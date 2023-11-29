# Build our application using a Go builder.
FROM golang:1.21.4 AS builder

WORKDIR /src/project
COPY . .
# need to install + build templ BEFORE compiling ze go.
RUN go install github.com/a-h/templ/cmd/templ@latest


# install tw?
RUN curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.3/tailwindcss-linux-x64
RUN chmod +x tailwindcss-linux-x64
RUN mv tailwindcss-linux-x64 tailwindcss
RUN yarn install
RUN ./tailwindcss -i ./styles/input.css -o ./static/css/styles.css --minify

# do our build stuff here.
RUN templ generate

# should prob read up on this
RUN go build -buildvcs=false -ldflags "-s -w -extldflags '-static'" -tags osusergo,netgo -o /usr/local/bin/go-htmlx ./cmd/main.go

# install goose.
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@latest



# CMD ["sleep", "infinity"]
# Our final Docker image stage starts here.
FROM alpine
COPY --from=flyio/litefs:0.5.8 /usr/local/bin/litefs /usr/local/bin/litefs
# Copy binaries from the previous build stages.
# may need to compile templ here too.
COPY --from=builder /usr/local/bin/go-htmlx /usr/local/bin/go-htmlx
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY migrations /migrations
COPY etc/litefs.yml /etc/litefs.yml

COPY --from=builder /src/project/static static
# install alpine deps
RUN apk add bash fuse3 sqlite ca-certificates curl gnupg

RUN mkdir -p /litefs/data

ENV GO_PORT="8081"
ENV PORT="8080"
ARG DOPPLER_TOKEN=''
ENV DOPPLER_TOKEN=''
RUN (curl -Ls --tlsv1.2 --proto "=https" --retry 3 https://cli.doppler.com/install.sh || wget -t 3 -qO- https://cli.doppler.com/install.sh) | sh


EXPOSE 8080
ENTRYPOINT litefs mount
