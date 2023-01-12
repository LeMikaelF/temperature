package main

import (
	"github.com/LeMikaelF/temperature/metrics"
	"github.com/LeMikaelF/temperature/server"
	"net/http"
)

func main() {
	metricsChan := make(chan float64, 1000)
	exporter := metrics.NewExporter(metricsChan)
	metricsServer := server.NewServer(metricsChan, exporter.LastReceived)
	go func() {
		err := exporter.Start()
		if err != nil {
			panic(err)
		}
	}()

	if err := http.ListenAndServe(":8080", http.HandlerFunc(metricsServer.Serve)); err != nil {
		panic(err)
	}
}
