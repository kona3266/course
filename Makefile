all: build

build: cli serv

cli:
	env CGO_ENABLED=0 go build -trimpath -o bin/cli ./cmd/client/main.go

serv:
	env CGO_ENABLED=0 go build -trimpath -o bin/serv ./cmd/server/main.go

clean:
	rm -f ./bin/cli
	rm -f ./bin/serv
