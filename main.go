package main

import (
	"flag"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var addr = flag.String("listen-address", ":9776", "The address to listen on for HTTP requests.")
var enableRuntimeMetrics = flag.Bool("runtime-metrics", false, "Enable prometheus runtime metrics.")
var logLevel = zap.LevelFlag("log-level", zap.WarnLevel, "The information level used for logging")
var host = flag.String("host", "", "The password used for logging into S-Miles Cloud.")
var port = flag.Int("port", 80, "The port .")

var log = initLog()

func initLog() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)
	config := zap.NewProductionEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	consoleEncoder := zapcore.NewConsoleEncoder(config)
	logCore := zapcore.NewCore(consoleEncoder, stdout, logLevel)
	logger := zap.New(logCore)

	return logger
}

func main() {
	defer log.Sync()

	flag.Parse()

	mainLog := log.Sugar().Named("main")

	if len(*host) == 0 {
		mainLog.Fatal("Host must not be empty.")
	}

	mainLog.Debug("Registering Prometheus metrics …")

	reg := prometheus.NewRegistry()

	reg.MustRegister(collectors.NewBuildInfoCollector())

	if *enableRuntimeMetrics {
		reg.MustRegister(collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.GoRuntimeMetricsRule{Matcher: regexp.MustCompile("/.*")}),
		))
	}

	reg.MustRegister(newMetrics())

	mainLog.Infof("Listening for HTTP requests on %s …", *addr)

	prometheusLog, _ := zap.NewStdLogAt(log.Named("prometheus"), zap.ErrorLevel)

	// Expose the registered metrics via HTTP.
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			ErrorLog: prometheusLog,
			// Opt into OpenMetrics to support exemplars.
			EnableOpenMetrics: true,
		},
	))
	mainLog.Fatal(http.ListenAndServe(*addr, nil))
}
