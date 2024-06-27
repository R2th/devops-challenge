package config

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/config"
)

type Config struct {
	Modules map[string]Module `yaml:"modules"`
}

type Module struct {
	Prober  string        `yaml:"prober,omitempty"`
	Timeout time.Duration `yaml:"timeout,omitempty"`
	RPC     RPCProbe      `yaml:"grpc,omitempty"`
}

type RPCProbe struct {
	Service             string           `yaml:"service,omitempty"`
	TLS                 bool             `yaml:"tls,omitempty"`
	TLSConfig           config.TLSConfig `yaml:"tls_config,omitempty"`
	IPProtocolFallback  bool             `yaml:"ip_protocol_fallback,omitempty"`
	PreferredIPProtocol string           `yaml:"preferred_ip_protocol,omitempty"`
}

type SafeConfig struct {
	sync.RWMutex
	C                   *Config
	configReloadSuccess prometheus.Gauge
	configReloadSeconds prometheus.Gauge
}

func NewSafeConfig(reg prometheus.Registerer) *SafeConfig {
	configReloadSuccess := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "chall2_exporter",
		Name:      "config_last_reload_successful",
		Help:      "Chall2 exporter config loaded successfully.",
	})

	configReloadSeconds := promauto.With(reg).NewGauge(prometheus.GaugeOpts{
		Namespace: "chall2_exporter",
		Name:      "config_last_reload_success_timestamp_seconds",
		Help:      "Timestamp of the last successful configuration reload.",
	})
	return &SafeConfig{C: &Config{}, configReloadSuccess: configReloadSuccess, configReloadSeconds: configReloadSeconds}
}
