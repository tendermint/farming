# Farming Plans

There are two different types of farming plans in the farming module. Where as a public farming plan can only be created through governance proposal, a private farming plan can be created with any account. The plan creator's account is used as distributing account FarmingPoolAddress that will be distributed to farmers automatically. 

In this documentation, there are sample data provided in JSON structure that are needed to test the functionality by using farming command line interfaces. 

## Table of Contetns

- [Bootstrap Local Network](#Boostrap)
- [Public Farming Plan](#Public-Farming-Plan)
  * [AddPublicFarmingFixedAmountPlan](#AddPublicFarmingFixedAmountPlan)
  * [AddPublicFarmingRatioPlan](#AddPublicFarmingRatioPlan)
  * [AddMultiplePublicPlans](#AddMultiplePublicPlans)
  * [UpdatePublicFarmingFixedAmountPlan](#UpdatePublicFarmingFixedAmountPlan)
  * [DeletePublicFarmingFixedAmountPlan](#DeletePublicFarmingFixedAmountPlan)
- [Private Farming Plan](#Private-Farming-Plan)
  * [PrivateFarmingFixedAmountPlan](#PrivateFarmingFixedAmountPlan)
  * [PrivateFarmingRatioPlan](#PrivateFarmingRatioPlan)

# Bootstrap

```bash
# Clone the project 
git clone https://github.com/tendermint/farming.git
cd cosmos-sdk
make install

# Configure variables
export BINARY=farmingd
export HOME_1=$HOME/.farmingapp
export CHAIN_ID=localnet
export VALIDATOR_1="struggle panic room apology luggage game screen wing want lazy famous eight robot picture wrap act uphold grab away proud music danger naive opinion"
export USER_1="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"
export GENESIS_COINS=10000000000stake,10000000000uatom,10000000000uusd

# Boostrap
$BINARY init $CHAIN_ID --chain-id $CHAIN_ID
echo $VALIDATOR_1 | $BINARY keys add val1 --keyring-backend test --recover
echo $USER_1 | $BINARY keys add user1 --keyring-backend test --recover
$BINARY add-genesis-account $($BINARY keys show val1 --keyring-backend test -a) $GENESIS_COINS
$BINARY add-genesis-account $($BINARY keys show user1 --keyring-backend test -a) $GENESIS_COINS
$BINARY gentx val1 100000000stake --chain-id $CHAIN_ID --keyring-backend test
$BINARY collect-gentxs

# Modify app.toml
sed -i '' 's/enable = false/enable = true/g' $HOME_1/config/app.toml
sed -i '' 's/swagger = false/swagger = true/g' $HOME_1/config/app.toml

# Modify governance proposal for testing purpose
sed -i '' 's%"amount": "10000000"%"amount": "1"%g' $HOME_1/config/genesis.json
sed -i '' 's%"quorum": "0.334000000000000000",%"quorum": "0.000000000000000001",%g' $HOME_1/config/genesis.json
sed -i '' 's%"threshold": "0.500000000000000000",%"threshold": "0.000000000000000001",%g' $HOME_1/config/genesis.json
sed -i '' 's%"voting_period": "172800s"%"voting_period": "60s"%g' $HOME_1/config/genesis.json

# Start
$BINARY start
```

# Public Farming Plan

## AddPublicFarmingFixedAmountPlan

```json
{
  "title": "Public Farming Plan",
  "description": "Are you ready to farm?",
  "add_request_proposals": [
    {
      "name": "First Public Farming Plan",
      "farming_pool_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "termination_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "staking_coin_weights": [
        {
          "denom": "PoolCoinDenom",
          "amount": "1.000000000000000000"
        }
      ],
      "start_time": "2021-07-15T08:41:21.662422Z",
      "end_time": "2022-07-16T08:41:21.662422Z",
      "epoch_amount": [
        {
          "denom": "uatom",
          "amount": "1"
        }
      ]
    }
  ]
}
```

## AddPublicFarmingRatioPlan

```json
{
  "title": "Public Farming Plan",
  "description": "Are you ready to farm?",
  "add_request_proposals": [
    {
      "name": "First Public Farming Plan",
      "farming_pool_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "termination_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "staking_coin_weights": [
        {
          "denom": "PoolCoinDenom",
          "amount": "1.000000000000000000"
        }
      ],
      "start_time": "2021-07-15T08:41:21.662422Z",
      "end_time": "2022-07-16T08:41:21.662422Z",
      "epoch_ratio": "1.000000000000000000"
    }
  ]
}
```

## AddMultiplePublicPlans

```json
{
  "title": "Public Farming Plan",
  "description": "Are you ready to farm?",
  "add_request_proposals": [
    {
      "name": "First Public Farming Plan",
      "farming_pool_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "termination_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "staking_coin_weights": [
        {
          "denom": "PoolCoinDenom",
          "amount": "1.000000000000000000"
        }
      ],
      "start_time": "2021-07-15T08:41:21.662422Z",
      "end_time": "2022-07-16T08:41:21.662422Z",
      "epoch_amount": [
        {
          "denom": "uatom",
          "amount": "1"
        }
      ]
    },
    {
      "name": "First Public Farming Plan",
      "farming_pool_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "termination_address": "cosmos1zaavvzxez0elundtn32qnk9lkm8kmcszzsv80v",
      "staking_coin_weights": [
        {
          "denom": "PoolCoinDenom",
          "amount": "1.000000000000000000"
        }
      ],
      "start_time": "2021-07-15T08:41:21.662422Z",
      "end_time": "2022-07-16T08:41:21.662422Z",
      "epoch_ratio": "1.000000000000000000"
    }
  ]
}
```

## UpdatePublicFarmingFixedAmountPlan

```json
{
  "title": "Let's Update the Farming Plan 1",
  "description": "FarmingPoolAddress needs to be changed",
  "update_request_proposals": [
    {
      "plan_id": 1,
      "farming_pool_address": "cosmos13w4ueuk80d3kmwk7ntlhp84fk0arlm3mqf0w08",
      "termination_address": "cosmos13w4ueuk80d3kmwk7ntlhp84fk0arlm3mqf0w08",
      "staking_coin_weights": [
        {
          "denom": "uatom",
          "amount": "1.000000000000000000"
        }
      ],
      "start_time": "2021-07-15T08:41:21.662422Z",
      "end_time": "2022-07-16T08:41:21.662422Z",
      "epoch_amount": [
        {
          "denom": "uatom",
          "amount": "1"
        }
      ]
    }
  ]
}
```

## DeletePublicFarmingFixedAmountPlan

```json
{
  "title": "Delete Public Farming Plan 1",
  "description": "This plan is no longer needed",
  "delete_request_proposals": [
    {
      "plan_id": 1
    }
  ]
}
```

# Private Farming Plan

## PrivateFarmingFixedAmountPlan

```json
{
	"name": "This Farming Plan intends to incentivize ATOM HODLERS!",
  "staking_coin_weights": [
	  {
	      "denom": "uatom",
	      "amount": "1.000000000000000000"
	  }
  ],
  "start_time": "2021-07-15T08:41:21.662422Z",
  "end_time": "2022-07-16T08:41:21.662422Z",
  "epoch_amount": [
    {
      "denom": "uatom",
      "amount": "1"
    }
  ]
}
```

## PrivateFarmingRatioPlan

```json
{
	"name": "This Farming Plan intends to incentivize ATOM HODLERS!",
  "staking_coin_weights": [
	  {
	      "denom": "uatom",
	      "amount": "1.000000000000000000"
	  }
  ],
  "start_time": "2021-07-15T08:41:21.662422Z",
  "end_time": "2022-07-16T08:41:21.662422Z",
  "epoch_ratio":"1.000000000000000000"
}
```