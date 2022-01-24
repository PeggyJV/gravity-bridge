#![allow(clippy::needless_question_mark)]
use crate::error::GravityError;
use ethers::prelude::*;
use ethers::types::Address as EthAddress;
use std::{panic, result::Result};

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

pub fn downcast_to_f64(input: U256) -> Option<f64> {
    match panic::catch_unwind(|| input.as_u128() as f64) {
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

pub fn format_eth_address(address: EthAddress) -> String {
    format!("0x{}", bytes_to_hex_str(address.as_bytes()))
}

pub fn format_eth_hash(hash: H256) -> String {
    format!("0x{}", bytes_to_hex_str(hash.as_bytes()))
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
        None => s,
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

pub fn vec_u8_to_fixed_32(v: Vec<u8>) -> Result<[u8; 32], GravityError> {
    if v.len() != 32 {
        return Err(GravityError::InvalidArgumentError(format!(
            "Error converting Vec<u8> to [u8; 32], length is not 32: {:?}",
            v
        )));
    }

    Ok(u8_slice_to_fixed_32(&v[..])?)
}

pub fn u8_slice_to_fixed_32(v: &[u8]) -> Result<[u8; 32], GravityError> {
    if v.len() != 32 {
        return Err(GravityError::InvalidArgumentError(format!(
            "Error converting &[u8] to [u8; 32], length is not 32: {:?}",
            v
        )));
    }

    let mut v_slice: [u8; 32] = Default::default();
    v_slice.copy_from_slice(v);
    Ok(v_slice)
}

#[test]
fn overflow_f32() {
    assert_eq!(downcast_to_f32(42.into()), Some(42f32));
    assert_eq!(downcast_to_f32(U256::MAX), None);
}

#[test]
fn overflow_u64() {
    assert_eq!(downcast_to_u64(42.into()), Some(42u64));
    assert_eq!(downcast_to_u64(U256::MAX), None);
}

#[test]
fn overflow_u128() {
    assert_eq!(downcast_to_u128(42.into()), Some(42u128));
    assert_eq!(downcast_to_u128(U256::MAX), None);
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

#[test]
fn decode_bytes() {
    assert_eq!(
        hex_str_to_bytes(&"deadbeef".to_owned()).expect("Unable to decode"),
        [222, 173, 190, 239]
    );
}

#[test]
fn decode_odd_amount_of_bytes() {
    assert_eq!(hex_str_to_bytes(&"f".to_owned()).unwrap(), vec![15]);
}

#[test]
fn bytes_raises_decode_error() {
    use crate::error::GravityError;

    let e = hex_str_to_bytes(&"\u{012345}deadbeef".to_owned()).unwrap_err();

    match e {
        GravityError::FromUtf8Error(_) => {}
        _ => panic!(),
    };
}

#[test]
fn bytes_raises_parse_error() {
    use crate::error::GravityError;

    let e = hex_str_to_bytes(&"Lorem ipsum".to_owned()).unwrap_err();
    match e {
        GravityError::ParseIntError(_) => {}
        _ => panic!(),
    }
}

#[test]
fn parse_prefixed_empty() {
    assert_eq!(
        hex_str_to_bytes(&"0x".to_owned()).unwrap(),
        Vec::<u8>::new()
    );
}

#[test]
fn parse_prefixed_non_empty() {
    assert_eq!(
        hex_str_to_bytes(&"0xdeadbeef".to_owned()).unwrap(),
        vec![0xde, 0xad, 0xbe, 0xef]
    );
}
