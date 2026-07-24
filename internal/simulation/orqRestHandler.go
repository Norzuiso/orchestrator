package simulation

import (
	"fmt"
	"net/http"
)

// Orchestrator
func (s *SimulationEngine) StartSimulation(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.StartSimualtion()
	fmt.Fprintf(w, "Simulation started.")
}

func (s *SimulationEngine) StopSimulation(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.PauseSimualtion()
	fmt.Fprintf(w, "Simulation stopped on \n\t epoch: \t%v \n\tphase:\t%v\n\t.", s.Orchestrator.Epoch, s.Orchestrator.GetPhase().GetName())
}

func (s *SimulationEngine) NextPhase(w http.ResponseWriter, req *http.Request) {
	oldPhase := s.Orchestrator.CurrentPhase.GetName()
	s.Orchestrator.CurrentPhase.NextPhase()
	newPhase := s.Orchestrator.CurrentPhase.GetName()

	fmt.Fprintf(w, "Phase updated. From: %v To: %v", oldPhase, newPhase)
}

func (s *SimulationEngine) NextEpoch(w http.ResponseWriter, req *http.Request) {
	s.Orchestrator.NextEpoch()
	fmt.Fprintf(w, "Epoch: %v", s.Orchestrator.Epoch)
}
