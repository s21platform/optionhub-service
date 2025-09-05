package infra

import (
	"context"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"

	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/optionhub-service/internal/config"
)

func MetricsInterceptor(metrics *pkg.Metrics) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler) (
		interface{}, error) {
		t := time.Now()
		method := strings.Join(strings.Split(strings.Trim(info.FullMethod, "/"), "/")[1:], "")
		metrics.Increment(method)

		ctx = context.WithValue(ctx, config.KeyMetrics, metrics)
		resp, err := handler(ctx, req)

		if err != nil {
			metrics.Increment(method + "_error")
		}

		metrics.Duration(time.Since(t).Milliseconds(), method)

		return resp, err
	}
}

func MetricsRequest(metrics *pkg.Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			metrics.Increment(r.Method)
			ctx = context.WithValue(ctx, config.KeyMetrics, metrics)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}