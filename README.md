# mhb
MQTT HTTP bridge. Intended to be compatible with Grafana streaming datasource, allowing any MQTT client to supply a data stream to Grafana.

Inspired by:
https://github.com/slim-bean/leafbus from Grafaconline 2020

https://github.com/ryantxu/streaming-json-datasource

## to-do
* Use msgpack or similar to carry flexible data types over MQTT, to work with decoded CAN bus data.
* Allow HTTP connection url/parameters to choose subscribed topic.
* Allow setting MQTT options via REST-ish API?
* Actually check it works with Grafana https://github.com/ryantxu/streaming-json-datasource
* MQTT TLS & https? (http should be via local connection or in stack/compose only right now)
