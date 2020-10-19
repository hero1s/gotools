package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hero1s/gotools/limit"
	"github.com/hero1s/gotools/log"
	"time"
)

type Engine struct {
	Gin  *gin.Engine
	port string
}

type Config struct {
	// http export port. :8080
	Port string

	// interface limit
	Limit *limit.Config
}

func InitServer(c *Config) *Engine {
	g := gin.Default()
	g.Use(cors())
	g.Use(gin.Recovery())
	if c.Limit != nil && c.Limit.Rate != 0 {
		g.Use(limit.NewLimiter(c.Limit).Limit())
	}
	engine := &Engine{Gin: g, port: c.Port}
	if !strings.Contains(strings.TrimSpace(c.Port), ":") {
		engine.port = ":" + c.Port
	}
	return engine
}

func (e *Engine) Start() {
	go func() {
		if err := e.Gin.Run(e.port); err != nil {
			panic(fmt.Sprintf("web server port(%s) run error(%+v).", e.port, err))
		}
	}()
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") // 请求头部
		if origin == "" {
			origin = c.Request.Host
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, X-CSRF-Token, authorization, sign, appid, ts")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			c.Header("Access-Control-Max-Age", "172800")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

// 日志记录到文件
func LogReq() gin.HandlerFunc {
	return func(c * gin.Context){
		// 开始时间
		startTime := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		endTime := time.Now()
		// 执行时间
		latencyTime := endTime.Sub(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// 请求路由
		reqUri := c.Request.RequestURI
		// 状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		// 日志格式
		log.Info("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}