# Build our application using a Go builder.
FROM golang:1.20 AS builder

WORKDIR /src/project
COPY . .
# should prob read up on this
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/go-htmlx ./main.go

# install goose.
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# CMD ["sleep", "infinity"]
# Our final Docker image stage starts here.
FROM alpine
COPY --from=flyio/litefs:0.54 /usr/local/bin/litefs /usr/local/bin/litefs
# Copy binaries from the previous build stages.
COPY --from=builder /usr/local/bin/go-htmlx /usr/local/bin/go-htmlx
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY etc/litefs.yml /etc/litefs.yml
# install alpine deps
RUN apk add bash fuse3 sqlite ca-certificates curl

RUN mkdir -p /litefs/data

EXPOSE 8080
ENTRYPOINT go-htmlx
