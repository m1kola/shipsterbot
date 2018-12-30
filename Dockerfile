# Build step
FROM golang:1.11 as build

WORKDIR /builddir


# Install dependencies
# This allows us to use Docker cache in most of the cases,
# so build happens faster
COPY Gopkg.lock Gopkg.toml Makefile ./
RUN make vendor

# Copy source files and compile
COPY . .

# Note that we need to disable CGO
# to be able to run a binary in Alpine linux
RUN CGO_ENABLED=0 make build


# Image build
FROM alpine:3.8
WORKDIR /app/bin/

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Copy a binary from the build step
COPY --from=build /builddir/shipsterbot .

# Define a command to run in a container
CMD ./shipsterbot migrate up && ./shipsterbot startbot telegram

EXPOSE 8443
