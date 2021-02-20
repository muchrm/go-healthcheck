FROM golang:1.15-buster AS build_base

ENV GO111MODULE=on

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM build_base AS server_builder
COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o go-healthcheck main.go 

FROM debian:buster-slim
RUN apt-get update && apt-get --no-install-recommends --no-install-suggests --yes --quiet install ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /go/src/app
COPY --from=server_builder /go/src/app/go-healthcheck /usr/bin
CMD [ "go-healthcheck" ]
