protogen:
	protoc --go_out=. --go-grpc_out=. ./api/optionhub.proto --experimental_allow_proto3_optional
	protoc --doc_out=. --doc_opt=markdown,GRPC_API.md ./api/optionhub.proto --experimental_allow_proto3_optional
