# Build our application using a Go builder.
FROM golang:1.20 AS builder

WORKDIR /src/project
COPY . .
# need to install + build templ BEFORE compiling ze go.


# should prob read up on this
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/go-htmlx ./cmd/main.go

# install goose.
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest



# CMD ["sleep", "infinity"]
# Our final Docker image stage starts here.
FROM alpine
COPY --from=flyio/litefs:0.5.7 /usr/local/bin/litefs /usr/local/bin/litefs
# Copy binaries from the previous build stages.
# may need to compile templ here too.
COPY --from=builder /usr/local/bin/go-htmlx /usr/local/bin/go-htmlx
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /migrations
COPY etc/litefs.yml /etc/litefs.yml
# install alpine deps
RUN apk add bash fuse3 sqlite ca-certificates curl

RUN mkdir -p /litefs/data

ENV GO_PORT="8081"
ENV PORT="8080"

EXPOSE 8080
ENTRYPOINT litefs mount
