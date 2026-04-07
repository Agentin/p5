package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"github.com/student/tech-ip-sem2/services/tasks/internal/grpcclient"
	"github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers"
	authMiddleware "github.com/student/tech-ip-sem2/services/tasks/internal/http/handlers/middleware"
	metricsMW "github.com/student/tech-ip-sem2/services/tasks/internal/http/middleware"
	"github.com/student/tech-ip-sem2/services/tasks/internal/repository"
	"github.com/student/tech-ip-sem2/shared/middleware"
)

func NewRouter(repo repository.TaskRepository, authClient *grpcclient.AuthClient, logger *zap.Logger) http.Handler {
	mux := http.NewServeMux()

	protected := authMiddleware.AuthMiddleware(authClient)

	// Основные CRUD
	mux.Handle("POST /v1/tasks", metricsMW.MetricsMiddleware("/v1/tasks")(protected(http.HandlerFunc(handlers.CreateTaskHandler(repo)))))
	mux.Handle("GET /v1/tasks", metricsMW.MetricsMiddleware("/v1/tasks")(protected(http.HandlerFunc(handlers.GetTasksHandler(repo)))))
	mux.Handle("GET /v1/tasks/{id}", metricsMW.MetricsMiddleware("/v1/tasks/:id")(protected(http.HandlerFunc(handlers.GetTaskHandler(repo)))))
	mux.Handle("PATCH /v1/tasks/{id}", metricsMW.MetricsMiddleware("/v1/tasks/:id")(protected(http.HandlerFunc(handlers.UpdateTaskHandler(repo)))))
	mux.Handle("DELETE /v1/tasks/{id}", metricsMW.MetricsMiddleware("/v1/tasks/:id")(protected(http.HandlerFunc(handlers.DeleteTaskHandler(repo)))))

	// Поиск (тоже защищён)
	mux.Handle("GET /v1/tasks/search", metricsMW.MetricsMiddleware("/v1/tasks/search")(protected(http.HandlerFunc(handlers.SearchTasksHandler(repo)))))

	// Метрики
	mux.Handle("GET /metrics", promhttp.Handler())

	handler := middleware.RequestIDMiddleware(mux)
	handler = middleware.HTTPAccessLogMiddleware(logger)(handler)
	return handler
}
