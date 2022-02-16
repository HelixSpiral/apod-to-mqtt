package main

import (
	"encoding/json"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/helixspiral/apod"
)

func main() {
	// Some initial setup with our environment
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	mqttTopic := os.Getenv("MQTT_TOPIC")
	apodApiKey := os.Getenv("APOD_API_KEY")

	// We always use EST since that's the timezone APOD updates in.
	locEST, err := time.LoadLocation("EST")
	if err != nil {
		panic(err)
	}

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Setup the Apod service
	Apod := apod.NewAPOD(apodApiKey)

	// Query the api for the picture of the day
	resp, err := Apod.Query(&apod.ApodQueryInput{
		Date: time.Now().In(locEST),
	})
	if err != nil {
		panic(err)
	}

	// Some initial setup for our MQTT client
	options := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientId)
	options.WriteTimeout = 20 * time.Second
	mqttClient := mqtt.NewClient(options)

	// Connect to the MQTT broker
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Jsonify our message so we have a []byte
	jsonMsg, err := json.Marshal(resp[0])
	if err != nil {
		panic(err)
	}

	// Send it and wait
	token := mqttClient.Publish(mqttTopic, 2, false, jsonMsg)
	_ = token.Wait()
	if token.Error() != nil {
		panic(err)
	}

	mqttClient.Disconnect(250)
}
