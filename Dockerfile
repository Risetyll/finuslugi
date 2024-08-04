FROM golang:1.22.5-bookworm

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /build/server ./cmd/server

WORKDIR /build

COPY api/certs /app/certs

EXPOSE 8443

CMD ["./server"]