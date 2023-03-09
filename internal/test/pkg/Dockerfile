FROM golang:latest
WORKDIR /app/internal/test/pkg
COPY go.mod .
RUN go mod download
COPY . .
RUN go test ./... -v
