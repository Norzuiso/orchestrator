package utils

type LogsProvider interface {
	WriteLogs(str string)
	WriteClientStream(clientId int64, isOpen bool)
}
