### 1. Get Open Oracle Contracts

```
git clone https://github.com/compound-finance/open-oracle.git
```

### 2. Install ioctl and solc

```
wget https://github.com/ethereum/solidity/releases/download/v0.6.10/solc-static-linux

chmod +x solc-static-linux

cp solc-static-linux /usr/local/bin/solc

curl https://raw.githubusercontent.com/iotexproject/iotex-core/master/install-cli.sh | sh -s "unstable"
```

### 3. Compile OpenOraclePriceData Contract

```
ioctl contract compile OpenOraclePriceData --abi-out OpenOraclePriceData.abi --bin-out OpenOraclePriceData.bin
```

### 4. Create and Set Private Key

```
ioctl account create
...

export PrivateKey=d0bd45f30f5efea...7a498
```
Get some testnet IOTX from https://faucet.iotex.io/ for the created account.


### 5. Run
```
./openoracle
```
