# State & Types

## Groups

Groups are simply aggregations of members.

### ID

They have an auto-incremented ID:
```go
type GroupID uint64
```

### Properties

#### `Members []Member`

An array of `Member` structs:

```go
// Member specifies a address and a weight for a group member
type Member struct {
	// The address of a group member. Can be another group or a contract
	Address sdk.AccAddress `json:"address"`
	// The integral weight of this member with respect to other members and the decision threshold
	Weight sdk.Int `json:"weight"`
}
```

#### `Memo string`

A single string memo

#### Indexes

#### `Member`

It should be possible to look-up groups by member

## Group Accounts

Group accounts associate a group with a decision policy

### ID

Group accounts are identified by an `sdk.AccAddress` generated
from an auto-incremented `uint64`.

*TODO:* How to encode group addresses?

### Properties

#### `Group GroupID`

The `GroupID`

### `DecisionPolicy DecisionPolicy`

An instance of:

```go
// DecisionPolicy allows for flexibility in decision policy based both on
// weights (the tally of yes, no, abstain, and veto votes) and time (via
// the block header proposalSubmitTime)
type DecisionPolicy interface {
	Allow(tally Tally, totalPower sdk.Int, header types.Header, proposalSubmitTime time.Time)
}
```

#### `Memo string`

A single string memo.

#### Indexes

#### `Group`

It should be possible to look-up group accounts by group ID.

## Proposals

## Votes




