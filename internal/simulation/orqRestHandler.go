package simulation

import (
	"fmt"
	"net/http"
)

// Orchestrator
func (s *SimulationEngine) StartSimulation(w http.ResponseWriter, req *http.Request) {
	nextState, _ := s.currentState.NextState()
	s.SetState(nextState)
	fmt.Fprintf(w, "Simulation started.")
}

func (s *SimulationEngine) StopSimulation(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.PauseSimualtion()
	fmt.Fprintf(w, "Simulation stopped on \n\t epoch: \t%v \n\tphase:\t%v\n\t.", s.Orchestrator.Epoch, s.currentState.GetStateName())
}

func (s *SimulationEngine) NextStateHttp(w http.ResponseWriter, req *http.Request) {
	newState, err := s.NextState()
	if err != nil {
		fmt.Fprintf(w, "%v", err)
	}

	fmt.Fprintf(w, "Phase updated. From: %v", newState)
}

func (s *SimulationEngine) NextEpoch(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.NextEpoch()
	fmt.Fprintf(w, "Epoch: %v", s.Orchestrator.Epoch)
}
