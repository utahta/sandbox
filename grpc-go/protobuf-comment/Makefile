build:
	protoc --go_out=plugins=grpc:. ./helloworld/helloworld.proto
	protoc --include_source_info --descriptor_set_out=./helloworld/FILE ./helloworld/helloworld.proto

