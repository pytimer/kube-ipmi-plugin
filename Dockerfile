# Build the kube-ipmi-plugin binary
FROM golang:1.12.6 as builder

# Copy in the go src
WORKDIR /go/src/github.com/pytimer/kube-ipmi-plugin
COPY main.go main.go
COPY pkg/    pkg/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o kube-ipmi-plugin github.com/pytimer/kube-ipmi-plugin

# Copy the kube-ipmi-plugin into a thin image
FROM alpine
WORKDIR /

RUN apk add --no-cache ipmitool

COPY --from=builder /go/src/github.com/pytimer/kube-ipmi-plugin/kube-ipmi-plugin .
ENTRYPOINT ["/kube-ipmi-plugin"]
