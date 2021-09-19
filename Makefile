run:
	go run cmd/joirecord.go

build:
	cd frontend && yarn build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/joirecord.go