package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/qeek-dev/quaistudio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/hykuan/k8s-client-example"
	k8sapi "github.com/hykuan/k8s-client-example/k8s-client/api/grpc"
	"github.com/hykuan/k8s-client-example/logger"
	"github.com/hykuan/k8s-client-example/models"
	"github.com/hykuan/k8s-client-example/models/api"
	grpcapi "github.com/hykuan/k8s-client-example/models/api/grpc"
	httpapi "github.com/hykuan/k8s-client-example/models/api/http"
)

const (
	defLogLevel   = "info"
	defHTTPPort   = "8182"
	defGRPCPort   = "8183"
	defSecret     = "users"
	defServerCert = ""
	defServerKey  = ""
	defK8sUrl     = "localhost:8181"
	envLogLevel   = "QS_MODELS_LOG_LEVEL"
	envHTTPPort   = "QS_MODELS_HTTP_PORT"
	envGRPCPort   = "QS_MODELS_GRPC_PORT"
	envSecret     = "QS_MODELS_SECRET"
	envServerCert = "QS_MODELS_SERVER_CERT"
	envServerKey  = "QS_MODELS_SERVER_KEY"
	envK8sUrl     = "QS_K8S_URL"
)

type config struct {
	logLevel   string
	httpPort   string
	grpcPort   string
	secret     string
	serverCert string
	serverKey  string
	k8sUrl     string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	conn := connectToK8sService(cfg.k8sUrl, logger)
	defer conn.Close()

	svc := newService(conn, logger)
	errs := make(chan error, 2)

	go startHTTPServer(svc, cfg.httpPort, cfg.serverCert, cfg.serverKey, logger, errs)
	go startGRPCServer(svc, cfg.grpcPort, cfg.serverCert, cfg.serverKey, logger, errs)

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	err = <-errs
	logger.Error(fmt.Sprintf("k8s-client service terminated: %s", err))
}

func loadConfig() config {
	return config{
		logLevel:   quaistudio.Env(envLogLevel, defLogLevel),
		httpPort:   quaistudio.Env(envHTTPPort, defHTTPPort),
		grpcPort:   quaistudio.Env(envGRPCPort, defGRPCPort),
		secret:     quaistudio.Env(envSecret, defSecret),
		serverCert: quaistudio.Env(envServerCert, defServerCert),
		serverKey:  quaistudio.Env(envServerKey, defServerKey),
		k8sUrl:     quaistudio.Env(envK8sUrl, defK8sUrl),
	}
}

func connectToK8sService(k8sAddr string, logger logger.Logger) *grpc.ClientConn {
	conn, err := grpc.Dial(k8sAddr, grpc.WithInsecure())
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to k8s service: %s", err))
		os.Exit(1)
	}
	return conn
}

func newService(conn *grpc.ClientConn, logger logger.Logger) models.Service {
	k8sClient := k8sapi.NewClient(conn)

	svc := models.New(k8sClient)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "models",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "models",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(svc models.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if certFile != "" || keyFile != "" {
		logger.Info(fmt.Sprintf("models service started using https, cert %s key %s, exposed port %s", certFile, keyFile, port))
		errs <- http.ListenAndServeTLS(p, certFile, keyFile, httpapi.MakeHandler(svc, logger))
	} else {
		logger.Info(fmt.Sprintf("models service started using http, exposed port %s", port))
		errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, logger))
	}
}

func startGRPCServer(svc models.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	var server *grpc.Server
	if certFile != "" || keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to load models certificates: %s", err))
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("models gRPC service started using https on port %s with cert %s key %s", port, certFile, keyFile))
		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		logger.Info(fmt.Sprintf("models gRPC service started using http on port %s", port))
		server = grpc.NewServer()
	}

	quai.RegisterModelServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("models gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
