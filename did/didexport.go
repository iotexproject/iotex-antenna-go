// Copyright (c) 2020 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package main

import (
	"C"
	"math/big"
)

const (
	failure = 0
	success = 1
)

//CeateDID returns transaction hash,transaction if success,error message
//export CeateDID
func CeateDID(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, id, didHash, url string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), failure, C.CString("gas price convert error")
	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	h, err := d.CreateDID(id, didHash, url)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	return C.CString(h), success, C.CString("")
}

//DeleteDID returns transaction hash,transaction if success,error message
//export DeleteDID
func DeleteDID(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, did string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), failure, C.CString("gas price convert error")

	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	h, err := d.DeleteDID(did)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	return C.CString(h), success, C.CString("")
}

//UpdateHash returns transaction hash,transaction if success,error message
//export UpdateHash
func UpdateHash(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, did, didHash string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), failure, C.CString("gas price convert error")
	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	h, err := d.UpdateHash(did, didHash)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	return C.CString(h), success, C.CString("")
}

//UpdateURI returns transaction hash,transaction if success,error message
//export UpdateURI
func UpdateURI(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, did, uri string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), failure, C.CString("gas price convert error")
	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	h, err := d.UpdateURI(did, uri)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	return C.CString(h), success, C.CString("")
}

//GetHash returns did hash,transaction if success,error message
//export GetHash
func GetHash(endpoint, contract, abiString, did string) (*C.char, uint64, *C.char) {
	d, err := NewDID(endpoint, "", contract, abiString, nil, 0)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	h, err := d.GetHash(did)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	return C.CString(h), success, C.CString("")
}

//GetURI returns did uri,transaction if success,error message
//export GetURI
func GetURI(endpoint, contract, abiString, did string) (*C.char, uint64, *C.char) {
	d, err := NewDID(endpoint, "", contract, abiString, nil, 0)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	uri, err := d.GetURI(did)
	if err != nil {
		return C.CString(""), failure, C.CString(err.Error())
	}
	return C.CString(uri), success, C.CString("")
}
func main() {}
