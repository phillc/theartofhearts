all: test install

install:
	go get "git.apache.org/thrift.git/lib/go/thrift"
	go build -o bin/my_agent

test:
	go test

run:
	time rake agents[1] &>logs/output.txt; tail -n 5 logs/output.txt
