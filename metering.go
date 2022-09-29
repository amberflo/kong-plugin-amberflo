package main

import (
	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/amberflo/metering-go"
	"github.com/xtgo/uuid"
	"log"
	"time"
)

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		log.Printf("Error starting plugin server: %s", err.Error())
	}
}

const Version = "0.1"
const Priority = 10

type Config struct {
	ApiKey         string `json:"apiKey"`
	MeterApiName   string `json:"meterApiName"`
	CustomerHeader string `json:"customerHeader"`
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	customerId, err := kong.Request.GetHeader(conf.CustomerHeader)
	if err != nil {
		log.Printf("Error reading header: '%s', %s", conf.CustomerHeader, err.Error())
		return
	}

	intervalSeconds := 30 * time.Second
	batchSize := 5
	debug := true

	meteringClient := metering.NewMeteringClient(
		conf.ApiKey,
		metering.WithBatchSize(batchSize),
		metering.WithIntervalSeconds(intervalSeconds),
		metering.WithDebug(debug),
	)

	meteringErr := meteringClient.Meter(&metering.MeterMessage{
		UniqueId:          uuid.NewRandom().String(),
		CustomerId:        customerId,
		MeterApiName:      conf.MeterApiName,
		MeterTimeInMillis: time.Now().UnixMilli(),
		MeterValue:        1,
	})
	if err != nil {
		log.Printf("Error metering request: %s", meteringErr.Error())
	}

	shutdownErr := meteringClient.Shutdown()
	if shutdownErr != nil {
		log.Printf("Error during shutdown: %s", shutdownErr.Error())
	}
}
