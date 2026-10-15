package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func TestRegisterMetricsRequiresBasicAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	engine := gin.New()
	registerMetrics(engine, prometheus.NewRegistry(), "/metrics", "metrics", "secret")

	for _, method := range []string{http.MethodGet, http.MethodHead} {
		t.Run(method+" without credentials", func(t *testing.T) {
			request := httptest.NewRequest(method, "/metrics", nil)
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)
			if response.Code != http.StatusUnauthorized {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
			}
		})

		t.Run(method+" with credentials", func(t *testing.T) {
			request := httptest.NewRequest(method, "/metrics", nil)
			request.SetBasicAuth("metrics", "secret")
			response := httptest.NewRecorder()
			engine.ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}
		})
	}
}
