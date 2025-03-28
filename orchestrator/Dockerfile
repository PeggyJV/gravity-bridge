# Reference: https://www.lpalmieri.com/posts/fast-rust-docker-builds/

FROM rust:1.81-bullseye AS cargo-chef-rust
RUN cargo install cargo-chef --version 0.1.62 --locked

FROM cargo-chef-rust AS planner
WORKDIR /app
# We only pay the installation cost once,
# it will be cached from the second build onwards
# To ensure a reproducible build consider pinning
# the cargo-chef version with `--version X.X.X`
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

FROM cargo-chef-rust AS cacher
WORKDIR /app
COPY --from=planner /app/recipe.json recipe.json
RUN cargo chef cook --release --recipe-path recipe.json

FROM cargo-chef-rust AS builder
WORKDIR /app
COPY . .
# Copy over the cached dependencies
COPY --from=cacher /app/target target
COPY --from=cacher /usr/local/cargo /usr/local/cargo
RUN cargo build --release --bin gorc

FROM debian:bullseye AS runtime
WORKDIR /app
COPY --from=builder /app/target/release/gorc /usr/local/bin
CMD ["gorc"]
