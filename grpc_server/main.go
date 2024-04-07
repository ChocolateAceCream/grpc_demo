package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sms/proto"
	"sms/service"
	"sms/transport"
	"syscall"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	logger := log.NewJSONLogger(os.Stdout)
	logger = log.With(logger, "timestamp", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	smsService := service.NewSMSService(logger)
	smsEndpoint := transport.MakeSMSEndpoint(smsService)
	gRPCServer := transport.NewGRPCServer(smsEndpoint, logger)
	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	gRPCListener, err := net.Listen("tcp", "localhost:9085")
	if err != nil {
		logger.Log("during", "listen", "err", err)
	}

	certificate, err := tls.LoadX509KeyPair("./certs2/server.crt", "./certs2/server.key")
	if err != nil {
		logger.Log("failed to load key pair: %s", err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("./certs2/ca.crt")
	if err != nil {
		logger.Log("could not read ca certificate: %s", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		logger.Log("failed to append ca certs", ok)
	}
	opts := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})),
	}

	go func() {
		// s := grpc.NewServer(grpc.Creds(creds))
		s := grpc.NewServer(opts...)
		proto.RegisterSmsServiceServer(s, gRPCServer)
		level.Info(logger).Log("msg", "Starting gRPC server at :9085")
		s.Serve(gRPCListener)
	}()

	level.Error(logger).Log("exit", <-errs)

}
