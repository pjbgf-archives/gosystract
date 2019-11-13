FROM golang:1.13-alpine AS build

LABEL repository="https://github.com/pjbgf/gosystract/"

WORKDIR /go/src/pjbgf/gosystract

RUN apk --update add git gcc
ADD . /go/src/pjbgf/gosystract

ENV GO111MODULE=on

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

RUN go build -ldflags "-w -X=github.com/pjbgf/gosystract/cmd/cli.gitcommit=$(git describe --tags --always)" -o /go/bin/gosystract

FROM alpine:latest
COPY --from=build /go/bin/gosystract /usr/bin
COPY --from=build /usr/local/go/bin/go /usr/bin
COPY --from=build /usr/local/go/pkg/tool/linux_amd64/objdump /usr/local/go/pkg/tool/linux_amd64/
CMD ["/gosystract"]