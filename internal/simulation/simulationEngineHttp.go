package simulation

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"

	pb "github.com/Norzuiso/protocol/gen/go/proto/orchestrator/v1"

	"github.com/Norzuiso/orchestrator/internal/models"
	"github.com/Norzuiso/orchestrator/internal/utils"
	"google.golang.org/protobuf/encoding/protojson"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *SimulationEngine) HttpConnect() {
	mux := http.NewServeMux()
	// msg
	mux.HandleFunc("POST /msg/all", s.SendMsgToAllClients)           // TODO
	mux.HandleFunc("POST /msg/client", s.SendMsgToClient)            // TODO
	mux.HandleFunc("POST /msg/clients/list", s.SendMsgToClientsList) // TODO
	mux.HandleFunc("POST /msg/client-to-client", s.RegisterClientToClientconnection)

	// clients
	mux.HandleFunc("GET /client/active", s.GetActiveClients)
	mux.HandleFunc("GET /client/clients-to-clients", s.GetAllClientToClientConnection)
	mux.HandleFunc("GET /client/client-to-client", s.GetClientToClientConnection)
	mux.HandleFunc("GET /client/open-streams", s.GetClientIdsWithOpenStreams)
	mux.HandleFunc("GET /client/info", s.GetClientInfoById)
	mux.HandleFunc("GET /client/open-streams/info", s.GetOpenStreamsClientsInfo)

	// simulation
	mux.HandleFunc("POST /simulation/start", s.StartSimulation)
	mux.HandleFunc("GET /simulation/pause", s.StopSimulation) // TODO
	mux.HandleFunc("GET /simulation/state/waiting-connections", s.WaitingConnections)
	mux.HandleFunc("GET /simulation/end", s.HttpEndSimulation)

	mux.HandleFunc("GET /simulation/next-phase", s.NextStateHttp) // TODO
	mux.HandleFunc("GET /simulation/next-epoch", s.NextEpoch)

	// Logs
	mux.HandleFunc("/logs/stream", s.StreamLogs)
	mux.HandleFunc("/client/open-streams/stream", s.StreamOpenStreamsClients)

	err := http.ListenAndServe(":8090", corsMiddleware(mux)) // TODO - Port Get it from config

	if err != nil {
		panic(err)
	}
}

func (s *SimulationEngine) StreamOpenStreamsClients(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Not supported streaming", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-req.Context().Done():
			return
		case l := <-s.OpenStreamChannel:
			fmt.Fprintf(w, "data: %s\n\n", l.String())
			flusher.Flush()
		}
	}
}

func (s *SimulationEngine) StreamLogs(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Not supported streaming", http.StatusInternalServerError)
		return
	}

	for {
		select {
		case <-req.Context().Done():
			return
		case l := <-s.LogChannel:
			log.Println(l)
			fmt.Fprintf(w, "data: %s\n\n", l)
			flusher.Flush()
		}
	}
}

func (s *SimulationEngine) Write(str string) {
	s.LogChannel <- str
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

func (s *SimulationEngine) GetClientIdsWithOpenStreams(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	keys := slices.Collect(maps.Keys(s.ClientService.ClientStreams))
	json.NewEncoder(w).Encode(keys)
}

func (s *SimulationEngine) GetClientToClientConnection(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	clientToClients, err := s.Storage.ClientToClientGet(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	data, err := protojson.Marshal(clientToClients)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func (s *SimulationEngine) GetOpenStreamsClientsInfo(w http.ResponseWriter, req *http.Request) {
	activeClientsIds := slices.Collect(maps.Keys(s.ClientService.ClientStreams))
	response := make([]*models.ClientInfo, 0)
	for _, id := range activeClientsIds {
		clientInfo, err := s.GetClientInfo(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		response = append(response, clientInfo)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *SimulationEngine) GetClientInfo(id int64) (*models.ClientInfo, error) {
	client, err := s.Storage.ActiveClientsGet(id)
	if err != nil {
		return nil, err
	}

	clientData, err := protojson.Marshal(client)
	if err != nil {
		return nil, err
	}

	clientToClients, err := s.Storage.ClientToClientGet(id)
	if err != nil {
		if err.Error() != "Client not found" {
			return nil, err
		}
		clientToClients = &pb.ClientConnectionList{}
	}
	ctcData, err := protojson.Marshal(clientToClients)
	if err != nil {
		return nil, err
	}

	_, hasStream := s.ClientService.ClientStreams[id]
	return &models.ClientInfo{
		HasOpenStream: hasStream,
		Client:        clientData,
		Connections:   ctcData,
	}, nil
}

func (s *SimulationEngine) GetClientInfoById(w http.ResponseWriter, req *http.Request) {
	idStr := req.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	response, err := s.GetClientInfo(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	s.Orchestrator.Epoch = init
	s.Orchestrator.MaxOfEpochs = end
	s.Orchestrator.SetSeed(rBody.Seed)
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
