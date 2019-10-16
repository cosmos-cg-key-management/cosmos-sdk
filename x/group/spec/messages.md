```go
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
type MsgCreateProposal struct {
	Proposer sdk.AccAddress `json:"proposer"`
	GroupAccount    sdk.AccAddress `json:"group"`
	Msgs     []sdk.Msg      `json:"msgs"`
	// Exec can be set to true in order to attempt to execute the proposal immediately
	// with no voting in a single transaction - this is useful for 1/N or 2/N multisig
	// key groups. Every signer of the MsgCreateProposal transaction is considered a yes
	// vote
	Exec bool `json:"exec,omitempty"`
}

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
