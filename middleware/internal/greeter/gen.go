package greeter

//go:generate protoc -I ./ ./greeter.proto --go_out=plugins=grpc:.
