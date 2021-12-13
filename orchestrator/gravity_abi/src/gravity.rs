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
            serde_json :: from_str ("{\n  \"_format\": \"hh-sol-artifact-1\",\n  \"contractName\": \"Gravity\",\n  \"sourceName\": \"contracts/Gravity.sol\",\n  \"abi\": [\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"_gravityId\",\n          \"type\": \"bytes32\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"_powerThreshold\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"internalType\": \"address[]\",\n          \"name\": \"_validators\",\n          \"type\": \"address[]\"\n        },\n        {\n          \"internalType\": \"uint256[]\",\n          \"name\": \"_powers\",\n          \"type\": \"uint256[]\"\n        }\n      ],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"constructor\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"BatchTimedOut\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"IncorrectCheckpoint\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"cumulativePower\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"powerThreshold\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"InsufficientPower\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"newNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"currentNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"InvalidBatchNonce\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"InvalidLogicCallFees\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"newNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"currentNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"InvalidLogicCallNonce\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"InvalidLogicCallTransfers\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"InvalidSendToCosmos\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"InvalidSignature\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"newNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"currentNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"InvalidValsetNonce\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"LogicCallTimedOut\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"MalformedBatch\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"MalformedCurrentValidatorSet\",\n      \"type\": \"error\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"MalformedNewValidatorSet\",\n      \"type\": \"error\"\n    },\n    {\n      \"anonymous\": false,\n      \"inputs\": [\n        {\n          \"indexed\": false,\n          \"internalType\": \"string\",\n          \"name\": \"_cosmosDenom\",\n          \"type\": \"string\"\n        },\n        {\n          \"indexed\": true,\n          \"internalType\": \"address\",\n          \"name\": \"_tokenContract\",\n          \"type\": \"address\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"string\",\n          \"name\": \"_name\",\n          \"type\": \"string\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"string\",\n          \"name\": \"_symbol\",\n          \"type\": \"string\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint8\",\n          \"name\": \"_decimals\",\n          \"type\": \"uint8\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_eventNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"ERC20DeployedEvent\",\n      \"type\": \"event\"\n    },\n    {\n      \"anonymous\": false,\n      \"inputs\": [\n        {\n          \"indexed\": false,\n          \"internalType\": \"bytes32\",\n          \"name\": \"_invalidationId\",\n          \"type\": \"bytes32\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_invalidationNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"bytes\",\n          \"name\": \"_returnData\",\n          \"type\": \"bytes\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_eventNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"LogicCallEvent\",\n      \"type\": \"event\"\n    },\n    {\n      \"anonymous\": false,\n      \"inputs\": [\n        {\n          \"indexed\": true,\n          \"internalType\": \"address\",\n          \"name\": \"_tokenContract\",\n          \"type\": \"address\"\n        },\n        {\n          \"indexed\": true,\n          \"internalType\": \"address\",\n          \"name\": \"_sender\",\n          \"type\": \"address\"\n        },\n        {\n          \"indexed\": true,\n          \"internalType\": \"bytes32\",\n          \"name\": \"_destination\",\n          \"type\": \"bytes32\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_amount\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_eventNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"SendToCosmosEvent\",\n      \"type\": \"event\"\n    },\n    {\n      \"anonymous\": false,\n      \"inputs\": [\n        {\n          \"indexed\": true,\n          \"internalType\": \"uint256\",\n          \"name\": \"_batchNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"indexed\": true,\n          \"internalType\": \"address\",\n          \"name\": \"_token\",\n          \"type\": \"address\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_eventNonce\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"TransactionBatchExecutedEvent\",\n      \"type\": \"event\"\n    },\n    {\n      \"anonymous\": false,\n      \"inputs\": [\n        {\n          \"indexed\": true,\n          \"internalType\": \"uint256\",\n          \"name\": \"_newValsetNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_eventNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256\",\n          \"name\": \"_rewardAmount\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"address\",\n          \"name\": \"_rewardToken\",\n          \"type\": \"address\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"address[]\",\n          \"name\": \"_validators\",\n          \"type\": \"address[]\"\n        },\n        {\n          \"indexed\": false,\n          \"internalType\": \"uint256[]\",\n          \"name\": \"_powers\",\n          \"type\": \"uint256[]\"\n        }\n      ],\n      \"name\": \"ValsetUpdatedEvent\",\n      \"type\": \"event\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"string\",\n          \"name\": \"_cosmosDenom\",\n          \"type\": \"string\"\n        },\n        {\n          \"internalType\": \"string\",\n          \"name\": \"_name\",\n          \"type\": \"string\"\n        },\n        {\n          \"internalType\": \"string\",\n          \"name\": \"_symbol\",\n          \"type\": \"string\"\n        },\n        {\n          \"internalType\": \"uint8\",\n          \"name\": \"_decimals\",\n          \"type\": \"uint8\"\n        }\n      ],\n      \"name\": \"deployERC20\",\n      \"outputs\": [],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"address\",\n          \"name\": \"_erc20Address\",\n          \"type\": \"address\"\n        }\n      ],\n      \"name\": \"lastBatchNonce\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"_invalidation_id\",\n          \"type\": \"bytes32\"\n        }\n      ],\n      \"name\": \"lastLogicCallNonce\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"address\",\n          \"name\": \"_tokenContract\",\n          \"type\": \"address\"\n        },\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"_destination\",\n          \"type\": \"bytes32\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"_amount\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"sendToCosmos\",\n      \"outputs\": [],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"state_gravityId\",\n      \"outputs\": [\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"\",\n          \"type\": \"bytes32\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"\",\n          \"type\": \"bytes32\"\n        }\n      ],\n      \"name\": \"state_invalidationMapping\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"internalType\": \"address\",\n          \"name\": \"\",\n          \"type\": \"address\"\n        }\n      ],\n      \"name\": \"state_lastBatchNonces\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"state_lastEventNonce\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"state_lastValsetCheckpoint\",\n      \"outputs\": [\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"\",\n          \"type\": \"bytes32\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"state_lastValsetNonce\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [],\n      \"name\": \"state_powerThreshold\",\n      \"outputs\": [\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"stateMutability\": \"view\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"components\": [\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"validators\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"powers\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"valsetNonce\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"rewardAmount\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"rewardToken\",\n              \"type\": \"address\"\n            }\n          ],\n          \"internalType\": \"struct ValsetArgs\",\n          \"name\": \"_currentValset\",\n          \"type\": \"tuple\"\n        },\n        {\n          \"components\": [\n            {\n              \"internalType\": \"uint8\",\n              \"name\": \"v\",\n              \"type\": \"uint8\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"r\",\n              \"type\": \"bytes32\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"s\",\n              \"type\": \"bytes32\"\n            }\n          ],\n          \"internalType\": \"struct ValSignature[]\",\n          \"name\": \"_sigs\",\n          \"type\": \"tuple[]\"\n        },\n        {\n          \"internalType\": \"uint256[]\",\n          \"name\": \"_amounts\",\n          \"type\": \"uint256[]\"\n        },\n        {\n          \"internalType\": \"address[]\",\n          \"name\": \"_destinations\",\n          \"type\": \"address[]\"\n        },\n        {\n          \"internalType\": \"uint256[]\",\n          \"name\": \"_fees\",\n          \"type\": \"uint256[]\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"_batchNonce\",\n          \"type\": \"uint256\"\n        },\n        {\n          \"internalType\": \"address\",\n          \"name\": \"_tokenContract\",\n          \"type\": \"address\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"_batchTimeout\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"submitBatch\",\n      \"outputs\": [],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"components\": [\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"validators\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"powers\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"valsetNonce\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"rewardAmount\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"rewardToken\",\n              \"type\": \"address\"\n            }\n          ],\n          \"internalType\": \"struct ValsetArgs\",\n          \"name\": \"_currentValset\",\n          \"type\": \"tuple\"\n        },\n        {\n          \"components\": [\n            {\n              \"internalType\": \"uint8\",\n              \"name\": \"v\",\n              \"type\": \"uint8\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"r\",\n              \"type\": \"bytes32\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"s\",\n              \"type\": \"bytes32\"\n            }\n          ],\n          \"internalType\": \"struct ValSignature[]\",\n          \"name\": \"_sigs\",\n          \"type\": \"tuple[]\"\n        },\n        {\n          \"components\": [\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"transferAmounts\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"transferTokenContracts\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"feeAmounts\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"feeTokenContracts\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"logicContractAddress\",\n              \"type\": \"address\"\n            },\n            {\n              \"internalType\": \"bytes\",\n              \"name\": \"payload\",\n              \"type\": \"bytes\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"timeOut\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"invalidationId\",\n              \"type\": \"bytes32\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"invalidationNonce\",\n              \"type\": \"uint256\"\n            }\n          ],\n          \"internalType\": \"struct LogicCallArgs\",\n          \"name\": \"_args\",\n          \"type\": \"tuple\"\n        }\n      ],\n      \"name\": \"submitLogicCall\",\n      \"outputs\": [],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"components\": [\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"validators\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"powers\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"valsetNonce\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"rewardAmount\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"rewardToken\",\n              \"type\": \"address\"\n            }\n          ],\n          \"internalType\": \"struct ValsetArgs\",\n          \"name\": \"_currentValset\",\n          \"type\": \"tuple\"\n        },\n        {\n          \"components\": [\n            {\n              \"internalType\": \"uint8\",\n              \"name\": \"v\",\n              \"type\": \"uint8\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"r\",\n              \"type\": \"bytes32\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"s\",\n              \"type\": \"bytes32\"\n            }\n          ],\n          \"internalType\": \"struct ValSignature[]\",\n          \"name\": \"_sigs\",\n          \"type\": \"tuple[]\"\n        },\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"_theHash\",\n          \"type\": \"bytes32\"\n        },\n        {\n          \"internalType\": \"uint256\",\n          \"name\": \"_powerThreshold\",\n          \"type\": \"uint256\"\n        }\n      ],\n      \"name\": \"testCheckValidatorSignatures\",\n      \"outputs\": [],\n      \"stateMutability\": \"pure\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"components\": [\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"validators\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"powers\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"valsetNonce\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"rewardAmount\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"rewardToken\",\n              \"type\": \"address\"\n            }\n          ],\n          \"internalType\": \"struct ValsetArgs\",\n          \"name\": \"_valsetArgs\",\n          \"type\": \"tuple\"\n        },\n        {\n          \"internalType\": \"bytes32\",\n          \"name\": \"_gravityId\",\n          \"type\": \"bytes32\"\n        }\n      ],\n      \"name\": \"testMakeCheckpoint\",\n      \"outputs\": [],\n      \"stateMutability\": \"pure\",\n      \"type\": \"function\"\n    },\n    {\n      \"inputs\": [\n        {\n          \"components\": [\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"validators\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"powers\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"valsetNonce\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"rewardAmount\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"rewardToken\",\n              \"type\": \"address\"\n            }\n          ],\n          \"internalType\": \"struct ValsetArgs\",\n          \"name\": \"_newValset\",\n          \"type\": \"tuple\"\n        },\n        {\n          \"components\": [\n            {\n              \"internalType\": \"address[]\",\n              \"name\": \"validators\",\n              \"type\": \"address[]\"\n            },\n            {\n              \"internalType\": \"uint256[]\",\n              \"name\": \"powers\",\n              \"type\": \"uint256[]\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"valsetNonce\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"uint256\",\n              \"name\": \"rewardAmount\",\n              \"type\": \"uint256\"\n            },\n            {\n              \"internalType\": \"address\",\n              \"name\": \"rewardToken\",\n              \"type\": \"address\"\n            }\n          ],\n          \"internalType\": \"struct ValsetArgs\",\n          \"name\": \"_currentValset\",\n          \"type\": \"tuple\"\n        },\n        {\n          \"components\": [\n            {\n              \"internalType\": \"uint8\",\n              \"name\": \"v\",\n              \"type\": \"uint8\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"r\",\n              \"type\": \"bytes32\"\n            },\n            {\n              \"internalType\": \"bytes32\",\n              \"name\": \"s\",\n              \"type\": \"bytes32\"\n            }\n          ],\n          \"internalType\": \"struct ValSignature[]\",\n          \"name\": \"_sigs\",\n          \"type\": \"tuple[]\"\n        }\n      ],\n      \"name\": \"updateValset\",\n      \"outputs\": [],\n      \"stateMutability\": \"nonpayable\",\n      \"type\": \"function\"\n    }\n  ],\n  \"bytecode\": \"0x60a0604052600060045560016005553480156200001b57600080fd5b5060405162003b5e38038062003b5e8339810160408190526200003e916200031d565b60016000558051825114620000665760405163c6617b7b60e01b815260040160405180910390fd5b6000805b8251811015620000c45782818151811062000089576200008962000410565b6020026020010151826200009e91906200043c565b915084821115620000af57620000c4565b80620000bb8162000457565b9150506200006a565b50838111620000f35760405162bfb6ab60e01b8152600481018290526024810185905260440160405180910390fd5b6200012f6040518060a001604052806060815260200160608152602001600081526020016000815260200160006001600160a01b031681525090565b506040805160a081018252848152602081018490526000918101829052606081018290526080810182905290620001678288620001cf565b60808890526006879055600181905560045460055460405192935090917f76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a91620001ba9160009081908b908b90620004ed565b60405180910390a2505050505050506200059c565b6000806918da1958dadc1bda5b9d60b21b60001b90506000838286604001518760000151886020015189606001518a608001516040516020016200021a97969594939291906200053c565b60408051808303601f19018152919052805160209091012095945050505050565b634e487b7160e01b600052604160045260246000fd5b604051601f8201601f191681016001600160401b03811182821017156200027c576200027c6200023b565b604052919050565b60006001600160401b03821115620002a057620002a06200023b565b5060051b60200190565b600082601f830112620002bc57600080fd5b81516020620002d5620002cf8362000284565b62000251565b82815260059290921b84018101918181019086841115620002f557600080fd5b8286015b84811015620003125780518352918301918301620002f9565b509695505050505050565b600080600080608085870312156200033457600080fd5b845160208087015160408801519296509450906001600160401b03808211156200035d57600080fd5b818801915088601f8301126200037257600080fd5b815162000383620002cf8262000284565b81815260059190911b8301840190848101908b831115620003a357600080fd5b938501935b82851015620003da5784516001600160a01b0381168114620003ca5760008081fd5b82529385019390850190620003a8565b60608b01519097509450505080831115620003f457600080fd5b50506200040487828801620002aa565b91505092959194509250565b634e487b7160e01b600052603260045260246000fd5b634e487b7160e01b600052601160045260246000fd5b6000821982111562000452576200045262000426565b500190565b60006000198214156200046e576200046e62000426565b5060010190565b600081518084526020808501945080840160005b83811015620004b05781516001600160a01b03168752958201959082019060010162000489565b509495945050505050565b600081518084526020808501945080840160005b83811015620004b057815187529582019590820190600101620004cf565b85815284602082015260018060a01b038416604082015260a0606082015260006200051c60a083018562000475565b8281036080840152620005308185620004bb565b98975050505050505050565b87815286602082015285604082015260e0606082015260006200056360e083018762000475565b8281036080840152620005778187620004bb565b60a084019590955250506001600160a01b039190911660c09091015295945050505050565b608051613591620005cd60003960008181610220015281816105700152818161060b01526109c501526135916000f3fe60806040523480156200001157600080fd5b5060043610620001145760003560e01c8063aca6b1c111620000a3578063df97174b116200006e578063df97174b1462000265578063e5a2b5d21462000288578063f2b533071462000292578063f7955637146200029c57600080fd5b8063aca6b1c114620001f9578063b56561fe1462000210578063bdda81d4146200021a578063c9d194d5146200024257600080fd5b80636941db9311620000e45780636941db93146200019e57806373b2054714620001b55780637dfb6f8614620001bf5780638690ff9814620001e257600080fd5b80629011531462000119578063010315251462000132578063011b217414620001495780631ffbe7f91462000187575b600080fd5b620001306200012a366004620018aa565b620002b3565b005b62000130620001433660046200192b565b620002c9565b620001756200015a36600462001990565b6001600160a01b031660009081526002602052604090205490565b60405190815260200160405180910390f35b6200013062000198366004620019ae565b620002e4565b62000130620001af36600462001bd3565b620004a2565b6200017560055481565b62000175620001d036600462001d6f565b60036020526000908152604090205481565b62000130620001f336600462001dd0565b62000835565b620001306200020a36600462001f09565b62000b88565b6200017560045481565b620001757f000000000000000000000000000000000000000000000000000000000000000081565b620001756200025336600462001d6f565b60009081526003602052604090205490565b620001756200027636600462001990565b60026020526000908152604090205481565b6200017560065481565b6200017560015481565b62000130620002ad36600462001ff3565b62000e72565b620002c2858585858562000f2a565b5050505050565b620002df620002d883620020a9565b8262001084565b505050565b60026000541415620003135760405162461bcd60e51b81526004016200030a906200215f565b60405180910390fd5b600260009081556040516370a0823160e01b81523060048201526001600160a01b038516906370a0823190602401602060405180830381865afa1580156200035f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000385919062002196565b90506200039e6001600160a01b038516333085620010f0565b6040516370a0823160e01b81523060048201526000906001600160a01b038616906370a0823190602401602060405180830381865afa158015620003e6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200040c919062002196565b90508181116200042f576040516321739d9b60e01b815260040160405180910390fd5b6005546200043f906001620021c6565b60055583336001600160a01b0387167fd7767894d73c589daeca9643f445f03d7be61aad2950c117e7cbff4176fca7e46200047b8686620021e1565b6005546040805192835260208301919091520160405180910390a450506001600055505050565b60026000541415620004c85760405162461bcd60e51b81526004016200030a906200215f565b600260005560c08101514310620004f25760405163bcf37c2560e01b815260040160405180910390fd5b61010081015160e0820151600090815260036020526040902054106200054f5761010081015160e082015160009081526003602052604090819020549051629427e960e11b8152600481019290925260248201526044016200030a565b6200055c84848462001163565b600154620005956200056e86620020a9565b7f000000000000000000000000000000000000000000000000000000000000000062001084565b14620005b45760405163723a340360e01b815260040160405180910390fd5b60208101515181515114620005dc57604051634298a95160e11b815260040160405180910390fd5b806060015151816040015151146200060757604051634829247960e01b815260040160405180910390fd5b60007f0000000000000000000000000000000000000000000000000000000000000000681b1bd9da58d0d85b1b60ba1b836000015184602001518560400151866060015187608001518860a001518960c001518a60e001518b6101000151604051602001620006819b9a99989796959493929190620022d0565b604051602081830303815290604052805190602001209050620006aa8585858460065462000f2a565b5061010081015160e08201516000908152600360205260408120919091555b8151518110156200074d5762000738826080015183600001518381518110620006f657620006f662002380565b60200260200101518460200151848151811062000717576200071762002380565b60200260200101516001600160a01b0316620011bc9092919063ffffffff16565b80620007448162002396565b915050620006c9565b5060006200076482608001518360a00151620011ee565b905060005b826040015151811015620007cb57620007b6338460400151838151811062000795576200079562002380565b60200260200101518560600151848151811062000717576200071762002380565b80620007c28162002396565b91505062000769565b50600554620007dc906001620021c6565b600581905560e08301516101008401516040517f7c2bb24f8e1b3725cb613d7f11ef97d9745cc97a0e40f730621c052d684077a19362000821939291869190620023b4565b60405180910390a150506001600055505050565b600260005414156200085b5760405162461bcd60e51b81526004016200030a906200215f565b600260008181556001600160a01b038416815260209190915260409020548311620008bf576001600160a01b0382166000908152600260205260409081902054905163f7f920ad60e01b81526200030a918591600401918252602082015260400190565b6001600160a01b038216600090815260026020526040902054620008e790620f4240620021c6565b8311156200092e576001600160a01b0382166000908152600260205260409081902054905163f7f920ad60e01b81526200030a918591600401918252602082015260400190565b8043106200094f576040516308b9266360e11b815260040160405180910390fd5b6200095c8c8c8c62001163565b6001546200096e6200056e8e620020a9565b146200098d5760405163723a340360e01b815260040160405180910390fd5b87861415806200099d5750878414155b15620009bc5760405163c1f97e3560e01b815260040160405180910390fd5b62000a3c8c8c8c7f00000000000000000000000000000000000000000000000000000000000000006f0e8e4c2dce6c2c6e8d2dedc84c2e8c6d60831b8e8e8e8e8e8e8e8e8e60405160200162000a1d9b9a999897969594939291906200245d565b6040516020818303038152906040528051906020012060065462000f2a565b6001600160a01b0382166000908152600260205260408120849055805b8981101562000b085762000ac889898381811062000a7b5762000a7b62002380565b905060200201602081019062000a92919062001990565b8c8c8481811062000aa75762000aa762002380565b90506020020135866001600160a01b0316620011bc9092919063ffffffff16565b86868281811062000add5762000add62002380565b905060200201358262000af19190620021c6565b91508062000aff8162002396565b91505062000a59565b5062000b1f6001600160a01b0384163383620011bc565b5060055462000b30906001620021c6565b60058190556040519081526001600160a01b0383169084907f02c7e81975f8edb86e2a0c038b7b86a49c744236abf0f6177ff5afc6986ab7089060200160405180910390a35050600160005550505050505050505050565b826040013584604001351162000bc1576040805163e0e8edf360e01b81528186013560048201529084013560248201526044016200030a565b62000bd46040840135620f4240620021c6565b8460400135111562000c09576040805163e0e8edf360e01b81528186013560048201529084013560248201526044016200030a565b62000c186020850185620024db565b905062000c268580620024db565b905014158062000c41575062000c3d8480620024db565b1590505b1562000c605760405163c01ba0ab60e01b815260040160405180910390fd5b62000c6d83838362001163565b6000805b62000c806020870187620024db565b905081101562000ce85762000c996020870187620024db565b8281811062000cac5762000cac62002380565b905060200201358262000cc09190620021c6565b915060065482111562000cd35762000ce8565b8062000cdf8162002396565b91505062000c71565b50600654811162000d1a5760065460405162bfb6ab60e01b81526200030a918391600401918252602082015260400190565b60015462000d2c6200056e86620020a9565b1462000d4b5760405163723a340360e01b815260040160405180910390fd5b600062000d5c6200056e87620020a9565b905062000d6f8585858460065462000f2a565b60018190556040860135600455600062000d9060a088016080890162001990565b6001600160a01b03161415801562000dab5750606086013515155b1562000dde5762000dde33606088013562000dcd60a08a0160808b0162001990565b6001600160a01b03169190620011bc565b60055462000dee906001620021c6565b60058190556040870135907f76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a90606089013562000e3260a08b0160808c0162001990565b62000e3e8b80620024db565b62000e4d60208e018e620024db565b60405162000e62979695949392919062002527565b60405180910390a2505050505050565b600030868686868660405162000e889062001835565b62000e9996959493929190620025a3565b604051809103906000f08015801562000eb6573d6000803e3d6000fd5b5060055490915062000eca90600162001239565b60058190556040516001600160a01b038316917f82fe3a4fa49c6382d0c085746698ddbbafe6c2bf61285b19410644b5b26287c79162000f18918c918c918c918c918c918c918c91620025f6565b60405180910390a25050505050505050565b6000805b62000f3a8780620024db565b9050811015620010505785858281811062000f595762000f5962002380565b62000f71926020606090920201908101915062002653565b60ff16156200103b5762000fd662000f8a8880620024db565b8381811062000f9d5762000f9d62002380565b905060200201602081019062000fb4919062001990565b8588888581811062000fca5762000fca62002380565b90506060020162001247565b62000ff457604051638baa579f60e01b815260040160405180910390fd5b620010036020880188620024db565b8281811062001016576200101662002380565b90506020020135826200102a9190620021c6565b9150828211156200103b5762001050565b80620010478162002396565b91505062000f2e565b508181116200107c5760405162bfb6ab60e01b815260048101829052602481018390526044016200030a565b505050505050565b6000806918da1958dadc1bda5b9d60b21b60001b90506000838286604001518760000151886020015189606001518a60800151604051602001620010cf979695949392919062002671565b60408051808303601f19018152919052805160209091012095945050505050565b6040516001600160a01b03808516602483015283166044820152606481018290526200115d9085906323b872dd60e01b906084015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152620012dc565b50505050565b620011726020840184620024db565b9050620011808480620024db565b90501415806200119d575080620011988480620024db565b905014155b15620002df5760405163c6617b7b60e01b815260040160405180910390fd5b6040516001600160a01b038316602482015260448101829052620002df90849063a9059cbb60e01b9060640162001125565b60606200123283836040518060400160405280601e81526020017f416464726573733a206c6f772d6c6576656c2063616c6c206661696c65640000815250620013b5565b9392505050565b6000620012328284620021c6565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c81018390526000908190605c0160408051601f1981840301815291905280516020918201209150620012be908290620012ae9086018662002653565b85602001358660400135620013ce565b6001600160a01b0316856001600160a01b0316149150509392505050565b600062001333826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316620013b59092919063ffffffff16565b805190915015620002df5780806020019051810190620013549190620026d1565b620002df5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016200030a565b6060620013c68484600085620013fa565b949350505050565b6000806000620013e1878787876200152c565b91509150620013f08162001621565b5095945050505050565b6060824710156200145d5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b60648201526084016200030a565b843b620014ad5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016200030a565b600080866001600160a01b03168587604051620014cb9190620026f5565b60006040518083038185875af1925050503d80600081146200150a576040519150601f19603f3d011682016040523d82523d6000602084013e6200150f565b606091505b509150915062001521828286620017f7565b979650505050505050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a083111562001565575060009050600362001618565b8460ff16601b141580156200157e57508460ff16601c14155b1562001591575060009050600462001618565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015620015e6573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116620016115760006001925092505062001618565b9150600090505b94509492505050565b600081600481111562001638576200163862002713565b1415620016425750565b600181600481111562001659576200165962002713565b1415620016a95760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016200030a565b6002816004811115620016c057620016c062002713565b1415620017105760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016200030a565b600381600481111562001727576200172762002713565b1415620017825760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016200030a565b600481600481111562001799576200179962002713565b1415620017f45760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016200030a565b50565b606083156200180857508162001232565b825115620018195782518084602001fd5b8160405162461bcd60e51b81526004016200030a919062002729565b610e1d806200273f83390190565b600060a082840312156200185657600080fd5b50919050565b60008083601f8401126200186f57600080fd5b5081356001600160401b038111156200188757600080fd5b602083019150836020606083028501011115620018a357600080fd5b9250929050565b600080600080600060808688031215620018c357600080fd5b85356001600160401b0380821115620018db57600080fd5b620018e989838a0162001843565b965060208801359150808211156200190057600080fd5b506200190f888289016200185c565b9699909850959660408101359660609091013595509350505050565b600080604083850312156200193f57600080fd5b82356001600160401b038111156200195657600080fd5b620019648582860162001843565b95602094909401359450505050565b80356001600160a01b03811681146200198b57600080fd5b919050565b600060208284031215620019a357600080fd5b620012328262001973565b600080600060608486031215620019c457600080fd5b620019cf8462001973565b95602085013595506040909401359392505050565b634e487b7160e01b600052604160045260246000fd5b60405161012081016001600160401b038111828210171562001a205762001a20620019e4565b60405290565b604051601f8201601f191681016001600160401b038111828210171562001a515762001a51620019e4565b604052919050565b60006001600160401b0382111562001a755762001a75620019e4565b5060051b60200190565b600082601f83011262001a9157600080fd5b8135602062001aaa62001aa48362001a59565b62001a26565b82815260059290921b8401810191818101908684111562001aca57600080fd5b8286015b8481101562001ae7578035835291830191830162001ace565b509695505050505050565b600082601f83011262001b0457600080fd5b8135602062001b1762001aa48362001a59565b82815260059290921b8401810191818101908684111562001b3757600080fd5b8286015b8481101562001ae75762001b4f8162001973565b835291830191830162001b3b565b600082601f83011262001b6f57600080fd5b81356001600160401b0381111562001b8b5762001b8b620019e4565b62001ba0601f8201601f191660200162001a26565b81815284602083860101111562001bb657600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000806060858703121562001bea57600080fd5b84356001600160401b038082111562001c0257600080fd5b62001c108883890162001843565b9550602087013591508082111562001c2757600080fd5b62001c35888389016200185c565b9095509350604087013591508082111562001c4f57600080fd5b90860190610120828903121562001c6557600080fd5b62001c6f620019fa565b82358281111562001c7f57600080fd5b62001c8d8a82860162001a7f565b82525060208301358281111562001ca357600080fd5b62001cb18a82860162001af2565b60208301525060408301358281111562001cca57600080fd5b62001cd88a82860162001a7f565b60408301525060608301358281111562001cf157600080fd5b62001cff8a82860162001af2565b60608301525062001d136080840162001973565b608082015260a08301358281111562001d2b57600080fd5b62001d398a82860162001b5d565b60a08301525060c083013560c082015260e083013560e08201526101009150818301358282015280935050505092959194509250565b60006020828403121562001d8257600080fd5b5035919050565b60008083601f84011262001d9c57600080fd5b5081356001600160401b0381111562001db457600080fd5b6020830191508360208260051b8501011115620018a357600080fd5b6000806000806000806000806000806000806101008d8f03121562001df457600080fd5b6001600160401b038d35111562001e0a57600080fd5b62001e198e8e358f0162001843565b9b506001600160401b0360208e0135111562001e3457600080fd5b62001e468e60208f01358f016200185c565b909b5099506001600160401b0360408e0135111562001e6457600080fd5b62001e768e60408f01358f0162001d89565b90995097506001600160401b0360608e0135111562001e9457600080fd5b62001ea68e60608f01358f0162001d89565b90975095506001600160401b0360808e0135111562001ec457600080fd5b62001ed68e60808f01358f0162001d89565b909550935060a08d0135925062001ef060c08e0162001973565b915060e08d013590509295989b509295989b509295989b565b6000806000806060858703121562001f2057600080fd5b84356001600160401b038082111562001f3857600080fd5b62001f468883890162001843565b9550602087013591508082111562001f5d57600080fd5b62001f6b8883890162001843565b9450604087013591508082111562001f8257600080fd5b5062001f91878288016200185c565b95989497509550505050565b60008083601f84011262001fb057600080fd5b5081356001600160401b0381111562001fc857600080fd5b602083019150836020828501011115620018a357600080fd5b803560ff811681146200198b57600080fd5b60008060008060008060006080888a0312156200200f57600080fd5b87356001600160401b03808211156200202757600080fd5b620020358b838c0162001f9d565b909950975060208a01359150808211156200204f57600080fd5b6200205d8b838c0162001f9d565b909750955060408a01359150808211156200207757600080fd5b50620020868a828b0162001f9d565b90945092506200209b90506060890162001fe1565b905092959891949750929550565b600060a08236031215620020bc57600080fd5b60405160a081016001600160401b038282108183111715620020e257620020e2620019e4565b816040528435915080821115620020f857600080fd5b620021063683870162001af2565b835260208501359150808211156200211d57600080fd5b506200212c3682860162001a7f565b6020830152506040830135604082015260608301356060820152620021546080840162001973565b608082015292915050565b6020808252601f908201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00604082015260600190565b600060208284031215620021a957600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b60008219821115620021dc57620021dc620021b0565b500190565b600082821015620021f657620021f6620021b0565b500390565b600081518084526020808501945080840160005b838110156200222d578151875295820195908201906001016200220f565b509495945050505050565b600081518084526020808501945080840160005b838110156200222d5781516001600160a01b0316875295820195908201906001016200224c565b60005b838110156200229057818101518382015260200162002276565b838111156200115d5750506000910152565b60008151808452620022bc81602086016020860162002273565b601f01601f19169290920160200192915050565b60006101608d83528c6020840152806040840152620022f28184018d620021fb565b9050828103606084015262002308818c62002238565b905082810360808401526200231e818b620021fb565b905082810360a084015262002334818a62002238565b6001600160a01b03891660c085015283810360e08501529050620023598188620022a2565b61010084019690965250506101208101929092526101409091015298975050505050505050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415620023ad57620023ad620021b0565b5060010190565b848152836020820152608060408201526000620023d56080830185620022a2565b905082606083015295945050505050565b81835260006001600160fb1b038311156200240057600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183526000602080850194508260005b858110156200222d576001600160a01b03620024498362001973565b16875295820195908201906001016200242d565b60006101008d83528c6020840152806040840152620024808184018c8e620023e6565b9050828103606084015262002497818a8c6200241d565b90508281036080840152620024ae81888a620023e6565b60a084019690965250506001600160a01b039290921660c083015260e09091015298975050505050505050565b6000808335601e19843603018112620024f357600080fd5b8301803591506001600160401b038211156200250e57600080fd5b6020019150600581901b3603821315620018a357600080fd5b87815286602082015260018060a01b038616604082015260a0606082015260006200255760a0830186886200241d565b82810360808401526200256c818587620023e6565b9a9950505050505050505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0387168152608060208201819052600090620025ca90830187896200257a565b8281036040840152620025df8186886200257a565b91505060ff83166060830152979650505050505050565b60a0815260006200260c60a083018a8c6200257a565b82810360208401526200262181898b6200257a565b90508281036040840152620026388187896200257a565b60ff9590951660608401525050608001529695505050505050565b6000602082840312156200266657600080fd5b620012328262001fe1565b87815286602082015285604082015260e0606082015260006200269860e083018762002238565b8281036080840152620026ac8187620021fb565b60a084019590955250506001600160a01b039190911660c09091015295945050505050565b600060208284031215620026e457600080fd5b815180151581146200123257600080fd5b600082516200270981846020870162002273565b9190910192915050565b634e487b7160e01b600052602160045260246000fd5b602081526000620012326020830184620022a256fe60806040526000196005553480156200001757600080fd5b5060405162000e1d38038062000e1d8339810160408190526200003a9162000311565b825183908390620000539060039060208501906200019e565b508051620000699060049060208401906200019e565b5050600680546001600160a81b031916600160a01b60ff8516026001600160a01b031916176001600160a01b03871617905550600554620000ac908590620000b6565b5050505062000419565b6001600160a01b038216620001115760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b8060026000828254620001259190620003b5565b90915550506001600160a01b0382166000908152602081905260408120805483929062000154908490620003b5565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b828054620001ac90620003dc565b90600052602060002090601f016020900481019282620001d057600085556200021b565b82601f10620001eb57805160ff19168380011785556200021b565b828001600101855582156200021b579182015b828111156200021b578251825591602001919060010190620001fe565b50620002299291506200022d565b5090565b5b808211156200022957600081556001016200022e565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200026c57600080fd5b81516001600160401b038082111562000289576200028962000244565b604051601f8301601f19908116603f01168101908282118183101715620002b457620002b462000244565b81604052838152602092508683858801011115620002d157600080fd5b600091505b83821015620002f55785820183015181830184015290820190620002d6565b83821115620003075760008385830101525b9695505050505050565b600080600080608085870312156200032857600080fd5b84516001600160a01b03811681146200034057600080fd5b60208601519094506001600160401b03808211156200035e57600080fd5b6200036c888389016200025a565b945060408701519150808211156200038357600080fd5b5062000392878288016200025a565b925050606085015160ff81168114620003aa57600080fd5b939692955090935050565b60008219821115620003d757634e487b7160e01b600052601160045260246000fd5b500190565b600181811c90821680620003f157607f821691505b602082108114156200041357634e487b7160e01b600052602260045260246000fd5b50919050565b6109f480620004296000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80635fd130a91161008c578063a457c2d711610066578063a457c2d7146101b6578063a9059cbb146101c9578063cbf0a64e146101dc578063dd62ed3e1461020757600080fd5b80635fd130a91461017057806370a082311461018557806395d89b41146101ae57600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd1461012b578063313ce5671461013e578063395093511461015d575b600080fd5b6100dc610240565b6040516100e99190610812565b60405180910390f35b610105610100366004610883565b6102d2565b60405190151581526020016100e9565b61011d6102e8565b6040519081526020016100e9565b6101056101393660046108ad565b610315565b600654600160a01b900460ff1660405160ff90911681526020016100e9565b61010561016b366004610883565b6103c4565b61018361017e3660046108e9565b610400565b005b61011d6101933660046108e9565b6001600160a01b031660009081526020819052604090205490565b6100dc61046a565b6101056101c4366004610883565b610479565b6101056101d7366004610883565b610512565b6006546101ef906001600160a01b031681565b6040516001600160a01b0390911681526020016100e9565b61011d61021536600461090b565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60606003805461024f9061093e565b80601f016020809104026020016040519081016040528092919081815260200182805461027b9061093e565b80156102c85780601f1061029d576101008083540402835291602001916102c8565b820191906000526020600020905b8154815290600101906020018083116102ab57829003601f168201915b5050505050905090565b60006102df33848461051f565b50600192915050565b6006546001600160a01b0316600090815260208190526040812054600554610310919061098f565b905090565b6000610322848484610643565b6001600160a01b0384166000908152600160209081526040808320338452909152902054828110156103ac5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084015b60405180910390fd5b6103b9853385840361051f565b506001949350505050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102df9185906103fb9086906109a6565b61051f565b6006546001600160a01b031633146104485760405162461bcd60e51b815260206004820152600b60248201526a4e6f74206772617669747960a81b60448201526064016103a3565b600680546001600160a01b0319166001600160a01b0392909216919091179055565b60606004805461024f9061093e565b3360009081526001602090815260408083206001600160a01b0386168452909152812054828110156104fb5760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084016103a3565b610508338585840361051f565b5060019392505050565b60006102df338484610643565b6001600160a01b0383166105815760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016103a3565b6001600160a01b0382166105e25760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016103a3565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6001600160a01b0383166106a75760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b60648201526084016103a3565b6001600160a01b0382166107095760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b60648201526084016103a3565b6001600160a01b038316600090815260208190526040902054818110156107815760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b60648201526084016103a3565b6001600160a01b038085166000908152602081905260408082208585039055918516815290812080548492906107b89084906109a6565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161080491815260200190565b60405180910390a350505050565b600060208083528351808285015260005b8181101561083f57858101830151858201604001528201610823565b81811115610851576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b038116811461087e57600080fd5b919050565b6000806040838503121561089657600080fd5b61089f83610867565b946020939093013593505050565b6000806000606084860312156108c257600080fd5b6108cb84610867565b92506108d960208501610867565b9150604084013590509250925092565b6000602082840312156108fb57600080fd5b61090482610867565b9392505050565b6000806040838503121561091e57600080fd5b61092783610867565b915061093560208401610867565b90509250929050565b600181811c9082168061095257607f821691505b6020821081141561097357634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156109a1576109a1610979565b500390565b600082198211156109b9576109b9610979565b50019056fea264697066735822122010756c7691f24fa76ef182760b9f99a094a8bf02470a418707f1de88a7d8d7b364736f6c634300080a0033a2646970667358221220e0cc737c893f4f801f620d372b01df0f8d6fe82c9585b72e83600f517f252d7a64736f6c634300080a0033\",\n  \"deployedBytecode\": \"0x60806040523480156200001157600080fd5b5060043610620001145760003560e01c8063aca6b1c111620000a3578063df97174b116200006e578063df97174b1462000265578063e5a2b5d21462000288578063f2b533071462000292578063f7955637146200029c57600080fd5b8063aca6b1c114620001f9578063b56561fe1462000210578063bdda81d4146200021a578063c9d194d5146200024257600080fd5b80636941db9311620000e45780636941db93146200019e57806373b2054714620001b55780637dfb6f8614620001bf5780638690ff9814620001e257600080fd5b80629011531462000119578063010315251462000132578063011b217414620001495780631ffbe7f91462000187575b600080fd5b620001306200012a366004620018aa565b620002b3565b005b62000130620001433660046200192b565b620002c9565b620001756200015a36600462001990565b6001600160a01b031660009081526002602052604090205490565b60405190815260200160405180910390f35b6200013062000198366004620019ae565b620002e4565b62000130620001af36600462001bd3565b620004a2565b6200017560055481565b62000175620001d036600462001d6f565b60036020526000908152604090205481565b62000130620001f336600462001dd0565b62000835565b620001306200020a36600462001f09565b62000b88565b6200017560045481565b620001757f000000000000000000000000000000000000000000000000000000000000000081565b620001756200025336600462001d6f565b60009081526003602052604090205490565b620001756200027636600462001990565b60026020526000908152604090205481565b6200017560065481565b6200017560015481565b62000130620002ad36600462001ff3565b62000e72565b620002c2858585858562000f2a565b5050505050565b620002df620002d883620020a9565b8262001084565b505050565b60026000541415620003135760405162461bcd60e51b81526004016200030a906200215f565b60405180910390fd5b600260009081556040516370a0823160e01b81523060048201526001600160a01b038516906370a0823190602401602060405180830381865afa1580156200035f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019062000385919062002196565b90506200039e6001600160a01b038516333085620010f0565b6040516370a0823160e01b81523060048201526000906001600160a01b038616906370a0823190602401602060405180830381865afa158015620003e6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906200040c919062002196565b90508181116200042f576040516321739d9b60e01b815260040160405180910390fd5b6005546200043f906001620021c6565b60055583336001600160a01b0387167fd7767894d73c589daeca9643f445f03d7be61aad2950c117e7cbff4176fca7e46200047b8686620021e1565b6005546040805192835260208301919091520160405180910390a450506001600055505050565b60026000541415620004c85760405162461bcd60e51b81526004016200030a906200215f565b600260005560c08101514310620004f25760405163bcf37c2560e01b815260040160405180910390fd5b61010081015160e0820151600090815260036020526040902054106200054f5761010081015160e082015160009081526003602052604090819020549051629427e960e11b8152600481019290925260248201526044016200030a565b6200055c84848462001163565b600154620005956200056e86620020a9565b7f000000000000000000000000000000000000000000000000000000000000000062001084565b14620005b45760405163723a340360e01b815260040160405180910390fd5b60208101515181515114620005dc57604051634298a95160e11b815260040160405180910390fd5b806060015151816040015151146200060757604051634829247960e01b815260040160405180910390fd5b60007f0000000000000000000000000000000000000000000000000000000000000000681b1bd9da58d0d85b1b60ba1b836000015184602001518560400151866060015187608001518860a001518960c001518a60e001518b6101000151604051602001620006819b9a99989796959493929190620022d0565b604051602081830303815290604052805190602001209050620006aa8585858460065462000f2a565b5061010081015160e08201516000908152600360205260408120919091555b8151518110156200074d5762000738826080015183600001518381518110620006f657620006f662002380565b60200260200101518460200151848151811062000717576200071762002380565b60200260200101516001600160a01b0316620011bc9092919063ffffffff16565b80620007448162002396565b915050620006c9565b5060006200076482608001518360a00151620011ee565b905060005b826040015151811015620007cb57620007b6338460400151838151811062000795576200079562002380565b60200260200101518560600151848151811062000717576200071762002380565b80620007c28162002396565b91505062000769565b50600554620007dc906001620021c6565b600581905560e08301516101008401516040517f7c2bb24f8e1b3725cb613d7f11ef97d9745cc97a0e40f730621c052d684077a19362000821939291869190620023b4565b60405180910390a150506001600055505050565b600260005414156200085b5760405162461bcd60e51b81526004016200030a906200215f565b600260008181556001600160a01b038416815260209190915260409020548311620008bf576001600160a01b0382166000908152600260205260409081902054905163f7f920ad60e01b81526200030a918591600401918252602082015260400190565b6001600160a01b038216600090815260026020526040902054620008e790620f4240620021c6565b8311156200092e576001600160a01b0382166000908152600260205260409081902054905163f7f920ad60e01b81526200030a918591600401918252602082015260400190565b8043106200094f576040516308b9266360e11b815260040160405180910390fd5b6200095c8c8c8c62001163565b6001546200096e6200056e8e620020a9565b146200098d5760405163723a340360e01b815260040160405180910390fd5b87861415806200099d5750878414155b15620009bc5760405163c1f97e3560e01b815260040160405180910390fd5b62000a3c8c8c8c7f00000000000000000000000000000000000000000000000000000000000000006f0e8e4c2dce6c2c6e8d2dedc84c2e8c6d60831b8e8e8e8e8e8e8e8e8e60405160200162000a1d9b9a999897969594939291906200245d565b6040516020818303038152906040528051906020012060065462000f2a565b6001600160a01b0382166000908152600260205260408120849055805b8981101562000b085762000ac889898381811062000a7b5762000a7b62002380565b905060200201602081019062000a92919062001990565b8c8c8481811062000aa75762000aa762002380565b90506020020135866001600160a01b0316620011bc9092919063ffffffff16565b86868281811062000add5762000add62002380565b905060200201358262000af19190620021c6565b91508062000aff8162002396565b91505062000a59565b5062000b1f6001600160a01b0384163383620011bc565b5060055462000b30906001620021c6565b60058190556040519081526001600160a01b0383169084907f02c7e81975f8edb86e2a0c038b7b86a49c744236abf0f6177ff5afc6986ab7089060200160405180910390a35050600160005550505050505050505050565b826040013584604001351162000bc1576040805163e0e8edf360e01b81528186013560048201529084013560248201526044016200030a565b62000bd46040840135620f4240620021c6565b8460400135111562000c09576040805163e0e8edf360e01b81528186013560048201529084013560248201526044016200030a565b62000c186020850185620024db565b905062000c268580620024db565b905014158062000c41575062000c3d8480620024db565b1590505b1562000c605760405163c01ba0ab60e01b815260040160405180910390fd5b62000c6d83838362001163565b6000805b62000c806020870187620024db565b905081101562000ce85762000c996020870187620024db565b8281811062000cac5762000cac62002380565b905060200201358262000cc09190620021c6565b915060065482111562000cd35762000ce8565b8062000cdf8162002396565b91505062000c71565b50600654811162000d1a5760065460405162bfb6ab60e01b81526200030a918391600401918252602082015260400190565b60015462000d2c6200056e86620020a9565b1462000d4b5760405163723a340360e01b815260040160405180910390fd5b600062000d5c6200056e87620020a9565b905062000d6f8585858460065462000f2a565b60018190556040860135600455600062000d9060a088016080890162001990565b6001600160a01b03161415801562000dab5750606086013515155b1562000dde5762000dde33606088013562000dcd60a08a0160808b0162001990565b6001600160a01b03169190620011bc565b60055462000dee906001620021c6565b60058190556040870135907f76d08978c024a4bf8cbb30c67fd78fcaa1827cbc533e4e175f36d07e64ccf96a90606089013562000e3260a08b0160808c0162001990565b62000e3e8b80620024db565b62000e4d60208e018e620024db565b60405162000e62979695949392919062002527565b60405180910390a2505050505050565b600030868686868660405162000e889062001835565b62000e9996959493929190620025a3565b604051809103906000f08015801562000eb6573d6000803e3d6000fd5b5060055490915062000eca90600162001239565b60058190556040516001600160a01b038316917f82fe3a4fa49c6382d0c085746698ddbbafe6c2bf61285b19410644b5b26287c79162000f18918c918c918c918c918c918c918c91620025f6565b60405180910390a25050505050505050565b6000805b62000f3a8780620024db565b9050811015620010505785858281811062000f595762000f5962002380565b62000f71926020606090920201908101915062002653565b60ff16156200103b5762000fd662000f8a8880620024db565b8381811062000f9d5762000f9d62002380565b905060200201602081019062000fb4919062001990565b8588888581811062000fca5762000fca62002380565b90506060020162001247565b62000ff457604051638baa579f60e01b815260040160405180910390fd5b620010036020880188620024db565b8281811062001016576200101662002380565b90506020020135826200102a9190620021c6565b9150828211156200103b5762001050565b80620010478162002396565b91505062000f2e565b508181116200107c5760405162bfb6ab60e01b815260048101829052602481018390526044016200030a565b505050505050565b6000806918da1958dadc1bda5b9d60b21b60001b90506000838286604001518760000151886020015189606001518a60800151604051602001620010cf979695949392919062002671565b60408051808303601f19018152919052805160209091012095945050505050565b6040516001600160a01b03808516602483015283166044820152606481018290526200115d9085906323b872dd60e01b906084015b60408051601f198184030181529190526020810180516001600160e01b03166001600160e01b031990931692909217909152620012dc565b50505050565b620011726020840184620024db565b9050620011808480620024db565b90501415806200119d575080620011988480620024db565b905014155b15620002df5760405163c6617b7b60e01b815260040160405180910390fd5b6040516001600160a01b038316602482015260448101829052620002df90849063a9059cbb60e01b9060640162001125565b60606200123283836040518060400160405280601e81526020017f416464726573733a206c6f772d6c6576656c2063616c6c206661696c65640000815250620013b5565b9392505050565b6000620012328284620021c6565b6040517f19457468657265756d205369676e6564204d6573736167653a0a3332000000006020820152603c81018390526000908190605c0160408051601f1981840301815291905280516020918201209150620012be908290620012ae9086018662002653565b85602001358660400135620013ce565b6001600160a01b0316856001600160a01b0316149150509392505050565b600062001333826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316620013b59092919063ffffffff16565b805190915015620002df5780806020019051810190620013549190620026d1565b620002df5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016200030a565b6060620013c68484600085620013fa565b949350505050565b6000806000620013e1878787876200152c565b91509150620013f08162001621565b5095945050505050565b6060824710156200145d5760405162461bcd60e51b815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f6044820152651c8818d85b1b60d21b60648201526084016200030a565b843b620014ad5760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016200030a565b600080866001600160a01b03168587604051620014cb9190620026f5565b60006040518083038185875af1925050503d80600081146200150a576040519150601f19603f3d011682016040523d82523d6000602084013e6200150f565b606091505b509150915062001521828286620017f7565b979650505050505050565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a083111562001565575060009050600362001618565b8460ff16601b141580156200157e57508460ff16601c14155b1562001591575060009050600462001618565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015620015e6573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116620016115760006001925092505062001618565b9150600090505b94509492505050565b600081600481111562001638576200163862002713565b1415620016425750565b600181600481111562001659576200165962002713565b1415620016a95760405162461bcd60e51b815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016200030a565b6002816004811115620016c057620016c062002713565b1415620017105760405162461bcd60e51b815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016200030a565b600381600481111562001727576200172762002713565b1415620017825760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604482015261756560f01b60648201526084016200030a565b600481600481111562001799576200179962002713565b1415620017f45760405162461bcd60e51b815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604482015261756560f01b60648201526084016200030a565b50565b606083156200180857508162001232565b825115620018195782518084602001fd5b8160405162461bcd60e51b81526004016200030a919062002729565b610e1d806200273f83390190565b600060a082840312156200185657600080fd5b50919050565b60008083601f8401126200186f57600080fd5b5081356001600160401b038111156200188757600080fd5b602083019150836020606083028501011115620018a357600080fd5b9250929050565b600080600080600060808688031215620018c357600080fd5b85356001600160401b0380821115620018db57600080fd5b620018e989838a0162001843565b965060208801359150808211156200190057600080fd5b506200190f888289016200185c565b9699909850959660408101359660609091013595509350505050565b600080604083850312156200193f57600080fd5b82356001600160401b038111156200195657600080fd5b620019648582860162001843565b95602094909401359450505050565b80356001600160a01b03811681146200198b57600080fd5b919050565b600060208284031215620019a357600080fd5b620012328262001973565b600080600060608486031215620019c457600080fd5b620019cf8462001973565b95602085013595506040909401359392505050565b634e487b7160e01b600052604160045260246000fd5b60405161012081016001600160401b038111828210171562001a205762001a20620019e4565b60405290565b604051601f8201601f191681016001600160401b038111828210171562001a515762001a51620019e4565b604052919050565b60006001600160401b0382111562001a755762001a75620019e4565b5060051b60200190565b600082601f83011262001a9157600080fd5b8135602062001aaa62001aa48362001a59565b62001a26565b82815260059290921b8401810191818101908684111562001aca57600080fd5b8286015b8481101562001ae7578035835291830191830162001ace565b509695505050505050565b600082601f83011262001b0457600080fd5b8135602062001b1762001aa48362001a59565b82815260059290921b8401810191818101908684111562001b3757600080fd5b8286015b8481101562001ae75762001b4f8162001973565b835291830191830162001b3b565b600082601f83011262001b6f57600080fd5b81356001600160401b0381111562001b8b5762001b8b620019e4565b62001ba0601f8201601f191660200162001a26565b81815284602083860101111562001bb657600080fd5b816020850160208301376000918101602001919091529392505050565b6000806000806060858703121562001bea57600080fd5b84356001600160401b038082111562001c0257600080fd5b62001c108883890162001843565b9550602087013591508082111562001c2757600080fd5b62001c35888389016200185c565b9095509350604087013591508082111562001c4f57600080fd5b90860190610120828903121562001c6557600080fd5b62001c6f620019fa565b82358281111562001c7f57600080fd5b62001c8d8a82860162001a7f565b82525060208301358281111562001ca357600080fd5b62001cb18a82860162001af2565b60208301525060408301358281111562001cca57600080fd5b62001cd88a82860162001a7f565b60408301525060608301358281111562001cf157600080fd5b62001cff8a82860162001af2565b60608301525062001d136080840162001973565b608082015260a08301358281111562001d2b57600080fd5b62001d398a82860162001b5d565b60a08301525060c083013560c082015260e083013560e08201526101009150818301358282015280935050505092959194509250565b60006020828403121562001d8257600080fd5b5035919050565b60008083601f84011262001d9c57600080fd5b5081356001600160401b0381111562001db457600080fd5b6020830191508360208260051b8501011115620018a357600080fd5b6000806000806000806000806000806000806101008d8f03121562001df457600080fd5b6001600160401b038d35111562001e0a57600080fd5b62001e198e8e358f0162001843565b9b506001600160401b0360208e0135111562001e3457600080fd5b62001e468e60208f01358f016200185c565b909b5099506001600160401b0360408e0135111562001e6457600080fd5b62001e768e60408f01358f0162001d89565b90995097506001600160401b0360608e0135111562001e9457600080fd5b62001ea68e60608f01358f0162001d89565b90975095506001600160401b0360808e0135111562001ec457600080fd5b62001ed68e60808f01358f0162001d89565b909550935060a08d0135925062001ef060c08e0162001973565b915060e08d013590509295989b509295989b509295989b565b6000806000806060858703121562001f2057600080fd5b84356001600160401b038082111562001f3857600080fd5b62001f468883890162001843565b9550602087013591508082111562001f5d57600080fd5b62001f6b8883890162001843565b9450604087013591508082111562001f8257600080fd5b5062001f91878288016200185c565b95989497509550505050565b60008083601f84011262001fb057600080fd5b5081356001600160401b0381111562001fc857600080fd5b602083019150836020828501011115620018a357600080fd5b803560ff811681146200198b57600080fd5b60008060008060008060006080888a0312156200200f57600080fd5b87356001600160401b03808211156200202757600080fd5b620020358b838c0162001f9d565b909950975060208a01359150808211156200204f57600080fd5b6200205d8b838c0162001f9d565b909750955060408a01359150808211156200207757600080fd5b50620020868a828b0162001f9d565b90945092506200209b90506060890162001fe1565b905092959891949750929550565b600060a08236031215620020bc57600080fd5b60405160a081016001600160401b038282108183111715620020e257620020e2620019e4565b816040528435915080821115620020f857600080fd5b620021063683870162001af2565b835260208501359150808211156200211d57600080fd5b506200212c3682860162001a7f565b6020830152506040830135604082015260608301356060820152620021546080840162001973565b608082015292915050565b6020808252601f908201527f5265656e7472616e637947756172643a207265656e7472616e742063616c6c00604082015260600190565b600060208284031215620021a957600080fd5b5051919050565b634e487b7160e01b600052601160045260246000fd5b60008219821115620021dc57620021dc620021b0565b500190565b600082821015620021f657620021f6620021b0565b500390565b600081518084526020808501945080840160005b838110156200222d578151875295820195908201906001016200220f565b509495945050505050565b600081518084526020808501945080840160005b838110156200222d5781516001600160a01b0316875295820195908201906001016200224c565b60005b838110156200229057818101518382015260200162002276565b838111156200115d5750506000910152565b60008151808452620022bc81602086016020860162002273565b601f01601f19169290920160200192915050565b60006101608d83528c6020840152806040840152620022f28184018d620021fb565b9050828103606084015262002308818c62002238565b905082810360808401526200231e818b620021fb565b905082810360a084015262002334818a62002238565b6001600160a01b03891660c085015283810360e08501529050620023598188620022a2565b61010084019690965250506101208101929092526101409091015298975050505050505050565b634e487b7160e01b600052603260045260246000fd5b6000600019821415620023ad57620023ad620021b0565b5060010190565b848152836020820152608060408201526000620023d56080830185620022a2565b905082606083015295945050505050565b81835260006001600160fb1b038311156200240057600080fd5b8260051b8083602087013760009401602001938452509192915050565b8183526000602080850194508260005b858110156200222d576001600160a01b03620024498362001973565b16875295820195908201906001016200242d565b60006101008d83528c6020840152806040840152620024808184018c8e620023e6565b9050828103606084015262002497818a8c6200241d565b90508281036080840152620024ae81888a620023e6565b60a084019690965250506001600160a01b039290921660c083015260e09091015298975050505050505050565b6000808335601e19843603018112620024f357600080fd5b8301803591506001600160401b038211156200250e57600080fd5b6020019150600581901b3603821315620018a357600080fd5b87815286602082015260018060a01b038616604082015260a0606082015260006200255760a0830186886200241d565b82810360808401526200256c818587620023e6565b9a9950505050505050505050565b81835281816020850137506000828201602090810191909152601f909101601f19169091010190565b6001600160a01b0387168152608060208201819052600090620025ca90830187896200257a565b8281036040840152620025df8186886200257a565b91505060ff83166060830152979650505050505050565b60a0815260006200260c60a083018a8c6200257a565b82810360208401526200262181898b6200257a565b90508281036040840152620026388187896200257a565b60ff9590951660608401525050608001529695505050505050565b6000602082840312156200266657600080fd5b620012328262001fe1565b87815286602082015285604082015260e0606082015260006200269860e083018762002238565b8281036080840152620026ac8187620021fb565b60a084019590955250506001600160a01b039190911660c09091015295945050505050565b600060208284031215620026e457600080fd5b815180151581146200123257600080fd5b600082516200270981846020870162002273565b9190910192915050565b634e487b7160e01b600052602160045260246000fd5b602081526000620012326020830184620022a256fe60806040526000196005553480156200001757600080fd5b5060405162000e1d38038062000e1d8339810160408190526200003a9162000311565b825183908390620000539060039060208501906200019e565b508051620000699060049060208401906200019e565b5050600680546001600160a81b031916600160a01b60ff8516026001600160a01b031916176001600160a01b03871617905550600554620000ac908590620000b6565b5050505062000419565b6001600160a01b038216620001115760405162461bcd60e51b815260206004820152601f60248201527f45524332303a206d696e7420746f20746865207a65726f206164647265737300604482015260640160405180910390fd5b8060026000828254620001259190620003b5565b90915550506001600160a01b0382166000908152602081905260408120805483929062000154908490620003b5565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b828054620001ac90620003dc565b90600052602060002090601f016020900481019282620001d057600085556200021b565b82601f10620001eb57805160ff19168380011785556200021b565b828001600101855582156200021b579182015b828111156200021b578251825591602001919060010190620001fe565b50620002299291506200022d565b5090565b5b808211156200022957600081556001016200022e565b634e487b7160e01b600052604160045260246000fd5b600082601f8301126200026c57600080fd5b81516001600160401b038082111562000289576200028962000244565b604051601f8301601f19908116603f01168101908282118183101715620002b457620002b462000244565b81604052838152602092508683858801011115620002d157600080fd5b600091505b83821015620002f55785820183015181830184015290820190620002d6565b83821115620003075760008385830101525b9695505050505050565b600080600080608085870312156200032857600080fd5b84516001600160a01b03811681146200034057600080fd5b60208601519094506001600160401b03808211156200035e57600080fd5b6200036c888389016200025a565b945060408701519150808211156200038357600080fd5b5062000392878288016200025a565b925050606085015160ff81168114620003aa57600080fd5b939692955090935050565b60008219821115620003d757634e487b7160e01b600052601160045260246000fd5b500190565b600181811c90821680620003f157607f821691505b602082108114156200041357634e487b7160e01b600052602260045260246000fd5b50919050565b6109f480620004296000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80635fd130a91161008c578063a457c2d711610066578063a457c2d7146101b6578063a9059cbb146101c9578063cbf0a64e146101dc578063dd62ed3e1461020757600080fd5b80635fd130a91461017057806370a082311461018557806395d89b41146101ae57600080fd5b806306fdde03146100d4578063095ea7b3146100f257806318160ddd1461011557806323b872dd1461012b578063313ce5671461013e578063395093511461015d575b600080fd5b6100dc610240565b6040516100e99190610812565b60405180910390f35b610105610100366004610883565b6102d2565b60405190151581526020016100e9565b61011d6102e8565b6040519081526020016100e9565b6101056101393660046108ad565b610315565b600654600160a01b900460ff1660405160ff90911681526020016100e9565b61010561016b366004610883565b6103c4565b61018361017e3660046108e9565b610400565b005b61011d6101933660046108e9565b6001600160a01b031660009081526020819052604090205490565b6100dc61046a565b6101056101c4366004610883565b610479565b6101056101d7366004610883565b610512565b6006546101ef906001600160a01b031681565b6040516001600160a01b0390911681526020016100e9565b61011d61021536600461090b565b6001600160a01b03918216600090815260016020908152604080832093909416825291909152205490565b60606003805461024f9061093e565b80601f016020809104026020016040519081016040528092919081815260200182805461027b9061093e565b80156102c85780601f1061029d576101008083540402835291602001916102c8565b820191906000526020600020905b8154815290600101906020018083116102ab57829003601f168201915b5050505050905090565b60006102df33848461051f565b50600192915050565b6006546001600160a01b0316600090815260208190526040812054600554610310919061098f565b905090565b6000610322848484610643565b6001600160a01b0384166000908152600160209081526040808320338452909152902054828110156103ac5760405162461bcd60e51b815260206004820152602860248201527f45524332303a207472616e7366657220616d6f756e74206578636565647320616044820152676c6c6f77616e636560c01b60648201526084015b60405180910390fd5b6103b9853385840361051f565b506001949350505050565b3360008181526001602090815260408083206001600160a01b038716845290915281205490916102df9185906103fb9086906109a6565b61051f565b6006546001600160a01b031633146104485760405162461bcd60e51b815260206004820152600b60248201526a4e6f74206772617669747960a81b60448201526064016103a3565b600680546001600160a01b0319166001600160a01b0392909216919091179055565b60606004805461024f9061093e565b3360009081526001602090815260408083206001600160a01b0386168452909152812054828110156104fb5760405162461bcd60e51b815260206004820152602560248201527f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f77604482015264207a65726f60d81b60648201526084016103a3565b610508338585840361051f565b5060019392505050565b60006102df338484610643565b6001600160a01b0383166105815760405162461bcd60e51b8152602060048201526024808201527f45524332303a20617070726f76652066726f6d20746865207a65726f206164646044820152637265737360e01b60648201526084016103a3565b6001600160a01b0382166105e25760405162461bcd60e51b815260206004820152602260248201527f45524332303a20617070726f766520746f20746865207a65726f206164647265604482015261737360f01b60648201526084016103a3565b6001600160a01b0383811660008181526001602090815260408083209487168084529482529182902085905590518481527f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925910160405180910390a3505050565b6001600160a01b0383166106a75760405162461bcd60e51b815260206004820152602560248201527f45524332303a207472616e736665722066726f6d20746865207a65726f206164604482015264647265737360d81b60648201526084016103a3565b6001600160a01b0382166107095760405162461bcd60e51b815260206004820152602360248201527f45524332303a207472616e7366657220746f20746865207a65726f206164647260448201526265737360e81b60648201526084016103a3565b6001600160a01b038316600090815260208190526040902054818110156107815760405162461bcd60e51b815260206004820152602660248201527f45524332303a207472616e7366657220616d6f756e7420657863656564732062604482015265616c616e636560d01b60648201526084016103a3565b6001600160a01b038085166000908152602081905260408082208585039055918516815290812080548492906107b89084906109a6565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef8460405161080491815260200190565b60405180910390a350505050565b600060208083528351808285015260005b8181101561083f57858101830151858201604001528201610823565b81811115610851576000604083870101525b50601f01601f1916929092016040019392505050565b80356001600160a01b038116811461087e57600080fd5b919050565b6000806040838503121561089657600080fd5b61089f83610867565b946020939093013593505050565b6000806000606084860312156108c257600080fd5b6108cb84610867565b92506108d960208501610867565b9150604084013590509250925092565b6000602082840312156108fb57600080fd5b61090482610867565b9392505050565b6000806040838503121561091e57600080fd5b61092783610867565b915061093560208401610867565b90509250929050565b600181811c9082168061095257607f821691505b6020821081141561097357634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b6000828210156109a1576109a1610979565b500390565b600082198211156109b9576109b9610979565b50019056fea264697066735822122010756c7691f24fa76ef182760b9f99a094a8bf02470a418707f1de88a7d8d7b364736f6c634300080a0033a2646970667358221220e0cc737c893f4f801f620d372b01df0f8d6fe82c9585b72e83600f517f252d7a64736f6c634300080a0033\",\n  \"linkReferences\": {},\n  \"deployedLinkReferences\": {}\n}\n") . expect ("invalid abi")
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
        #[doc = "Calls the contract's `submitBatch` (0x8690ff98) function"]
        pub fn submit_batch(
            &self,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
            amounts: ::std::vec::Vec<ethers::core::types::U256>,
            destinations: ::std::vec::Vec<ethers::core::types::Address>,
            fees: ::std::vec::Vec<ethers::core::types::U256>,
            batch_nonce: ethers::core::types::U256,
            token_contract: ethers::core::types::Address,
            batch_timeout: ethers::core::types::U256,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
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
        #[doc = "Calls the contract's `submitLogicCall` (0x6941db93) function"]
        pub fn submit_logic_call(
            &self,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
            args: LogicCallArgs,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([105, 65, 219, 147], (current_valset, sigs, args))
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `testCheckValidatorSignatures` (0x00901153) function"]
        pub fn test_check_validator_signatures(
            &self,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
            the_hash: [u8; 32],
            power_threshold: ethers::core::types::U256,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash(
                    [0, 144, 17, 83],
                    (current_valset, sigs, the_hash, power_threshold),
                )
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `testMakeCheckpoint` (0x01031525) function"]
        pub fn test_make_checkpoint(
            &self,
            valset_args: ValsetArgs,
            gravity_id: [u8; 32],
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([1, 3, 21, 37], (valset_args, gravity_id))
                .expect("method not found (this should never happen)")
        }
        #[doc = "Calls the contract's `updateValset` (0xaca6b1c1) function"]
        pub fn update_valset(
            &self,
            new_valset: ValsetArgs,
            current_valset: ValsetArgs,
            sigs: ::std::vec::Vec<ValSignature>,
        ) -> ethers::contract::builders::ContractCall<M, ()> {
            self.0
                .method_hash([172, 166, 177, 193], (new_valset, current_valset, sigs))
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
        ethers :: contract :: EthDisplay,
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
        ethers :: contract :: EthDisplay,
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
        pub return_data: ethers::core::types::Bytes,
        pub event_nonce: ethers::core::types::U256,
    }
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthEvent,
        ethers :: contract :: EthDisplay,
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
        ethers :: contract :: EthDisplay,
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
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethevent(
        name = "ValsetUpdatedEvent",
        abi = "ValsetUpdatedEvent(uint256,uint256,uint256,address,address[],uint256[])"
    )]
    pub struct ValsetUpdatedEventFilter {
        #[ethevent(indexed)]
        pub new_valset_nonce: ethers::core::types::U256,
        pub event_nonce: ethers::core::types::U256,
        pub reward_amount: ethers::core::types::U256,
        pub reward_token: ethers::core::types::Address,
        pub validators: Vec<ethers::core::types::Address>,
        pub powers: Vec<ethers::core::types::U256>,
    }
    #[derive(Debug, Clone, PartialEq, Eq, ethers :: contract :: EthAbiType)]
    pub enum GravityEvents {
        Erc20DeployedEventFilter(Erc20DeployedEventFilter),
        LogicCallEventFilter(LogicCallEventFilter),
        SendToCosmosEventFilter(SendToCosmosEventFilter),
        TransactionBatchExecutedEventFilter(TransactionBatchExecutedEventFilter),
        ValsetUpdatedEventFilter(ValsetUpdatedEventFilter),
    }
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
    impl ::std::fmt::Display for GravityEvents {
        fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
            match self {
                GravityEvents::Erc20DeployedEventFilter(element) => element.fmt(f),
                GravityEvents::LogicCallEventFilter(element) => element.fmt(f),
                GravityEvents::SendToCosmosEventFilter(element) => element.fmt(f),
                GravityEvents::TransactionBatchExecutedEventFilter(element) => element.fmt(f),
                GravityEvents::ValsetUpdatedEventFilter(element) => element.fmt(f),
            }
        }
    }
    #[doc = "Container type for all input parameters for the `deployERC20`function with signature `deployERC20(string,string,string,uint8)` and selector `[247, 149, 86, 55]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "deployERC20", abi = "deployERC20(string,string,string,uint8)")]
    pub struct DeployERC20Call {
        pub cosmos_denom: String,
        pub name: String,
        pub symbol: String,
        pub decimals: u8,
    }
    #[doc = "Container type for all input parameters for the `lastBatchNonce`function with signature `lastBatchNonce(address)` and selector `[1, 27, 33, 116]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "lastBatchNonce", abi = "lastBatchNonce(address)")]
    pub struct LastBatchNonceCall {
        pub erc_20_address: ethers::core::types::Address,
    }
    #[doc = "Container type for all input parameters for the `lastLogicCallNonce`function with signature `lastLogicCallNonce(bytes32)` and selector `[201, 209, 148, 213]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "lastLogicCallNonce", abi = "lastLogicCallNonce(bytes32)")]
    pub struct LastLogicCallNonceCall {
        pub invalidation_id: [u8; 32],
    }
    #[doc = "Container type for all input parameters for the `sendToCosmos`function with signature `sendToCosmos(address,bytes32,uint256)` and selector `[31, 251, 231, 249]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "sendToCosmos", abi = "sendToCosmos(address,bytes32,uint256)")]
    pub struct SendToCosmosCall {
        pub token_contract: ethers::core::types::Address,
        pub destination: [u8; 32],
        pub amount: ethers::core::types::U256,
    }
    #[doc = "Container type for all input parameters for the `state_gravityId`function with signature `state_gravityId()` and selector `[189, 218, 129, 212]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "state_gravityId", abi = "state_gravityId()")]
    pub struct StateGravityIdCall;
    #[doc = "Container type for all input parameters for the `state_invalidationMapping`function with signature `state_invalidationMapping(bytes32)` and selector `[125, 251, 111, 134]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(
        name = "state_invalidationMapping",
        abi = "state_invalidationMapping(bytes32)"
    )]
    pub struct StateInvalidationMappingCall(pub [u8; 32]);
    #[doc = "Container type for all input parameters for the `state_lastBatchNonces`function with signature `state_lastBatchNonces(address)` and selector `[223, 151, 23, 75]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "state_lastBatchNonces", abi = "state_lastBatchNonces(address)")]
    pub struct StateLastBatchNoncesCall(pub ethers::core::types::Address);
    #[doc = "Container type for all input parameters for the `state_lastEventNonce`function with signature `state_lastEventNonce()` and selector `[115, 178, 5, 71]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "state_lastEventNonce", abi = "state_lastEventNonce()")]
    pub struct StateLastEventNonceCall;
    #[doc = "Container type for all input parameters for the `state_lastValsetCheckpoint`function with signature `state_lastValsetCheckpoint()` and selector `[242, 181, 51, 7]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(
        name = "state_lastValsetCheckpoint",
        abi = "state_lastValsetCheckpoint()"
    )]
    pub struct StateLastValsetCheckpointCall;
    #[doc = "Container type for all input parameters for the `state_lastValsetNonce`function with signature `state_lastValsetNonce()` and selector `[181, 101, 97, 254]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "state_lastValsetNonce", abi = "state_lastValsetNonce()")]
    pub struct StateLastValsetNonceCall;
    #[doc = "Container type for all input parameters for the `state_powerThreshold`function with signature `state_powerThreshold()` and selector `[229, 162, 181, 210]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(name = "state_powerThreshold", abi = "state_powerThreshold()")]
    pub struct StatePowerThresholdCall;
    #[doc = "Container type for all input parameters for the `submitBatch`function with signature `submitBatch((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],uint256[],address[],uint256[],uint256,address,uint256)` and selector `[134, 144, 255, 152]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(
        name = "submitBatch",
        abi = "submitBatch((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],uint256[],address[],uint256[],uint256,address,uint256)"
    )]
    pub struct SubmitBatchCall {
        pub current_valset: ValsetArgs,
        pub sigs: ::std::vec::Vec<ValSignature>,
        pub amounts: ::std::vec::Vec<ethers::core::types::U256>,
        pub destinations: ::std::vec::Vec<ethers::core::types::Address>,
        pub fees: ::std::vec::Vec<ethers::core::types::U256>,
        pub batch_nonce: ethers::core::types::U256,
        pub token_contract: ethers::core::types::Address,
        pub batch_timeout: ethers::core::types::U256,
    }
    #[doc = "Container type for all input parameters for the `submitLogicCall`function with signature `submitLogicCall((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],(uint256[],address[],uint256[],address[],address,bytes,uint256,bytes32,uint256))` and selector `[105, 65, 219, 147]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
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
    #[doc = "Container type for all input parameters for the `testCheckValidatorSignatures`function with signature `testCheckValidatorSignatures((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],bytes32,uint256)` and selector `[0, 144, 17, 83]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(
        name = "testCheckValidatorSignatures",
        abi = "testCheckValidatorSignatures((address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[],bytes32,uint256)"
    )]
    pub struct TestCheckValidatorSignaturesCall {
        pub current_valset: ValsetArgs,
        pub sigs: ::std::vec::Vec<ValSignature>,
        pub the_hash: [u8; 32],
        pub power_threshold: ethers::core::types::U256,
    }
    #[doc = "Container type for all input parameters for the `testMakeCheckpoint`function with signature `testMakeCheckpoint((address[],uint256[],uint256,uint256,address),bytes32)` and selector `[1, 3, 21, 37]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
    )]
    #[ethcall(
        name = "testMakeCheckpoint",
        abi = "testMakeCheckpoint((address[],uint256[],uint256,uint256,address),bytes32)"
    )]
    pub struct TestMakeCheckpointCall {
        pub valset_args: ValsetArgs,
        pub gravity_id: [u8; 32],
    }
    #[doc = "Container type for all input parameters for the `updateValset`function with signature `updateValset((address[],uint256[],uint256,uint256,address),(address[],uint256[],uint256,uint256,address),(uint8,bytes32,bytes32)[])` and selector `[172, 166, 177, 193]`"]
    #[derive(
        Clone,
        Debug,
        Default,
        Eq,
        PartialEq,
        ethers :: contract :: EthCall,
        ethers :: contract :: EthDisplay,
        serde :: Deserialize,
        serde :: Serialize,
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
    #[derive(Debug, Clone, PartialEq, Eq, ethers :: contract :: EthAbiType)]
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
    impl ethers::core::abi::AbiDecode for GravityCalls {
        fn decode(data: impl AsRef<[u8]>) -> Result<Self, ethers::core::abi::AbiError> {
            if let Ok(decoded) =
                <DeployERC20Call as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::DeployERC20(decoded));
            }
            if let Ok(decoded) =
                <LastBatchNonceCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::LastBatchNonce(decoded));
            }
            if let Ok(decoded) =
                <LastLogicCallNonceCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::LastLogicCallNonce(decoded));
            }
            if let Ok(decoded) =
                <SendToCosmosCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::SendToCosmos(decoded));
            }
            if let Ok(decoded) =
                <StateGravityIdCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::StateGravityId(decoded));
            }
            if let Ok(decoded) =
                <StateInvalidationMappingCall as ethers::core::abi::AbiDecode>::decode(
                    data.as_ref(),
                )
            {
                return Ok(GravityCalls::StateInvalidationMapping(decoded));
            }
            if let Ok(decoded) =
                <StateLastBatchNoncesCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::StateLastBatchNonces(decoded));
            }
            if let Ok(decoded) =
                <StateLastEventNonceCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::StateLastEventNonce(decoded));
            }
            if let Ok(decoded) =
                <StateLastValsetCheckpointCall as ethers::core::abi::AbiDecode>::decode(
                    data.as_ref(),
                )
            {
                return Ok(GravityCalls::StateLastValsetCheckpoint(decoded));
            }
            if let Ok(decoded) =
                <StateLastValsetNonceCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::StateLastValsetNonce(decoded));
            }
            if let Ok(decoded) =
                <StatePowerThresholdCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::StatePowerThreshold(decoded));
            }
            if let Ok(decoded) =
                <SubmitBatchCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::SubmitBatch(decoded));
            }
            if let Ok(decoded) =
                <SubmitLogicCallCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::SubmitLogicCall(decoded));
            }
            if let Ok(decoded) =
                <TestCheckValidatorSignaturesCall as ethers::core::abi::AbiDecode>::decode(
                    data.as_ref(),
                )
            {
                return Ok(GravityCalls::TestCheckValidatorSignatures(decoded));
            }
            if let Ok(decoded) =
                <TestMakeCheckpointCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::TestMakeCheckpoint(decoded));
            }
            if let Ok(decoded) =
                <UpdateValsetCall as ethers::core::abi::AbiDecode>::decode(data.as_ref())
            {
                return Ok(GravityCalls::UpdateValset(decoded));
            }
            Err(ethers::core::abi::Error::InvalidData.into())
        }
    }
    impl ethers::core::abi::AbiEncode for GravityCalls {
        fn encode(self) -> Vec<u8> {
            match self {
                GravityCalls::DeployERC20(element) => element.encode(),
                GravityCalls::LastBatchNonce(element) => element.encode(),
                GravityCalls::LastLogicCallNonce(element) => element.encode(),
                GravityCalls::SendToCosmos(element) => element.encode(),
                GravityCalls::StateGravityId(element) => element.encode(),
                GravityCalls::StateInvalidationMapping(element) => element.encode(),
                GravityCalls::StateLastBatchNonces(element) => element.encode(),
                GravityCalls::StateLastEventNonce(element) => element.encode(),
                GravityCalls::StateLastValsetCheckpoint(element) => element.encode(),
                GravityCalls::StateLastValsetNonce(element) => element.encode(),
                GravityCalls::StatePowerThreshold(element) => element.encode(),
                GravityCalls::SubmitBatch(element) => element.encode(),
                GravityCalls::SubmitLogicCall(element) => element.encode(),
                GravityCalls::TestCheckValidatorSignatures(element) => element.encode(),
                GravityCalls::TestMakeCheckpoint(element) => element.encode(),
                GravityCalls::UpdateValset(element) => element.encode(),
            }
        }
    }
    impl ::std::fmt::Display for GravityCalls {
        fn fmt(&self, f: &mut ::std::fmt::Formatter<'_>) -> ::std::fmt::Result {
            match self {
                GravityCalls::DeployERC20(element) => element.fmt(f),
                GravityCalls::LastBatchNonce(element) => element.fmt(f),
                GravityCalls::LastLogicCallNonce(element) => element.fmt(f),
                GravityCalls::SendToCosmos(element) => element.fmt(f),
                GravityCalls::StateGravityId(element) => element.fmt(f),
                GravityCalls::StateInvalidationMapping(element) => element.fmt(f),
                GravityCalls::StateLastBatchNonces(element) => element.fmt(f),
                GravityCalls::StateLastEventNonce(element) => element.fmt(f),
                GravityCalls::StateLastValsetCheckpoint(element) => element.fmt(f),
                GravityCalls::StateLastValsetNonce(element) => element.fmt(f),
                GravityCalls::StatePowerThreshold(element) => element.fmt(f),
                GravityCalls::SubmitBatch(element) => element.fmt(f),
                GravityCalls::SubmitLogicCall(element) => element.fmt(f),
                GravityCalls::TestCheckValidatorSignatures(element) => element.fmt(f),
                GravityCalls::TestMakeCheckpoint(element) => element.fmt(f),
                GravityCalls::UpdateValset(element) => element.fmt(f),
            }
        }
    }
    impl ::std::convert::From<DeployERC20Call> for GravityCalls {
        fn from(var: DeployERC20Call) -> Self {
            GravityCalls::DeployERC20(var)
        }
    }
    impl ::std::convert::From<LastBatchNonceCall> for GravityCalls {
        fn from(var: LastBatchNonceCall) -> Self {
            GravityCalls::LastBatchNonce(var)
        }
    }
    impl ::std::convert::From<LastLogicCallNonceCall> for GravityCalls {
        fn from(var: LastLogicCallNonceCall) -> Self {
            GravityCalls::LastLogicCallNonce(var)
        }
    }
    impl ::std::convert::From<SendToCosmosCall> for GravityCalls {
        fn from(var: SendToCosmosCall) -> Self {
            GravityCalls::SendToCosmos(var)
        }
    }
    impl ::std::convert::From<StateGravityIdCall> for GravityCalls {
        fn from(var: StateGravityIdCall) -> Self {
            GravityCalls::StateGravityId(var)
        }
    }
    impl ::std::convert::From<StateInvalidationMappingCall> for GravityCalls {
        fn from(var: StateInvalidationMappingCall) -> Self {
            GravityCalls::StateInvalidationMapping(var)
        }
    }
    impl ::std::convert::From<StateLastBatchNoncesCall> for GravityCalls {
        fn from(var: StateLastBatchNoncesCall) -> Self {
            GravityCalls::StateLastBatchNonces(var)
        }
    }
    impl ::std::convert::From<StateLastEventNonceCall> for GravityCalls {
        fn from(var: StateLastEventNonceCall) -> Self {
            GravityCalls::StateLastEventNonce(var)
        }
    }
    impl ::std::convert::From<StateLastValsetCheckpointCall> for GravityCalls {
        fn from(var: StateLastValsetCheckpointCall) -> Self {
            GravityCalls::StateLastValsetCheckpoint(var)
        }
    }
    impl ::std::convert::From<StateLastValsetNonceCall> for GravityCalls {
        fn from(var: StateLastValsetNonceCall) -> Self {
            GravityCalls::StateLastValsetNonce(var)
        }
    }
    impl ::std::convert::From<StatePowerThresholdCall> for GravityCalls {
        fn from(var: StatePowerThresholdCall) -> Self {
            GravityCalls::StatePowerThreshold(var)
        }
    }
    impl ::std::convert::From<SubmitBatchCall> for GravityCalls {
        fn from(var: SubmitBatchCall) -> Self {
            GravityCalls::SubmitBatch(var)
        }
    }
    impl ::std::convert::From<SubmitLogicCallCall> for GravityCalls {
        fn from(var: SubmitLogicCallCall) -> Self {
            GravityCalls::SubmitLogicCall(var)
        }
    }
    impl ::std::convert::From<TestCheckValidatorSignaturesCall> for GravityCalls {
        fn from(var: TestCheckValidatorSignaturesCall) -> Self {
            GravityCalls::TestCheckValidatorSignatures(var)
        }
    }
    impl ::std::convert::From<TestMakeCheckpointCall> for GravityCalls {
        fn from(var: TestMakeCheckpointCall) -> Self {
            GravityCalls::TestMakeCheckpoint(var)
        }
    }
    impl ::std::convert::From<UpdateValsetCall> for GravityCalls {
        fn from(var: UpdateValsetCall) -> Self {
            GravityCalls::UpdateValset(var)
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
        pub payload: ethers::core::types::Bytes,
        pub time_out: ethers::core::types::U256,
        pub invalidation_id: [u8; 32],
        pub invalidation_nonce: ethers::core::types::U256,
    }
    #[doc = "`ValSignature(uint8,bytes32,bytes32)`"]
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
    pub struct ValSignature {
        pub v: u8,
        pub r: [u8; 32],
        pub s: [u8; 32],
    }
    #[doc = "`ValsetArgs(address[],uint256[],uint256,uint256,address)`"]
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
    pub struct ValsetArgs {
        pub validators: Vec<ethers::core::types::Address>,
        pub powers: Vec<ethers::core::types::U256>,
        pub valset_nonce: ethers::core::types::U256,
        pub reward_amount: ethers::core::types::U256,
        pub reward_token: ethers::core::types::Address,
    }
}
