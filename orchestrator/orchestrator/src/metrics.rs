use std::{convert::TryInto, net};

use axum::prelude::*;
use ethers::prelude::*;
use hyper::Server;
use lazy_static::lazy_static;
use prometheus::*;

pub async fn metrics_main_loop(addr: &net::SocketAddr) {
    let get_metrics = || async {
        let mut buffer = Vec::new();
        let encoder = TextEncoder::new();
        let metric_families = prometheus::gather();
        encoder.encode(&metric_families, &mut buffer).unwrap();
        String::from_utf8(buffer.clone()).unwrap()
    };

    let app = route("/", get(get_metrics));

    info!("metrics listening on {}", addr);
    Server::bind(addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

// Counters
lazy_static! {
    pub static ref COSMOS_UNAVAILABLE: IntCounter = register_int_counter!(opts!(
        "cosmos_unavailable",
        "cosmos chain was unavailable",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    pub static ref ETHEREUM_EVENT_CHECK_FAILURES: IntCounter = register_int_counter!(opts!(
        "ethereum_event_check_failures",
        "ethereum events could not be retrieved",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    pub static ref ETHEREUM_UNAVAILABLE: IntCounter = register_int_counter!(opts!(
        "ethereum_unavailable",
        "ethereum chain was unavailable",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    pub static ref UNSIGNED_BATCH_FAILURES: IntCounter = register_int_counter!(opts!(
        "unsigned_batch_failures",
        "unsigned batches could not be retrieved",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    pub static ref UNSIGNED_LOGIC_CALL_FAILURES: IntCounter = register_int_counter!(opts!(
        "unsigned_logic_call_failures",
        "unsigned logic calls could not be retrieved",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    pub static ref UNSIGNED_VALSET_FAILURES: IntCounter = register_int_counter!(opts!(
        "unsigned_valset_failures",
        "unsigned valsets could not be retrieved",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
}

// Gauges (guarded by setters)
lazy_static! {
    static ref COSMOS_BLOCK_HEIGHT: IntGauge = register_int_gauge!(opts!(
        "cosmos_block_height",
        "cosmos chain block height",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    static ref COSMOS_LAST_EVENT_NONCE: IntGauge = register_int_gauge!(opts!(
        "cosmos_last_event_nonce",
        "last event nonce committed by this validator",
        labels! {"chain" => "cosmos"}
    ))
    .unwrap();
    static ref ETHEREUM_BLOCK_HEIGHT: IntGauge = register_int_gauge!(opts!(
        "ethereum_block_height",
        "ethereum chain block height",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_CHECK_FOR_EVENTS_END_BLOCK: IntGauge = register_int_gauge!(opts!(
        "ethereum_check_for_events_end_block",
        "end of block range being scanned ethereum events",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_CHECK_FOR_EVENTS_STARTING_BLOCK: IntGauge = register_int_gauge!(opts!(
        "ethereum_check_for_events_starting_block",
        "start of block range being scanned ethereum events",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_BATCH_EVENT: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_batch_event",
        "last event_nonce for a batch event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_BATCH_NONCE: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_batch_nonce",
        "last batch_nonce for a batch event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_DEPOSIT_BLOCK: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_deposit_block",
        "last block height for a deposit event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_DEPOSIT_EVENT: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_deposit_event",
        "last event_nonce for a deposit event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_ERC20_BLOCK: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_erc20_block",
        "last block height for an erc20 event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_ERC20_EVENT: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_erc20_event",
        "last event_nonce for a erc20 event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_EVENT_NONCE: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_event_nonce",
        "last event nonce on the gravity contract",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_LOGIC_CALL_EVENT: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_logic_call_event",
        "last event_nonce for a logic call event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_LOGIC_CALL_NONCE: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_logic_call_nonce",
        "last logic_call_nonce for a logic call event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_VALSET_EVENT: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_valset_event",
        "last event_nonce for a valset event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_LAST_VALSET_NONCE: IntGauge = register_int_gauge!(opts!(
        "ethereum_last_valset_nonce",
        "last valset_nonce for a valset event",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
    static ref ETHEREUM_BAL: IntGauge = register_int_gauge!(opts!(
        "ethereum_balance",
        "orchestrator_key current ethereum_balance",
        labels! {"chain" => "ethereum"}
    ))
    .unwrap();
}

pub fn set_cosmos_block_height(v: u64) {
    set_u64(&COSMOS_BLOCK_HEIGHT, v)
}

pub fn set_cosmos_last_event_nonce(v: u64) {
    set_u64(&COSMOS_LAST_EVENT_NONCE, v);
}

pub fn set_ethereum_block_height(v: u64) {
    set_u64(&ETHEREUM_BLOCK_HEIGHT, v);
}

pub fn set_ethereum_check_for_events_end_block(v: u64) {
    set_u64(&ETHEREUM_CHECK_FOR_EVENTS_END_BLOCK, v);
}

pub fn set_ethereum_check_for_events_starting_block(v: u64) {
    set_u64(&ETHEREUM_CHECK_FOR_EVENTS_STARTING_BLOCK, v);
}

pub fn set_ethereum_last_batch_event(v: U256) {
    set_u256(&ETHEREUM_LAST_BATCH_EVENT, v);
    set_u256(&ETHEREUM_LAST_EVENT_NONCE, v);
}

pub fn set_ethereum_last_batch_nonce(v: U256) {
    set_u256(&ETHEREUM_LAST_BATCH_NONCE, v);
}

pub fn set_ethereum_last_deposit_block(v: U256) {
    set_u256(&ETHEREUM_LAST_DEPOSIT_BLOCK, v);
}

pub fn set_ethereum_last_deposit_event(v: U256) {
    set_u256(&ETHEREUM_LAST_DEPOSIT_EVENT, v);
    set_u256(&ETHEREUM_LAST_EVENT_NONCE, v);
}

pub fn set_ethereum_last_erc20_block(v: U256) {
    set_u256(&ETHEREUM_LAST_ERC20_BLOCK, v);
}

pub fn set_ethereum_last_erc20_event(v: U256) {
    set_u256(&ETHEREUM_LAST_ERC20_EVENT, v);
    set_u256(&ETHEREUM_LAST_EVENT_NONCE, v);
}

pub fn set_ethereum_last_logic_call_event(v: U256) {
    set_u256(&ETHEREUM_LAST_LOGIC_CALL_EVENT, v);
    set_u256(&ETHEREUM_LAST_EVENT_NONCE, v);
}

pub fn set_ethereum_last_logic_call_nonce(v: U256) {
    set_u256(&ETHEREUM_LAST_LOGIC_CALL_NONCE, v);
}

pub fn set_ethereum_last_valset_event(v: U256) {
    set_u256(&ETHEREUM_LAST_VALSET_EVENT, v);
    set_u256(&ETHEREUM_LAST_EVENT_NONCE, v);
}

pub fn set_ethereum_last_valset_nonce(v: U256) {
    set_u256(&ETHEREUM_LAST_VALSET_NONCE, v);
}

pub fn set_ethereum_bal(v: U256) {
    set_u256(&ETHEREUM_BAL, v);
}

fn set_u64(gauge: &IntGauge, value: u64) {
    let v = value.try_into().unwrap_or(-1);
    if v > gauge.get() {
        gauge.set(v);
    }
}

fn set_u256(gauge: &IntGauge, value: U256) {
    let v = value.to_string().parse().unwrap_or(-1);
    if v > gauge.get() {
        gauge.set(v);
    }
}
