GOPATH=$(shell pwd)

all:
	env GOPATH=${GOPATH} go get -u github.com/jteeuwen/go-bindata/...
	env GOPATH=${GOPATH} bin/go-bindata -o src/go-cowsay/bindata.go src/go-cowsay/cows
	env GOPATH=${GOPATH} go build -o cowsay go-cowsay
	env GOPATH=${GOPATH} go build -o cowthink go-cowsay

clean:
	rm -f -- cowsay
	rm -f -- cowthink
	rm -f -- go-cowsay
	rm -f -- src/go-cowsay/bindata.go
	rm -rf -- bin
	rm -rf -- pkg
	rm -rf -- src/github.com
	rm -rf -- src/golang.com

fmt:
	env GOPATH=${GOPATH} go fmt go-cowsay

install:
	sudo cp cowsay /usr/local/bin/
