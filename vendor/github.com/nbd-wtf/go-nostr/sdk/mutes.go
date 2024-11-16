package sdk

import "context"

type MuteList = GenericList[Follow]

func (sys *System) FetchMuteList(ctx context.Context, pubkey string) MuteList {
	ml, _ := fetchGenericList[Follow](sys, ctx, pubkey, 10000, parseFollow, nil, false)
	return ml
}
