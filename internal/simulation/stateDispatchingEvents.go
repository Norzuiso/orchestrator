package simulation

type StateDispatchingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateDispatchingEvents) WaitingConnections() error {
	return nil
}

func (w *StateDispatchingEvents) ReadMsg() error {
	return nil
}
func (w *StateDispatchingEvents) RequestMsg() error {
	return nil
}

func NewStateDispatchingEvents(s *SimulationEngine) *StateDispatchingEvents {
	return &StateDispatchingEvents{SimulationEngine: s}
}
