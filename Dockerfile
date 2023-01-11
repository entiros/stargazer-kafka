# Docker builder
FROM golang:1.19.5 as builder

# Set destination for COPY
WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Download Go modules
RUN go mod download

# Test
#RUN go test -v ./...

# Build
RUN go build -o /stargazer-kafka ./cmd/stargazer-kafka/main.go


# Docker run
FROM ubuntu

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy executable from builder
COPY --from=builder /stargazer-kafka /stargazer-kafka

# Configuration directory and default configuration file
RUN mkdir -p /configs/

COPY ./configs/config_example.yaml /configs/config.yml

# Run
CMD [ "/stargazer-kafka", "/configs/config.yml" ]
