version=1.0.3

run:
	go run cmd/jrecord.go

build:
	cd frontend && yarn build
	docker build --platform linux/amd64 -t registry.cn-hongkong.aliyuncs.com/joi/jrecord:${version} .
