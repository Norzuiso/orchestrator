package simulation

type StateCollectingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateCollectingEvents) WaitingConnections() error {
	return nil
}

func (w *StateCollectingEvents) ReadMsg() error {
	return nil
}
func (w *StateCollectingEvents) RequestMsg() error {
	return nil
}

func NewStateCollectingEvents(s *SimulationEngine) *StateCollectingEvents {
	return &StateCollectingEvents{SimulationEngine: s}
}
