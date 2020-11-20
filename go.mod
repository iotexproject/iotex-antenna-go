module github.com/iotexproject/iotex-antenna-go/v2

go 1.14

require (
	github.com/aristanetworks/goarista v0.0.0-20190531155855-fef20d617fa7 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/ethereum/go-ethereum v1.8.27
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/iotexproject/go-pkgs v0.1.4
	github.com/iotexproject/iotex-address v0.2.3
	github.com/iotexproject/iotex-proto v0.4.3
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.4.0
	google.golang.org/grpc v1.27.0
)

replace github.com/ethereum/go-ethereum => github.com/iotexproject/go-ethereum v0.3.1
