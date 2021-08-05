use std::convert::TryInto;

use axum::prelude::*;
use hyper::Server;
use lazy_static::lazy_static;
use prometheus::*;

lazy_static! {
    //
    static ref COSMOS_UNAVAILABLE: IntCounter = register_int_counter!(opts!(
        "cosmos_unavailable",
        "cosmos chain was unavailable",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    //
    static ref COSMOS_BLOCK_HEIGHT: IntGauge = register_int_gauge!(opts!(
        "cosmos_block_height",
        "cosmos chain block height",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    //
    static ref COSMOS_LAST_EVENT_NONCE: IntGauge = register_int_gauge!(opts!(
        "cosmos_last_event_nonce",
        "last event nonce committed by this validator",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    //
    static ref ETHEREUM_LAST_EVENT_NONCE: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_event_nonce",
        "last event nonce on the gravity contract",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    //
    static ref ETHEREUM_EVENT_CHECK_FAILURES: IntCounter = register_int_counter!(opts!(
        "ethereum_event_check_failures",
        "ethereum events could not be retrieved",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    //
    static ref ETHEREUM_UNAVAILABLE: IntCounter = register_int_counter!(opts!(
        "ethereum_unavailable",
        "ethereum chain was unavailable",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    //
    static ref ETHEREUM_BLOCK_HEIGHT: IntGauge = register_int_gauge!(opts!(
        "ethereum_block_height",
        "ethereum chain block height",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    //
    static ref ETHEREUM_CHECK_FOR_EVENTS_STARTING_BLOCK: IntGauge = register_int_gauge!(opts!(
        "ethereum_check_for_events_starting_block",
        "start of block range being scanned ethereum events",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    //
    static ref ETHEREUM_CHECK_FOR_EVENTS_END_BLOCK: IntGauge = register_int_gauge!(opts!(
        "ethereum_check_for_events_end_block",
        "end of block range being scanned ethereum events",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    //
    static ref UNSIGNED_BATCH_FAILURES: IntCounter = register_int_counter!(opts!(
        "unsigned_batch_failures",
        "unsigned batches could not be retrieved",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    //
    static ref UNSIGNED_LOGIC_CALL_FAILURES: IntCounter = register_int_counter!(opts!(
        "unsigned_logic_call_failures",
        "unsigned logic calls could not be retrieved",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    //
    static ref UNSIGNED_VALSET_FAILURES: IntCounter = register_int_counter!(opts!(
        "unsigned_valset_failures",
        "unsigned valsets could not be retrieved",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
}

pub fn increment_cosmos_unavailable() {
    COSMOS_UNAVAILABLE.inc();
}

pub fn increment_ethereum_unavailable() {
    ETHEREUM_UNAVAILABLE.inc();
}

pub fn increment_ethereum_event_check_failures() {
    ETHEREUM_EVENT_CHECK_FAILURES.inc();
}

pub fn increment_unsigned_batch_failures() {
    UNSIGNED_BATCH_FAILURES.inc();
}

pub fn increment_unsigned_logic_call_failures() {
    UNSIGNED_LOGIC_CALL_FAILURES.inc();
}

pub fn increment_unsigned_valset_failures() {
    UNSIGNED_VALSET_FAILURES.inc();
}

pub fn set_cosmos_block_height(v: u64) {
    COSMOS_BLOCK_HEIGHT.set(match v.try_into() {
        Ok(v) => v,
        Err(_) => -1,
    });
}

pub fn set_cosmos_last_event_nonce(v: u64) {
    COSMOS_LAST_EVENT_NONCE.set(match v.try_into() {
        Ok(v) => v,
        Err(_) => -1,
    });
}

pub fn set_ethereum_last_event_nonce(v: clarity::Uint256) {
    ETHEREUM_LAST_EVENT_NONCE.set(match v.to_str_radix(10).parse() {
        Ok(v) => v,
        Err(_) => -1,
    });
}

pub fn set_ethereum_block_height(v: clarity::Uint256) {
    ETHEREUM_BLOCK_HEIGHT.set(match v.to_str_radix(10).parse() {
        Ok(v) => v,
        Err(_) => -1,
    });
}

pub fn set_ethereum_check_for_events_starting_block(v: clarity::Uint256) {
    ETHEREUM_CHECK_FOR_EVENTS_STARTING_BLOCK.set(match v.to_str_radix(10).parse() {
        Ok(v) => v,
        Err(_) => -1,
    });
}

pub fn set_ethereum_check_for_events_end_block(v: clarity::Uint256) {
    ETHEREUM_CHECK_FOR_EVENTS_END_BLOCK.set(match v.to_str_radix(10).parse() {
        Ok(v) => v,
        Err(_) => -1,
    });
}

pub async fn metrics_main_loop() {
    let get_metrics = || async {
        let mut buffer = Vec::new();
        let encoder = TextEncoder::new();
        let metric_families = prometheus::gather();
        encoder.encode(&metric_families, &mut buffer).unwrap();
        String::from_utf8(buffer.clone()).unwrap()
    };

    let app = route("/", get(get_metrics));

    // TODO(Levi) accept config for SocketAddr
    let addr = std::net::SocketAddr::from(([127, 0, 0, 1], 3000));
    info!("metrics listening on {}", addr);
    Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}
