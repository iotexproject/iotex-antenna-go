### Install solc and ioctl

```
wget https://github.com/ethereum/solidity/releases/download/v0.5.16/solc-static-linux

chmod +x solc-static-linux

cp solc-static-linux /usr/local/bin/solc

curl https://raw.githubusercontent.com/iotexproject/iotex-core/master/install-cli.sh | sh -s "unstable"
```
### Compile contract using ioctl

```
ioctl contract compile XRC20 --abi-out XRC20.abi --bin-out XRC20.bin```