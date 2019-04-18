#!/bin/bash
apt-get update
apt-get install openssl -y
apt-get install libssl-dev -y
apt-get install autoconf -y
apt-get install default-jre -y
apt-get install build-essential automake libtool pkg-config libffi-dev python-dev python-pip -y
cd vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1
./autogen.sh
./configure --disable-shared --with-pic --with-bignum=no --enable-module-recovery --disable-jni
make
make install