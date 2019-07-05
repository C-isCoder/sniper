default:clean rpc doc
	@go build

clean:
	@git clean -x -f -d
# 接口代码生成
rpc:
	@find rpc -name '*.proto' -exec protoc --twirp_out=. --go_out=. {} \;
# 开发运行服务。
dev:
	@go run main.go server
# 生成 markdown 接口文档
doc:
	@find rpc -name '*.proto' -exec protoc --markdown_out=path_prefix=/api:. {} \;

.PHONY: test rpc
