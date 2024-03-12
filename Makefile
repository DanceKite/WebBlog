.PHONY: all build run gotool clean help

BINARY="WebBlog"

all: gotool build

build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}.exe

run:
	@go run ./main.go 

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	@if [ -f ${BINARY}.exe ] ; then rm ${BINARY}.exe ; fi

help:
	@echo "make - 格式化 Go 代码并编译生成二进制文件"
	@echo "make build - 编译 Go 代码，生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make gotool - 运行 Go 工具 'fmt' 和 'vet'"
	@echo "make clean - 移除二进制文件 vim swap files"
	@echo "make help - 查看帮助"