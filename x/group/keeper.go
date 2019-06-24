package group

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/delegation"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bech32"
)

type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	accountKeeper auth.AccountKeeper
	dispatcher    delegation.Keeper
}

func NewKeeper(groupStoreKey sdk.StoreKey, cdc *codec.Codec, accountKeeper auth.AccountKeeper, dispatcher delegation.Keeper) Keeper {
	return Keeper{
		groupStoreKey,
		cdc,
		accountKeeper,
		dispatcher,
	}
}

type GroupAccount struct {
	*auth.BaseAccount
}

func (acc *GroupAccount) SetPubKey(pubKey crypto.PubKey) error {
	return fmt.Errorf("cannot set a PubKey on a Group account")
}

var (
	keyNewGroupID    = []byte("newGroupID")
	keyNewProposalID = []byte("newProposalID")
)

func KeyGroupID(id sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("g/%x", id))
}

func KeyGroupIDByMemberAddress(addr sdk.AccAddress, id sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("g/%x/%x", addr, id))
}

func KeyProposalsByGroupID(groupID sdk.AccAddress, proposalID ProposalID) []byte {
	return []byte(fmt.Sprintf("p/%x/%x", groupID, proposalID))
}

func KeyProposal(id ProposalID) []byte {
	return []byte(fmt.Sprintf("p/%x", id))
}

func (keeper Keeper) GetGroupInfo(ctx sdk.Context, id sdk.AccAddress) (info Group, err sdk.Error) {
	if len(id) < 1 || id[0] != 'G' {
		return info, sdk.ErrUnknownRequest("Not a valid group")
	}
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyGroupID(id))
	if bz == nil {
		return info, sdk.ErrUnknownRequest("Not found")
	}
	info = Group{}
	marshalErr := keeper.cdc.UnmarshalBinaryBare(bz, &info)
	if marshalErr != nil {
		return info, sdk.ErrUnknownRequest(marshalErr.Error())
	}
	return info, nil
}

// GetGroups gets all groups
func (keeper Keeper) GetGroups(ctx sdk.Context) []Group {
	prefix := fmt.Sprintf("g/")
	prefixBytes := []byte(prefix)
	store := ctx.KVStore(keeper.storeKey)
	var groups []Group
	iter := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var group Group
		// keeper.cdc.MustUnmarshalBinaryBare(iter.Value(), &group)
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iter.Value(), &group)
		groups = append(groups, group)
	}

	return groups
}

// GetGroupsByMemberAddress get groups that I'm a member of
// key: g/%x/%x
// g/[member address]/[group id] -> [group id]
func (keeper Keeper) GetGroupsByMemberAddress(ctx sdk.Context, memberAddr sdk.AccAddress) []Group {
	prefix := fmt.Sprintf("g/%x/", memberAddr)
	prefixBytes := []byte(prefix)
	store := ctx.KVStore(keeper.storeKey)
	var groups []Group
	iter := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		groupID := iter.Value()
		group, err := keeper.GetGroupInfo(ctx, groupID)
		if err != nil {
			panic(err)
		}
		groups = append(groups, group)
	}

	return groups
}

func (keeper Keeper) GetProposalsByGroupID(ctx sdk.Context, groupID sdk.AccAddress) []Proposal {
	prefix := fmt.Sprintf("p/%x/", groupID)
	prefixBytes := []byte(prefix)
	store := ctx.KVStore(keeper.storeKey)
	var proposals []Proposal
	iter := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var proposal Proposal
		keeper.cdc.MustUnmarshalBinaryBare(iter.Value(), &proposal)
		proposals = append(proposals, proposal)
	}

	return proposals
}

func addrFromUint64(id uint64) sdk.AccAddress {
	addr := make([]byte, binary.MaxVarintLen64+1)
	addr[0] = 'G'
	n := binary.PutUvarint(addr[1:], id)
	return addr[:n+1]
}

func (keeper Keeper) getNewGroupId(ctx sdk.Context) sdk.AccAddress {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(keyNewGroupID)
	var groupId uint64 = 0
	if bz != nil {
		keeper.cdc.MustUnmarshalBinaryBare(bz, &groupId)
	}
	bz = keeper.cdc.MustMarshalBinaryBare(groupId + 1)
	store.Set(keyNewGroupID, bz)
	return addrFromUint64(groupId)
}

