# Step 1: Modules caching
FROM golang:1.22.1-alpine3.18 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22.1-alpine3.18 as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./crypto_building ./cmd/

# Step 3: Final
FROM alpine 
WORKDIR /app
COPY --from=builder /app/config/dev.conf.yaml conf.yaml 
COPY --from=builder /app/config/dev.env .env
COPY --from=builder /app/migrations/ migrations/
COPY --from=builder /app/crypto_building init

# Command to run the executable
CMD ["./init"]

