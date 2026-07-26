package simulation

type StateRequestingEvents struct {
	SimulationEngine *SimulationEngine
}

func (w *StateRequestingEvents) WaitingConnections() error {
	return nil
}

func (w *StateRequestingEvents) ReadMsg() error {
	return nil
}
func (w *StateRequestingEvents) RequestMsg() error {
	return nil
}

func NewStateRequestingEvents(s *SimulationEngine) *StateRequestingEvents {
	return &StateRequestingEvents{SimulationEngine: s}
}
