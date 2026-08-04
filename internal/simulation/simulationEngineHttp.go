package simulation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	"google.golang.org/protobuf/encoding/protojson"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

func (s *SimulationEngine) HttpConnect() {
	// msg
	http.HandleFunc("POST /msg/all", s.SendMsgToAllClients)           // TODO
	http.HandleFunc("POST /msg/client", s.SendMsgToClient)            // TODO
	http.HandleFunc("POST /msg/clients/list", s.SendMsgToClientsList) // TODO
	http.HandleFunc("POST /msg/client-to-client", s.RegisterClientToClientconnection)

	// clients
	http.HandleFunc("GET /client/active", s.GetActiveClients)
	http.HandleFunc("GET /client/client-to-client", s.GetAllClientToClientConnection)

	// simulation
	http.HandleFunc("POST /simulation/start", s.StartSimulation)
	http.HandleFunc("GET /simulation/pause", s.StopSimulation) // TODO
	http.HandleFunc("GET /simulation/state/waiting-connections", s.WaitingConnections)
	http.HandleFunc("GET /simulation/end", s.HttpEndSimulation)

	http.HandleFunc("GET /simulation/next-phase", s.NextStateHttp) // TODO
	http.HandleFunc("GET /simulation/next-epoch", s.NextEpoch)

	err := http.ListenAndServe(":8090", nil) // TODO - Port Get it from config

	if err != nil {
		panic(err)
	}
}

func (s *SimulationEngine) GetActiveClients(w http.ResponseWriter, req *http.Request) {
	activeClients, err := s.Storage.GetAllActiveClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(activeClients)
}

func (s *SimulationEngine) GetAllClientToClientConnection(w http.ResponseWriter, req *http.Request) {
	clientToClients, err := s.Storage.GetAllClientToClientConnection()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("["))

	for i, ctc := range clientToClients {
		if i > 0 {
			w.Write([]byte(","))
		}
		data, err := protojson.Marshal(ctc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	}

	w.Write([]byte("]"))
}

func (s *SimulationEngine) GetClientToClientConnection(w http.ResponseWriter, req *http.Request) {

}

func (s *SimulationEngine) GetClientStatus(w http.ResponseWriter, req *http.Request) {

}

func (s *SimulationEngine) RegisterClientToClientconnection(w http.ResponseWriter, req *http.Request) {
	rBody := &pb.RegisterConnectionRequest{}
	body, _ := io.ReadAll(req.Body)

	if err := protojson.Unmarshal(body, rBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := s.ClientService.RegisterConnection(req.Context(), rBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "%s", res.String())
}

func (s *SimulationEngine) SendMsgToClient(w http.ResponseWriter, req *http.Request) {
	rBody := &models.MsgToClient{}

	utils.ParseBody(req, rBody)
	client, ok := s.ClientService.ClientStreams[rBody.ClientId]
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

// Orchestrator
func (s *SimulationEngine) StartSimulation(w http.ResponseWriter, req *http.Request) {
	rBody := &models.StartSimualtionReq{}
	utils.ParseBody(req, rBody)

	init := rBody.InitEpoch
	end := rBody.EndEpoch

	if end <= init {
		w.Header().Set("Content-type", "pkgplication/json")
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintln(w, "Init epoch need to be greater than end epoch")
		return
	}

	if init < 0 {
		w.Header().Set("Content-type", "pkgplication/json")
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprintln(w, "Init epoch need to be greater than 0")
		return
	}

	s.Orchestrator.Epoch = (init - 1)
	s.Orchestrator.MaxOfEpochs = end
	s.Orchestrator.Seed = rBody.Seed
	s.Orchestrator.StepsMode = rBody.StepsMode
	if _, err := fmt.Fprintf(w, "Simulation started."); err != nil {
		return
	}
	s.NextState()
}

func (s *SimulationEngine) HttpEndSimulation(w http.ResponseWriter, req *http.Request) {
	s.EndSimulation()
	fmt.Fprintf(w, "End of simulation")
}

func (s *SimulationEngine) WaitingConnections(w http.ResponseWriter, req *http.Request) {
	s.CurrentState = s.StateWaitingConnections
	s.CurrentState.StartState()
	fmt.Fprintf(w, "Waiting connections")
}

func (s *SimulationEngine) StopSimulation(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.PauseSimualtion()
	fmt.Fprintf(w, "Simulation stopped on \n\t epoch: \t%v \n\tphase:\t%v\n\t.", s.Orchestrator.Epoch, s.CurrentState.GetStateName())
}

func (s *SimulationEngine) NextStateHttp(w http.ResponseWriter, req *http.Request) {
	s.NextState()
	fmt.Fprintf(w, "State updated. Current State: %v", s.CurrentState.GetStateName())
}

func (s *SimulationEngine) NextEpoch(w http.ResponseWriter, req *http.Request) {
	s.StartEpoch()
	fmt.Fprintf(w, "Epoch: %v", s.Orchestrator.Epoch)
}
