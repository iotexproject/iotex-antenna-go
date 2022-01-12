package iotex

import (
	"crypto/tls"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewDefaultGRPCConn creates a default grpc connection. With tls and retry.
func NewDefaultGRPCConn(endpoint string) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Second)),
		grpc_retry.WithMax(3),
	}
	return grpc.Dial(endpoint,
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
}

// NewGRPCConnWithoutTLS creates a default grpc connection, without tls
func NewGRPCConnWithoutTLS(endpoint string) (*grpc.ClientConn, error) {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(100 * time.Second)),
		grpc_retry.WithMax(3),
	}
	return grpc.Dial(endpoint,
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
		grpc.WithInsecure())
}
