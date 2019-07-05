package demoserver1

import (
	"context"

	pb "sniper/rpc/demo/v1"
)

// Server 实现 /twirp/demo.v1.Demo 服务
type Server struct{}

// Hello 实现 /twirp/demo.v1.Demo/Hello 接口
func (s *Server) Hello(ctx context.Context, req *pb.Req) (resp *pb.Resp, err error) {
	// FIXME 请开始你的表演
	return
}
