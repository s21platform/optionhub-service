protogen:
# Для использования нужно установить:
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	protoc --go_out=. --go-grpc_out=. ./api/optionhub.proto --experimental_allow_proto3_optional
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/optionhub.proto --experimental_allow_proto3_optional

codegen:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
	oapi-codegen -generate chi-server -package api api/schema.yaml > internal/generated/server.gen.go
	oapi-codegen -generate types -package api api/schema.yaml > internal/generated/models.gen.go