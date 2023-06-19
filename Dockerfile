FROM golang:1.20-alpine
WORKDIR /workspace
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o bin/ ./cmd/api/
ENTRYPOINT [ "bin/api" ]


