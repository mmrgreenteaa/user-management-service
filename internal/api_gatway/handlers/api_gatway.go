package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mmrgreenteaa/user-management-service/internal/api_gatway/redis"
	authpb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/auth"
	usermepb "github.com/mmrgreenteaa/user-management-service/internal/gen/proto/user_manegement"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ApiGatway struct {
	UserMenedger usermepb.UserManagementClient
	AuthServis   authpb.AuthClient
	Rdb          *redis.DB
	logger       *slog.Logger
	Srv          *http.Server
}

type Metrics struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestDuration *prometheus.HistogramVec
}

type ApiGatwayConfig struct {
	Ip         string        `mapstructure:"listen_addr"`
	AuthIp     string        `mapstructure:"authIp"`
	UserMengIP string        `mapstructure:"user_managerIp"`
	Rcfg       redis.RCongif `mapstructure:"redis"`
}

func NewApiGatwaty(confg *ApiGatwayConfig) (*ApiGatway, error) {

	r := redis.Connect(&confg.Rcfg)

	aconn, err := grpc.NewClient(confg.AuthIp, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect auth server: %w", err)
	}
	aClient := authpb.NewAuthClient(aconn)

	uConn, err := grpc.NewClient(confg.UserMengIP, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect user menegment: %w", err)
	}
	uClient := usermepb.NewUserManagementClient(uConn)
	return &ApiGatway{
		AuthServis:   aClient,
		UserMenedger: uClient,
		logger:       slog.Default(),
		Rdb:          r,
	}, nil

}

func (apgt ApiGatway) DoMetrics(m *Metrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())

		path := c.FullPath()
		if path == "" {
			path = "not_found"
		}

		m.HttpRequestsTotal.WithLabelValues(c.Request.Method,  status).Inc()
		m.HttpRequestDuration.WithLabelValues(c.Request.Method, path).Observe(duration)
	}
}

func NewMenrics() *Metrics {

	HttpRequestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gin_http_requests_total",
			Help: "Общее количество обработанных HTTP запросов",
		},
		[]string{"method", "status"},
	)

	HttpRequestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gin_http_request_duration_seconds",
			Help:    "Время обработки запроса в секундах",
			Buckets: []float64{.01, .025, .05, .1, .25, .5, 1, 2.5, 5},
		},
		[]string{"method", "endpoint"},
	)

	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(HttpRequestsTotal)
	return &Metrics{
		HttpRequestsTotal:   HttpRequestsTotal,
		HttpRequestDuration: HttpRequestDuration,
	}
}
