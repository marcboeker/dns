build-linux-amd64:
	CGO_ENABLED=1 CC="x86_64-unknown-linux-gnu-cc" GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dns cmd/server/main.go

build:
	CGO_ENABLED=1 go build -ldflags "-s -w" -o dns cmd/server/main.go

run:
	go run cmd/server/main.go

update-blocklist:
	curl -o blocklist.txt https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts
	go run cmd/blocker/main.go blocklist.txt blocklist.db
	rm blocklist.txt