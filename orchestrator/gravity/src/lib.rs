//! This crate contains various components and utilities for interacting with the Gravity Cosmos module. Primarily
//! Extensions to Althea's 'deep_space' Cosmos transaction library to allow it to send Gravity module specific messages
//! parse Gravity module specific endpoints and generally interact with the multitude of Gravity specific functionality
//! that's part of the Cosmos module.

#[macro_use]
extern crate log;

pub mod build;
pub mod ethereum;
pub mod query;
pub mod send;
pub mod utils;
use std::str::FromStr;

pub use deep_space;

// Re-export ethereum functions for backward compatibility
pub use ethereum::{one_eth, one_eth_f32};
use num256::Uint256;

#[derive(Debug, Clone, PartialEq, Eq, PartialOrd, Ord, Hash)]
pub struct Coin {
    pub amount: Uint256,
    pub denom: String,
}

impl Coin {
    pub fn new(amount: Uint256, denom: String) -> Self {
        Self { amount, denom }
    }
}

impl From<Coin> for deep_space::Coin {
    fn from(value: Coin) -> Self {
        deep_space::Coin {
            amount: value.amount,
            denom: value.denom,
        }
    }
}

impl From<deep_space::Coin> for Coin {
    fn from(value: deep_space::Coin) -> Self {
        Coin {
            amount: value.amount,
            denom: value.denom,
        }
    }
}

impl From<Coin> for cosmos_sdk_proto_althea::cosmos::base::v1beta1::Coin {
    fn from(value: Coin) -> Self {
        cosmos_sdk_proto_althea::cosmos::base::v1beta1::Coin {
            amount: value.amount.to_string(),
            denom: value.denom,
        }
    }
}

impl From<cosmos_sdk_proto_althea::cosmos::base::v1beta1::Coin> for Coin {
    fn from(value: cosmos_sdk_proto_althea::cosmos::base::v1beta1::Coin) -> Self {
        Coin {
            amount: Uint256::from_str(&value.amount).unwrap(),
            denom: value.denom,
        }
    }
}

impl From<Coin> for cosmos_sdk_proto::cosmos::base::v1beta1::Coin {
    fn from(value: Coin) -> Self {
        cosmos_sdk_proto::cosmos::base::v1beta1::Coin {
            amount: value.amount.to_string(),
            denom: value.denom,
        }
    }
}

impl From<cosmos_sdk_proto::cosmos::base::v1beta1::Coin> for Coin {
    fn from(value: cosmos_sdk_proto::cosmos::base::v1beta1::Coin) -> Self {
        Coin {
            amount: Uint256::from_str(&value.amount).unwrap(),
            denom: value.denom,
        }
    }
}
