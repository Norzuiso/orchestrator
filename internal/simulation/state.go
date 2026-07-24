package simulation

type State interface {
	WaitConnections() error
}
