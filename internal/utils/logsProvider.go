package utils

type LogsProvider interface {
	WriteLogs(str string)
	WriteClientStream(clientId int64, isActive bool)
}
