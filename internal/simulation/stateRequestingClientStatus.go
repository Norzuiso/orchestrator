package simulation

type StateRequestingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (w *StateRequestingClientStatus) WaitingConnections() error {
	return nil
}

func (w *StateRequestingClientStatus) ReadMsg() error {
	return nil
}
func (w *StateRequestingClientStatus) RequestMsg() error {
	return nil
}

func NewStateRequestingClientStatus(s *SimulationEngine) *StateRequestingClientStatus {
	return &StateRequestingClientStatus{SimulationEngine: s}
}
