// nolint
// autogenerated code using github.com/rigelrozanski/multitool
// aliases generated for the following subdirectories:
// ALIASGEN: github.com/cosmos/cosmos-sdk/x/fee_grant/internal/keeper
// ALIASGEN: github.com/cosmos/cosmos-sdk/x/fee_grant/internal/types
package fee_grant

import (
	"github.com/cosmos/cosmos-sdk/x/fee_grant/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/fee_grant/internal/types"
)

const (
	QueryGetFeeAllowances = keeper.QueryGetFeeAllowances
	DefaultCodespace      = types.DefaultCodespace
	ModuleName            = types.ModuleName
	StoreKey              = types.StoreKey
	RouterKey             = types.RouterKey
	QuerierRoute          = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper                   = keeper.NewKeeper
	NewQuerier                  = keeper.NewQuerier
	RegisterCodec               = types.RegisterCodec
	ExpiresAtTime               = types.ExpiresAtTime
	ExpiresAtHeight             = types.ExpiresAtHeight
	ClockDuration               = types.ClockDuration
	BlockDuration               = types.BlockDuration
	FeeAllowanceKey             = types.FeeAllowanceKey
	FeeAllowancePrefixByGrantee = types.FeeAllowancePrefixByGrantee
	NewMsgGrantFeeAllowance     = types.NewMsgGrantFeeAllowance
	NewMsgRevokeFeeAllowance    = types.NewMsgRevokeFeeAllowance

	// variable aliases
	ModuleCdc             = types.ModuleCdc
	ErrFeeLimitExceeded   = types.ErrFeeLimitExceeded
	ErrFeeLimitExpired    = types.ErrFeeLimitExpired
	ErrInvalidDuration    = types.ErrInvalidDuration
	FeeAllowanceKeyPrefix = types.FeeAllowanceKeyPrefix
)

type (
	Keeper                = keeper.Keeper
	BasicFeeAllowance     = types.BasicFeeAllowance
	ExpiresAt             = types.ExpiresAt
	Duration              = types.Duration
	FeeAllowanceGrant     = types.FeeAllowanceGrant
	MsgGrantFeeAllowance  = types.MsgGrantFeeAllowance
	MsgRevokeFeeAllowance = types.MsgRevokeFeeAllowance
	PeriodicFeeAllowance  = types.PeriodicFeeAllowance
)