# Docker builder
FROM golang:1.18 as builder

# Set destination for COPY
WORKDIR /app

COPY . .

# Download Go modules
RUN go mod download

# Build
RUN go build -o /stargazer-kafka ./cmd/stargazer-kafka/main.go


# Docker run
FROM golang:1.18

# Copy executable from builder
COPY --from=builder /stargazer-kafka /stargazer-kafka

# Configuration directory and default configuration file
RUN mkdir -p /configs/
COPY ./configs/config.yml /configs/config.yml

# Run
CMD [ "/stargazer-kafka", "/configs/config.yml" ]
