FROM golang:1.21 as builder

RUN mkdir -p /opt/hap-ui
WORKDIR /opt/hap-ui

COPY go.mod go.sum ./
RUN go mod download

COPY ./ /opt/hap-ui
RUN go build .


FROM golang:1.21

RUN mkdir -p /opt/hap-ui
WORKDIR /opt/hap-ui

VOLUME ["/opt/hap-ui/db"]
VOLUME ["/opt/hap-ui/conf"]

ENV HAPUICONFIG="/opt/hap-ui/conf/config.toml"

COPY --from=builder /opt/hap-ui/hap-ui /opt/hap-ui/hap-ui

CMD ./hap-ui
