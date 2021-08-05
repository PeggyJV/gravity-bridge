use axum::prelude::*;
use hyper::Server;
use lazy_static::lazy_static;
use prometheus::{self, register_int_counter, Encoder, IntCounter, TextEncoder};

lazy_static! {
    static ref HIGH_FIVE_COUNTER: IntCounter =
        register_int_counter!("highfives", "Number of high fives received").unwrap();
}

pub async fn metrics_main_loop() {
    let get_root = || async {
        HIGH_FIVE_COUNTER.inc();
        "Hello, World!\n"
    };

    let get_metrics = || async {
        let mut buffer = Vec::new();
        let encoder = TextEncoder::new();
        let metric_families = prometheus::gather();
        encoder.encode(&metric_families, &mut buffer).unwrap();
        String::from_utf8(buffer.clone()).unwrap()
    };

    let app = route("/", get(get_root)).route("/metrics", get(get_metrics));

    // TODO(Levi) accept config for SocketAddr
    let addr = std::net::SocketAddr::from(([127, 0, 0, 1], 3000));
    info!("metrics listening on {}", addr);
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
