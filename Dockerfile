FROM golang:1.16
WORKDIR /work
COPY . .
RUN go build -o jrecord cmd/jrecord.go

FROM ubuntu:latest
RUN apt-get update && apt-get install -y ca-certificates
WORKDIR /work
COPY --from=0 ./work/jrecord .
COPY ./frontend/dist/ ./frontend/dist/
CMD ["./jrecord"]
