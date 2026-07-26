package simulation

type StateAwaitingClientStatus struct {
	SimulationEngine *SimulationEngine
}

func (w *StateAwaitingClientStatus) WaitingConnections() error {
	return nil
}

func (w *StateAwaitingClientStatus) ReadMsg() error {
	return nil
}
func (w *StateAwaitingClientStatus) RequestMsg() error {
	return nil
}

func NewStateAwaitingClientStatus(s *SimulationEngine) *StateAwaitingClientStatus {
	return &StateAwaitingClientStatus{SimulationEngine: s}
}
