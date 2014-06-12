all: test build

build:
	go build -o bin/my_agent

test:
	go test

run:
	time rake agents[1] &>logs/output.txt; tail -n 10 logs/output.txt; less logs/output.txt
