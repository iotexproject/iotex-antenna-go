package service

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

type IotexService interface {
	Connect() error
	AuthClient() iotex.AuthedClient
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

func NewIotexService(accountPrivate, endpoint string, secure bool) IotexService {
	return &iotexService{
		endpoint:       endpoint,
		secure:         secure,
		accountPrivate: accountPrivate,
	}
}

func (s *iotexService) Connect() (err error) {
	return s.connect()
}

func (s *iotexService) AuthClient() iotex.AuthedClient {
	return s.authedClient
}
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
		s.authedClient = iotex.NewAuthedClient(iotexapi.NewAPIServiceClient(s.grpcConn), creator)
	}

	s.readOnlyClient = iotex.NewReadOnlyClient(iotexapi.NewAPIServiceClient(s.grpcConn))
	return
}
