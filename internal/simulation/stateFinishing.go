package simulation

type StateFinishing struct {
	SimulationEngine *SimulationEngine
}

func (w *StateFinishing) WaitingConnections() error {
	return nil
}

func (w *StateFinishing) ReadMsg() error {
	return nil
}
func (w *StateFinishing) RequestMsg() error {
	return nil
}

func NewStateFinishing(s *SimulationEngine) *StateFinishing {
	return &StateFinishing{SimulationEngine: s}
}
