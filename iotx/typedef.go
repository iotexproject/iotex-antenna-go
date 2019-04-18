// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package iotx

type TransferRequest struct {
	From     string
	To       string
	Value    string
	Payload  string
	GasLimit string
	GasPrice string
}
type ContractRequest struct {
	From   string
	Amount string
	// contract bytecode
	Data     []byte
	GasLimit string
	GasPrice string
}
