package group

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	clientrest "github.com/cosmos/cosmos-sdk/client/rest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes registers staking-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	// Withdraw all delegator rewards
	r.HandleFunc(
		"/group/create",
		createdGroupHandlerFn(cliCtx),
	).Methods("POST")
}

type createGroupReq struct {
	BaseReq rest.BaseReq `json:"base_req"`
	Members []string     `json:"members"`
}

func createdGroupHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createGroupReq

		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		var members []Member
		for _, memberStr := range req.Members {
			memberAddr, _ := sdk.AccAddressFromBech32(memberStr)
			member := Member{
				Address: memberAddr,
				Weight:  sdk.NewInt(10),
			}
			members = append(members, member)
		}

		signer := cliCtx.GetFromAddress()
		msg := NewMsgCreateGroup(signer, nil, members, sdk.NewInt(10), "")

		clientrest.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/group/groups_by_member/{memberAddr}",
		memberGroupsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/group/proposals_by_group_id/{groupId}",
		groupProposalsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/group/groups",
		groupsHandlerFn(cliCtx),
	).Methods("GET")
}
func groupsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", "group", "groups")

		res, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func memberGroupsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		memberAddr := vars["memberAddr"]
		route := fmt.Sprintf("custom/%s/%s", "group", "groups_by_member")

		decodedAddr, _ := sdk.AccAddressFromBech32(memberAddr)
		params := QueryGroupsByMemberParams{
			Address: decodedAddr,
		}

		bz, _ := cliCtx.Codec.MarshalJSON(params)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func groupProposalsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		memberAddr := vars["groupId"]
		route := fmt.Sprintf("custom/%s/%s", "group", "proposals_by_group_id")

		decodedAddr, _ := sdk.AccAddressFromBech32(memberAddr)
		params := QueryProposalsByGroupIDrParams{
			Address: decodedAddr,
		}

		bz, _ := cliCtx.Codec.MarshalJSON(params)
		res, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
