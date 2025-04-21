FROM --platform=linux/amd64 quay.io/frrouting/frr:8.5.4

RUN apk update && apk add --no-cache \
    wget \
    git \
    build-base \
    ca-certificates \
    bash

RUN wget https://go.dev/dl/go1.24.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.24.2.linux-amd64.tar.gz && \
    rm go1.24.2.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin
ENV GOPATH=/go
ENV PATH=$PATH:$GOPATH/bin

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" "$GOPATH/pkg"

RUN go install github.com/air-verse/air@latest

WORKDIR /app

RUN mkdir -p /app/src /app/local /app/protobufSource /app/tmp

COPY ./dockerfile/.air.toml .


#VOLUME ["/app"]

COPY dockerfile/files/tempClient /usr/local/bin/tempClient
COPY dockerfile/files/start_frr_dev /usr/bin/start_frr

RUN chmod +x /usr/bin/start_frr
RUN chmod +x /usr/local/bin/tempClient

CMD ["/usr/bin/start_frr"]
