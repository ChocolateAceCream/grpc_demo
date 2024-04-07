package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	pb "grpc-client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cert, err := tls.LoadX509KeyPair("certs2/client.crt", "certs2/client.key")
	if err != nil {
		log.Fatalf("could not load client key pair: %s", err)
	}
	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("certs2/ca.crt")
	if err != nil {
		log.Fatalf("could not read ca certificate: %s", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatal("failed to append ca certs")
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      certPool,
	}

	// conn, err := grpc.NewClient("localhost:9085", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// Set up a connection to the server.
	conn, err := grpc.Dial("sms.easydelivery.ltd:9085", grpc.WithTransportCredentials(credentials.NewTLS(config)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSmsServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SendSMS(ctx, &pb.SendSMSRequest{
		Receiver: "123456789",
		Message:  "Hello",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Resp: %s", r.GetResp())
	log.Printf("Code: %v", r.GetCode())
}
