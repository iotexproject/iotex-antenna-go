#include <stdio.h>
#include <string.h>
#include <unistd.h>
#include "../didlib.h"

int main ()
{
    char privateKey[] = "414efa99dfac6f4095d6954713fb0085268d400d6a05a8ae8a69b5b1c10b4bed";
    char updatedHash[] = "414efa99dfac6f4095d6954713fb0085268d400d6a05a8ae8a69b5b1c10eeeee";
	char endpoint[] = "api.testnet.iotex.one:443";
	char contract[] = "io1zgs5gqjl679qlj4gqqpa9t329r8f5gr8xc9lr0";
	char did[] = "did:io:0x0ddfC506136fb7c050Cc2E9511eccD81b15e7426";
	char abi[] = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"didString\",\"type\":\"string\"}],\"name\":\"CreateDID\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"didString\",\"type\":\"string\"}],\"name\":\"DeleteDID\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"didString\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"UpdateHash\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"didString\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"UpdateURI\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"id\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"createDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"did\",\"type\":\"string\"}],\"name\":\"deleteDID\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"dids\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"exist\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"did\",\"type\":\"string\"}],\"name\":\"getHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"string\",\"name\":\"did\",\"type\":\"string\"}],\"name\":\"getURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"did\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"}],\"name\":\"updateHash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"string\",\"name\":\"did\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"uri\",\"type\":\"string\"}],\"name\":\"updateURI\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]";
    char gasPrice[] = "1000000000000";
    char uri[] = "urixxx";
    char updatedUri[] = "urixxxupdated";
    char id[] = "";
    GoString PrivateKey;
    PrivateKey.p = privateKey;
    PrivateKey.n = strlen(privateKey);

    GoString UpdatedHash;
    UpdatedHash.p = updatedHash;
    UpdatedHash.n = strlen(updatedHash);

    GoString Endpoint;
    Endpoint.p = endpoint;
    Endpoint.n = strlen(endpoint);
    GoString Contract;
    Contract.p = contract;
    Contract.n = strlen(contract);
    GoString Abi;
    Abi.p = abi;
    Abi.n = strlen(abi);
    GoString Did;
    Did.p = did;
    Did.n = strlen(did);
    GoString GasPrice;
    GasPrice.p = gasPrice;
    GasPrice.n = strlen(gasPrice);
    GoString Uri;
    Uri.p = uri;
    Uri.n = strlen(uri);
    GoString Id;
    Id.p = id;
    Id.n = strlen(id);
    GoString UpdatedUri;
    UpdatedUri.p = updatedUri;
    UpdatedUri.n = strlen(updatedUri);
    //CeateDID
    struct CeateDID_return CeateDIDRet=CeateDID(Endpoint,PrivateKey,Contract,Abi,GasPrice,1000000,Id,PrivateKey,Uri);
    printf("CeateDID %s %lld %s\n",CeateDIDRet.r0,CeateDIDRet.r1,CeateDIDRet.r2);
    sleep(10);

    //GetHash
    struct GetHash_return GetHashRet=GetHash(Endpoint,Contract,Abi,Did);
    printf("GetHash %s %lld %s\n",GetHashRet.r0,GetHashRet.r1,GetHashRet.r2);
    //GetUri
    struct GetUri_return GetUriRet=GetUri(Endpoint,Contract,Abi,Did);
    printf("GetUri %s %lld %s\n",GetUriRet.r0,GetUriRet.r1,GetUriRet.r2);

    //UpdateHash
    struct UpdateHash_return UpdateHashRet=UpdateHash(Endpoint,PrivateKey,Contract,Abi,GasPrice,1000000,Id,UpdatedHash);
    printf("UpdateHash %s %lld %s\n",UpdateHashRet.r0,UpdateHashRet.r1,UpdateHashRet.r2);
    sleep(10);
    //UpdateUri
    struct UpdateUri_return UpdateUriRet=UpdateUri(Endpoint,PrivateKey,Contract,Abi,GasPrice,1000000,Id,UpdatedUri);
    printf("UpdateUri %s %lld %s\n",UpdateUriRet.r0,UpdateUriRet.r1,UpdateUriRet.r2);
    sleep(10);

    printf("after update\n");
    //GetHash
    GetHashRet=GetHash(Endpoint,Contract,Abi,Did);
    printf("GetHash %s %lld %s\n",GetHashRet.r0,GetHashRet.r1,GetHashRet.r2);
    //GetUri
    GetUriRet=GetUri(Endpoint,Contract,Abi,Did);
    printf("GetUri %s %lld %s\n",GetUriRet.r0,GetUriRet.r1,GetUriRet.r2);

    //DeleteDID
    struct DeleteDID_return DeleteDIDRet=DeleteDID(Endpoint,PrivateKey,Contract,Abi,GasPrice,1000000,Id);
    printf("DeleteDID %s %lld %s\n",DeleteDIDRet.r0,DeleteDIDRet.r1,DeleteDIDRet.r2);
    sleep(10);

    printf("after delete\n");
    //GetHash
    GetHashRet=GetHash(Endpoint,Contract,Abi,Did);
    printf("GetHash %s %lld %s\n",GetHashRet.r0,GetHashRet.r1,GetHashRet.r2);
    //GetUri
    GetUriRet=GetUri(Endpoint,Contract,Abi,Did);
    printf("GetUri %s %lld %s\n",GetUriRet.r0,GetUriRet.r1,GetUriRet.r2);

    return 0;
}