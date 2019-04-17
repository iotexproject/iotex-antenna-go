#!/bin/bash
apt-get update
sudo apt-get install openssl -y
sudo apt-get install libssl-dev -y
sudo apt-get install autoconf -y
sudo apt-get install default-jre -y
sudo apt-get install build-essential automake libtool pkg-config libffi-dev python-dev python-pip -y
cd vendor/github.com/ethereum/go-ethereum/crypto/secp256k1/libsecp256k1
./autogen.sh
./configure --disable-shared --with-pic --with-bignum=no --enable-module-recovery --disable-jni
make
sudo make install