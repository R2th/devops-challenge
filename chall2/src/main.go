package main

import (
	"chall2/exporter-sd/config"
	"chall2/exporter-sd/probe"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	versioncollector "github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var (
	sc            = config.NewSafeConfig(prometheus.DefaultRegisterer)
	routePrefix   = kingpin.Flag("web.route-prefix", "").PlaceHolder("<path>").String()
	timeoutOffset = kingpin.Flag("timeout-offset", "").Default("0.5").Float64()
	externalURL   = kingpin.Flag("web.external-url", "").PlaceHolder("<url>").String()
	toolkitFlags  = webflag.AddFlags(kingpin.CommandLine, ":9115")

	moduleUnknownCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "chall2_module_unknown_total",
		Help: "Count of unknown modules requested by probes",
	})
)

func init() {
	prometheus.MustRegister(versioncollector.NewCollector("chall2_exporter"))
}

func main() {
	os.Exit(run())
}

func run() int {

	logger := promlog.New(&promlog.Config{})

	beURL, err := computeExternalURL(*externalURL, ":9915")
	if err != nil {
		panic(err)
	}

	// Default -web.route-prefix to path of -web.external-url.
	if *routePrefix == "" {
		*routePrefix = beURL.Path
	}

	// routePrefix must always be at least '/'.
	*routePrefix = "/" + strings.Trim(*routePrefix, "/")
	// routePrefix requires path to have trailing "/" in order
	// for browsers to interpret the path-relative path correctly, instead of stripping it.
	if *routePrefix != "/" {
		*routePrefix = *routePrefix + "/"
	}

	http.HandleFunc(path.Join(*routePrefix, "/probe"), func(w http.ResponseWriter, r *http.Request) {
		sc.Lock()
		conf := sc.C
		sc.Unlock()
		probe.Handler(w, r, conf, *timeoutOffset, nil, moduleUnknownCounter)
	})

	srv := &http.Server{}
	srvc := make(chan struct{})
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
			// level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
			close(srvc)
		}
	}()

	for {
		select {
		case <-term:
			// level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			return 0
		case <-srvc:
			return 1
		}
	}
}

func startsOrEndsWithQuote(s string) bool {
	return strings.HasPrefix(s, "\"") || strings.HasPrefix(s, "'") ||
		strings.HasSuffix(s, "\"") || strings.HasSuffix(s, "'")
}

// computeExternalURL computes a sanitized external URL from a raw input. It infers unset
// URL parts from the OS and the given listen address.
func computeExternalURL(u, listenAddr string) (*url.URL, error) {
	if u == "" {
		hostname, err := os.Hostname()
		if err != nil {
			return nil, err
		}
		_, port, err := net.SplitHostPort(listenAddr)
		if err != nil {
			return nil, err
		}
		u = fmt.Sprintf("http://%s:%s/", hostname, port)
	}

	fmt.Printf("%v\n", u)

	if startsOrEndsWithQuote(u) {
		return nil, errors.New("URL must not begin or end with quotes")
	}

	eu, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	ppref := strings.TrimRight(eu.Path, "/")
	if ppref != "" && !strings.HasPrefix(ppref, "/") {
		ppref = "/" + ppref
	}
	eu.Path = ppref

	return eu, nil
}
