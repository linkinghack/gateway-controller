package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/linkinghack/gateway-controller/config"
	"github.com/linkinghack/gateway-controller/pkg/kit/daemonwaiter"
	"github.com/linkinghack/gateway-controller/pkg/server/apiserver"
	"github.com/linkinghack/gateway-controller/pkg/server/xds"
	"github.com/spf13/cobra"
)

// serverRootCmd represents the base command when called without any subcommands
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start to serve gateway control API",
	Long:  `Start to serve data plane gateways with dynamic registered configs`,
	Run:   startServe,
}

func init() {
	serveCmd.Flags().String("webConfig.listenAddr", "", "指定服务器监听地址，如 0.0.0.0:5500")
}

func startServe(cmd *cobra.Command, args []string) {
	conf := config.GetGlobalConfig()
	apiServer := apiserver.NewGWControlAPIServer(&conf.ServerConfig)
	dw := daemonwaiter.GetDefaultDaemonWaiter()

	// HTTP API server
	go func() {
		err := apiServer.Start()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// XDS Server
	xdsServer := xds.NewXdsServerWithGlobalConfig(dw.GetContext())
	dw.AddAndStart(xdsServer)

	// Start DaemonWaiter
	go dw.Run()

	// graceful shut down
	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 2)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so does not need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 7 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()

	// Stop DaemonWaiter
	dw.Stop()

	if err := apiServer.Stop(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
}
