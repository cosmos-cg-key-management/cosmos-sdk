# State & Types

## Groups

Groups are simply aggregations of members.

### Key Value Store Layout

| Key                 | Description              | Type              |
|---------------------|--------------------|--------------------------|
| `g/<group>/desc`  | Group description    | `string`  |
| `g/<group>/<member>`  | Member's voting power    | `sdk.Int`  |
| `g/<group>/<member>/desc`  | Member description    | `string`  |
| `g/<group>/owner`  | Group owner    | `sdk.AccAdress`  |
| `g/<group>/totalPower`  | Group's computed total power    | `sdk.Int`  |
| `mg/<member>/<group>`  |  Member -> group reverse index  | empty |
| `og/<owner>/<group>`  |  Owner -> group reverse index  | empty |

### Key Types

| Key                 | Type             |
|---------------------|------------------|
| `<group>`  | `GroupID`  |
| `<owner>`  | `sdk.AccAddress`  |
| `<member>`  | `sdk.AccAddress`  |

### Custom Types

```go
type GroupID []byte
```

`GroupID` is generated from an auto-incrementing `uint64` preprended with the
prefix `0`. This prefix allows other group composition mechanisms in the future,
specifically via token ownership rather than group membership.

## Group Accounts

Group accounts associate a group with a decision policy.

### Key-Value Store Layout

| Key                 | Description              | Type              |
|---------------------|--------------------|--------------------------|
| `a/<group-account>/description`  | Group account description    | `string`  |
| `a/<group-account>/group`  | Group account's underlying group    | `GroupID`  |
| `a/<group-account>/decisionPolicy`  | Group account's decision policy    | `DecisionPolicy`  |
| `a/<group-account>/owner`  | Group account's owner | `sdk.AccAddress`  |
| `ga/<group>/<group-account>`  | Group -> group account reverse index  | empty |
| `oa/<owner>/<group-account>`  | Owner -> group account reverse index  | empty |

### Key Types

| Key                 | Type             | Description |
|---------------------|------------------| ------------|
| `<group-account>`  | `GroupID`  | Generated from an auto-incremented `uint64` |
| `<group>`  | `GroupID`  | | 
| `<owner>`  | `sdk.AccAddress`  | |

*TODO:* How to encode group account addresses?

### Custom Types

```go
type Tally struct {
	YesCount sdk.Int
	NoCount sdk.Int
	AbstainCount sdk.Int
	VetoCount sdk.Int
}
__
// DecisionPolicy allows for flexibility in decision policy based both on
// powers (the tally of yes, no, abstain, and veto votes) and time since voting
// started
type DecisionPolicy interface {
	Allow(tally Tally, totalPower sdk.Int, timeSinceVotingStart time.Duration)
}
```

## Proposals

### Key Value Store Layout

| Key                 | Description              | Type              |
|---------------------|--------------------|--------------------------|
| `p/<proposal>/desc`  | Proposal description    | `string`  |
| `p/<proposal>/ga`  | Proposal's group account    | `sdk.AccAddress`  |
| `p/<proposal>/msgs`  | Messages that will be run if the proposal succeeds    | `[]sdk.Msg`  |
| `p/<proposal>/proposer`  | Account that proposed the proposal    | `sdk.AccAddress`  |
| `p/<proposal>/votingStart`  | When voting started    | `time.Time`  |
| `p/<proposal>/<voter>/vote`  | A voter's vote on the proposal    | `Vote`  |
| `p/<proposal>/<voter>/comment`  | A voter's comment on their vote | `string`  |
| `vp/<voter>/<proposal>`  | Voter -> proposal reverse look-up | empty  |
| `ap/<group-account>/<proposal>`  | Group account -> proposal reverse look-up | empty  |
| `pp/<proposer>/<proposal>`  | Proposer -> proposal reverse look-up | empty  |

### Key Type

| Key                 | Type             |------|
|---------------------|------------------|------|
| `<proposal>`  | `ProposalID`  | auto-incremented `uint64`
| `<voter>`  | `sdk.AccAddress`  |
| `<proposer>`  | `sdk.AccAddress`  |
| `<group-account>`  | `sdk.AccAddress`  |

### Custom Types

```go
type Vote int

const (
	No Vote = 0
	Yes = 1
	Abstain = 2
	Veto = 3
)
```
