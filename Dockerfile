# Build our application using a Go builder.
FROM golang:1.20 AS builder

WORKDIR /src/project
COPY . .
# should prob read up on this
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/go-htmlx ./main.go

# Our final Docker image stage starts here.
FROM alpine

# Copy binaries from the previous build stages.
COPY --from=builder /usr/local/bin/go-htmlx /usr/local/bin/go-htmlx

# install goose to run migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# install alpine deps
RUN apk add bash fuse3 sqlite ca-certificates curl

ENTRYPOINT go-htmlx