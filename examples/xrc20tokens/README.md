### 1. Install ioctl and solc

```
wget https://github.com/ethereum/solidity/releases/download/v0.5.16/solc-static-linux

chmod +x solc-static-linux

cp solc-static-linux /usr/local/bin/solc

curl https://raw.githubusercontent.com/iotexproject/iotex-core/master/install-cli.sh | sh -s "unstable"
```

### 2. Compile XRC20 Contract using ioctl

```
ioctl contract compile XRC20 --abi-out XRC20.abi --bin-out XRC20.bin
```


### 3. Create and Set Private Key

```
ioctl account create
...

export PrivateKey=d0bd45f30f5efea...7a498
```

### 4. Run

```
./xrc20tokens
```
