FROM golang:1.18.1-alpine AS gobuilder
ARG VERSION=master
ARG OS=linux
ARG ARCH=amd64

ENV CGO_ENABLED=0
ENV GOOS=${OS}
ENV GOARCH=${ARCH}

RUN mkdir /webhook-go && chmod 700 /webhook-go
WORKDIR /webhook-go

RUN apk add --no-cache curl
RUN curl -s -L "https://github.com/voxpupuli/webhook-go/archive/${VERSION}.tar.gz" | tar --exclude .gitignore --strip 1 -xz &&\
    go build .

FROM puppet/puppet-agent:7.16.0

EXPOSE 8080

RUN /opt/puppetlabs/puppet/bin/gem install r10k

USER puppet

COPY --from=gobuilder /webhook-go/webhook-go /webhook-go
ENTRYPOINT [ "/webhook-go" ]
