FROM lukemathwalker/cargo-chef:latest AS planner

WORKDIR /opt/virtual-ac

COPY . .
RUN cargo chef prepare --recipe-path recipe.json


FROM rust:1.72 as builder

RUN mkdir -p /opt/virtual-ac
WORKDIR /opt/virtual-ac

RUN apt-get update && \
    apt-get install -y cmake && \
    cargo install cargo-chef --locked

COPY --from=planner /opt/virtual-ac/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

COPY . .
RUN cargo build -r


FROM debian:bookworm-slim

RUN mkdir -p /opt/virtual-ac
WORKDIR /opt/virtual-ac

VOLUME ["/opt/virtual-ac/conf"]

ENV VIRTUALACCONFIG="/opt/virtual-ac/conf/config.toml"

COPY --from=builder /opt/virtual-ac/target/release/virtual-ac /opt/virtual-ac/virtual-ac

CMD ./virtual-ac
