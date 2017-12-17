# Build step
FROM golang:1.9 as build

WORKDIR /builddir


# Install dependencies
# This allows us to use Docker cache in most of the cases,
# so build happens faster
COPY Gopkg.lock Gopkg.toml Makefile ./
RUN make vendor

# Copy source files and compile
COPY . .
RUN make build


# Image build
FROM golang:1.9
WORKDIR /app/bin/

# Copy a binary from the build step
COPY --from=build /builddir/shipsterbot .

# Define a command to run in a container
CMD ./shipsterbot migrate up && ./shipsterbot startbot telegram

EXPOSE 8443
