use crate::error::GravityError;
use ethers::prelude::*;
use std::panic;

pub fn downcast_to_f32(input: U256) -> Option<f32> {
    // technically the max value of u128 is larger than f32, but
    // in practicality this won't matter for any of the cases we
    // would care about downcasting from a U256, and Rust will
    // gracefully saturate the cast
    match panic::catch_unwind(|| input.as_u128() as f32) {
        Ok(downcasted) => Some(downcasted),
        Err(_) => None,
    }
}

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

pub fn bytes_to_hex_str(bytes: &[u8]) -> String {
    bytes
        .iter()
        .map(|b| format!("{:0>2x?}", b))
        .fold(String::new(), |acc, x| acc + &x)
}

pub fn hex_str_to_bytes(s: &str) -> Result<Vec<u8>, GravityError> {
    let s = match s.strip_prefix("0x") {
        Some(s) => s,
        None => &s,
    };
    let bytes = s
        .as_bytes()
        .chunks(2)
        .map::<Result<u8, GravityError>, _>(|ch| {
            let str = String::from_utf8(ch.to_vec())?;
            let byte = u8::from_str_radix(&str, 16)?;

            Ok(byte)
        })
        .collect::<Result<Vec<_>, _>>()?;

    Ok(bytes)
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

#[test]
fn encode_bytes() {
    assert_eq!(bytes_to_hex_str(&[0xf]), "0f".to_owned());
    assert_eq!(bytes_to_hex_str(&[0xff]), "ff".to_owned());
    assert_eq!(
        bytes_to_hex_str(&[0xde, 0xad, 0xbe, 0xef]),
        "deadbeef".to_owned()
    );
}