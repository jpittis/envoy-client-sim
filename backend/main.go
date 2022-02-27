package main

import (
	"context"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	pb "github.com/jpittis/envoy-client-sim/backend/proto"
)

const (
	numClients           = 2
	sleepBetweenRequests = 1 * time.Second
)

// Mostly so I can run the binary outside of docker-compose without it crashing.
var (
	defaultEndpoints = []string{"10081", "10082"}

	success = promauto.NewCounter(prometheus.CounterOpts{
		Name: "backend_success_total",
		Help: "A count of successful requests",
	})
	failure = promauto.NewCounter(prometheus.CounterOpts{
		Name: "backend_failure_total",
		Help: "A count of failed requests",
	})
	latency = promauto.NewSummary(prometheus.SummaryOpts{
		Name:       "backend_latency",
		Help:       "A summary of request latency",
		Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
	})
)

func main() {
	endpoints := defaultEndpoints
	// Source endpoint config from a shared file so that they can be coordinated with the
	// generated Envoy config.
	buf, err := ioutil.ReadFile("/etc/endpoints.txt")
	if err != nil {
		log.Println("Error:", err)
	} else {
		endpoints = strings.Split(string(buf), ",")
	}
	log.Println("Endpoints:", endpoints)
	// Spawn one gRPC server per endpoint to simulate multiple backends.
	for _, port := range endpoints {
		go func(port string) {
			err := listen("127.0.0.1:"+port, port)
			if err != nil {
				log.Fatal(err)
			}
		}(port)
	}
	// Spawn one or more gRPC clients, each serially generating load controlled by the
	// sleep time between requests.
	for i := 0; i < numClients; i++ {
		go func() {
			client, conn, err := connect("127.0.0.1:10080")
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()

			for {
				start := time.Now()
				timer := prometheus.NewTimer(latency)
				rep, err := client.Get(context.Background(), &pb.GetRequest{})
				duration := time.Since(start)
				timer.ObserveDuration()
				if err != nil {
					failure.Inc()
					log.Printf("Failure! (duration=%s)", duration)
				} else {
					success.Inc()
					log.Printf("Success! (name=%s, duration=%s)", rep.Name, duration)
				}
				withJitter := time.Duration(rand.Int63n(int64(sleepBetweenRequests) * 2))
				time.Sleep(withJitter)
			}
		}()
	}
	// Block on running a prometheus metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("127.0.0.1:2112", nil))
}

func listen(addr, name string) error {
	li, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterBackendServer(grpcServer, &backend{name})
	return grpcServer.Serve(li)
}

type backend struct {
	name string
}

func (b *backend) Get(
	ctx context.Context,
	req *pb.GetRequest,
) (*pb.GetResponse, error) {
	return &pb.GetResponse{Name: b.name}, nil
}

func connect(addr string) (pb.BackendClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	return pb.NewBackendClient(conn), conn, nil
}
