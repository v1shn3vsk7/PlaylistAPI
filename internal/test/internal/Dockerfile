FROM golang:latest
WORKDIR /app/internal/test/internal
COPY go.mod .
RUN go mod download
COPY . .
RUN go test ./... -v
