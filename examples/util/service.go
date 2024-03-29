// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package util

import (
	"crypto/tls"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"

	"github.com/iotexproject/iotex-proto/golang/iotexapi"

	"github.com/iotexproject/iotex-antenna-go/v2/account"
	"github.com/iotexproject/iotex-antenna-go/v2/iotex"
)

// IotexService is the IotexService interface
type IotexService interface {
	// Connect connect to iotex server
	Connect() error
	// AuthClient is the client with private key
	AuthClient() iotex.AuthedClient
	// ReadOnlyClient is the client without private key
	ReadOnlyClient() iotex.ReadOnlyClient
}

type iotexService struct {
	sync.RWMutex
	endpoint       string
	secure         bool
	accountPrivate string

	grpcConn       *grpc.ClientConn
	authedClient   iotex.AuthedClient
	readOnlyClient iotex.ReadOnlyClient
}

// NewIotexService returns IotexService
func NewIotexService(accountPrivate, endpoint string, secure bool) IotexService {
	return &iotexService{
		endpoint:       endpoint,
		secure:         secure,
		accountPrivate: accountPrivate,
	}
}

// Connect connect to iotex server
func (s *iotexService) Connect() (err error) {
	return s.connect()
}

// AuthClient is the client with private key
func (s *iotexService) AuthClient() iotex.AuthedClient {
	return s.authedClient
}

// AuthClient is the client without private key
func (s *iotexService) ReadOnlyClient() iotex.ReadOnlyClient {
	return s.readOnlyClient
}

func (s *iotexService) connect() (err error) {
	s.Lock()
	defer s.Unlock()
	// Check if the existing connection is good.
	if s.grpcConn != nil && s.grpcConn.GetState() != connectivity.Shutdown {
		return
	}
	opts := []grpc.DialOption{}
	if s.secure {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{})))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	s.grpcConn, err = grpc.Dial(s.endpoint, opts...)
	if err != nil {
		return
	}
	if s.accountPrivate != "" {
		creator, err := account.HexStringToAccount(s.accountPrivate)
		if err != nil {
			return err
		}
		s.authedClient = iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(s.grpcConn), 1, creator)
	}

	s.readOnlyClient = iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(s.grpcConn))
	return
}
