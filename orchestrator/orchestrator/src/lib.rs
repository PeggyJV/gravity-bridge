//! Orchestrator is a sort of specialized relayer for the Gravity bridge that runs on every validator.
//! Things this library is responsible for
//!   * Performing all the Ethereum signing required to submit updates and generate batches
//!   * Progressing the validator set update generation process.
//!   * Observing events on the Ethereum chain and submitting oracle messages for validator consensus
//! Things this library needs
//!   * Access to the validators signing Ethereum key
//!   * Access to the validators Cosmos key
//!   * Access to an Cosmos chain RPC server
//!   * Access to an Ethereum chain RPC server

pub mod ethereum_event_watcher;
pub mod get_with_retry;
pub mod main_loop;
pub mod metrics;
pub mod oracle_resync;

#[macro_use]
extern crate log;
extern crate prometheus;
