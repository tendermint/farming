---
Title: Localnet
Description: A tutorial of how to build `farmingd` and bootstrap local network.
---

## Get farming module source code

```bash
# Use git to clone farming module source code and install `farmingd`
git clone https://github.com/tendermint/farming.git
cd farming
make install
```

## Start a blockchain with Starport

Use [Starport CLI](https://docs.starport.network/cli/#starport-chain-serve) to start a local blockchain with automatic reloading. You can configure custom settings in [config.yml](../../../config.yml).

```bash
starport chain serve
```

## Start a blockchain with commands

The following commands are used to bootstrap a single chain with a single validator in your local machine. Copy the commands and run them in your terminal.

```bash
# Configure variables
export BINARY=farmingd
export HOME_FARMINGAPP=$HOME/.farmingapp
export CHAIN_ID=localnet
export VALIDATOR_1="struggle panic room apology luggage game screen wing want lazy famous eight robot picture wrap act uphold grab away proud music danger naive opinion"
export USER_1="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"
export USER_2="fuel obscure melt april direct second usual hair leave hobby beef bacon solid drum used law mercy worry fat super must ritual bring faculty"
export VALIDATOR_1_GENESIS_COINS=10000000000stake,10000000000uatom,10000000000uusd
export USER_1_GENESIS_COINS=10000000000stake,10000000000uatom,10000000000uusd
export USER_2_GENESIS_COINS=10000000000stake,10000000000poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4

# Initialize chain and craete gentx for a single validator
$BINARY init $CHAIN_ID --chain-id $CHAIN_ID
echo $VALIDATOR_1 | $BINARY keys add val1 --keyring-backend test --recover
echo $USER_1 | $BINARY keys add user1 --keyring-backend test --recover
echo $USER_2 | $BINARY keys add user2 --keyring-backend test --recover
$BINARY add-genesis-account $($BINARY keys show val1 --keyring-backend test -a) $VALIDATOR_1_GENESIS_COINS
$BINARY add-genesis-account $($BINARY keys show user1 --keyring-backend test -a) $USER_1_GENESIS_COINS
$BINARY add-genesis-account $($BINARY keys show user2 --keyring-backend test -a) $USER_2_GENESIS_COINS
$BINARY gentx val1 100000000stake --chain-id $CHAIN_ID --keyring-backend test
$BINARY collect-gentxs

# Check platform
platform='unknown'
unamestr=`uname`
if [ "$unamestr" = 'Linux' ]; then
   platform='linux'
fi

if [ $platform = 'linux' ]; then
	sed -i 's/enable = false/enable = true/g' $HOME_BUDGETAPP/config/app.toml
	sed -i 's/swagger = false/swagger = true/g' $HOME_BUDGETAPP/config/app.toml
	sed -i 's%"amount": "10000000"%"amount": "1"%g' $HOME_BUDGETAPP/config/genesis.json
    # (Optional) Modify governance proposal for testing public plan proposal
	sed -i 's%"quorum": "0.334000000000000000",%"quorum": "0.000000000000000001",%g' $HOME_BUDGETAPP/config/genesis.json
	sed -i 's%"threshold": "0.500000000000000000",%"threshold": "0.000000000000000001",%g' $HOME_BUDGETAPP/config/genesis.json
	sed -i 's%"voting_period": "172800s"%"voting_period": "30s"%g' $HOME_BUDGETAPP/config/genesis.json
else
	sed -i '' 's/enable = false/enable = true/g' $HOME_BUDGETAPP/config/app.toml
	sed -i '' 's/swagger = false/swagger = true/g' $HOME_BUDGETAPP/config/app.toml
	sed -i '' 's%"amount": "10000000"%"amount": "1"%g' $HOME_BUDGETAPP/config/genesis.json
    # (Optional) Modify governance proposal for testing public plan proposal
	sed -i '' 's%"quorum": "0.334000000000000000",%"quorum": "0.000000000000000001",%g' $HOME_BUDGETAPP/config/genesis.json
	sed -i '' 's%"threshold": "0.500000000000000000",%"threshold": "0.000000000000000001",%g' $HOME_BUDGETAPP/config/genesis.json
	sed -i '' 's%"voting_period": "172800s"%"voting_period": "30s"%g' $HOME_BUDGETAPP/config/genesis.json
fi

# Start
$BINARY start
```
