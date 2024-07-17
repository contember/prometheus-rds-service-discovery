FROM golang:latest as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

# Build the Go application (you can adjust the output binary name as needed)
RUN CGO_ENABLED=0 GOOS=linux go build -o service-discovery ./cmd/service-discovery

CMD ["./service-discovery"]

FROM alpine:latest as runtime

WORKDIR /app

# Copy the binary from the builder image to the final image
COPY --from=builder /app/service-discovery .

CMD ["./service-discovery"]
