pub use gravity::*;
/// This module was auto-generated with ethers-rs Abigen.
/// More information at: <https://github.com/gakonst/ethers-rs>
#[allow(
    clippy::enum_variant_names,
    clippy::too_many_arguments,
    clippy::upper_case_acronyms,
    clippy::type_complexity,
    dead_code,
    non_camel_case_types
)]
pub mod gravity {
    const _: () = {
        ::core::include_bytes!("../Gravity.json",);
    };
    #[allow(deprecated)]
    fn __abi() -> ::ethers::core::abi::Abi {
        ::ethers::core::abi::ethabi::Contract {
            constructor: ::core::option::Option::Some(::ethers::core::abi::ethabi::Constructor {
                inputs: ::std::vec![
                    ::ethers::core::abi::ethabi::Param {
                        name: ::std::borrow::ToOwned::to_owned("_gravityId"),
                        kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                        internal_type: ::core::option::Option::Some(
                            ::std::borrow::ToOwned::to_owned("bytes32"),
                        ),
                    },
                    ::ethers::core::abi::ethabi::Param {
                        name: ::std::borrow::ToOwned::to_owned("_powerThreshold"),
                        kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                        internal_type: ::core::option::Option::Some(
                            ::std::borrow::ToOwned::to_owned("uint256"),
                        ),
                    },
                    ::ethers::core::abi::ethabi::Param {
                        name: ::std::borrow::ToOwned::to_owned("_validators"),
                        kind: ::ethers::core::abi::ethabi::ParamType::Array(
                            ::std::boxed::Box::new(::ethers::core::abi::ethabi::ParamType::Address,),
                        ),
                        internal_type: ::core::option::Option::Some(
                            ::std::borrow::ToOwned::to_owned("address[]"),
                        ),
                    },
                    ::ethers::core::abi::ethabi::Param {
                        name: ::std::borrow::ToOwned::to_owned("_powers"),
                        kind: ::ethers::core::abi::ethabi::ParamType::Array(
                            ::std::boxed::Box::new(::ethers::core::abi::ethabi::ParamType::Uint(
                                256usize
                            ),),
                        ),
                        internal_type: ::core::option::Option::Some(
                            ::std::borrow::ToOwned::to_owned("uint256[]"),
                        ),
                    },
                ],
            }),
            functions: ::core::convert::From::from([
                (
                    ::std::borrow::ToOwned::to_owned("deployERC20"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("deployERC20"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_cosmosDenom"),
                                kind: ::ethers::core::abi::ethabi::ParamType::String,
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("string"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_name"),
                                kind: ::ethers::core::abi::ethabi::ParamType::String,
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("string"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_symbol"),
                                kind: ::ethers::core::abi::ethabi::ParamType::String,
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("string"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_decimals"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(8usize),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint8"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::NonPayable,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("lastBatchNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("lastBatchNonce"),
                        inputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::borrow::ToOwned::to_owned("_erc20Address"),
                            kind: ::ethers::core::abi::ethabi::ParamType::Address,
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("address"),
                            ),
                        },],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("lastLogicCallNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("lastLogicCallNonce"),
                        inputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::borrow::ToOwned::to_owned("_invalidation_id"),
                            kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("bytes32"),
                            ),
                        },],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("sendToCosmos"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("sendToCosmos"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_tokenContract"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("address"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_destination"),
                                kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("bytes32"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_amount"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::NonPayable,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_gravityId"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_gravityId"),
                        inputs: ::std::vec![],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("bytes32"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_invalidationMapping"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_invalidationMapping",),
                        inputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("bytes32"),
                            ),
                        },],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_lastBatchNonces"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_lastBatchNonces",),
                        inputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Address,
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("address"),
                            ),
                        },],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_lastEventNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_lastEventNonce",),
                        inputs: ::std::vec![],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_lastValsetCheckpoint"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_lastValsetCheckpoint",),
                        inputs: ::std::vec![],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("bytes32"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_lastValsetNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_lastValsetNonce",),
                        inputs: ::std::vec![],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("state_powerThreshold"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("state_powerThreshold",),
                        inputs: ::std::vec![],
                        outputs: ::std::vec![::ethers::core::abi::ethabi::Param {
                            name: ::std::string::String::new(),
                            kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                            internal_type: ::core::option::Option::Some(
                                ::std::borrow::ToOwned::to_owned("uint256"),
                            ),
                        },],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::View,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("submitBatch"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("submitBatch"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_currentValset"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValsetArgs"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_sigs"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                            ::ethers::core::abi::ethabi::ParamType::Uint(8usize),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                        ],),
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValSignature[]"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_amounts"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256[]"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_destinations"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Address,
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("address[]"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_fees"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256[]"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_batchNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_tokenContract"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("address"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_batchTimeout"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::NonPayable,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("submitLogicCall"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("submitLogicCall"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_currentValset"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValsetArgs"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_sigs"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                            ::ethers::core::abi::ethabi::ParamType::Uint(8usize),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                        ],),
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValSignature[]"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_args"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                    ::ethers::core::abi::ethabi::ParamType::Bytes,
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct LogicCallArgs"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::NonPayable,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("testCheckValidatorSignatures"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("testCheckValidatorSignatures",),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_currentValset"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValsetArgs"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_sigs"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                            ::ethers::core::abi::ethabi::ParamType::Uint(8usize),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                        ],),
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValSignature[]"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_theHash"),
                                kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("bytes32"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_powerThreshold"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::Pure,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("testMakeCheckpoint"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("testMakeCheckpoint"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_valsetArgs"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValsetArgs"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_gravityId"),
                                kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("bytes32"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::Pure,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("updateValset"),
                    ::std::vec![::ethers::core::abi::ethabi::Function {
                        name: ::std::borrow::ToOwned::to_owned("updateValset"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_newValset"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValsetArgs"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_currentValset"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Address,
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Array(
                                        ::std::boxed::Box::new(
                                            ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                        ),
                                    ),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ::ethers::core::abi::ethabi::ParamType::Address,
                                ],),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValsetArgs"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("_sigs"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Tuple(::std::vec![
                                            ::ethers::core::abi::ethabi::ParamType::Uint(8usize),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                            ::ethers::core::abi::ethabi::ParamType::FixedBytes(
                                                32usize
                                            ),
                                        ],),
                                    ),
                                ),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("struct ValSignature[]"),
                                ),
                            },
                        ],
                        outputs: ::std::vec![],
                        constant: ::core::option::Option::None,
                        state_mutability: ::ethers::core::abi::ethabi::StateMutability::NonPayable,
                    },],
                ),
            ]),
            events: ::core::convert::From::from([
                (
                    ::std::borrow::ToOwned::to_owned("ERC20DeployedEvent"),
                    ::std::vec![::ethers::core::abi::ethabi::Event {
                        name: ::std::borrow::ToOwned::to_owned("ERC20DeployedEvent"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_cosmosDenom"),
                                kind: ::ethers::core::abi::ethabi::ParamType::String,
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_tokenContract"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_name"),
                                kind: ::ethers::core::abi::ethabi::ParamType::String,
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_symbol"),
                                kind: ::ethers::core::abi::ethabi::ParamType::String,
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_decimals"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(8usize),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_eventNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                        ],
                        anonymous: false,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("LogicCallEvent"),
                    ::std::vec![::ethers::core::abi::ethabi::Event {
                        name: ::std::borrow::ToOwned::to_owned("LogicCallEvent"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_invalidationId"),
                                kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_invalidationNonce",),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_returnData"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Bytes,
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_eventNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                        ],
                        anonymous: false,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("SendToCosmosEvent"),
                    ::std::vec![::ethers::core::abi::ethabi::Event {
                        name: ::std::borrow::ToOwned::to_owned("SendToCosmosEvent"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_tokenContract"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_sender"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_destination"),
                                kind: ::ethers::core::abi::ethabi::ParamType::FixedBytes(32usize,),
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_amount"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_eventNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                        ],
                        anonymous: false,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("TransactionBatchExecutedEvent"),
                    ::std::vec![::ethers::core::abi::ethabi::Event {
                        name: ::std::borrow::ToOwned::to_owned("TransactionBatchExecutedEvent",),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_batchNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_token"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_eventNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                        ],
                        anonymous: false,
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("ValsetUpdatedEvent"),
                    ::std::vec![::ethers::core::abi::ethabi::Event {
                        name: ::std::borrow::ToOwned::to_owned("ValsetUpdatedEvent"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_newValsetNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: true,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_eventNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_rewardAmount"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_rewardToken"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Address,
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_validators"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Address,
                                    ),
                                ),
                                indexed: false,
                            },
                            ::ethers::core::abi::ethabi::EventParam {
                                name: ::std::borrow::ToOwned::to_owned("_powers"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Array(
                                    ::std::boxed::Box::new(
                                        ::ethers::core::abi::ethabi::ParamType::Uint(256usize),
                                    ),
                                ),
                                indexed: false,
                            },
                        ],
                        anonymous: false,
                    },],
                ),
            ]),
            errors: ::core::convert::From::from([
                (
                    ::std::borrow::ToOwned::to_owned("BatchTimedOut"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("BatchTimedOut"),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("IncorrectCheckpoint"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("IncorrectCheckpoint",),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InsufficientPower"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InsufficientPower"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("cumulativePower"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("powerThreshold"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidBatchNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidBatchNonce"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("newNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("currentNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidLogicCallFees"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidLogicCallFees",),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidLogicCallNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidLogicCallNonce",),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("newNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("currentNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidLogicCallTransfers"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidLogicCallTransfers",),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidSendToCosmos"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidSendToCosmos",),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidSignature"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidSignature"),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("InvalidValsetNonce"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("InvalidValsetNonce"),
                        inputs: ::std::vec![
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("newNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                            ::ethers::core::abi::ethabi::Param {
                                name: ::std::borrow::ToOwned::to_owned("currentNonce"),
                                kind: ::ethers::core::abi::ethabi::ParamType::Uint(256usize,),
                                internal_type: ::core::option::Option::Some(
                                    ::std::borrow::ToOwned::to_owned("uint256"),
                                ),
                            },
                        ],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("LogicCallTimedOut"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("LogicCallTimedOut"),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("MalformedBatch"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("MalformedBatch"),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("MalformedCurrentValidatorSet"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("MalformedCurrentValidatorSet",),
                        inputs: ::std::vec![],
                    },],
                ),
                (
                    ::std::borrow::ToOwned::to_owned("MalformedNewValidatorSet"),
                    ::std::vec![::ethers::core::abi::ethabi::AbiError {
                        name: ::std::borrow::ToOwned::to_owned("MalformedNewValidatorSet",),
                        inputs: ::std::vec![],
                    },],
                ),
            ]),
            receive: false,
            fallback: false,
        }
    }
    ///The parsed JSON ABI of the contract.
    pub static GRAVITY_ABI: ::ethers::contract::Lazy<::ethers::core::abi::Abi> =
        ::ethers::contract::Lazy::new(__abi);
    pub struct Gravity<M>(::ethers::contract::Contract<M>);
    impl<M> ::core::clone::Clone for Gravity<M> {
        fn clone(&self) -> Self {
            Self(::core::clone::Clone::clone(&self.0))
        }
    }
    impl<M> ::core::ops::Deref for Gravity<M> {
        type Target = ::ethers::contract::Contract<M>;
        fn deref(&self) -> &Self::Target {
            &self.0
        }
    }
    impl<M> ::core::ops::DerefMut for Gravity<M> {
        fn deref_mut(&mut self) -> &mut Self::Target {
            &mut self.0
        }
    }
    impl<M> ::core::fmt::Debug for Gravity<M> {
        fn fmt(&self, f: &mut ::core::fmt::Formatter<'_>) -> ::core::fmt::Result {
            f.debug_tuple(::core::stringify!(Gravity))
                .field(&self.address())
                .finish()
        }
    }
    impl<M: ::ethers::providers::Middleware> Gravity<M> {
        /// Creates a new contract instance with the specified `ethers` client at
        /// `address`. The contract derefs to a `ethers::Contract` object.
        pub fn new<T: Into<::ethers::core::types::Address>>(
            address: T,
            client: ::std::sync::Arc<M>,
        ) -> Self {
            Self(::ethers::contract::Contract::new(
                address.into(),
                GRAVITY_ABI.clone(),
                client,
            ))
        }
        ///Calls the contract's `deployERC20` (0xf7955637) function
        pub fn deploy_erc20(
            &self,
            cosmos_denom: ::std::string::String,
            name: ::std::string::String,
            symbol: ::std::string::String,
            decimals: u8,
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([247, 149, 86, 55], (cosmos_denom, name, symbol, decimals))
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `lastBatchNonce` (0x011b2174) function
        pub fn last_batch_nonce(
            &self,
            erc_20_address: ::ethers::core::types::Address,
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([1, 27, 33, 116], erc_20_address)
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `lastLogicCallNonce` (0xc9d194d5) function
        pub fn last_logic_call_nonce(
            &self,
            invalidation_id: [u8; 32],
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([201, 209, 148, 213], invalidation_id)
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `sendToCosmos` (0x1ffbe7f9) function
        pub fn send_to_cosmos(
            &self,
            token_contract: ::ethers::core::types::Address,
            destination: [u8; 32],
            amount: ::ethers::core::types::U256,
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([31, 251, 231, 249], (token_contract, destination, amount))
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_gravityId` (0xbdda81d4) function
        pub fn state_gravity_id(&self) -> ::ethers::contract::builders::ContractCall<M, [u8; 32]> {
            self.0
                .method_hash([189, 218, 129, 212], ())
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_invalidationMapping` (0x7dfb6f86) function
        pub fn state_invalidation_mapping(
            &self,
            p0: [u8; 32],
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([125, 251, 111, 134], p0)
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_lastBatchNonces` (0xdf97174b) function
        pub fn state_last_batch_nonces(
            &self,
            p0: ::ethers::core::types::Address,
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([223, 151, 23, 75], p0)
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_lastEventNonce` (0x73b20547) function
        pub fn state_last_event_nonce(
            &self,
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([115, 178, 5, 71], ())
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_lastValsetCheckpoint` (0xf2b53307) function
        pub fn state_last_valset_checkpoint(
            &self,
        ) -> ::ethers::contract::builders::ContractCall<M, [u8; 32]> {
            self.0
                .method_hash([242, 181, 51, 7], ())
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_lastValsetNonce` (0xb56561fe) function
        pub fn state_last_valset_nonce(
            &self,
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([181, 101, 97, 254], ())
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `state_powerThreshold` (0xe5a2b5d2) function
        pub fn state_power_threshold(
            &self,
        ) -> ::ethers::contract::builders::ContractCall<M, ::ethers::core::types::U256> {
            self.0
                .method_hash([229, 162, 181, 210], ())
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `submitBatch` (0x8690ff98) function
        pub fn submit_batch(
            &self,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
            amounts: ::std::vec::Vec<::ethers::core::types::U256>,
            destinations: ::std::vec::Vec<::ethers::core::types::Address>,
            fees: ::std::vec::Vec<::ethers::core::types::U256>,
            batch_nonce: ::ethers::core::types::U256,
            token_contract: ::ethers::core::types::Address,
            batch_timeout: ::ethers::core::types::U256,
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [134, 144, 255, 152],
                    (
                        current_valset,
                        sigs,
                        amounts,
                        destinations,
                        fees,
                        batch_nonce,
                        token_contract,
                        batch_timeout,
                    ),
                )
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `submitLogicCall` (0x6941db93) function
        pub fn submit_logic_call(
            &self,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
            args: LogicCallArgs,
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([105, 65, 219, 147], (current_valset, sigs, args))
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `testCheckValidatorSignatures` (0x00901153) function
        pub fn test_check_validator_signatures(
            &self,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
            the_hash: [u8; 32],
            power_threshold: ::ethers::core::types::U256,
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [0, 144, 17, 83],
                    (current_valset, sigs, the_hash, power_threshold),
                )
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `testMakeCheckpoint` (0x01031525) function
        pub fn test_make_checkpoint(
            &self,
            valset_args: ValsetArgs,
            gravity_id: [u8; 32],
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([1, 3, 21, 37], (valset_args, gravity_id))
                .expect("method not found (this should never happen)")
        }
        ///Calls the contract's `updateValset` (0xaca6b1c1) function
        pub fn update_valset(
            &self,
            new_valset: ValsetArgs,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
        ) -> ::ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([172, 166, 177, 193], (new_valset, current_valset, sigs))
                .expect("method not found (this should never happen)")
        }
        ///Gets the contract's `ERC20DeployedEvent` event
        pub fn erc20_deployed_event_filter(
            &self,
        ) -> ::ethers::contract::builders::Event<::std::sync::Arc<M>, M, Erc20DeployedEventFilter>
        {
            self.0.event()
        }
        ///Gets the contract's `LogicCallEvent` event
        pub fn logic_call_event_filter(
            &self,
        ) -> ::ethers::contract::builders::Event<::std::sync::Arc<M>, M, LogicCallEventFilter>
        {
            self.0.event()
        }
        ///Gets the contract's `SendToCosmosEvent` event
        pub fn send_to_cosmos_event_filter(
            &self,
        ) -> ::ethers::contract::builders::Event<::std::sync::Arc<M>, M, SendToCosmosEventFilter>
        {
            self.0.event()
        }
        ///Gets the contract's `TransactionBatchExecutedEvent` event
        pub fn transaction_batch_executed_event_filter(
            &self,
        ) -> ::ethers::contract::builders::Event<
            ::std::sync::Arc<M>,
            M,
            TransactionBatchExecutedEventFilter,
        > {
            self.0.event()
        }
        ///Gets the contract's `ValsetUpdatedEvent` event
        pub fn valset_updated_event_filter(
            &self,
        ) -> ::ethers::contract::builders::Event<::std::sync::Arc<M>, M, ValsetUpdatedEventFilter>
        {
            self.0.event()
        }
        /// Returns an `Event` builder for all the events of this contract.
        pub fn events(
            &self,
        ) -> ::ethers::contract::builders::Event<::std::sync::Arc<M>, M, GravityEvents> {
            self.0
                .event_with_filter(::core::default::Default::default())
        }
    }
    impl<M: ::ethers::providers::Middleware> From<::ethers::contract::Contract<M>> for Gravity<M> {
        fn from(contract: ::ethers::contract::Contract<M>) -> Self {
            Self::new(contract.address(), contract.client())
        }
    }
    ///Custom Error type `BatchTimedOut` with signature `BatchTimedOut()` and selector `0x11724cc6`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "BatchTimedOut", abi = "BatchTimedOut()")]
    pub struct BatchTimedOut;
    ///Custom Error type `IncorrectCheckpoint` with signature `IncorrectCheckpoint()` and selector `0x723a3403`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "IncorrectCheckpoint", abi = "IncorrectCheckpoint()")]
    pub struct IncorrectCheckpoint;
    ///Custom Error type `InsufficientPower` with signature `InsufficientPower(uint256,uint256)` and selector `0x00bfb6ab`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "InsufficientPower", abi = "InsufficientPower(uint256,uint256)")]
    pub struct InsufficientPower {
        pub cumulative_power: ::ethers::core::types::U256,
        pub power_threshold: ::ethers::core::types::U256,
    }
    ///Custom Error type `InvalidBatchNonce` with signature `InvalidBatchNonce(uint256,uint256)` and selector `0xf7f920ad`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "InvalidBatchNonce", abi = "InvalidBatchNonce(uint256,uint256)")]
    pub struct InvalidBatchNonce {
        pub new_nonce: ::ethers::core::types::U256,
        pub current_nonce: ::ethers::core::types::U256,
    }
    ///Custom Error type `InvalidLogicCallFees` with signature `InvalidLogicCallFees()` and selector `0x48292479`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "InvalidLogicCallFees", abi = "InvalidLogicCallFees()")]
    pub struct InvalidLogicCallFees;
    ///Custom Error type `InvalidLogicCallNonce` with signature `InvalidLogicCallNonce(uint256,uint256)` and selector `0x01284fd2`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(
        name = "InvalidLogicCallNonce",
        abi = "InvalidLogicCallNonce(uint256,uint256)"
    )]
    pub struct InvalidLogicCallNonce {
        pub new_nonce: ::ethers::core::types::U256,
        pub current_nonce: ::ethers::core::types::U256,
    }
    ///Custom Error type `InvalidLogicCallTransfers` with signature `InvalidLogicCallTransfers()` and selector `0x853152a2`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(
        name = "InvalidLogicCallTransfers",
        abi = "InvalidLogicCallTransfers()"
    )]
    pub struct InvalidLogicCallTransfers;
    ///Custom Error type `InvalidSendToCosmos` with signature `InvalidSendToCosmos()` and selector `0x21739d9b`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "InvalidSendToCosmos", abi = "InvalidSendToCosmos()")]
    pub struct InvalidSendToCosmos;
    ///Custom Error type `InvalidSignature` with signature `InvalidSignature()` and selector `0x8baa579f`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "InvalidSignature", abi = "InvalidSignature()")]
    pub struct InvalidSignature;
    ///Custom Error type `InvalidValsetNonce` with signature `InvalidValsetNonce(uint256,uint256)` and selector `0xe0e8edf3`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(
        name = "InvalidValsetNonce",
        abi = "InvalidValsetNonce(uint256,uint256)"
    )]
    pub struct InvalidValsetNonce {
        pub new_nonce: ::ethers::core::types::U256,
        pub current_nonce: ::ethers::core::types::U256,
    }
    ///Custom Error type `LogicCallTimedOut` with signature `LogicCallTimedOut()` and selector `0xbcf37c25`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "LogicCallTimedOut", abi = "LogicCallTimedOut()")]
    pub struct LogicCallTimedOut;
    ///Custom Error type `MalformedBatch` with signature `MalformedBatch()` and selector `0xc1f97e35`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "MalformedBatch", abi = "MalformedBatch()")]
    pub struct MalformedBatch;
    ///Custom Error type `MalformedCurrentValidatorSet` with signature `MalformedCurrentValidatorSet()` and selector `0xc6617b7b`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(
        name = "MalformedCurrentValidatorSet",
        abi = "MalformedCurrentValidatorSet()"
    )]
    pub struct MalformedCurrentValidatorSet;
    ///Custom Error type `MalformedNewValidatorSet` with signature `MalformedNewValidatorSet()` and selector `0xc01ba0ab`
    #[derive(
        Clone,
        ::ethers::contract::EthError,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[etherror(name = "MalformedNewValidatorSet", abi = "MalformedNewValidatorSet()")]
    pub struct MalformedNewValidatorSet;
    ///Container type for all of the contract's custom errors
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        serde::Deserialize,
        serde::Serialize,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub enum GravityErrors {
        BatchTimedOut(BatchTimedOut),
        IncorrectCheckpoint(IncorrectCheckpoint),
        InsufficientPower(InsufficientPower),
        InvalidBatchNonce(InvalidBatchNonce),
        InvalidLogicCallFees(InvalidLogicCallFees),
        InvalidLogicCallNonce(InvalidLogicCallNonce),
        InvalidLogicCallTransfers(InvalidLogicCallTransfers),
        InvalidSendToCosmos(InvalidSendToCosmos),
        InvalidSignature(InvalidSignature),
        InvalidValsetNonce(InvalidValsetNonce),
        LogicCallTimedOut(LogicCallTimedOut),
        MalformedBatch(MalformedBatch),
        MalformedCurrentValidatorSet(MalformedCurrentValidatorSet),
        MalformedNewValidatorSet(MalformedNewValidatorSet),
        /// The standard solidity revert string, with selector
        /// Error(string) -- 0x08c379a0
        RevertString(::std::string::String),
    }
    impl ::ethers::core::abi::AbiDecode for GravityErrors {
        fn decode(
            data: impl AsRef<[u8]>,
        ) -> ::core::result::Result<Self, ::ethers::core::abi::AbiError> {
            let data = data.as_ref();
            if let Ok(decoded) =
                <::std::string::String as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::RevertString(decoded));
            }
            if let Ok(decoded) = <BatchTimedOut as ::ethers::core::abi::AbiDecode>::decode(data) {
                return Ok(Self::BatchTimedOut(decoded));
            }
            if let Ok(decoded) =
                <IncorrectCheckpoint as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::IncorrectCheckpoint(decoded));
            }
            if let Ok(decoded) = <InsufficientPower as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InsufficientPower(decoded));
            }
            if let Ok(decoded) = <InvalidBatchNonce as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidBatchNonce(decoded));
            }
            if let Ok(decoded) =
                <InvalidLogicCallFees as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidLogicCallFees(decoded));
            }
            if let Ok(decoded) =
                <InvalidLogicCallNonce as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidLogicCallNonce(decoded));
            }
            if let Ok(decoded) =
                <InvalidLogicCallTransfers as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidLogicCallTransfers(decoded));
            }
            if let Ok(decoded) =
                <InvalidSendToCosmos as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidSendToCosmos(decoded));
            }
            if let Ok(decoded) = <InvalidSignature as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidSignature(decoded));
            }
            if let Ok(decoded) =
                <InvalidValsetNonce as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::InvalidValsetNonce(decoded));
            }
            if let Ok(decoded) = <LogicCallTimedOut as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::LogicCallTimedOut(decoded));
            }
            if let Ok(decoded) = <MalformedBatch as ::ethers::core::abi::AbiDecode>::decode(data) {
                return Ok(Self::MalformedBatch(decoded));
            }
            if let Ok(decoded) =
                <MalformedCurrentValidatorSet as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::MalformedCurrentValidatorSet(decoded));
            }
            if let Ok(decoded) =
                <MalformedNewValidatorSet as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::MalformedNewValidatorSet(decoded));
            }
            Err(::ethers::core::abi::Error::InvalidData.into())
        }
    }
    impl ::ethers::core::abi::AbiEncode for GravityErrors {
        fn encode(self) -> ::std::vec::Vec<u8> {
            match self {
                Self::BatchTimedOut(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::IncorrectCheckpoint(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::InsufficientPower(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::InvalidBatchNonce(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::InvalidLogicCallFees(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::InvalidLogicCallNonce(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::InvalidLogicCallTransfers(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::InvalidSendToCosmos(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::InvalidSignature(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::InvalidValsetNonce(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::LogicCallTimedOut(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::MalformedBatch(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::MalformedCurrentValidatorSet(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::MalformedNewValidatorSet(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::RevertString(s) => ::ethers::core::abi::AbiEncode::encode(s),
            }
        }
    }
    impl ::ethers::contract::ContractRevert for GravityErrors {
        fn valid_selector(selector: [u8; 4]) -> bool {
            match selector {
                [0x08, 0xc3, 0x79, 0xa0] => true,
                _ if selector == <BatchTimedOut as ::ethers::contract::EthError>::selector() => {
                    true
                }
                _ if selector
                    == <IncorrectCheckpoint as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <InsufficientPower as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <InvalidBatchNonce as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <InvalidLogicCallFees as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <InvalidLogicCallNonce as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <InvalidLogicCallTransfers as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <InvalidSendToCosmos as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector == <InvalidSignature as ::ethers::contract::EthError>::selector() => {
                    true
                }
                _ if selector
                    == <InvalidValsetNonce as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector
                    == <LogicCallTimedOut as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ if selector == <MalformedBatch as ::ethers::contract::EthError>::selector() => {
                    true
                }
                _ if selector
                    == <MalformedCurrentValidatorSet as ::ethers::contract::EthError>::selector(
                    ) =>
                {
                    true
                }
                _ if selector
                    == <MalformedNewValidatorSet as ::ethers::contract::EthError>::selector() =>
                {
                    true
                }
                _ => false,
            }
        }
    }
    impl ::core::fmt::Display for GravityErrors {
        fn fmt(&self, f: &mut ::core::fmt::Formatter<'_>) -> ::core::fmt::Result {
            match self {
                Self::BatchTimedOut(element) => ::core::fmt::Display::fmt(element, f),
                Self::IncorrectCheckpoint(element) => ::core::fmt::Display::fmt(element, f),
                Self::InsufficientPower(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidBatchNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidLogicCallFees(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidLogicCallNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidLogicCallTransfers(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidSendToCosmos(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidSignature(element) => ::core::fmt::Display::fmt(element, f),
                Self::InvalidValsetNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::LogicCallTimedOut(element) => ::core::fmt::Display::fmt(element, f),
                Self::MalformedBatch(element) => ::core::fmt::Display::fmt(element, f),
                Self::MalformedCurrentValidatorSet(element) => {
                    ::core::fmt::Display::fmt(element, f)
                }
                Self::MalformedNewValidatorSet(element) => ::core::fmt::Display::fmt(element, f),
                Self::RevertString(s) => ::core::fmt::Display::fmt(s, f),
            }
        }
    }
    impl ::core::convert::From<::std::string::String> for GravityErrors {
        fn from(value: String) -> Self {
            Self::RevertString(value)
        }
    }
    impl ::core::convert::From<BatchTimedOut> for GravityErrors {
        fn from(value: BatchTimedOut) -> Self {
            Self::BatchTimedOut(value)
        }
    }
    impl ::core::convert::From<IncorrectCheckpoint> for GravityErrors {
        fn from(value: IncorrectCheckpoint) -> Self {
            Self::IncorrectCheckpoint(value)
        }
    }
    impl ::core::convert::From<InsufficientPower> for GravityErrors {
        fn from(value: InsufficientPower) -> Self {
            Self::InsufficientPower(value)
        }
    }
    impl ::core::convert::From<InvalidBatchNonce> for GravityErrors {
        fn from(value: InvalidBatchNonce) -> Self {
            Self::InvalidBatchNonce(value)
        }
    }
    impl ::core::convert::From<InvalidLogicCallFees> for GravityErrors {
        fn from(value: InvalidLogicCallFees) -> Self {
            Self::InvalidLogicCallFees(value)
        }
    }
    impl ::core::convert::From<InvalidLogicCallNonce> for GravityErrors {
        fn from(value: InvalidLogicCallNonce) -> Self {
            Self::InvalidLogicCallNonce(value)
        }
    }
    impl ::core::convert::From<InvalidLogicCallTransfers> for GravityErrors {
        fn from(value: InvalidLogicCallTransfers) -> Self {
            Self::InvalidLogicCallTransfers(value)
        }
    }
    impl ::core::convert::From<InvalidSendToCosmos> for GravityErrors {
        fn from(value: InvalidSendToCosmos) -> Self {
            Self::InvalidSendToCosmos(value)
        }
    }
    impl ::core::convert::From<InvalidSignature> for GravityErrors {
        fn from(value: InvalidSignature) -> Self {
            Self::InvalidSignature(value)
        }
    }
    impl ::core::convert::From<InvalidValsetNonce> for GravityErrors {
        fn from(value: InvalidValsetNonce) -> Self {
            Self::InvalidValsetNonce(value)
        }
    }
    impl ::core::convert::From<LogicCallTimedOut> for GravityErrors {
        fn from(value: LogicCallTimedOut) -> Self {
            Self::LogicCallTimedOut(value)
        }
    }
    impl ::core::convert::From<MalformedBatch> for GravityErrors {
        fn from(value: MalformedBatch) -> Self {
            Self::MalformedBatch(value)
        }
    }
    impl ::core::convert::From<MalformedCurrentValidatorSet> for GravityErrors {
        fn from(value: MalformedCurrentValidatorSet) -> Self {
            Self::MalformedCurrentValidatorSet(value)
        }
    }
    impl ::core::convert::From<MalformedNewValidatorSet> for GravityErrors {
        fn from(value: MalformedNewValidatorSet) -> Self {
            Self::MalformedNewValidatorSet(value)
        }
    }
    #[derive(
        Clone,
        ::ethers::contract::EthEvent,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethevent(
        name = "ERC20DeployedEvent",
        abi = "ERC20DeployedEvent(string,address,string,string,uint8,uint256)"
    )]
    pub struct Erc20DeployedEventFilter {
        pub cosmos_denom: ::std::string::String,
        #[ethevent(indexed)]
        pub token_contract: ::ethers::core::types::Address,
        pub name: ::std::string::String,
        pub symbol: ::std::string::String,
        pub decimals: u8,
        pub event_nonce: ::ethers::core::types::U256,
    }
    #[derive(
        Clone,
        ::ethers::contract::EthEvent,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethevent(
        name = "LogicCallEvent",
        abi = "LogicCallEvent(bytes32,uint256,bytes,uint256)"
    )]
    pub struct LogicCallEventFilter {
        pub invalidation_id: [u8; 32],
        pub invalidation_nonce: ::ethers::core::types::U256,
        pub return_data: ::ethers::core::types::Bytes,
        pub event_nonce: ::ethers::core::types::U256,
    }
    #[derive(
        Clone,
        ::ethers::contract::EthEvent,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethevent(
        name = "SendToCosmosEvent",
        abi = "SendToCosmosEvent(address,address,bytes32,uint256,uint256)"
    )]
    pub struct SendToCosmosEventFilter {
        #[ethevent(indexed)]
        pub token_contract: ::ethers::core::types::Address,
        #[ethevent(indexed)]
        pub sender: ::ethers::core::types::Address,
        #[ethevent(indexed)]
        pub destination: [u8; 32],
        pub amount: ::ethers::core::types::U256,
        pub event_nonce: ::ethers::core::types::U256,
    }
    #[derive(
        Clone,
        ::ethers::contract::EthEvent,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethevent(
        name = "TransactionBatchExecutedEvent",
        abi = "TransactionBatchExecutedEvent(uint256,address,uint256)"
    )]
    pub struct TransactionBatchExecutedEventFilter {
        #[ethevent(indexed)]
        pub batch_nonce: ::ethers::core::types::U256,
        #[ethevent(indexed)]
        pub token: ::ethers::core::types::Address,
        pub event_nonce: ::ethers::core::types::U256,
    }
    #[derive(
        Clone,
        ::ethers::contract::EthEvent,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethevent(
        name = "ValsetUpdatedEvent",
        abi = "ValsetUpdatedEvent(uint256,uint256,uint256,address,address[],uint256[])"
    )]
    pub struct ValsetUpdatedEventFilter {
        #[ethevent(indexed)]
        pub new_valset_nonce: ::ethers::core::types::U256,
        pub event_nonce: ::ethers::core::types::U256,
        pub reward_amount: ::ethers::core::types::U256,
        pub reward_token: ::ethers::core::types::Address,
        pub validators: ::std::vec::Vec<::ethers::core::types::Address>,
        pub powers: ::std::vec::Vec<::ethers::core::types::U256>,
    }
    ///Container type for all of the contract's events
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        serde::Deserialize,
        serde::Serialize,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub enum GravityEvents {
        Erc20DeployedEventFilter(Erc20DeployedEventFilter),
        LogicCallEventFilter(LogicCallEventFilter),
        SendToCosmosEventFilter(SendToCosmosEventFilter),
        TransactionBatchExecutedEventFilter(TransactionBatchExecutedEventFilter),
        ValsetUpdatedEventFilter(ValsetUpdatedEventFilter),
    }
    impl ::ethers::contract::EthLogDecode for GravityEvents {
        fn decode_log(
            log: &::ethers::core::abi::RawLog,
        ) -> ::core::result::Result<Self, ::ethers::core::abi::Error> {
            if let Ok(decoded) = Erc20DeployedEventFilter::decode_log(log) {
                return Ok(GravityEvents::Erc20DeployedEventFilter(decoded));
            }
            if let Ok(decoded) = LogicCallEventFilter::decode_log(log) {
                return Ok(GravityEvents::LogicCallEventFilter(decoded));
            }
            if let Ok(decoded) = SendToCosmosEventFilter::decode_log(log) {
                return Ok(GravityEvents::SendToCosmosEventFilter(decoded));
            }
            if let Ok(decoded) = TransactionBatchExecutedEventFilter::decode_log(log) {
                return Ok(GravityEvents::TransactionBatchExecutedEventFilter(decoded));
            }
            if let Ok(decoded) = ValsetUpdatedEventFilter::decode_log(log) {
                return Ok(GravityEvents::ValsetUpdatedEventFilter(decoded));
            }
            Err(::ethers::core::abi::Error::InvalidData)
        }
    }
    impl ::core::fmt::Display for GravityEvents {
        fn fmt(&self, f: &mut ::core::fmt::Formatter<'_>) -> ::core::fmt::Result {
            match self {
                Self::Erc20DeployedEventFilter(element) => ::core::fmt::Display::fmt(element, f),
                Self::LogicCallEventFilter(element) => ::core::fmt::Display::fmt(element, f),
                Self::SendToCosmosEventFilter(element) => ::core::fmt::Display::fmt(element, f),
                Self::TransactionBatchExecutedEventFilter(element) => {
                    ::core::fmt::Display::fmt(element, f)
                }
                Self::ValsetUpdatedEventFilter(element) => ::core::fmt::Display::fmt(element, f),
            }
        }
    }
    impl ::core::convert::From<Erc20DeployedEventFilter> for GravityEvents {
        fn from(value: Erc20DeployedEventFilter) -> Self {
            Self::Erc20DeployedEventFilter(value)
        }
    }
    impl ::core::convert::From<LogicCallEventFilter> for GravityEvents {
        fn from(value: LogicCallEventFilter) -> Self {
            Self::LogicCallEventFilter(value)
        }
    }
    impl ::core::convert::From<SendToCosmosEventFilter> for GravityEvents {
        fn from(value: SendToCosmosEventFilter) -> Self {
            Self::SendToCosmosEventFilter(value)
        }
    }
    impl ::core::convert::From<TransactionBatchExecutedEventFilter> for GravityEvents {
        fn from(value: TransactionBatchExecutedEventFilter) -> Self {
            Self::TransactionBatchExecutedEventFilter(value)
        }
    }
    impl ::core::convert::From<ValsetUpdatedEventFilter> for GravityEvents {
        fn from(value: ValsetUpdatedEventFilter) -> Self {
            Self::ValsetUpdatedEventFilter(value)
        }
    }
    ///Container type for all input parameters for the `deployERC20` function with signature `deployERC20(string,string,string,uint8)` and selector `0xf7955637`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "deployERC20", abi = "deployERC20(string,string,string,uint8)")]
    pub struct DeployERC20Call {
        pub cosmos_denom: ::std::string::String,
        pub name: ::std::string::String,
        pub symbol: ::std::string::String,
        pub decimals: u8,
    }
    ///Container type for all input parameters for the `lastBatchNonce` function with signature `lastBatchNonce(address)` and selector `0x011b2174`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "lastBatchNonce", abi = "lastBatchNonce(address)")]
    pub struct LastBatchNonceCall {
        pub erc_20_address: ::ethers::core::types::Address,
    }
    ///Container type for all input parameters for the `lastLogicCallNonce` function with signature `lastLogicCallNonce(bytes32)` and selector `0xc9d194d5`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "lastLogicCallNonce", abi = "lastLogicCallNonce(bytes32)")]
    pub struct LastLogicCallNonceCall {
        pub invalidation_id: [u8; 32],
    }
    ///Container type for all input parameters for the `sendToCosmos` function with signature `sendToCosmos(address,bytes32,uint256)` and selector `0x1ffbe7f9`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "sendToCosmos", abi = "sendToCosmos(address,bytes32,uint256)")]
    pub struct SendToCosmosCall {
        pub token_contract: ::ethers::core::types::Address,
        pub destination: [u8; 32],
        pub amount: ::ethers::core::types::U256,
    }
    ///Container type for all input parameters for the `state_gravityId` function with signature `state_gravityId()` and selector `0xbdda81d4`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "state_gravityId", abi = "state_gravityId()")]
    pub struct StateGravityIdCall;
    ///Container type for all input parameters for the `state_invalidationMapping` function with signature `state_invalidationMapping(bytes32)` and selector `0x7dfb6f86`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "state_invalidationMapping",
        abi = "state_invalidationMapping(bytes32)"
    )]
    pub struct StateInvalidationMappingCall(pub [u8; 32]);
    ///Container type for all input parameters for the `state_lastBatchNonces` function with signature `state_lastBatchNonces(address)` and selector `0xdf97174b`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "state_lastBatchNonces", abi = "state_lastBatchNonces(address)")]
    pub struct StateLastBatchNoncesCall(pub ::ethers::core::types::Address);
    ///Container type for all input parameters for the `state_lastEventNonce` function with signature `state_lastEventNonce()` and selector `0x73b20547`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "state_lastEventNonce", abi = "state_lastEventNonce()")]
    pub struct StateLastEventNonceCall;
    ///Container type for all input parameters for the `state_lastValsetCheckpoint` function with signature `state_lastValsetCheckpoint()` and selector `0xf2b53307`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "state_lastValsetCheckpoint",
        abi = "state_lastValsetCheckpoint()"
    )]
    pub struct StateLastValsetCheckpointCall;
    ///Container type for all input parameters for the `state_lastValsetNonce` function with signature `state_lastValsetNonce()` and selector `0xb56561fe`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "state_lastValsetNonce", abi = "state_lastValsetNonce()")]
    pub struct StateLastValsetNonceCall;
    ///Container type for all input parameters for the `state_powerThreshold` function with signature `state_powerThreshold()` and selector `0xe5a2b5d2`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(name = "state_powerThreshold", abi = "state_powerThreshold()")]
    pub struct StatePowerThresholdCall;
    ///Container type for all input parameters for the `submitBatch` function with signature `submitBatch((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],uint256[],address[],uint256[],uint256,address,uint256)` and selector `0x8690ff98`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "submitBatch",
        abi = "submitBatch((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],uint256[],address[],uint256[],uint256,address,uint256)"
    )]
    pub struct SubmitBatchCall {
        pub current_valset: ValsetArgs,
        pub sigs: ::std::vec::Vec<ValSignature>,
        pub amounts: ::std::vec::Vec<::ethers::core::types::U256>,
        pub destinations: ::std::vec::Vec<::ethers::core::types::Address>,
        pub fees: ::std::vec::Vec<::ethers::core::types::U256>,
        pub batch_nonce: ::ethers::core::types::U256,
        pub token_contract: ::ethers::core::types::Address,
        pub batch_timeout: ::ethers::core::types::U256,
    }
    ///Container type for all input parameters for the `submitLogicCall` function with signature `submitLogicCall((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],(uint256[],address[],uint256[],address[],address,bytes,uint256,bytes32,uint256))` and selector `0x6941db93`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "submitLogicCall",
        abi = "submitLogicCall((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],(uint256[],address[],uint256[],address[],address,bytes,uint256,bytes32,uint256))"
    )]
    pub struct SubmitLogicCallCall {
        pub current_valset: ValsetArgs,
        pub sigs: ::std::vec::Vec<ValSignature>,
        pub args: LogicCallArgs,
    }
    ///Container type for all input parameters for the `testCheckValidatorSignatures` function with signature `testCheckValidatorSignatures((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],bytes32,uint256)` and selector `0x00901153`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "testCheckValidatorSignatures",
        abi = "testCheckValidatorSignatures((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],bytes32,uint256)"
    )]
    pub struct TestCheckValidatorSignaturesCall {
        pub current_valset: ValsetArgs,
        pub sigs: ::std::vec::Vec<ValSignature>,
        pub the_hash: [u8; 32],
        pub power_threshold: ::ethers::core::types::U256,
    }
    ///Container type for all input parameters for the `testMakeCheckpoint` function with signature `testMakeCheckpoint((address[],uint256[],uint256,uint256,address),bytes32)` and selector `0x01031525`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "testMakeCheckpoint",
        abi = "testMakeCheckpoint((address[],uint256[],uint256,uint256,address),bytes32)"
    )]
    pub struct TestMakeCheckpointCall {
        pub valset_args: ValsetArgs,
        pub gravity_id: [u8; 32],
    }
    ///Container type for all input parameters for the `updateValset` function with signature `updateValset((address[],uint256[],uint256,uint256,address),(address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[])` and selector `0xaca6b1c1`
    #[derive(
        Clone,
        ::ethers::contract::EthCall,
        ::ethers::contract::EthDisplay,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    #[ethcall(
        name = "updateValset",
        abi = "updateValset((address[],uint256[],uint256,uint256,address),(address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[])"
    )]
    pub struct UpdateValsetCall {
        pub new_valset: ValsetArgs,
        pub current_valset: ValsetArgs,
        pub sigs: ::std::vec::Vec<ValSignature>,
    }
    ///Container type for all of the contract's call
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        serde::Deserialize,
        serde::Serialize,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub enum GravityCalls {
        DeployERC20(DeployERC20Call),
        LastBatchNonce(LastBatchNonceCall),
        LastLogicCallNonce(LastLogicCallNonceCall),
        SendToCosmos(SendToCosmosCall),
        StateGravityId(StateGravityIdCall),
        StateInvalidationMapping(StateInvalidationMappingCall),
        StateLastBatchNonces(StateLastBatchNoncesCall),
        StateLastEventNonce(StateLastEventNonceCall),
        StateLastValsetCheckpoint(StateLastValsetCheckpointCall),
        StateLastValsetNonce(StateLastValsetNonceCall),
        StatePowerThreshold(StatePowerThresholdCall),
        SubmitBatch(SubmitBatchCall),
        SubmitLogicCall(SubmitLogicCallCall),
        TestCheckValidatorSignatures(TestCheckValidatorSignaturesCall),
        TestMakeCheckpoint(TestMakeCheckpointCall),
        UpdateValset(UpdateValsetCall),
    }
    impl ::ethers::core::abi::AbiDecode for GravityCalls {
        fn decode(
            data: impl AsRef<[u8]>,
        ) -> ::core::result::Result<Self, ::ethers::core::abi::AbiError> {
            let data = data.as_ref();
            if let Ok(decoded) = <DeployERC20Call as ::ethers::core::abi::AbiDecode>::decode(data) {
                return Ok(Self::DeployERC20(decoded));
            }
            if let Ok(decoded) =
                <LastBatchNonceCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::LastBatchNonce(decoded));
            }
            if let Ok(decoded) =
                <LastLogicCallNonceCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::LastLogicCallNonce(decoded));
            }
            if let Ok(decoded) = <SendToCosmosCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::SendToCosmos(decoded));
            }
            if let Ok(decoded) =
                <StateGravityIdCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StateGravityId(decoded));
            }
            if let Ok(decoded) =
                <StateInvalidationMappingCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StateInvalidationMapping(decoded));
            }
            if let Ok(decoded) =
                <StateLastBatchNoncesCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StateLastBatchNonces(decoded));
            }
            if let Ok(decoded) =
                <StateLastEventNonceCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StateLastEventNonce(decoded));
            }
            if let Ok(decoded) =
                <StateLastValsetCheckpointCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StateLastValsetCheckpoint(decoded));
            }
            if let Ok(decoded) =
                <StateLastValsetNonceCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StateLastValsetNonce(decoded));
            }
            if let Ok(decoded) =
                <StatePowerThresholdCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::StatePowerThreshold(decoded));
            }
            if let Ok(decoded) = <SubmitBatchCall as ::ethers::core::abi::AbiDecode>::decode(data) {
                return Ok(Self::SubmitBatch(decoded));
            }
            if let Ok(decoded) =
                <SubmitLogicCallCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::SubmitLogicCall(decoded));
            }
            if let Ok(decoded) =
                <TestCheckValidatorSignaturesCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::TestCheckValidatorSignatures(decoded));
            }
            if let Ok(decoded) =
                <TestMakeCheckpointCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::TestMakeCheckpoint(decoded));
            }
            if let Ok(decoded) = <UpdateValsetCall as ::ethers::core::abi::AbiDecode>::decode(data)
            {
                return Ok(Self::UpdateValset(decoded));
            }
            Err(::ethers::core::abi::Error::InvalidData.into())
        }
    }
    impl ::ethers::core::abi::AbiEncode for GravityCalls {
        fn encode(self) -> Vec<u8> {
            match self {
                Self::DeployERC20(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::LastBatchNonce(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::LastLogicCallNonce(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::SendToCosmos(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::StateGravityId(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::StateInvalidationMapping(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::StateLastBatchNonces(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::StateLastEventNonce(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::StateLastValsetCheckpoint(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::StateLastValsetNonce(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::StatePowerThreshold(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::SubmitBatch(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::SubmitLogicCall(element) => ::ethers::core::abi::AbiEncode::encode(element),
                Self::TestCheckValidatorSignatures(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::TestMakeCheckpoint(element) => {
                    ::ethers::core::abi::AbiEncode::encode(element)
                }
                Self::UpdateValset(element) => ::ethers::core::abi::AbiEncode::encode(element),
            }
        }
    }
    impl ::core::fmt::Display for GravityCalls {
        fn fmt(&self, f: &mut ::core::fmt::Formatter<'_>) -> ::core::fmt::Result {
            match self {
                Self::DeployERC20(element) => ::core::fmt::Display::fmt(element, f),
                Self::LastBatchNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::LastLogicCallNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::SendToCosmos(element) => ::core::fmt::Display::fmt(element, f),
                Self::StateGravityId(element) => ::core::fmt::Display::fmt(element, f),
                Self::StateInvalidationMapping(element) => ::core::fmt::Display::fmt(element, f),
                Self::StateLastBatchNonces(element) => ::core::fmt::Display::fmt(element, f),
                Self::StateLastEventNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::StateLastValsetCheckpoint(element) => ::core::fmt::Display::fmt(element, f),
                Self::StateLastValsetNonce(element) => ::core::fmt::Display::fmt(element, f),
                Self::StatePowerThreshold(element) => ::core::fmt::Display::fmt(element, f),
                Self::SubmitBatch(element) => ::core::fmt::Display::fmt(element, f),
                Self::SubmitLogicCall(element) => ::core::fmt::Display::fmt(element, f),
                Self::TestCheckValidatorSignatures(element) => {
                    ::core::fmt::Display::fmt(element, f)
                }
                Self::TestMakeCheckpoint(element) => ::core::fmt::Display::fmt(element, f),
                Self::UpdateValset(element) => ::core::fmt::Display::fmt(element, f),
            }
        }
    }
    impl ::core::convert::From<DeployERC20Call> for GravityCalls {
        fn from(value: DeployERC20Call) -> Self {
            Self::DeployERC20(value)
        }
    }
    impl ::core::convert::From<LastBatchNonceCall> for GravityCalls {
        fn from(value: LastBatchNonceCall) -> Self {
            Self::LastBatchNonce(value)
        }
    }
    impl ::core::convert::From<LastLogicCallNonceCall> for GravityCalls {
        fn from(value: LastLogicCallNonceCall) -> Self {
            Self::LastLogicCallNonce(value)
        }
    }
    impl ::core::convert::From<SendToCosmosCall> for GravityCalls {
        fn from(value: SendToCosmosCall) -> Self {
            Self::SendToCosmos(value)
        }
    }
    impl ::core::convert::From<StateGravityIdCall> for GravityCalls {
        fn from(value: StateGravityIdCall) -> Self {
            Self::StateGravityId(value)
        }
    }
    impl ::core::convert::From<StateInvalidationMappingCall> for GravityCalls {
        fn from(value: StateInvalidationMappingCall) -> Self {
            Self::StateInvalidationMapping(value)
        }
    }
    impl ::core::convert::From<StateLastBatchNoncesCall> for GravityCalls {
        fn from(value: StateLastBatchNoncesCall) -> Self {
            Self::StateLastBatchNonces(value)
        }
    }
    impl ::core::convert::From<StateLastEventNonceCall> for GravityCalls {
        fn from(value: StateLastEventNonceCall) -> Self {
            Self::StateLastEventNonce(value)
        }
    }
    impl ::core::convert::From<StateLastValsetCheckpointCall> for GravityCalls {
        fn from(value: StateLastValsetCheckpointCall) -> Self {
            Self::StateLastValsetCheckpoint(value)
        }
    }
    impl ::core::convert::From<StateLastValsetNonceCall> for GravityCalls {
        fn from(value: StateLastValsetNonceCall) -> Self {
            Self::StateLastValsetNonce(value)
        }
    }
    impl ::core::convert::From<StatePowerThresholdCall> for GravityCalls {
        fn from(value: StatePowerThresholdCall) -> Self {
            Self::StatePowerThreshold(value)
        }
    }
    impl ::core::convert::From<SubmitBatchCall> for GravityCalls {
        fn from(value: SubmitBatchCall) -> Self {
            Self::SubmitBatch(value)
        }
    }
    impl ::core::convert::From<SubmitLogicCallCall> for GravityCalls {
        fn from(value: SubmitLogicCallCall) -> Self {
            Self::SubmitLogicCall(value)
        }
    }
    impl ::core::convert::From<TestCheckValidatorSignaturesCall> for GravityCalls {
        fn from(value: TestCheckValidatorSignaturesCall) -> Self {
            Self::TestCheckValidatorSignatures(value)
        }
    }
    impl ::core::convert::From<TestMakeCheckpointCall> for GravityCalls {
        fn from(value: TestMakeCheckpointCall) -> Self {
            Self::TestMakeCheckpoint(value)
        }
    }
    impl ::core::convert::From<UpdateValsetCall> for GravityCalls {
        fn from(value: UpdateValsetCall) -> Self {
            Self::UpdateValset(value)
        }
    }
    ///Container type for all return fields from the `lastBatchNonce` function with signature `lastBatchNonce(address)` and selector `0x011b2174`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct LastBatchNonceReturn(pub ::ethers::core::types::U256);
    ///Container type for all return fields from the `lastLogicCallNonce` function with signature `lastLogicCallNonce(bytes32)` and selector `0xc9d194d5`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct LastLogicCallNonceReturn(pub ::ethers::core::types::U256);
    ///Container type for all return fields from the `state_gravityId` function with signature `state_gravityId()` and selector `0xbdda81d4`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StateGravityIdReturn(pub [u8; 32]);
    ///Container type for all return fields from the `state_invalidationMapping` function with signature `state_invalidationMapping(bytes32)` and selector `0x7dfb6f86`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StateInvalidationMappingReturn(pub ::ethers::core::types::U256);
    ///Container type for all return fields from the `state_lastBatchNonces` function with signature `state_lastBatchNonces(address)` and selector `0xdf97174b`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StateLastBatchNoncesReturn(pub ::ethers::core::types::U256);
    ///Container type for all return fields from the `state_lastEventNonce` function with signature `state_lastEventNonce()` and selector `0x73b20547`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StateLastEventNonceReturn(pub ::ethers::core::types::U256);
    ///Container type for all return fields from the `state_lastValsetCheckpoint` function with signature `state_lastValsetCheckpoint()` and selector `0xf2b53307`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StateLastValsetCheckpointReturn(pub [u8; 32]);
    ///Container type for all return fields from the `state_lastValsetNonce` function with signature `state_lastValsetNonce()` and selector `0xb56561fe`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StateLastValsetNonceReturn(pub ::ethers::core::types::U256);
    ///Container type for all return fields from the `state_powerThreshold` function with signature `state_powerThreshold()` and selector `0xe5a2b5d2`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct StatePowerThresholdReturn(pub ::ethers::core::types::U256);
    ///`LogicCallArgs(uint256[],address[],uint256[],address[],address,bytes,uint256,bytes32,uint256)`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct LogicCallArgs {
        pub transfer_amounts: ::std::vec::Vec<::ethers::core::types::U256>,
        pub transfer_token_contracts: ::std::vec::Vec<::ethers::core::types::Address>,
        pub fee_amounts: ::std::vec::Vec<::ethers::core::types::U256>,
        pub fee_token_contracts: ::std::vec::Vec<::ethers::core::types::Address>,
        pub logic_contract_address: ::ethers::core::types::Address,
        pub payload: ::ethers::core::types::Bytes,
        pub time_out: ::ethers::core::types::U256,
        pub invalidation_id: [u8; 32],
        pub invalidation_nonce: ::ethers::core::types::U256,
    }
    ///`ValSignature(uint8,bytes32,bytes32)`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct ValSignature {
        pub v: u8,
        pub r: [u8; 32],
        pub s: [u8; 32],
    }
    ///`ValsetArgs(address[],uint256[],uint256,uint256,address)`
    #[derive(
        Clone,
        ::ethers::contract::EthAbiType,
        ::ethers::contract::EthAbiCodec,
        serde::Deserialize,
        serde::Serialize,
        Default,
        Debug,
        PartialEq,
        Eq,
        Hash,
    )]
    pub struct ValsetArgs {
        pub validators: ::std::vec::Vec<::ethers::core::types::Address>,
        pub powers: ::std::vec::Vec<::ethers::core::types::U256>,
        pub valset_nonce: ::ethers::core::types::U256,
        pub reward_amount: ::ethers::core::types::U256,
        pub reward_token: ::ethers::core::types::Address,
    }
}
