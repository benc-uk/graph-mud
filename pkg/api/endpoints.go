// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Optional extra endpoints you may want to add to your API
// ----------------------------------------------------------------------------

package api

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/elastic/go-sysinfo"
	"github.com/go-chi/chi/v5"
	metrics "github.com/m8as/go-chi-metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Status struct {
	Service    string `json:"service"`
	Healthy    bool   `json:"healthy"`
	Version    string `json:"version"`
	BuildInfo  string `json:"buildInfo"`
	Hostname   string `json:"hostname"`
	OS         string `json:"os"`
	Arch       string `json:"architecture"`
	CPU        int    `json:"cpuCount"`
	GoVersion  string `json:"goVersion"`
	ClientAddr string `json:"clientAddress"`
	ServerHost string `json:"serverHost"`
	Uptime     string `json:"uptime"`
}

// AddOKEndpoint adds an endpoint that respond 200 when hitting it
func (b *Base) AddOKEndpoint(r chi.Router, path string) {
	log.Printf("### 🍏 API: 200 OK endpoint at: %s", "/"+path)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})
}

// AddMetrics adds Prometheus metrics to the router
func (b *Base) AddMetricsEndpoint(r chi.Router, path string) {
	log.Printf("### 🔬 API: metrics endpoint at: %s", "/"+path)

	r.Use(metrics.SetRequestDuration)
	r.Use(metrics.IncRequestCount)
	r.Handle("/"+path, promhttp.Handler())
}

// AddHealth adds a health check endpoint to the API
func (b *Base) AddHealthEndpoint(r chi.Router, path string) {
	log.Printf("### 💚 API: health endpoint at: %s", "/"+path)

	r.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		if b.Healthy {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK: Service is healthy"))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte("Error: Service is not healthy"))
		}
	})
}

// AddStatus adds a status & info endpoint to the API
func (b *Base) AddStatusEndpoint(r chi.Router, path string) {
	log.Printf("### 🔮 API: status endpoint at: %s", "/"+path)

	r.HandleFunc("/"+path, func(w http.ResponseWriter, r *http.Request) {
		host, _ := sysinfo.Host()
		host.Info().Uptime()

		status := Status{
			Service:    b.ServiceName,
			Healthy:    b.Healthy,
			Version:    b.Version,
			BuildInfo:  b.BuildInfo,
			Hostname:   host.Info().Hostname,
			GoVersion:  runtime.Version(),
			OS:         runtime.GOOS,
			Arch:       runtime.GOARCH,
			CPU:        runtime.NumCPU(),
			ClientAddr: r.RemoteAddr,
			ServerHost: r.Host,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(status)
	})
}
