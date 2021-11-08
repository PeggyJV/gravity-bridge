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

#[test]
fn test_downcast_to_u64() {
    let mut i = 0u64;
    while i < 100_000 {
        assert_eq!(i, downcast_to_u64(i.into()).unwrap());
        i += 1
    }
    let mut i: u64 = std::u32::MAX.into();
    i -= 100;
    let end = i + 100_000;
    while i < end {
        assert_eq!(i, downcast_to_u64(i.into()).unwrap());
        i += 1
    }
}

#[test]
fn test_downcast_to_u128() {
    let mut i = 0u128;
    while i < 100_000 {
        assert_eq!(i, downcast_to_u128(i.into()).unwrap());
        i += 1
    }
    let mut i: u128 = std::u64::MAX.into();
    i -= 100;
    let end = i + 100_000;
    while i < end {
        assert_eq!(i, downcast_to_u128(i.into()).unwrap());
        i += 1
    }
}