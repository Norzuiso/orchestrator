package simulation

import (
	"fmt"
	"net/http"
)

// Orchestrator
func (s *SimulationEngine) StartSimulation(w http.ResponseWriter, req *http.Request) {
	if _, err := fmt.Fprintf(w, "Simulation started."); err != nil {
		return
	}
	s.NextState()
}

func (s *SimulationEngine) EndSimulation(w http.ResponseWriter, req *http.Request) {
	s.CurrentState = s.StateEnd
	s.CurrentState.StartState()
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
	s.Orchestrator.NextEpoch()
	fmt.Fprintf(w, "Epoch: %v", s.Orchestrator.Epoch)
}
