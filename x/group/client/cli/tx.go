package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/delegate"
	"github.com/cosmos/cosmos-sdk/x/group"
	"github.com/spf13/cobra"
	"strings"
)

func membersFromArray(arr []string) []group.Member {
	n := len(arr)
	res := make([]group.Member, n)
	for i := 0; i < n; i++ {
		strs := strings.Split(arr[i], "=")
		if len(strs) <= 0 {
			panic("empty array")
		}
		acc, err := sdk.AccAddressFromBech32(strs[0])
		if err != nil {
			panic(err)
		}
		mem := group.Member{
			Address: acc,
		}
		if len(strs) == 2 {
			var ok bool
			mem.Weight, ok = sdk.NewIntFromString(strs[1])
			if !ok {
				panic(fmt.Errorf("invalid weight: %s", strs[i]))
			}
		} else {
			mem.Weight = sdk.NewInt(1)
		}
		res[i] = mem
	}
	return res
}

func GetCmdCreateGroup(cdc *codec.Codec) *cobra.Command {
	var threshold int64
	var members []string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create an group",
		//Args:  cobra.MinimumNArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			info := group.Group{
				Members:           membersFromArray(members),
				DecisionThreshold: sdk.NewInt(threshold),
			}

			msg := group.NewMsgCreateGroup(info, account)
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int64Var(&threshold, "decision-threshold", 1, "Decision threshold")
	cmd.Flags().StringArrayVar(&members, "members", []string{}, "Members")

	return cmd
}

type ActionCreator func(cmd *cobra.Command, args []string) (delegate.Action, error)

func GetCmdPropose(cdc *codec.Codec, actionCreator ActionCreator) *cobra.Command {
	var exec bool

	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			action, err := actionCreator(cmd, args)

			if err != nil {
				return err
			}

			msg := group.MsgCreateProposal{
				Proposer: account,
				Action:   action,
				Exec:     exec,
			}
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().BoolVar(&exec, "exec", false, "try to execute the proposal immediately")
	return cmd
}

func getRunVote(cdc *codec.Codec, approve bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

		txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

		if err := cliCtx.EnsureAccountExists(); err != nil {
			return err
		}

		account := cliCtx.GetFromAddress()

		id := group.MustDecodeProposalIDBech32(args[0])

		msg := group.MsgVote{
			ProposalId: id,
			Voter:      account,
			Vote:       approve,
		}
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}

		cliCtx.PrintResponse = true

		return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
	}
}

func GetCmdApprove(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "approve [ID]",
		Short: "vote to approve a proposal",
		Args:  cobra.ExactArgs(1),
		RunE:  getRunVote(cdc, true),
	}
}

func GetCmdUnapprove(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unapprove [ID]",
		Short: "vote to un-approve a proposal that you have previously approved",
		Args:  cobra.ExactArgs(1),
		RunE:  getRunVote(cdc, false),
	}
}

func GetCmdTryExec(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "try-exec [ID]",
		Short: "try to execute the proposal (will fail if not enough signers have approved it)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			id := group.MustDecodeProposalIDBech32(args[0])

			msg := group.MsgTryExecuteProposal{
				ProposalId: id,
				Signer:     account,
			}
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}

func GetCmdWithdraw(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw [ID]",
		Short: "withdraw a proposer that you previously proposed",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			account := cliCtx.GetFromAddress()

			id := group.MustDecodeProposalIDBech32(args[0])

			msg := group.MsgWithdrawProposal{
				ProposalId: id,
				Proposer:   account,
			}
			err := msg.ValidateBasic()
			if err != nil {
				return err
			}

			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
}
