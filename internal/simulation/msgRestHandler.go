package simulation

import (
	"fmt"
	"net/http"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
)

func (s *SimulationEngine) GetActiveClients(w http.ResponseWriter, req *http.Request) {

}

func (s *SimulationEngine) GetAllClientToClientConnection(w http.ResponseWriter, req *http.Request) {

}

func (s *SimulationEngine) GetClientToClientConnection(w http.ResponseWriter, req *http.Request) {

}

func (s *SimulationEngine) GetClientStatus(w http.ResponseWriter, req *http.Request) {

}

func (s *SimulationEngine) SendMsgToClient(w http.ResponseWriter, req *http.Request) {
	rBody := &models.MsgToClient{}

	utils.ParseBody(req, rBody)
	client, ok := s.GrpcServer.ClientStreams[rBody.ClientId]
	if !ok {
		w.Header().Set("Content-type", "pkgplication/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, "Correct request, but Client %v does not exist or does not have an active connection.", client)

		return
	}

	// rMsg := rBody.Msg.MsgRequestToMessage(s.Orchestrator.Epoch, 0)
	// client.MsgHandler.SendMsg(rMsg)

	w.Header().Set("Content-type", "pkgplication/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Msg send it")
}

func (s *SimulationEngine) SendMsgToAllClients(w http.ResponseWriter, req *http.Request) {
	rBody := &models.MsgToAllClients{}

	utils.ParseBody(req, rBody)

	// for _, client := range s.GrpcServer.ClientStreams {
	// 	rMsg := rBody.Msg.MsgRequestToMessage(s.Orchestrator.Epoch, 0)
	// 	err := client.MsgHandler.SendMsg(rMsg)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Error on client: %v.\nerror:%v", client, err)
	// 	}
	// }
	w.Header().Set("Content-type", "pkgplication/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Msg send it")
}

func (s *SimulationEngine) SendMsgToClientsList(w http.ResponseWriter, req *http.Request) {
	rBody := &models.MsgToClientsList{}
	utils.ParseBody(req, rBody)

	// for _, rClient := range rBody.ClientIds {
	// 	client, ok := s.GrpcServer.ClientStreams[rClient]
	// 	if !ok {
	// 		continue
	// 	}
	// 	rMsg := rBody.Msg.MsgRequestToMessage(s.Orchestrator.Epoch, 0)
	// 	err := client.MsgHandler.SendMsg(rMsg)
	// 	if err != nil {
	// 		fmt.Fprintf(w, "Error on client: %v.\nerror:%v", client, err)
	// 	}
	// }
	w.Header().Set("Content-type", "pkgplication/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Msg send it")
}