func (keeper Keeper) CreateGroup(ctx sdk.Context, members []Member, decisionThreshold sdk.Int, owner sdk.AccAddress, memo string) (sdk.AccAddress, sdk.Error) {
	id := keeper.getNewGroupId(ctx)
	store := ctx.KVStore(keeper.storeKey)
	// iterate through members in group and add member <-> group id association
	for _, member := range members {
		key := KeyGroupIDByMemberAddress(member.Address, id)
		store.Set(key, id)
	}
	acct := &GroupAccount{
		BaseAccount: &auth.BaseAccount{
			Address: id,
		},
	}
	existingAcc := keeper.accountKeeper.GetAccount(ctx, id)
	if existingAcc != nil {
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("account with address %s already exists", id.String()))
	}
	keeper.accountKeeper.SetAccount(ctx, acct)

	return id, nil
}

func (keeper Keeper) Authorize(ctx sdk.Context, group sdk.AccAddress, signers []sdk.AccAddress) bool {
	info, err := keeper.GetGroupInfo(ctx, group)
	if err != nil {
		return false
	}
	ctx.GasMeter().ConsumeGas(10, "group auth")
	return keeper.AuthorizeGroupInfo(ctx, &info, signers)
}

func (keeper Keeper) AuthorizeGroupInfo(ctx sdk.Context, info *Group, signers []sdk.AccAddress) bool {
	voteCount := sdk.NewInt(0)
	sigThreshold := info.DecisionThreshold

	nMembers := len(info.Members)
	nSigners := len(signers)
	for i := 0; i < nMembers; i++ {
		mem := info.Members[i]
		// TODO Use a hash map to optimize this
		for j := 0; j < nSigners; j++ {
			ctx.GasMeter().ConsumeGas(10, "check addr")
			if bytes.Compare(mem.Address, signers[j]) == 0 || keeper.Authorize(ctx, mem.Address, signers) {
				voteCount = voteCount.Add(mem.Weight)
				diff := voteCount.Sub(sigThreshold)
				if diff.IsZero() || diff.IsPositive() {
					return true
				}
				break
			}
		}
	}
	return false
}

const (
	Bech32Prefix = "proposal"
)

func mustEncodeProposalIDBech32(id ProposalID) string {
	bz := make([]byte, binary.MaxVarintLen64+1)
	n := binary.PutUvarint(bz[0:], uint64(id))
	str, err := bech32.ConvertAndEncode(Bech32Prefix, bz[:n])
	if err != nil {
		panic(err)
	}
	return str
}

func MustDecodeProposalIDBech32(bech string) ProposalID {
	hrp, data, err := bech32.DecodeAndConvert(bech)
	if err != nil {
		panic(err)
	}
	if hrp != Bech32Prefix {
		panic(fmt.Sprintf("Expected bech32 prefix %s", Bech32Prefix))
	}
	id, err := binary.ReadUvarint(bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	return ProposalID(id)
}

func (keeper Keeper) getNewProposalId(ctx sdk.Context) ProposalID {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(keyNewProposalID)
	var id uint64 = 0
	if bz != nil {
		keeper.cdc.MustUnmarshalBinaryBare(bz, &id)
	}
	bz = keeper.cdc.MustMarshalBinaryBare(id + 1)
	store.Set(keyNewProposalID, bz)
	return ProposalID(id)
}

func (keeper Keeper) Propose(ctx sdk.Context, proposer sdk.AccAddress, group sdk.AccAddress, msgs []sdk.Msg) (ProposalID, sdk.Result) {
	id := keeper.getNewProposalId(ctx)

	prop := Proposal{
		Group:     group,
		Proposer:  proposer,
		Msgs:      msgs,
		Approvers: []sdk.AccAddress{proposer},
	}

	keeper.storeProposal(ctx, id, &prop)

	res := sdk.Result{}
	res.Tags = res.Tags.
		AppendTag("proposal.id", mustEncodeProposalIDBech32(id))
	return id, res
}

func (keeper Keeper) storeProposal(ctx sdk.Context, id ProposalID, proposal *Proposal) {
	store := ctx.KVStore(keeper.storeKey)
	bz, err := keeper.cdc.MarshalBinaryBare(proposal)
	if err != nil {
		panic(err)
	}

	store.Set(KeyProposal(id), bz)
	store.Set(KeyProposalsByGroupID(proposal.Group, id), bz)
}

func (keeper Keeper) GetProposal(ctx sdk.Context, id ProposalID) (proposal *Proposal, err sdk.Error) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(KeyProposal(id))
	proposal = &Proposal{}
	marshalErr := keeper.cdc.UnmarshalBinaryBare(bz, proposal)
	if marshalErr != nil {
		return proposal, sdk.ErrUnknownRequest(marshalErr.Error())
	}
	return proposal, nil
}

