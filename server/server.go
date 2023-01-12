package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	metrics         chan<- float64
	initMetricState func() float64
}

type temperatureStruct struct {
	Temperature float64 `json:"temperature"`
}

func NewServer(metrics chan<- float64, initMetricState func() float64) *Server {
	return &Server{metrics: metrics, initMetricState: initMetricState}
}

func (s Server) Serve(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/temperature" && r.Method == http.MethodGet {
		s.handleGetTemperature(w, r)
	} else if r.URL.Path == "/temperature" && r.Method == http.MethodPost {
		s.handlePostTemperature(w, r)
	} else if r.URL.Path == "/" && r.Method == http.MethodGet {
		http.ServeFile(w, r, "server/index.html")
	} else if r.URL.Path == "/favicon.ico" {
		//no-op
	} else {
		log.Printf("received unexpected request to path %s", r.URL.Path)
	}
}

func (s Server) handlePostTemperature(w http.ResponseWriter, r *http.Request) {
	var temperature temperatureStruct
	if err := json.NewDecoder(r.Body).Decode(&temperature); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("receiving metric %f", temperature.Temperature)
	s.metric(temperature.Temperature)
}

func (s Server) metric(metric float64) {
	s.metrics <- metric
}

func (s Server) handleGetTemperature(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(temperatureStruct{s.initMetricState()})
	if err != nil {
		log.Printf("handleGetTemperature marshaling temperature: %v", err)
		return
	}

	_, err = w.Write(payload)
	if err != nil {
		log.Printf("handleGetTemperature writing response body: %v", err)
	}
}
