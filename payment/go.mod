module github.com/pinai4/microservices-course-project/payment

go 1.24.0

replace github.com/pinai4/microservices-course-project/shared => ../shared

require (
	github.com/google/uuid v1.6.0
	github.com/pinai4/microservices-course-project/shared v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.75.1
)

require (
	go.opentelemetry.io/otel/sdk/metric v1.38.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
)
