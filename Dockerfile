FROM golang:latest
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -a /app/cmd/apiserver/
ENTRYPOINT exec go run cmd/apiserver/main.go