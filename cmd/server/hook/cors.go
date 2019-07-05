package hook

import (
	"context"
	"github.com/bilibili/twirp"
)

// NewCors 解决跨域问题
func NewCors() *twirp.ServerHooks {
	return &twirp.ServerHooks{
		RequestReceived: func(ctx context.Context) (context.Context, error) {
			_ = twirp.SetHTTPResponseHeader(ctx, "Access-Control-Allow-Origin", "*")
			_ = twirp.SetHTTPResponseHeader(ctx, "Access-Control-Allow-Headers", "Content-Type")
			_ = twirp.SetHTTPResponseHeader(ctx, "content-type", "application/json")
			return ctx, nil
		},
	}
}
