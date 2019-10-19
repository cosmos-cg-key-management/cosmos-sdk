```go
// Groups get their own GroupID
type MsgCreateGroup struct {
	// The Owner of the group is allowed to change the group structure. A group account
	// can own a group in order for the group to be able to manage its own members
	Owner  sdk.AccAddress `json:"owner"`
	// The members of the group and their associated weight
	Members []Member `json:"members,omitempty"`
	// TODO maybe make this something more specific to a domain name or a claim on identity? or Info leave it generic
	Description string `json:"Description,omitempty"`
    DecisionPolicy DecisionPolicy `json:"decision_policy"`
}

type MsgUpdateGroupOwner struct {
	Owner  sdk.AccAddress `json:"owner"`
	Group  sdk.AccAddress `json:"group"`
	NewOwner  sdk.AccAddress `json:"new_owner"`
}

type MsgUpdateGroupDecisionPolicy struct {
	Owner  sdk.AccAddress `json:"owner"`
	Group  sdk.AccAddress `json:"group"`
    DecisionPolicy DecisionPolicy `json:"decision_policy"`
}

type MsgUpdateGroupDescription struct {
	Owner  sdk.AccAddress `json:"owner"`
	Group  sdk.AccAddress `json:"group"`
	Description string `json:"Description,omitempty"`
}

type MsgUpdateGroupMembers struct {
	Owner  sdk.AccAddress `json:"owner"`
	Group  sdk.AccAddress `json:"group"`
	Members []Member `json:"members,omitempty"`
}

// group accounts get their own sdk.AccAddress
type MsgCreateGroupPolicy struct {
	Owner          sdk.AccAddress `json:"owner"`
	Group          sdk.AccAddress `json:"group"`
	DecisionPolicy DecisionPolicy `json:"decision_policy"`
	Description string `json:"Description,omitempty"`
}

// group accounts get their own sdk.AccAddress
type MsgGroupPolicyGrant struct {
	Owner          sdk.AccAddress `json:"owner"`
    Policy         GroupPolicyID
    Capability     Capability    
	Expiration time.Time      `json:"expiration"`
}

type MsgGroupPolicyRevoke struct {
	Owner          sdk.AccAddress `json:"owner"`
    Policy         GroupPolicyID
    MsgType sdk.Msg        `json:"msg_type"`
}

type MsgUpdateGroupPolicy struct {
	Owner  sdk.AccAddress `json:"owner"`
    Policy         GroupPolicyID
    DecisionPolicy DecisionPolicy `json:"decision_policy"`
}

type MsgUpdateGroupPolicyDescription struct {
	Owner  sdk.AccAddress `json:"owner"`
    Policy         GroupPolicyID
    DecisionPolicy DecisionPolicy `json:"decision_policy"`
}

type MsgDeleteGroupPolicy struct {
	Owner  sdk.AccAddress `json:"owner"`
    Policy         GroupPolicyID
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

type MsgDelegateVote struct {
	Delegator  sdk.AccAddress `json:"delegator"`
	Delegate  sdk.AccAddress `json:"delegate"`
}

type MsgUndelegateVote struct {
	Delegator  sdk.AccAddress `json:"delegator"`
}

type MsgDeposit struct {
	ProposalID ProposalID     `json:"proposal_id"`
	Depositor  sdk.AccAddress `json:"depositor"`
    Deposit sdk.Coins
}
```
