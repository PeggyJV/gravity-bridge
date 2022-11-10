# Reference: https://www.lpalmieri.com/posts/fast-rust-docker-builds/

# Collin: There is some strange dependency conflict in the docker build that for some
# reason doesn't occur when using vendored dependencies.
FROM rust:1 as cargo-chef-rust
# RUN cargo install cargo-chef

FROM cargo-chef-rust as planner
WORKDIR app
# We only pay the installation cost once,
# it will be cached from the second build onwards
# To ensure a reproducible build consider pinning
# the cargo-chef version with `--version X.X.X`
COPY . .
RUN cargo vendor
# RUN cargo chef prepare --recipe-path recipe.json

FROM cargo-chef-rust as cacher
WORKDIR app
COPY . .
# COPY --from=planner /app/recipe.json recipe.json
COPY --from=planner /usr/local/cargo /usr/local/cargo
RUN cargo build --release --bin gorc

# FROM cargo-chef-rust as builder
# WORKDIR app
# COPY . .
# # Copy over the cached dependencies
# COPY --from=cacher /app/target target
# COPY --from=cacher /usr/local/cargo /usr/local/cargo
# RUN cargo build --release --bin gorc

FROM cargo-chef-rust as runtime
WORKDIR app
COPY --from=cacher /app/target/release/gorc /usr/local/bin
CMD gorc
