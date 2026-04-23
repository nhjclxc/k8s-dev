package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Recovery 恢复中间件，捕获 panic
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取堆栈信息
				stack := projectStack(fmt.Sprintf("%v", err))

				// 记录日志
				slog.Error("发生 panic",
					"error", err,
					"stack", stack,
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				// 如果启用了 tracing，记录错误到 span
				span := trace.SpanFromContext(c.Request.Context())
				if span.SpanContext().IsValid() {
					span.SetStatus(codes.Error, "panic recovered")
					span.RecordError(err.(error))
				}

				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "服务器内部错误",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

func projectStack(err string) string {
	fullStack := debug.Stack()
	lines := strings.Split(string(fullStack), "\n")

	var filtered []string
	filtered = append(filtered, "panic reason:"+err)
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// 只保留包含你项目路径的堆栈
		if strings.Contains(line, "gin_kubelet") && strings.Contains(line, ".go:") {
			// 上一行一般是函数签名，保留
			if i > 0 {
				filtered = append(filtered, lines[i-1])
			}
			filtered = append(filtered, line)
		}
	}

	return strings.Join(filtered, "\n")
}
