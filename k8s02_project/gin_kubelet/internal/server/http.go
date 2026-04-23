package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	httpmw "gin_kubelet/internal/middleware/http"
	"gin_kubelet/pkg/app"
)

// HTTPServer HTTP 服务器
type HTTPServer struct {
	server *http.Server
	app    *app.App
}

// NewHTTPServer 创建 HTTP 服务器实例
func NewHTTPServer(application *app.App) *HTTPServer {
	return &HTTPServer{
		app: application,
	}
}

// Start 启动 HTTP 服务器
func (s *HTTPServer) Start() error {
	cfg := s.app.Config

	// 设置 Gin 模式
	if !cfg.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建 Gin 引擎
	r := gin.New()

	// 注册中间件
	r.Use(httpmw.Recovery())

	r.Use(httpmw.Logger())

	// 启用 pprof
	if cfg.App.Debug {
		pprof.Register(r)
	}

	// 注册健康检查路由
	s.registerHealthRoutes(r)

	// 注册业务路由
	s.registerRoutes(r)

	// 创建 HTTP 服务器
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	slog.Info("HTTP 服务器启动", "port", cfg.HTTP.Port)

	// 启动服务器
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP 服务器启动失败: %w", err)
	}

	return nil
}

// Stop 停止 HTTP 服务器
func (s *HTTPServer) Stop(ctx context.Context) error {
	slog.Info("正在关闭 HTTP 服务器...")
	if s.server != nil {
		return s.server.Shutdown(ctx)
	}
	return nil
}

// registerHealthRoutes 注册健康检查路由
func (s *HTTPServer) registerHealthRoutes(r *gin.Engine) {

	// 完整健康检查（包含详细信息）
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})
}

// registerRoutes 注册业务路由
func (s *HTTPServer) registerRoutes(r *gin.Engine) {

}
