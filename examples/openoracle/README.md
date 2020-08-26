**Prepare the contracts as below before the deployement and invocation**

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
ioctl contract compile OpenOraclePriceData --abi-out OpenOraclePriceData.abi --bin-out OpenOraclePriceData.bin```
