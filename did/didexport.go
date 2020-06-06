package main

import (
	"C"
	"math/big"
)

const (
	Failure = 0
	Success = 1
)

//CeateDID returns transaction hash,transaction if success,error message
//export CeateDID
func CeateDID(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, id, didHash, url string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), Failure, C.CString("gas price convert error")
	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	h, err := d.CreateDID(id, didHash, url)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	return C.CString(h), Success, C.CString("")
}

//DeleteDID returns transaction hash,transaction if success,error message
//export DeleteDID
func DeleteDID(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, did string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), Failure, C.CString("gas price convert error")

	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	h, err := d.DeleteDID(did)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	return C.CString(h), Success, C.CString("")
}

//UpdateHash returns transaction hash,transaction if success,error message
//export UpdateHash
func UpdateHash(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, did, didHash string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), Failure, C.CString("gas price convert error")
	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	h, err := d.UpdateHash(did, didHash)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	return C.CString(h), Success, C.CString("")
}

//UpdateUri returns transaction hash,transaction if success,error message
//export UpdateUri
func UpdateUri(endpoint, privateKey, contract, abiString, gasPrice string, gasLimit uint64, did, uri string) (*C.char, uint64, *C.char) {
	gp, ok := new(big.Int).SetString(gasPrice, 10)
	if !ok {
		return C.CString(""), Failure, C.CString("gas price convert error")
	}
	d, err := NewDID(endpoint, privateKey, contract, abiString, gp, gasLimit)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	h, err := d.UpdateUri(did, uri)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	return C.CString(h), Success, C.CString("")
}

//GetHash returns did hash,transaction if success,error message
//export GetHash
func GetHash(endpoint, contract, abiString, did string) (*C.char, uint64, *C.char) {
	d, err := NewDID(endpoint, "", contract, abiString, nil, 0)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	h, err := d.GetHash(did)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	return C.CString(h), Success, C.CString("")
}

//GetUri returns did uri,transaction if success,error message
//export GetUri
func GetUri(endpoint, contract, abiString, did string) (*C.char, uint64, *C.char) {
	d, err := NewDID(endpoint, "", contract, abiString, nil, 0)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	uri, err := d.GetUri(did)
	if err != nil {
		return C.CString(""), Failure, C.CString(err.Error())
	}
	return C.CString(uri), Success, C.CString("")
}
func main() {}
