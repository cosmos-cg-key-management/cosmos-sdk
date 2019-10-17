# State & Types

## Groups

Groups are simply aggregations of members.

### ID

```go
type GroupID []byte
```

`GroupID` is generated from an auto-incrementing `uint64` preprended with the
prefix `0`. This prefix allows other group composition mechanisms in the future,
specifically via token ownership rather than group membership.

### Properties

#### `Members []Member`

An array of `Member` structs:

```go
// Member specifies a address and a power for a group member
type Member struct {
	// The address of a group member. Can be another group or a contract
	Address sdk.AccAddress `json:"address"`
	// The integral power of this member with respect to other members and the decision threshold
	Power sdk.Int `json:"power"`
}
```

#### `Owner sdk.ACcAddress`

Owner is the account which "owns" the group and has the permission to 
add and remove memers.

#### `Memo string`

A single string memo.

#### Indexes

#### `Member`

It should be possible to look-up groups by member address.

#### `Owner`

It should be possible to look-up groups by owner address.

#### `TotalPower sdk.Int`

The sum total power of all member powers should be cached for quick tallying.


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
type Tally struct {
	YesCount sdk.Int
	NoCount sdk.Int
	AbstainCount sdk.Int
	VetoCount sdk.Int
}
__
// DecisionPolicy allows for flexibility in decision policy based both on
// powers (the tally of yes, no, abstain, and veto votes) and time (via
// the block header proposalSubmitTime)
type DecisionPolicy interface {
	Allow(tally Tally, totalPower sdk.Int, header types.Header, submittedTime time.Time, submittedHeight int64)
}
```

#### `Owner sdk.ACcAddress`

Owner is the account which "owns" the group account and has the permission to
change its `DecisionPolicy`. It should be left `nil` if the group account
"owns" itself.

#### `Memo string`

A single string memo.

#### Indexes

#### `Group`

It should be possible to look-up group accounts by group ID.

#### `Owner`

## Proposals

### ID

Proposals get an auto-incremented ID:
```go
type Proposal uint64
```

### Properties

#### `Proposer sdk.AccAddress`

#### `GroupAccount sdk.AccAddress`

#### `Msgs []sdk.Msg`

#### `SubmittedHeight int64`

#### `SubmittedTime time.Time`

#### `Votes []Vote`

```go
type Choice int

const (
	No Choice = iota
	Yes
	Abstain
	Veto
)

type Vote struct {
  Voter sdk.AccAddress
  Choice Choice
}
```

### Indexes

### `GroupAccount`
### `Proposer`
### `SuvmittedTime`
### `Voter`

It should be possible to look up proposal votes by voter address


