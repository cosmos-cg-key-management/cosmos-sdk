package group

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	groupTxCmd := &cobra.Command{
		Use:   ModuleName,
		Short: "Group transactions subcommands",
	}

	groupTxCmd.AddCommand(client.PostCommands(
		GetCmdApprove(cdc),
		GetCmdCreateGroup(cdc),
		GetCmdPropose(cdc),
		GetCmdTryExec(cdc),
		GetCmdUnapprove(cdc),
		GetCmdWithdraw(cdc),
	)...)

	return groupTxCmd
}

func membersFromArray(arr []string) []Member {
	n := len(arr)
	res := make([]Member, n)
	for i := 0; i < n; i++ {
		strs := strings.Split(arr[i], "=")
		if len(strs) <= 0 {
			panic("empty array")
		}
		acc, err := sdk.AccAddressFromBech32(strs[0])
		if err != nil {
			panic(err)
		}
		mem := Member{
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
    var memo string
	var ownerStr string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create an group",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))

			account := cliCtx.GetFromAddress()

			var ownerAddr sdk.AccAddress
			var err error
			if len(ownerStr) != 0 {
				ownerAddr, err = sdk.AccAddressFromBech32(ownerStr)
				if err != nil {
					return err
				}
			}

			msg := NewMsgCreateGroup(account, ownerAddr, membersFromArray(members), sdk.NewInt(threshold), memo)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			cliCtx.PrintResponse = true

			return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int64Var(&threshold, "decision-threshold", 1, "Decision threshold")
	cmd.Flags().StringArrayVar(&members, "members", []string{}, "Members")
	cmd.Flags().StringVar(&memo, "memo", "", "Memo")
	cmd.Flags().StringVar(&ownerStr, "owner", "", "The group owner who can change the group composition, leave blank for the group to own itself")

	return cmd
}

func GetCmdPropose(cdc *codec.Codec) *cobra.Command {
	var exec bool

	cmd := &cobra.Command{
		Use:   "propose [group-id] [proposal-json]",
		Short: "make a proposal, where proposal-json is any Msg encoded as json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)

			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			account := cliCtx.GetFromAddress()

			group, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var proposal sdk.Msg
			err = cdc.UnmarshalJSON([]byte(args[1]), &proposal)
			if err != nil {
				return err
			}

			msg := MsgCreateProposal{
				Proposer: account,
				Group:    group,
				Msgs:     []sdk.Msg{proposal},
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

		id := MustDecodeProposalIDBech32(args[0])

		msg := MsgVote{
			ProposalID: id,
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

			id := MustDecodeProposalIDBech32(args[0])

			msg := MsgTryExecuteProposal{
				ProposalID: id,
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

			id := MustDecodeProposalIDBech32(args[0])

			msg := MsgWithdrawProposal{
				ProposalID: id,
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
