FROM golang:latest
WORKDIR /PlaylistAPI/internal/test/pkg
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
CMD ["go", "test", "./test/...", "-v"]