func (keeper Keeper) Vote(ctx sdk.Context, proposalId ProposalID, voter sdk.AccAddress, yesNo bool) sdk.Result {
	proposal, err := keeper.GetProposal(ctx, proposalId)

	if err != nil {
		return sdk.Result{
			Code: sdk.CodeUnknownRequest,
			Log:  "can't find proposal",
		}
	}

	var newVotes []sdk.AccAddress
	votes := proposal.Approvers
	nVotes := len(votes)

	if yesNo {
		newVotes = make([]sdk.AccAddress, nVotes+1)
		for i := 0; i < nVotes; i++ {
			oldVoter := votes[i]
			if bytes.Equal(voter, oldVoter) {
				// Already voted YES
				return sdk.Result{
					Code: sdk.CodeUnknownRequest,
					Log:  "already voted yes",
				}
			}
			newVotes[i] = oldVoter
		}
		newVotes[nVotes] = voter
	} else {
		newVotes = make([]sdk.AccAddress, nVotes)
		didntVote := true
		j := 0
		for i := 0; i < nVotes; i++ {
			oldVoter := votes[i]
			if bytes.Equal(voter, oldVoter) {
				didntVote = false
			} else {
				newVotes[j] = oldVoter
				j++
			}
		}
		if didntVote {
			return sdk.Result{
				Code: sdk.CodeUnknownRequest,
				Log:  "didn't vote yes previously",
			}
		}
		if j != nVotes-1 {
			panic("unexpected vote count")
		}
		newVotes = newVotes[:j]
	}

	newProp := Proposal{
		Proposer:  proposal.Proposer,
		Msgs:      proposal.Msgs,
		Approvers: newVotes,
	}

	keeper.storeProposal(ctx, proposalId, &newProp)

	return sdk.Result{Code: sdk.CodeOK,
		Tags: sdk.EmptyTags().
			AppendTag("proposal.id", mustEncodeProposalIDBech32(proposalId)),
	}
}

func (keeper Keeper) TryExecute(ctx sdk.Context, proposalId ProposalID) sdk.Result {
	proposal, err := keeper.GetProposal(ctx, proposalId)

	if err != nil {
		return sdk.ErrUnknownRequest("can't find proposal").Result()
	}

	if !keeper.Authorize(ctx, proposal.Group, proposal.Approvers) {
		return sdk.ErrUnauthorized("proposal failed").Result()
	}

	res := keeper.dispatcher.DispatchActions(ctx, proposal.Group, proposal.Msgs)

	if res.Code == sdk.CodeOK {
		store := ctx.KVStore(keeper.storeKey)
		store.Delete(KeyProposal(proposalId))
	}

	return res
}

func (keeper Keeper) Withdraw(ctx sdk.Context, proposalId ProposalID, proposer sdk.AccAddress) sdk.Result {
	proposal, err := keeper.GetProposal(ctx, proposalId)

	if err != nil {
		return sdk.Result{
			Code: sdk.CodeUnknownRequest,
			Log:  "can't find proposal",
		}
	}

	if !bytes.Equal(proposer, proposal.Proposer) {
		return sdk.Result{
			Code: sdk.CodeUnauthorized,
			Log:  "you didn't propose this",
		}
	}

	store := ctx.KVStore(keeper.storeKey)
	store.Delete(KeyProposal(proposalId))

	return sdk.Result{Code: sdk.CodeOK,
		Tags: sdk.EmptyTags().
			AppendTag("proposal.id", mustEncodeProposalIDBech32(proposalId)),
	}
}
