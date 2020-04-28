# BUILD ENVIRONMENT

FROM golang:alpine AS builder

# All these steps will be cached
WORKDIR $GOPATH/src/github.com/ropenttd/openttd_exporter
COPY go.mod .
COPY go.sum .

# Get dependencies - will also be cached if we won't change mod/sum
RUN go mod download

# Then copy the rest of this source code
COPY . .

# And build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/openttd_exporter

# END BUILD ENVIRONMENT
# DEPLOY ENVIRONMENT

FROM scratch
MAINTAINER duck. <me@duck.me.uk>

# Copy the executable
COPY --from=builder /go/bin/openttd_exporter /openttd_exporter

# Finally, let's run
ENTRYPOINT ["/openttd_exporter"]
