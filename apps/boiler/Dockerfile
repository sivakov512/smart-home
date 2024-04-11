FROM lukemathwalker/cargo-chef:latest AS planner

WORKDIR /opt/boiler

COPY . .
RUN cargo chef prepare --recipe-path recipe.json


FROM rust:1.77 as builder

RUN mkdir -p /opt/boiler
WORKDIR /opt/boiler

RUN apt-get update && \
    apt-get install -y cmake && \
    cargo install cargo-chef --locked

COPY --from=planner /opt/boiler/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

COPY . .
RUN cargo build -r


FROM debian:bookworm-slim

RUN mkdir -p /opt/boiler
WORKDIR /opt/boiler

RUN apt-get update && \
    apt-get install -y openssl

VOLUME ["/opt/boiler/conf"]

ENV BOILERCONFIG="/opt/boiler/conf/config.toml"

COPY --from=builder /opt/boiler/target/release/boiler /opt/boiler/boiler

CMD ./boiler
