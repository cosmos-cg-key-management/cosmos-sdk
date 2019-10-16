```go
type GroupID uint64

// Groups get their own GroupID
type MsgCreateGroup struct {
	Signer sdk.AccAddress `json:"signer"`
	// The Owner of the group is allowed to change the group structure. A group account
	// can own a group in order for the group to be able to manage its own members
	Owner  sdk.AccAddress `json:"owner"`
	// The members of the group and their associated weight
	Members []Member `json:"members,omitempty"`
	// TODO maybe make this something more specific to a domain name or a claim on identity? or Info leave it generic
	Memo string `json:"memo,omitempty"`
}

// group accounts get their own sdk.AccAddress
type MsgCreateGroupAccount struct {
	Signer         sdk.AccAddress `json:"signer"`
	// The Owner of a group account is allowed to change the DecisionPolicy. This can be left nil 
	// in order for the group account to "own" itself
	Owner          sdk.AccAddress `json:"owner"`
	Group          GroupID        `json:"group"`
	DecisionPolicy DecisionPolicy `json:"decision_policy"`
	Memo           string         `json:"memo,omitempty"`
}

type Tally struct {
	YesCount sdk.Int
	NoCount sdk.Int
	AbstainCount sdk.Int
	VetoCount sdk.Int
}

// DecisionPolicy allows for flexibility in decision policy based both on
// weights (the tally of yes, no, abstain, and veto votes) and time (via
// the block header proposalSubmitTime)
type DecisionPolicy interface {
	Allow(tally Tally, totalPower sdk.Int, header types.Header, proposalSubmitTime time.Time)
}

type ThresholdDecisionPolicy struct {
	// Specifies the number of votes that must be accumulated in order for a decision to be made by the group.
	// A member gets as many votes as is indicated by their Weight field.
	// A big integer is used here to avoid any potential vulnerabilities from overflow errors
	// where large weight and threshold values are used.
	DecisionThreshold sdk.Int `json:"decision_threshold"`
}

type PercentageDecisionPolicy struct {
	Percent sdk.Dec `json:"percent"`
}

// A member specifies a address and a weight for a group member
type Member struct {
	// The address of a group member. Can be another group or a contract
	Address sdk.AccAddress `json:"address"`
	// The integral weight of this member with respect to other members and the decision threshold
	Weight sdk.Int `json:"weight"`
}
```

Because a group has its own `sdk.AccAddress` group members can also be other
groups so that groups can be nested.

## Proposals

Groups can execute any authorized action on the blockchain using their group
`sdk.AccAddress` by approving proposals.

Proposals have a simple creation, voting, and execution behavior based on the
following messages:

```go
type MsgCreateProposal struct {
	Proposer sdk.AccAddress `json:"proposer"`
	Group    sdk.AccAddress `json:"group"`
	Msgs     []sdk.Msg      `json:"msgs"`
	// Exec can be set to true in order to attempt to execute the proposal immediately
	// with no voting in a single transaction - this is useful for 1/N or 2/N multisig
	// key groups. Every signer of the MsgCreateProposal transaction is considered a yes
	// vote
	Exec bool `json:"exec,omitempty"`
}

type Vote int

const (
	No Vote = iota
	Yes
	Abstain
	Veto
)

type MsgVote struct {
	ProposalID ProposalID     `json:"proposal_id"`
	// Voters must sign this transaction
	Voters     []sdk.AccAddress `json:"voters"`
	Vote       Vote           `json:"vote"`
}

type MsgTryExecuteProposal struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Signer     sdk.AccAddress `json:"signer"`
}
```

As one can see proposal `Msgs` can be any `sdk.Msg`. The execution of these
messages is handled by the `delegation` module which checks whether or not
the group is authorized to execute the provided `Msg` and routes these
messages back to the `BaseApp` `Router` if so.

Groups can use proposals to update their group members and decision policies.

