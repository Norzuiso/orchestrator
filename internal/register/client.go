package register

import pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

type ConnectedClients struct {
	ActiveClients map[string]*pb.Client
}
