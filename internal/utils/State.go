package utils

type State interface {
	ReadMsg() error
	RequestMsg() error
}
