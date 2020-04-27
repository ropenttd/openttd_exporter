# BUILD ENVIRONMENT

FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/ropenttd/openttd_exporter
COPY . .

# Fetch dependencies
RUN go get -d -v

# Build the binary
RUN go build -o /go/bin/openttd_exporter

# END BUILD ENVIRONMENT
# DEPLOY ENVIRONMENT

FROM scratch
MAINTAINER duck. <me@duck.me.uk>

# Copy the executable
COPY --from=builder /go/bin/openttd_exporter /openttd_exporter

# Finally, let's run
ENTRYPOINT ["/openttd_exporter"]
