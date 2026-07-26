package simulation

type StateWaitingConnections struct {
	SimulationEngine *SimulationEngine
}

func (w *StateWaitingConnections) WaitingConnections() error {
	return nil
}

func (w *StateWaitingConnections) ReadMsg() error {
	return nil
}
func (w *StateWaitingConnections) RequestMsg() error {
	return nil
}

func NewStateWaitingConnections(s *SimulationEngine) *StateWaitingConnections {
	return &StateWaitingConnections{SimulationEngine: s}
}
