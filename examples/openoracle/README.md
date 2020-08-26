### Get contract first

```
git clone https://github.com/compound-finance/open-oracle.git
```

### Install solc and ioctl

```
wget https://github.com/ethereum/solidity/releases/download/v0.6.10/solc-static-linux

chmod +x solc-static-linux

cp solc-static-linux /usr/local/bin/solc

curl https://raw.githubusercontent.com/iotexproject/iotex-core/master/install-cli.sh | sh -s "unstable"
```
### Compile OpenOraclePriceData using ioctl in open-oracle/contracts folder

```
ioctl contract compile OpenOraclePriceData --abi-out OpenOraclePriceData.abi --bin-out OpenOraclePriceData.bin```