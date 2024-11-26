package main

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/helixspiral/apod"
)

func init() {
	logLevel := &slog.LevelVar{}
	if logLevelEnv := os.Getenv("LOG_LEVEL"); logLevelEnv != "" {
		switch logLevelEnv {
		case "ERROR":
			logLevel.Set(slog.LevelError)
		case "WARN":
			logLevel.Set(slog.LevelWarn)
		case "INFO":
			logLevel.Set(slog.LevelInfo)
		case "DEBUG":
			logLevel.Set(slog.LevelDebug)
		default:
			logLevel.Set(slog.LevelInfo)
		}
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// Replace the default logger
	slog.SetDefault(logger)
}

func main() {
	// Some initial setup with our environment
	mqttBroker := os.Getenv("MQTT_BROKER")
	mqttClientId := os.Getenv("MQTT_CLIENT_ID")
	mqttTopic := os.Getenv("MQTT_TOPIC")
	apodApiKey := os.Getenv("APOD_API_KEY")

	slog.Info("Starting the Astronomy Picture of the Day to MQTT service", "mqttBroker", mqttBroker, "mqttClientId", mqttClientId, "mqttTopic", mqttTopic)

	// We always use EST since that's the timezone APOD updates in.
	locEST, err := time.LoadLocation("EST")
	if err != nil {
		slog.Error("Failed to load timezone information", "error", err)

		os.Exit(1)
	}

	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)

	// Setup the Apod service
	Apod := apod.NewAPOD(apodApiKey)

	// Query the api for the picture of the day
	resp, err := Apod.Query(&apod.ApodQueryInput{
		Date:   time.Now().In(locEST),
		Thumbs: true,
	})
	if err != nil {
		slog.Error("Failed to query the APOD API", "error", err)

		os.Exit(1)
	}

	slog.Info("Successfully queried the APOD API", "title", resp[0].Title, "date", resp[0].Date)

	// Some initial setup for our MQTT client
	options := mqtt.NewClientOptions().AddBroker(mqttBroker).SetClientID(mqttClientId)
	options.WriteTimeout = 20 * time.Second
	mqttClient := mqtt.NewClient(options)

	// Connect to the MQTT broker
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		slog.Error("Failed to connect to the MQTT broker", "error", token.Error())

		os.Exit(1)
	}

	// Jsonify our message so we have a []byte
	jsonMsg, err := json.Marshal(resp[0])
	if err != nil {
		slog.Error("Failed to marshal the APOD response", "error", err)

		os.Exit(1)
	}

	// Send it and wait
	token := mqttClient.Publish(mqttTopic, 2, false, jsonMsg)
	_ = token.Wait()
	if token.Error() != nil {
		panic(token.Error())
	}

	mqttClient.Disconnect(250)
}
