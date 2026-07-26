package simulation

type StateStoringClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (w *StateStoringClientStatus) WaitingConnections() error {
	return nil
}

func (w *StateStoringClientStatus) ReadMsg() error {
	return nil
}
func (w *StateStoringClientStatus) RequestMsg() error {
	return nil
}

func NewStateStoringClientStatus(s *SimulationEngine) *StateStoringClientStatus {
	return &StateStoringClientStatus{SimulationEngine: s}
}
