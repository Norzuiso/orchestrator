package orchestrator

type Orchestrator struct {
	Epoch float32

	conectedClients int // 10
	doneClients     int // 4
	pendingClients  int // 6
}

func NewOrquestrator() *Orchestrator {
	o := &Orchestrator{}
	o.Epoch = 0
	return o
}

func (o *Orchestrator) StartSimualtion() {
	o.Epoch = 0
}
func (o *Orchestrator) PauseSimualtion() {
	o.Epoch = 0
}

func (o *Orchestrator) NextEpoch() {
	o.Epoch = +1
}
