package orchestrator

import (
	Phase "github.com/Norzuiso/orchestrator/internal/orchestrator/phase"
)

type Orchestrator struct {
	Epoch float32

	CurrentPhase    *Phase.Phase // 0 : Hasta que doneClients == connectedClients y pendingClients == 0 se hace cambio de phase
	conectedClients int          // 10
	doneClients     int          // 4
	pendingClients  int          // 6
}

func NewOrquestrator() *Orchestrator {
	o := &Orchestrator{}
	o.CurrentPhase = Phase.NewPhase().WaitingConnectionPhase()
	o.Epoch = 0
	return o
}

func (o *Orchestrator) StartSimualtion() {
	o.CurrentPhase.RequestEventPhase()
	o.Epoch = 0
}
func (o *Orchestrator) PauseSimualtion() {
	o.CurrentPhase.RequestEventPhase()
	o.Epoch = 0
}
func (o *Orchestrator) GetPhase() *Phase.Phase {
	return o.CurrentPhase
}

func (o *Orchestrator) NextEpoch() {
	o.Epoch = +1
}
