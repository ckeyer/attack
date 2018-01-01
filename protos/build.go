//go:generate sh -c "set -ex; protoc -I. -I/usr/local/include -I${GOPATH}/src -I${GOPATH}/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=plugins=grpc:. *.proto"

package protos
