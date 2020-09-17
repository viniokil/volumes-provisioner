# Default to Go 1.13
ARG GO_VERSION=1.13
FROM golang:${GO_VERSION}-alpine as builder

# Necessary to run 'go get' and to compile the linked binary
RUN apk add git musl-dev

COPY . /go/src/github.com/viniokil/volumes-provisioner

WORKDIR /go/src/github.com/viniokil/volumes-provisioner

# ENV GO111MODULE=on

# build & install server
RUN go get -u ./... && \
        CGO_ENABLED=0 \
        GOOS=linux \
        GOARCH=amd64 \
        go build -a -ldflags '-w' -o /go/bin/volumes-provisioner

FROM scratch AS final
LABEL maintainer="Valerii Kravets <viniokil@gmail.com>"

COPY --from=builder /go/bin/volumes-provisioner /volumes-provisioner

ENTRYPOINT ["/volumes-provisioner"]
