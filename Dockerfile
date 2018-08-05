ARG ARCH
# Build go binary
FROM golang AS gobuild

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY ./ /go/src/github.com/meyskens/dispatch.sh-applogs
WORKDIR /go/src/github.com/meyskens/dispatch.sh-applogs

ARG GOARCH
ARG GOARM

RUN dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} GOARM=${GOARM} go build -a -installsuffix cgo ./

# Set up deinitive image
ARG ARCH
FROM multiarch/alpine:${ARCH}-edge

RUN apk add --no-cache ca-certificates
COPY --from=gobuild /go/src/github.com/meyskens/dispatch.sh-applogs/dispatch.sh-applogs /usr/local/bin/dispatch.sh-applogs

CMD dispatch.sh-applogs