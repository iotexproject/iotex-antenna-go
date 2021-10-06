module github.com/iotexproject/iotex-antenna-go/v2

go 1.16

require (
	github.com/btcsuite/btcd v0.21.0-beta // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/ethereum/go-ethereum v1.10.4
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/iotexproject/go-pkgs v0.1.5-0.20210604060651-be5ee19f2575
	github.com/iotexproject/iotex-address v0.2.4
	github.com/iotexproject/iotex-proto v0.5.0
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	google.golang.org/grpc v1.33.1
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.4.0
