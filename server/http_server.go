package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"taleteller/constants"
	"taleteller/logger"
	"taleteller/middleware"

	"github.com/gorilla/handlers"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	mdw "github.com/slok/go-http-metrics/middleware"
	middlewarestd "github.com/slok/go-http-metrics/middleware/std"
)

func startHTTPServer(dependencies Dependencies) (err error) {
	//port := config.AppPort()
	//metricsPort := config.MetricsPort()
	addr := fmt.Sprintf(":%s", strconv.Itoa(8001))
	//metricsAddr := fmt.Sprintf(":%s", strconv.Itoa(metricsPort))

	muxRouter := initRouter(dependencies)
	// metrics middleware
	metricsMdw := mdw.New(mdw.Config{
		Recorder: metrics.NewRecorder(metrics.Config{}),
	})

	headersOk := handlers.AllowedHeaders([]string{constants.ContentType, constants.Authorization, constants.VerificationToken, constants.UserAgent, constants.TimeZone})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	corsHandler := handlers.CORS(headersOk, originsOk, methodsOk)(muxRouter)
	logHandler := middleware.RequestLoggerHandler(corsHandler)
	handler := middlewarestd.Handler("", metricsMdw, logHandler)

	// Serve metrics.
	//logger.Infof(context.Background(), "serving metrics at: %s", metricsAddr)
	//go http.ListenAndServe(metricsAddr, promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	logger.Infof(context.TODO(), "starting API server on %s", addr)
	server.ListenAndServe()
	return
}
