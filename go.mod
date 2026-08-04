module github.com/Norzuiso/orchestrator

go 1.26.1

replace github.com/Norzuiso/protocol => ../protocol

require (
	github.com/Norzuiso/protocol v0.0.0-00010101000000-000000000000
	github.com/grpc-ecosystem/go-grpc-middleware/v2 v2.3.3
	google.golang.org/grpc v1.81.1
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	modernc.org/libc v1.74.4 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

require (
	go.etcd.io/bbolt v1.5.0
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260226221140-a57be14db171 // indirect
	modernc.org/sqlite v1.56.0
)
