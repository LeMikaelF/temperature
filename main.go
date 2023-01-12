package main

import (
	"github.com/LeMikaelF/temperature/metrics"
	"github.com/LeMikaelF/temperature/server"
	"net/http"
)

func main() {
	metricsChan := make(chan float64)
	exporter := metrics.NewExporter(metricsChan)
	metricsServer := server.NewServer(metricsChan, exporter.LastPushed)
	go exporter.Start()

	if err := http.ListenAndServe(":8080", http.HandlerFunc(metricsServer.Serve)); err != nil {
		panic(err)
	}
}
