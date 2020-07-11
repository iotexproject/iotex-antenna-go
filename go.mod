module github.com/iotexproject/iotex-antenna-go/v2

require (
	github.com/aristanetworks/goarista v0.0.0-20190531155855-fef20d617fa7 // indirect
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/ethereum/go-ethereum v1.8.27
	github.com/gogo/protobuf v1.2.1
	github.com/golang/mock v1.4.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/iotexproject/go-pkgs v0.1.2-0.20200212033110-8fa5cf96fc1b
	github.com/iotexproject/iotex-address v0.2.1
	github.com/iotexproject/iotex-proto v0.3.0
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	golang.org/x/net v0.0.0-20191204025024-5ee1b9f4859a // indirect
	google.golang.org/genproto v0.0.0-20191203220235-3fa9dbf08042 // indirect
	google.golang.org/grpc v1.21.0
)

replace github.com/ethereum/go-ethereum v1.8.27 => github.com/iotexproject/go-ethereum v0.1.0

go 1.13
