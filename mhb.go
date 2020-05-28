package main

import (
	"net/http"
	"time"
	"fmt"
	"os"
	"log"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/eclipse/paho.mqtt.golang"
	"encoding/json"
)

type Data struct {
	Name      string  `json:"name"`
	//Timestamp int64   `json:"ts"`
	Val       float64 `json:"val"`
	//bytes     []byte
}

var ch chan Data = make(chan Data)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	val, _ := strconv.ParseFloat(string(msg.Payload()), 64)
	ch <- Data{Name: msg.Topic(), Val: val}
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://test.mosquitto.org:1883").SetClientID("mhbtest")
	opts.SetKeepAlive(30 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(3 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Println("starting router")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	cors := cors.New(cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	AllowCredentials: true,
	MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	fmt.Println("started router")
	r.Get("/t/{topic}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("request")
		topic := chi.URLParam(r, "topic")
		fmt.Println(r.RequestURI, r.Header)
		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		if token := c.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}		
		for {
			select {
			case d := <- ch:
				err := enc.Encode(d)
				fmt.Println(d)
				if err != nil {
					log.Println("Failed to marshal data object to json stream:", err)
				}
				//w.Write([]byte(d))
				flusher.Flush()
			case <-r.Context().Done():
				log.Println("Client connection closed")
				if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
					fmt.Println(token.Error())
					os.Exit(1)
				}
				return
			}
		}		
	})

	http.ListenAndServe(":3333", r)
}