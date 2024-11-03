FROM golang:1.23-alpine AS builder

WORKDIR /src

# Copy src folder
COPY ./src .

# Install build dependencies (ca-certificates will be copied to runner image)
RUN apk add --update git ca-certificates

# graphql generate
RUN go run github.com/99designs/gqlgen generate

# Fetch dependencies
RUN go get

# Build into an executable
RUN go build -ldflags "-s -w" -o /dist/backend

# From an empty docker image, copy just the executable
FROM scratch

COPY --from=builder /dist/backend /usr/bin/backend

# copy the ca-certificate.crt from the build stage
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Set gin to use release mode (instead of the default debug mode)
ENV GIN_MODE=release

ENTRYPOINT ["/usr/bin/backend"]