package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	HeaderUserId        = "user_id"
	HeaderOperationType = "operation_type"
)

func main() {
	r := gin.Default()

	// Gauge registration
	lastRequestReceivedTime := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "last_request_received_time",
		Help: "Time when the last request was processed",
	}, []string{HeaderUserId, HeaderOperationType})
	err := prometheus.Register(lastRequestReceivedTime)
	handleErr(err)

	// Metrics handler
	r.GET("/metrics", func(c *gin.Context) {
		handler := promhttp.Handler()
		handler.ServeHTTP(c.Writer, c.Request)
	})

	// Middleware to set lastRequestReceivedTime for all requests
	middleware := func(context *gin.Context) {
		lastRequestReceivedTime.With(prometheus.Labels{
			HeaderUserId:        context.GetHeader(HeaderUserId),
			HeaderOperationType: context.GetHeader(HeaderOperationType),
		}).SetToCurrentTime()
	}

	// Request handler
	r.GET("/data", middleware, func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "success"})
	})

	err = r.Run(":8080")
	handleErr(err)
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
