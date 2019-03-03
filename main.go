package main

import (
	"fmt"

    "github.com/iotexproject/iotex-core/protogen/iotextypes"
)

func main() {
    var blk iotextypes.Block
	fmt.Printf("Hello Antenna %s\n", blk.Header.ProducerPubkey)
}
