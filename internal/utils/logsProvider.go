package utils

type LogsProvider interface {
	WriteLogs(str string)
}
