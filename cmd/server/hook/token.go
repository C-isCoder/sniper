package hook

import (
	"context"
	"github.com/bilibili/twirp"
	"strings"
	"sniper/util/auth"
	"sniper/util/log"
)

func NewToken() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			r, ok := twirp.Request(ctx)
			if !ok {
				return ctx, twirp.InternalError("get request error")
			}
			aut := r.Header.Get("Authorization")
			if aut == "" {
				return ctx, auth.NotTokenError
			}
			log.Get(ctx).Debugln("aut: ", aut)
			arr := strings.Split(aut, " ")
			if len(arr) != 2 && arr[0] != "Bearer" {
				return ctx, auth.FailTokenError
			}
			token := arr[1]
			log.Get(ctx).Debugln("token -> ", token)
			ok, err := auth.JWT(token).VerifyToken()
			if ok {
				return ctx, nil
			} else {
				return ctx, err
			}
		},
	}
}
