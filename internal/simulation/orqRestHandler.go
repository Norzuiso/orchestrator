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

func (s *SimulationEngine) NextPhase(w http.ResponseWriter, req *http.Request) {
	oldState := s.currentState.GetStateName()
	nextState, err := s.currentState.NextState()
	if err != nil {
		fmt.Fprintf(w, "Error  %v", err)
	}
	s.SetState(nextState)
	newState := nextState.GetStateName()

	fmt.Fprintf(w, "Phase updated. From: %v To: %v", oldState, newState)
}

func (s *SimulationEngine) NextEpoch(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.NextEpoch()
	fmt.Fprintf(w, "Epoch: %v", s.Orchestrator.Epoch)
}
