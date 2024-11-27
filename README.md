APOD to MQTT
---

This is a simple application that uses the [APOD API Wrapper](https://github.com/helixspiral/apod) to get the Astronomy Picture of the Day from NASA, and then publishes that to an MQTT topic.

This application is largely environment agnostic and relies on environment variables to determine where it sends data.

The variables used are:

* MQTT_BROKER - The broker for your MQTT server
* MQTT_CLIENT_ID - The client ID to use when connecting to your MQTT server
* MQTT_TOPIC - The MQTT topic to publish to
* APOD_API_KEY - The API key to use for the NASA APOD API
* APOD_API_DOMAIN - The domain of the API to use, defaults to the NASA provided API

Build with Docker
---

We use the Docker buildx feature to build multiple architectures: `docker buildx build --platform linux/amd64,linux/arm64 -t ghcr.io/helixspiral/apod-to-mqtt:latest .`

If all you need is your arch you can omit the platform specific stuff and just do a normal docker build.

Kubernetes setup
---

We've provided an example config map and cronjob that can be used to run this. You'll also need to create a k8s secret with the `APOD_API_KEY` in the same namespace.

If the image repo being used is private you'll also need to provide k8s with credentials. You can do this with a secret: `kubectl create secret docker-registry <name> --docker-server=<server> --docker-username=<username> --docker-email=<email> --docker-password=<password> -n <namespace>`

Usage
---

To run this you'll need to have an MQTT broker setup to pass messages to, and at least one service that can handle receiving the messages from MQTT.

If you don't have something to consume MQTT messages you can use the Discord bot provided here: https://github.com/HelixSpiral/apod-mqtt-to-discord