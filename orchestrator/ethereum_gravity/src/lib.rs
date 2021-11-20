//! This crate contains various components and utilities for interacting with the Gravity Ethereum contract.

use ethers::types::U256;

#[macro_use]
extern crate log;

pub mod deploy_erc20;
pub mod erc20_utils;
pub mod logic_call;
pub mod send_to_cosmos;
pub mod submit_batch;
pub mod types;
pub mod utils;
pub mod valset_update;

pub fn one_eth() -> U256 {
    1000000000000000000u128.into()
}

pub fn one_eth_f32() -> f32 {
    1000000000000000000u128 as f32
}
