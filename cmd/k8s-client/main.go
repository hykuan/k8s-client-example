package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/qeek-dev/quaistudio"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/hykuan/k8s-client-example"
	"github.com/hykuan/k8s-client-example/k8s-client"
	"github.com/hykuan/k8s-client-example/k8s-client/api"
	grpcapi "github.com/hykuan/k8s-client-example/k8s-client/api/grpc"
	httpapi "github.com/hykuan/k8s-client-example/k8s-client/api/http"
	"github.com/hykuan/k8s-client-example/logger"
)

const (
	defLogLevel   = "info"
	defHTTPPort   = "8180"
	defGRPCPort   = "8181"
	defSecret     = "users"
	defServerCert = ""
	defServerKey  = ""
	envLogLevel   = "QS_USERS_LOG_LEVEL"
	envHTTPPort   = "QS_USERS_HTTP_PORT"
	envGRPCPort   = "QS_USERS_GRPC_PORT"
	envSecret     = "QS_USERS_SECRET"
	envServerCert = "QS_USERS_SERVER_CERT"
	envServerKey  = "QS_USERS_SERVER_KEY"
)

type config struct {
	logLevel   string
	httpPort   string
	grpcPort   string
	secret     string
	serverCert string
	serverKey  string
}

func main() {
	cfg := loadConfig()

	logger, err := logger.New(os.Stdout, cfg.logLevel)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	svc := newService(clientset, logger)
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
	}
}

func newService(clientSet *kubernetes.Clientset, logger logger.Logger) k8s_client.Service {
	svc := k8s_client.New(clientSet)
	svc = api.LoggingMiddleware(svc, logger)
	svc = api.MetricsMiddleware(
		svc,
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "users",
			Subsystem: "api",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"}),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "users",
			Subsystem: "api",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, []string{"method"}),
	)
	return svc
}

func startHTTPServer(svc k8s_client.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	if certFile != "" || keyFile != "" {
		logger.Info(fmt.Sprintf("k8s-client service started using https, cert %s key %s, exposed port %s", certFile, keyFile, port))
		errs <- http.ListenAndServeTLS(p, certFile, keyFile, httpapi.MakeHandler(svc, logger))
	} else {
		logger.Info(fmt.Sprintf("k8s-client service started using http, exposed port %s", port))
		errs <- http.ListenAndServe(p, httpapi.MakeHandler(svc, logger))
	}
}

func startGRPCServer(svc k8s_client.Service, port string, certFile string, keyFile string, logger logger.Logger, errs chan error) {
	p := fmt.Sprintf(":%s", port)
	listener, err := net.Listen("tcp", p)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to listen on port %s: %s", port, err))
	}

	var server *grpc.Server
	if certFile != "" || keyFile != "" {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to load users certificates: %s", err))
			os.Exit(1)
		}
		logger.Info(fmt.Sprintf("k8s-client gRPC service started using https on port %s with cert %s key %s", port, certFile, keyFile))
		server = grpc.NewServer(grpc.Creds(creds))
	} else {
		logger.Info(fmt.Sprintf("k8s-client gRPC service started using http on port %s", port))
		server = grpc.NewServer()
	}

	quai.RegisterK8SClientServiceServer(server, grpcapi.NewServer(svc))
	logger.Info(fmt.Sprintf("k8s-client gRPC service started, exposed port %s", port))
	errs <- server.Serve(listener)
}
