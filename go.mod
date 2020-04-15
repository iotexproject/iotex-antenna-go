module github.com/iotexproject/iotex-antenna-go/v2

require (
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/ethereum/go-ethereum v1.8.27
	github.com/gogo/protobuf v1.2.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/iotexproject/go-pkgs v0.1.1
	github.com/iotexproject/iotex-address v0.2.1
	github.com/iotexproject/iotex-proto v0.2.6-0.20200409230611-748f6ab69ca5
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.3.0
	google.golang.org/grpc v1.20.1
)

replace github.com/ethereum/go-ethereum v1.8.27 => github.com/iotexproject/go-ethereum v0.1.0

go 1.13
