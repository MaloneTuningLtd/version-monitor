FROM golang:alpine AS build-env
RUN apk --no-cache add git
COPY . /go/src/github.com/MaloneTuningLtd/version-monitor
RUN cd /go/src/github.com/MaloneTuningLtd/version-monitor \
  && go get \
  && go build -o version-monitor

FROM alpine:latest
RUN apk add --no-cache ca-certificates
COPY --from=build-env /go/src/github.com/MaloneTuningLtd/version-monitor/version-monitor /usr/local/bin/version-monitor
ENTRYPOINT /usr/local/bin/version-monitor
