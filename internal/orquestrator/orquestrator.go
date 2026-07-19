package orquestrator

import (
	"net/http"

	Phase "github.com/Norzuiso/orchestrator/internal/orquestrator/phase"
)

type Orquestrator struct {
	epoch float64

	currentPhase    Phase.Phase // 0 : Hasta que doneClients == connectedClients y pendingClients == 0 se hace cambio de phase
	conectedClients int         // 10
	doneClients     int         // 4
	pendingClients  int         // 6

	serverAddr string
	httpSrv    *http.Server
}

// Orquestrator phases match with message type
// "WAITING_CONNECTIONS" =
//    	MESSAGE_TYPE_OPEN_STREAM = 5;
/**
RequestEvents =
	Orquestrator send: 	MESSAGE_TYPE_RECOLECT_EVENTS
	Client send: 		MESSAGE_TYPE_SEND_EVENTS

    ==================
	What happen when the client sends the event before to be requested?
	-------------------------------------------------------------------
	Client send: 		MESSAGE_TYPE_SEND_EVENTS
	Orquestrator send: 	MESSAGE_TYPE_RECOLECT_EVENTS
	-------------------------------------------------------------------

	Once the orquestrator recolect an event from a client. it sends the event
    to other clients that apply/react
    MESSAGE_TYPE_APPLY_EVENT =4;
RecolectStatus =
	Client send: 	MESSAGE_TYPE_SEND_EVENTS =  Rejected.
												Message type not allow it in phase RecolectStatus

**/

func NewOrquestrator(serverAddr string) *Orquestrator {
	o := &Orquestrator{serverAddr: serverAddr}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /simulations/start", o.handleStart)
	mux.HandleFunc("GET /simulations/stop", o.handleStop)

	o.httpSrv = &http.Server{
		Addr:    serverAddr,
		Handler: mux,
	}
	return o
}

func (o *Orquestrator) ListenAndServe() error {
	return o.httpSrv.ListenAndServe()
}

func (o *Orquestrator) handleStop(w http.ResponseWriter, r *http.Request) {
	o.currentPhase.PausePhase()
}

func (o *Orquestrator) handleStart(w http.ResponseWriter, r *http.Request) {
	o.StartSimualtion()
}

func (o *Orquestrator) StartSimualtion() {
	o.currentPhase.RequestEventPhase()
	o.epoch = 0
}

func (o *Orquestrator) GetPhase() *Phase.Phase {
	return &o.currentPhase
}

func (o *Orquestrator) ChangePhase() {

}

func (o *Orquestrator) ChangeEpoch() {

}
