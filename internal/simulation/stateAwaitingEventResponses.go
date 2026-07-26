package simulation

type StateAwaitingEventResponses struct {
	SimulationEngine *SimulationEngine
}

func (w *StateAwaitingEventResponses) WaitingConnections() error {
	return nil
}

func (w *StateAwaitingEventResponses) ReadMsg() error {
	return nil
}
func (w *StateAwaitingEventResponses) RequestMsg() error {
	return nil
}

func NewStateAwaitingEventResponses(s *SimulationEngine) *StateAwaitingEventResponses {
	return &StateAwaitingEventResponses{SimulationEngine: s}
}
