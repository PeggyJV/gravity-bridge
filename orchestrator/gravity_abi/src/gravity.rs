pub use gravity_mod::*;
#[allow(clippy::too_many_arguments)]
mod gravity_mod {
    #![allow(clippy::enum_variant_names)]
    #![allow(dead_code)]
    #![allow(clippy::type_complexity)]
    #![allow(unused_imports)]
    use ethers::contract::{
        builders::{ContractCall, Event},
        Contract, Lazy,
    };
    use ethers::core::{
        abi::{Abi, Detokenize, InvalidOutputType, Token, Tokenizable},
        types::*,
    };
    use ethers::providers::Middleware;
    #[doc = "Gravity was auto-generated with ethers-rs Abigen. More information at: https://github.com/gakonst/ethers-rs"]
    use std::sync::Arc;
    pub static GRAVITY_ABI: ethers::contract::Lazy<ethers::core::abi::Abi> =
        ethers::contract::Lazy::new(|| {
            serde_json :: from_str ("[\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_gravityId\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_powerThreshold\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_validators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_powers\",\n        \"type\": \"uint256[]\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"constructor\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"_cosmosDenom\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"_name\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"_symbol\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint8\",\n        \"name\": \"_decimals\",\n        \"type\": \"uint8\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"ERC20DeployedEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes32\",\n        \"name\": \"_invalidationId\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_invalidationNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes\",\n        \"name\": \"_returnData\",\n        \"type\": \"bytes\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"LogicCallEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_sender\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"bytes32\",\n        \"name\": \"_destination\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_amount\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"SendToCosmosEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"uint256\",\n        \"name\": \"_batchNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_token\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"TransactionBatchExecutedEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"uint256\",\n        \"name\": \"_newValsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"address[]\",\n        \"name\": \"_validators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_powers\",\n        \"type\": \"uint256[]\"\n      }\n    ],\n    \"name\": \"ValsetUpdatedEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_cosmosDenom\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_name\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_symbol\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"uint8\",\n        \"name\": \"_decimals\",\n        \"type\": \"uint8\"\n      }\n    ],\n    \"name\": \"deployERC20\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address\",\n        \"name\": \"_erc20Address\",\n        \"type\": \"address\"\n      }\n    ],\n    \"name\": \"lastBatchNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_invalidation_id\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"lastLogicCallNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_destination\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"sendToCosmos\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_gravityId\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"state_invalidationMapping\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address\",\n        \"name\": \"\",\n        \"type\": \"address\"\n      }\n    ],\n    \"name\": \"state_lastBatchNonces\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_lastEventNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_lastValsetCheckpoint\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_lastValsetNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_powerThreshold\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_currentValidators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_currentPowers\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_currentValsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint8[]\",\n        \"name\": \"_v\",\n        \"type\": \"uint8[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_r\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_s\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_amounts\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_destinations\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_fees\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_batchNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_batchTimeout\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"submitBatch\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_currentValidators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_currentPowers\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_currentValsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint8[]\",\n        \"name\": \"_v\",\n        \"type\": \"uint8[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_r\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_s\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"transferAmounts\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"transferTokenContracts\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"feeAmounts\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"feeTokenContracts\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"logicContractAddress\",\n            \"type\": \"address\"\n          },\n          {\n            \"internalType\": \"bytes\",\n            \"name\": \"payload\",\n            \"type\": \"bytes\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"timeOut\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"invalidationId\",\n            \"type\": \"bytes32\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"invalidationNonce\",\n            \"type\": \"uint256\"\n          }\n        ],\n        \"internalType\": \"struct LogicCallArgs\",\n        \"name\": \"_args\",\n        \"type\": \"tuple\"\n      }\n    ],\n    \"name\": \"submitLogicCall\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_currentValidators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_currentPowers\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint8[]\",\n        \"name\": \"_v\",\n        \"type\": \"uint8[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_r\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_s\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_theHash\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_powerThreshold\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"testCheckValidatorSignatures\",\n    \"outputs\": [],\n    \"stateMutability\": \"pure\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_validators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_powers\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_valsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_gravityId\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"testMakeCheckpoint\",\n    \"outputs\": [],\n    \"stateMutability\": \"pure\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_newValidators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_newPowers\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_newValsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_currentValidators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_currentPowers\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_currentValsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint8[]\",\n        \"name\": \"_v\",\n        \"type\": \"uint8[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_r\",\n        \"type\": \"bytes32[]\"\n      },\n      {\n        \"internalType\": \"bytes32[]\",\n        \"name\": \"_s\",\n        \"type\": \"bytes32[]\"\n      }\n    ],\n    \"name\": \"updateValset\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  }\n]\n") . expect ("invalid abi")
        });
    #[derive(Clone)]
    pub struct Gravity<M>(ethers::contract::Contract<M>);
    impl<M> std::ops::Deref for Gravity<M> {
        type Target = ethers::contract::Contract<M>;
        fn deref(&self) -> &Self::Target {
            &self.0
        }
    }
    impl<M: ethers::providers::Middleware> std::fmt::Debug for Gravity<M> {
        fn fmt(&self, f: &mut std::fmt::Formatter) -> std::fmt::Result {
            f.debug_tuple(stringify!(Gravity))
                .field(&self.address())
                .finish()
        }
    }
    impl<'a, M: ethers::providers::Middleware> Gravity<M> {
        #[doc = r" Creates a new contract instance with the specified `ethers`"]
        #[doc = r" client at the given `Address`. The contract derefs to a `ethers::Contract`"]
        #[doc = r" object"]
        pub fn new<T: Into<ethers::core::types::Address>>(
            address: T,
            client: ::std::sync::Arc<M>,
        ) -> Self {
            let contract =
                ethers::contract::Contract::new(address.into(), GRAVITY_ABI.clone(), client);
            Self(contract)
        }
        #[doc = "Calls the contract's `deployERC20` (0xf7955637) function"]
        pub fn deploy_erc20(
            &self,
            cosmos_denom: String,
            name: String,
            symbol: String,
            decimals: u8,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([247, 149, 86, 55], (cosmos_denom, name, symbol, decimals))
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `lastBatchNonce` (0x011b2174) function"]
        pub fn last_batch_nonce(
            &self,
            erc_20_address: ethers::core::types::Address,
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([1, 27, 33, 116], erc_20_address)
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `lastLogicCallNonce` (0xc9d194d5) function"]
        pub fn last_logic_call_nonce(
            &self,
            invalidation_id: [u8; 32],
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([201, 209, 148, 213], invalidation_id)
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `sendToCosmos` (0x1ffbe7f9) function"]
        pub fn send_to_cosmos(
            &self,
            token_contract: ethers::core::types::Address,
            destination: [u8; 32],
            amount: ethers::core::types::U256,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([31, 251, 231, 249], (token_contract, destination, amount))
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_gravityId` (0xbdda81d4) function"]
        pub fn state_gravity_id(&self) -> ethers::contract::builders::ContractCall<M, [u8; 32]> {
            self.0
                .method_hash([189, 218, 129, 212], ())
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_invalidationMapping` (0x7dfb6f86) function"]
        pub fn state_invalidation_mapping(
            &self,
            p0: [u8; 32],
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([125, 251, 111, 134], p0)
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_lastBatchNonces` (0xdf97174b) function"]
        pub fn state_last_batch_nonces(
            &self,
            p0: ethers::core::types::Address,
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([223, 151, 23, 75], p0)
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_lastEventNonce` (0x73b20547) function"]
        pub fn state_last_event_nonce(
            &self,
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([115, 178, 5, 71], ())
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_lastValsetCheckpoint` (0xf2b53307) function"]
        pub fn state_last_valset_checkpoint(
            &self,
        ) -> ethers::contract::builders::ContractCall<M, [u8; 32]> {
            self.0
                .method_hash([242, 181, 51, 7], ())
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_lastValsetNonce` (0xb56561fe) function"]
        pub fn state_last_valset_nonce(
            &self,
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([181, 101, 97, 254], ())
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `state_powerThreshold` (0xe5a2b5d2) function"]
        pub fn state_power_threshold(
            &self,
        ) -> ethers::contract::builders::ContractCall<M, ethers::core::types::U256> {
            self.0
                .method_hash([229, 162, 181, 210], ())
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `submitBatch` (0x83b435db) function"]
        pub fn submit_batch(
            &self,
            current_validators: ::std::vec::Vec<ethers::core::types::Address>,
            current_powers: ::std::vec::Vec<ethers::core::types::U256>,
            current_valset_nonce: ethers::core::types::U256,
            v: ::std::vec::Vec<u8>,
            r: ::std::vec::Vec<[u8; 32]>,
            s: ::std::vec::Vec<[u8; 32]>,
            amounts: ::std::vec::Vec<ethers::core::types::U256>,
            destinations: ::std::vec::Vec<ethers::core::types::Address>,
            fees: ::std::vec::Vec<ethers::core::types::U256>,
            batch_nonce: ethers::core::types::U256,
            token_contract: ethers::core::types::Address,
            batch_timeout: ethers::core::types::U256,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [131, 180, 53, 219],
                    (
                        current_validators,
                        current_powers,
                        current_valset_nonce,
                        v,
                        r,
                        s,
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
        #[doc = "Calls the contract's `submitLogicCall` (0x0c246c82) function"]
        pub fn submit_logic_call(
            &self,
            current_validators: ::std::vec::Vec<ethers::core::types::Address>,
            current_powers: ::std::vec::Vec<ethers::core::types::U256>,
            current_valset_nonce: ethers::core::types::U256,
            v: ::std::vec::Vec<u8>,
            r: ::std::vec::Vec<[u8; 32]>,
            s: ::std::vec::Vec<[u8; 32]>,
            args: LogicCallArgs,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [12, 36, 108, 130],
                    (
                        current_validators,
                        current_powers,
                        current_valset_nonce,
                        v,
                        r,
                        s,
                        args,
                    ),
                )
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `testCheckValidatorSignatures` (0xdb7c4e57) function"]
        pub fn test_check_validator_signatures(
            &self,
            current_validators: ::std::vec::Vec<ethers::core::types::Address>,
            current_powers: ::std::vec::Vec<ethers::core::types::U256>,
            v: ::std::vec::Vec<u8>,
            r: ::std::vec::Vec<[u8; 32]>,
            s: ::std::vec::Vec<[u8; 32]>,
            the_hash: [u8; 32],
            power_threshold: ethers::core::types::U256,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [219, 124, 78, 87],
                    (
                        current_validators,
                        current_powers,
                        v,
                        r,
                        s,
                        the_hash,
                        power_threshold,
                    ),
                )
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `testMakeCheckpoint` (0xc227c30b) function"]
        pub fn test_make_checkpoint(
            &self,
            validators: ::std::vec::Vec<ethers::core::types::Address>,
            powers: ::std::vec::Vec<ethers::core::types::U256>,
            valset_nonce: ethers::core::types::U256,
            gravity_id: [u8; 32],
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [194, 39, 195, 11],
                    (validators, powers, valset_nonce, gravity_id),
                )
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `updateValset` (0xe3cb9f62) function"]
        pub fn update_valset(
            &self,
            new_validators: ::std::vec::Vec<ethers::core::types::Address>,
            new_powers: ::std::vec::Vec<ethers::core::types::U256>,
            new_valset_nonce: ethers::core::types::U256,
            current_validators: ::std::vec::Vec<ethers::core::types::Address>,
            current_powers: ::std::vec::Vec<ethers::core::types::U256>,
            current_valset_nonce: ethers::core::types::U256,
            v: ::std::vec::Vec<u8>,
            r: ::std::vec::Vec<[u8; 32]>,
            s: ::std::vec::Vec<[u8; 32]>,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [227, 203, 159, 98],
                    (
                        new_validators,
                        new_powers,
                        new_valset_nonce,
                        current_validators,
                        current_powers,
                        current_valset_nonce,
                        v,
                        r,
                        s,
                    ),
                )
                .expect("method not found (this should never happen)")
        }
        #[doc = "Gets the contract's `ERC20DeployedEvent` event"]
        pub fn erc20_deployed_event_filter(
            &self,
        ) -> ethers::contract::builders::Event<M, Erc20DeployedEventFilter> {
            self.0.event()
        }
        #[doc = "Gets the contract's `LogicCallEvent` event"]
        pub fn logic_call_event_filter(
            &self,
        ) -> ethers::contract::builders::Event<M, LogicCallEventFilter> {
            self.0.event()
        }
        #[doc = "Gets the contract's `SendToCosmosEvent` event"]
        pub fn send_to_cosmos_event_filter(
            &self,
        ) -> ethers::contract::builders::Event<M, SendToCosmosEventFilter> {
            self.0.event()
        }
        #[doc = "Gets the contract's `TransactionBatchExecutedEvent` event"]
        pub fn transaction_batch_executed_event_filter(
            &self,
        ) -> ethers::contract::builders::Event<M, TransactionBatchExecutedEventFilter> {
            self.0.event()
        }
        #[doc = "Gets the contract's `ValsetUpdatedEvent` event"]
        pub fn valset_updated_event_filter(
            &self,
        ) -> ethers::contract::builders::Event<M, ValsetUpdatedEventFilter> {
            self.0.event()
        }
        #[doc = r" Returns an [`Event`](#ethers_contract::builders::Event) builder for all events of this contract"]
        pub fn events(&self) -> ethers::contract::builders::Event<M, GravityEvents> {
            self.0.event_with_filter(Default::default())
        }
    }
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthEvent,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethevent(
        name = "ERC20DeployedEvent",
        abi = "ERC20DeployedEvent(string,address,string,string,uint8,uint256)"
    )]
    pub struct Erc20DeployedEventFilter {
        pub cosmos_denom: String,
        #[ethevent(indexed)]
        pub token_contract: ethers::core::types::Address,
        pub name: String,
        pub symbol: String,
        pub decimals: u8,
        pub event_nonce: ethers::core::types::U256,
    }
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthEvent,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethevent(
        name = "LogicCallEvent",
        abi = "LogicCallEvent(bytes32,uint256,bytes,uint256)"
    )]
    pub struct LogicCallEventFilter {
        pub invalidation_id: [u8; 32],
        pub invalidation_nonce: ethers::core::types::U256,
        pub return_data: Vec<u8>,
        pub event_nonce: ethers::core::types::U256,
    }
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthEvent,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethevent(
        name = "SendToCosmosEvent",
        abi = "SendToCosmosEvent(address,address,bytes32,uint256,uint256)"
    )]
    pub struct SendToCosmosEventFilter {
        #[ethevent(indexed)]
        pub token_contract: ethers::core::types::Address,
        #[ethevent(indexed)]
        pub sender: ethers::core::types::Address,
        #[ethevent(indexed)]
        pub destination: [u8; 32],
        pub amount: ethers::core::types::U256,
        pub event_nonce: ethers::core::types::U256,
    }
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthEvent,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethevent(
        name = "TransactionBatchExecutedEvent",
        abi = "TransactionBatchExecutedEvent(uint256,address,uint256)"
    )]
    pub struct TransactionBatchExecutedEventFilter {
        #[ethevent(indexed)]
        pub batch_nonce: ethers::core::types::U256,
        #[ethevent(indexed)]
        pub token: ethers::core::types::Address,
        pub event_nonce: ethers::core::types::U256,
    }
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthEvent,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethevent(
        name = "ValsetUpdatedEvent",
        abi = "ValsetUpdatedEvent(uint256,uint256,address[],uint256[])"
    )]
    pub struct ValsetUpdatedEventFilter {
        #[ethevent(indexed)]
        pub new_valset_nonce: ethers::core::types::U256,
        pub event_nonce: ethers::core::types::U256,
        pub validators: Vec<ethers::core::types::Address>,
        pub powers: Vec<ethers::core::types::U256>,
    }
    #[derive(Debug, Clone, PartialEq, Eq)]
    pub enum GravityEvents {
        Erc20DeployedEventFilter(Erc20DeployedEventFilter),
        LogicCallEventFilter(LogicCallEventFilter),
        SendToCosmosEventFilter(SendToCosmosEventFilter),
        TransactionBatchExecutedEventFilter(TransactionBatchExecutedEventFilter),
        ValsetUpdatedEventFilter(ValsetUpdatedEventFilter),
    }
    impl ethers::core::abi::Tokenizable for GravityEvents {
        fn from_token(
            token: ethers::core::abi::Token,
        ) -> Result<Self, ethers::core::abi::InvalidOutputType>
        where
            Self: Sized,
        {
            if let Ok(decoded) = Erc20DeployedEventFilter::from_token(token.clone()) {
                return Ok(GravityEvents::Erc20DeployedEventFilter(decoded));
            }
            if let Ok(decoded) = LogicCallEventFilter::from_token(token.clone()) {
                return Ok(GravityEvents::LogicCallEventFilter(decoded));
            }
            if let Ok(decoded) = SendToCosmosEventFilter::from_token(token.clone()) {
                return Ok(GravityEvents::SendToCosmosEventFilter(decoded));
            }
            if let Ok(decoded) = TransactionBatchExecutedEventFilter::from_token(token.clone()) {
                return Ok(GravityEvents::TransactionBatchExecutedEventFilter(decoded));
            }
            if let Ok(decoded) = ValsetUpdatedEventFilter::from_token(token.clone()) {
                return Ok(GravityEvents::ValsetUpdatedEventFilter(decoded));
            }
            Err(ethers::core::abi::InvalidOutputType(
                "Failed to decode all event variants".to_string(),
            ))
        }
        fn into_token(self) -> ethers::core::abi::Token {
            match self {
                GravityEvents::Erc20DeployedEventFilter(element) => element.into_token(),
                GravityEvents::LogicCallEventFilter(element) => element.into_token(),
                GravityEvents::SendToCosmosEventFilter(element) => element.into_token(),
                GravityEvents::TransactionBatchExecutedEventFilter(element) => element.into_token(),
                GravityEvents::ValsetUpdatedEventFilter(element) => element.into_token(),
            }
        }
    }
    impl ethers::core::abi::TokenizableItem for GravityEvents {}
    impl ethers::contract::EthLogDecode for GravityEvents {
        fn decode_log(log: &ethers::core::abi::RawLog) -> Result<Self, ethers::core::abi::Error>
        where
            Self: Sized,
        {
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
            Err(ethers::core::abi::Error::InvalidData)
        }
    }
    #[doc = "`LogicCallArgs(uint256[],address[],uint256[],address[],address,bytes,uint256,bytes32,uint256)`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthAbiType,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    pub struct LogicCallArgs {
        pub transfer_amounts: Vec<ethers::core::types::U256>,
        pub transfer_token_contracts: Vec<ethers::core::types::Address>,
        pub fee_amounts: Vec<ethers::core::types::U256>,
        pub fee_token_contracts: Vec<ethers::core::types::Address>,
        pub logic_contract_address: ethers::core::types::Address,
        pub payload: Vec<u8>,
        pub time_out: ethers::core::types::U256,
        pub invalidation_id: [u8; 32],
        pub invalidation_nonce: ethers::core::types::U256,
    }
}
