arch = $(shell dpkg --print-architecture)


all: bindata
	go install

update: bindata
	go clean
	go install

bindata:
	go-bindata -pkg="drops" js/

l32:
	GOOS=linux GOARCH=386 go build -o bin/$(filename)_i386

l64:
	GOOS=linux GOARCH=amd64 go build -o bin/$(filename)_amd64

watch:
	CompileDaemon -build="make" -pattern="(.+\\.go|.+\\.c|.+\\.js)$$" -exclude="bindata.go"