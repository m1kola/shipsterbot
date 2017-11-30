# Build step
FROM golang:1.9 as build

RUN mkdir -p /go/src/github.com/m1kola/shipsterbot
WORKDIR /go/src/github.com/m1kola/shipsterbot

# Install dependencies
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep
COPY Gopkg.lock Gopkg.toml ./
RUN dep ensure -vendor-only

# Copy source files and compile
COPY . .
RUN go build -o app

# Image build
FROM golang:1.9
WORKDIR /app/bin/

# Copy a binary from the build step
COPY --from=build /go/src/github.com/m1kola/shipsterbot/app .

# Define a command to run in a container
CMD ./app migrate up && ./app startbot telegram

EXPOSE 8080
