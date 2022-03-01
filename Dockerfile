FROM golang:1.17-alpine as builder

WORKDIR /service
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build  -o ./build/account-service ./cmd/api/main.go

FROM alpine:latest

WORKDIR /service
COPY ./configs/docker.yaml ./configs/main.yaml
COPY .env .
COPY --from=builder /service/build/account-service ./

ENTRYPOINT ["./account-service"]
