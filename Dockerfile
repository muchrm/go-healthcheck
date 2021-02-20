FROM golang:1.15-stretch AS build_base

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS server_builder
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build main.go

FROM debian:stretch-slim
RUN apt-get update && apt-get --no-install-recommends --no-install-suggests --yes --quiet install ca-certificates
WORKDIR /go/src/app
COPY --from=server_builder /go/src/app/main /go/src/app
ENTRYPOINT [ "./main" ]
