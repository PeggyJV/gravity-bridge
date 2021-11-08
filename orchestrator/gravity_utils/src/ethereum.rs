use ethers::prelude::*;
use std::panic;

pub fn downcast_to_u64(input: U256) -> Option<u64> {
    match panic::catch_unwind(|| input.as_u64()) {
        Ok(downcasted) => Some(downcasted),
        Err(_) => None,
    }
}

pub fn downcast_to_u128(input: U256) -> Option<u128> {
    match panic::catch_unwind(|| input.as_u128()) {
        Ok(downcasted) => Some(downcasted),
        Err(_) => None,
    }
}