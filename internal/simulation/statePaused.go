package simulation

type StatePaused struct {
	SimulationEngine *SimulationEngine
}

func (w *StatePaused) WaitingConnections() error {
	return nil
}

func (w *StatePaused) ReadMsg() error {
	return nil
}
func (w *StatePaused) RequestMsg() error {
	return nil
}

func NewStatePaused(s *SimulationEngine) *StatePaused {
	return &StatePaused{SimulationEngine: s}
}
