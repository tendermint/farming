<!-- order: 7 -->

# Events

The farming module emits the following events:

## EndBlocker

| Type                  | Attribute Key         | Attribute Value           |
| --------------------- | --------------------- | ------------------------- |
|   plan_terminated  |                 | {}    |
|   rewards_allocated  |                 | {}    |

## Handlers

### MsgCreateRatioPlan

| Type                      | Attribute Key    | Attribute Value |
| ------------------------- | ---------------- | --------------- |
| create_fixed_amount_plan  | plan_id          | {planID}        |

### MsgCreateRatioPlan

| Type                      | Attribute Key    | Attribute Value |
| ------------------------- | ---------------- | --------------- |
| create_ratio_plan  | plan_id          | {planID}        |

### MsgStake

| Type                      | Attribute Key    | Attribute Value |
| ------------------------- | ---------------- | --------------- |
| stake  | plan_id          | {planID}        |

### MsgUnstake

| Type                      | Attribute Key    | Attribute Value |
| ------------------------- | ---------------- | --------------- |
| unstake  | plan_id          | {planID}        |

### MsgHarvest

| Type                      | Attribute Key    | Attribute Value |
| ------------------------- | ---------------- | --------------- |
| harvest  | plan_id          | {planID}        |
### MsgAdvanceEpoch

This message is for testing purpose. It is only available when you build `farmingd` binary by `make install-testing` command.
