FROM golang:1.23-alpine AS builder

WORKDIR /src

# Copy src folder
COPY ./src .

# Fetch dependencies
RUN go get

# Build into an executable
RUN go build -ldflags "-s -w" -o /dist/backend

# From an empty docker image, copy just the executable
FROM scratch

COPY --from=builder /dist/backend /usr/bin/backend

# Set gin to use release mode (instead of the default debug mode)
ENV GIN_MODE=release

ENTRYPOINT ["/usr/bin/backend"]