.PHONY: build
build: pingpong_server pingpong_client

pingpong_server: server/server.go pb/pingpong.pb.go
	CGO_ENABLED=0 go build -o $@ $<

pingpong_client: client/client.go pb/pingpong.pb.go
	CGO_ENABLED=0 go build -o $@ $<

.PHONY: proto
proto:
	protoc -I ${GOPATH}/src -I ./vendor -I . \
		--go_out=plugins=grpc:. \
		--go_opt=paths=source_relative \
		pb/pingpong.proto

gogoproto:
	protoc -I ${GOPATH}/src -I ./vendor -I . \
		--gogo_out=plugins=grpc,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,paths=source_relative:. \
		pb/pingpong.proto

clean:
	go clean
	rm -f pingpong_server pingpong_client
