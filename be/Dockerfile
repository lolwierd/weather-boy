FROM golang:1.24 AS builder

WORKDIR /app

RUN apt-get update -y && apt-get install gh

RUN --mount=type=secret,id=GITHUB_TOKEN \
  gh auth login --with-token < /run/secrets/GITHUB_TOKEN && \
  gh auth setup-git

COPY go.mod go.sum ./

RUN go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/weatherboy/be

FROM alpine:3.20

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
