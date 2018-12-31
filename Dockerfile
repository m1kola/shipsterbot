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

# Install app dependencies
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

# Download, verify and install a pre-built DB migration tool
RUN wget -q https://github.com/golang-migrate/migrate/releases/download/v4.1.0/migrate.linux-amd64.tar.gz \
    && echo  "56546df1fcd708e981b965676c4930a23aa05543e662c681d700599315e5553d  migrate.linux-amd64.tar.gz" | sha256sum -c - \
    && tar -xzf migrate.linux-amd64.tar.gz && mv migrate.linux-amd64 /usr/local/bin/migrate && rm migrate.linux-amd64.tar.gz

# Copy a binary from the build step
COPY --from=build /builddir/shipsterbot .

# Define a command to run in a container
CMD ./shipsterbot startbot telegram
