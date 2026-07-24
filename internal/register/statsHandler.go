package register

import (
	"context"

	"google.golang.org/grpc/stats"
)

type StatsHandler struct{}

func (h *StatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	return ctx
}

func (h *StatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {

}

func (h *StatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}
func (h *StatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {

}
