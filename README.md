APOD to MQTT
---

This is a simple application that uses the [APOD API Wrapper](https://github.com/helixspiral/apod) to get the Astronomy Picture of the Day from NASA, and then publishes that to an MQTT topic.

Usage
---

This application is largely environment agnostic and relies on environment variables to determine where it sends data.

The variables used are:

* MQTT_BROKER - The broker for your MQTT server
* MQTT_CLIENT_ID - The client ID to use when connecting to your MQTT server
* MQTT_TOPIC - The MQTT topic to publish to
* APOD_API_KEY - The API key to use for the NASA APOD API