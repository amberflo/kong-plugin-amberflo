package main

import (
	"log"
	"strings"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/amberflo/metering-go/v2"
	"github.com/xtgo/uuid"
)

// See https://docs.konghq.com/gateway/latest/plugin-development/pluginserver/go/
func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		log.Printf("Error starting plugin server: %s", err.Error())
	}

	// Handle graceful shutdown of metering client
	if client != nil {
		err := client.Shutdown()
		if err != nil {
			log.Printf("Error shutting down metering client: %s", err.Error())
		}
	}
}

const Version = "0.2"
const Priority = 10

type Config struct {
	ApiKey           string            `json:"apiKey"`           // required
	MeterApiName     string            `json:"meterApiName"`     // required
	CustomerHeader   string            `json:"customerHeader"`   // required; get the customer id from this header
	IntervalSeconds  int               `json:"intervalSeconds"`  // send the meter record batch every x seconds
	BatchSize        int               `json:"batchSize"`        // send the meter record batch when it reaches this size
	Debug            bool              `json:"debug"`            // passed to the amberflo metering client
	MethodDimension  string            `json:"methodDimension"`  // name of the dimension for the request method
	HostDimension    string            `json:"hostDimension"`    // name of the dimension for the target url host
	RouteDimension   string            `json:"routeDimension"`   // name of the dimension for the route name
	ServiceDimension string            `json:"serviceDimension"` // name of the dimension for the service name
	DimensionHeaders map[string]string `json:"dimensionHeaders"` // map of "dimension name" to "header name", for inclusion in the meter record
	Replacements     map[string]string `json:"replacements"`     // map of "old" to "new" values for transforming dimension values
}

func New() interface{} {
	return &Config{
		IntervalSeconds: 1,
		BatchSize:       10,
		Debug:           false,
		Replacements:    map[string]string{"/": ":"},
	}
}

// Kong plugin method, see:
// - https://docs.konghq.com/gateway/3.0.x/plugin-development/pluginserver/go/#phase-handlers
// - https://docs.konghq.com/gateway/3.0.x/plugin-development/custom-logic/#available-contexts
func (conf Config) Access(kong *pdk.PDK) {

	// Get the customer id. Abort if there is no id.
	customerId := getHeader(kong, conf.CustomerHeader)
	if customerId == "" {
		log.Print("ERROR: missing customer id")
		return
	}

	// The plugin currently assumes a fixed meterApiName
	meterApiName := conf.MeterApiName

	dimensions := make(map[string]string)

	// Get dimensions from request headers, if configured.
	for name, header := range conf.DimensionHeaders {
		value := getHeader(kong, header)
		if value != "" {
			// Apply Config.Replacements on the header value, so it becomes a
			// valid dimension value.
			dimensions[name] = apply(conf.Replacements, value)
		}
	}

	// Add the request method as a dimension, if configured
	if conf.MethodDimension != "" {
		method, err := kong.Request.GetMethod()
		if err != nil {
			log.Printf("ERROR getting request method: %s", err.Error())
		} else if method != "" {
			dimensions[conf.MethodDimension] = method
		}
	}

	// Add the url host as a dimension, if configured
	if conf.HostDimension != "" {
		host, err := kong.Request.GetHost()
		if err != nil {
			log.Printf("ERROR getting request url host: %s", err.Error())
		} else if host != "" {
			dimensions[conf.HostDimension] = host
		}
	}

	// Add the service name as a dimension, if configured
	if conf.ServiceDimension != "" {
		service, err := kong.Router.GetService()
		if err != nil {
			log.Printf("ERROR getting request service: %s", err.Error())
		} else {
			dimensions[conf.ServiceDimension] = service.Name
		}
	}

	// Add the route name as a dimension, if configured
	if conf.RouteDimension != "" {
		route, err := kong.Router.GetRoute()
		if err != nil {
			log.Printf("ERROR getting request route: %s", err.Error())
		} else {
			dimensions[conf.RouteDimension] = route.Name
		}
	}

	// Meter the request
	conf.meter(customerId, meterApiName, dimensions)
}

// Meter a request. This will actually queue the meter record for asynchronous
// batching.
func (conf *Config) meter(customerId string, meterApiName string, dimensions map[string]string) {
	client := conf.getClient()

	err := client.Meter(&metering.MeterMessage{
		UniqueId:          uuid.NewRandom().String(),
		CustomerId:        customerId,
		MeterApiName:      meterApiName,
		MeterTimeInMillis: time.Now().UnixMilli(),
		MeterValue:        1,
		Dimensions:        dimensions,
	})
	if err != nil {
		log.Printf("ERROR metering request: %s", err.Error())
	}
}

// We'll use the same metering client for all requests, to benefit of its
// batching behavior.
var client *metering.Metering

// Get the metering client, initializing it if necessary.
func (conf *Config) getClient() *metering.Metering {

	if client == nil {
		intervalSeconds := time.Second * time.Duration(conf.IntervalSeconds)
		batchSize := conf.BatchSize
		debug := conf.Debug

		client = metering.NewMeteringClient(
			conf.ApiKey,
			metering.WithBatchSize(batchSize),
			metering.WithIntervalSeconds(intervalSeconds),
			metering.WithDebug(debug),
		)
	}

	return client
}

// Helper method. Get a request header, log errors. Returns a possibly empty string.
func getHeader(kong *pdk.PDK, header string) string {
	value, err := kong.Request.GetHeader(header)
	if err != nil {
		log.Printf("ERROR reading header '%s': %s", header, err.Error())
		return ""
	}
	return value
}

// Helper method. Apply the Config.Replacements on a dimension value.
func apply(replacements map[string]string, value string) string {
	for old, new := range replacements {
		value = strings.ReplaceAll(value, old, new)
	}
	return value
}
