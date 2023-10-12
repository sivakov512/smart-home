FROM lukemathwalker/cargo-chef:latest AS planner

WORKDIR /opt/heater

COPY . .
RUN cargo chef prepare --recipe-path recipe.json


FROM rust:1.73 as builder

RUN mkdir -p /opt/heater
WORKDIR /opt/heater

RUN apt-get update && \
    apt-get install -y cmake && \
    cargo install cargo-chef --locked

COPY --from=planner /opt/heater/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

COPY . .
RUN cargo build -r


FROM debian:bookworm-slim

RUN mkdir -p /opt/heater
WORKDIR /opt/heater

RUN apt-get update && \
    apt-get install -y openssl

VOLUME ["/opt/heater/conf"]

ENV VIRTUALACCONFIG="/opt/heater/conf/config.toml"

COPY --from=builder /opt/heater/target/release/heater /opt/heater/heater

CMD ./heater

