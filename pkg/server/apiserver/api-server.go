package apiserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/linkinghack/gateway-controller/config"
	"github.com/linkinghack/gateway-controller/pkg/log"
)

type ControlPlaneAPIServer struct {
	listenAddress string
	ginRouter     *gin.Engine
	httpServer    *http.Server
}

func NewGWControlAPIServer(config *config.ControlPlaneAPIServerConfig) *ControlPlaneAPIServer {
	gin.SetMode(config.GinMode)
	if config.GinMode == "debug" {
		gin.ForceConsoleColor()
	}
	gin.DefaultWriter = io.MultiWriter(os.Stdout) // TODO add gin log to other log collector

	ginEngine := gin.New()
	ginEngine.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	ginEngine.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	ginEngine.Use(cors.New(cors.Config{
		AllowOrigins:     config.CorsOrigins,
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	provisionerAPIServer := ControlPlaneAPIServer{
		ginRouter:     ginEngine,
		listenAddress: config.ListenAddr,
		httpServer: &http.Server{
			Addr:    config.ListenAddr,
			Handler: ginEngine,
		},
	}

	return &provisionerAPIServer
}

func (s *ControlPlaneAPIServer) Start() error {
	// 1. register all http routes to the gin router
	s.registerRoutes()

	fmt.Printf("Provisioner HTTP-APIServer trying listen on %s\n", s.listenAddress)
	return s.httpServer.ListenAndServe()
}

func (s *ControlPlaneAPIServer) Stop(ctx context.Context) error {
	logger := log.GetSpecificLogger("ProvisionerAPIServer.Stop")
	logger.Infoln("Shutting down ProvisionerAPIServer http server")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		logger.WithField("err", err.Error()).Error("Error shutting down ProvisionerAPIServer.httpServer")
		return err
	}
	return nil
}
