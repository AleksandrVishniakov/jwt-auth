FROM golang:1.23-alpine AS build

WORKDIR /go/src/jwt-auth

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd/
COPY internal ./internal/

RUN go build -o ../../bin/app ./cmd/app/main.go

FROM alpine
WORKDIR /go

COPY config.yaml ./config.yaml
COPY --from=build /go/bin/app /bin/app

# CMD ["app"]