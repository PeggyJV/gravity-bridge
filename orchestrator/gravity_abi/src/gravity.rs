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
            serde_json :: from_str ("[\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_gravityId\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_powerThreshold\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_validators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_powers\",\n        \"type\": \"uint256[]\"\n      }\n    ],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"constructor\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"BatchTimedOut\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"IncorrectCheckpoint\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"cumulativePower\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"powerThreshold\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"InsufficientPower\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"newNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"currentNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"InvalidBatchNonce\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"InvalidLogicCallFees\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"newNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"currentNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"InvalidLogicCallNonce\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"InvalidLogicCallTransfers\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"InvalidSendToCosmos\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"InvalidSignature\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"newNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"currentNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"InvalidValsetNonce\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"LogicCallTimedOut\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"MalformedBatch\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"MalformedCurrentValidatorSet\",\n    \"type\": \"error\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"MalformedNewValidatorSet\",\n    \"type\": \"error\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"_cosmosDenom\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"_name\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"_symbol\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint8\",\n        \"name\": \"_decimals\",\n        \"type\": \"uint8\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"ERC20DeployedEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes32\",\n        \"name\": \"_invalidationId\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_invalidationNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"bytes\",\n        \"name\": \"_returnData\",\n        \"type\": \"bytes\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"LogicCallEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_sender\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"bytes32\",\n        \"name\": \"_destination\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_amount\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"SendToCosmosEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"uint256\",\n        \"name\": \"_batchNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": true,\n        \"internalType\": \"address\",\n        \"name\": \"_token\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"TransactionBatchExecutedEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": true,\n        \"internalType\": \"uint256\",\n        \"name\": \"_newValsetNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_eventNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"_rewardAmount\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"address\",\n        \"name\": \"_rewardToken\",\n        \"type\": \"address\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"address[]\",\n        \"name\": \"_validators\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_powers\",\n        \"type\": \"uint256[]\"\n      }\n    ],\n    \"name\": \"ValsetUpdatedEvent\",\n    \"type\": \"event\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_cosmosDenom\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_name\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"string\",\n        \"name\": \"_symbol\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"uint8\",\n        \"name\": \"_decimals\",\n        \"type\": \"uint8\"\n      }\n    ],\n    \"name\": \"deployERC20\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address\",\n        \"name\": \"_erc20Address\",\n        \"type\": \"address\"\n      }\n    ],\n    \"name\": \"lastBatchNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_invalidation_id\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"lastLogicCallNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_destination\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"sendToCosmos\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_gravityId\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"state_invalidationMapping\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"address\",\n        \"name\": \"\",\n        \"type\": \"address\"\n      }\n    ],\n    \"name\": \"state_lastBatchNonces\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_lastEventNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_lastValsetCheckpoint\",\n    \"outputs\": [\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_lastValsetNonce\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [],\n    \"name\": \"state_powerThreshold\",\n    \"outputs\": [\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"stateMutability\": \"view\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"components\": [\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"validators\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"powers\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"valsetNonce\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"rewardAmount\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"rewardToken\",\n            \"type\": \"address\"\n          }\n        ],\n        \"internalType\": \"struct ValsetArgs\",\n        \"name\": \"_currentValset\",\n        \"type\": \"tuple\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"uint8\",\n            \"name\": \"v\",\n            \"type\": \"uint8\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"r\",\n            \"type\": \"bytes32\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"s\",\n            \"type\": \"bytes32\"\n          }\n        ],\n        \"internalType\": \"struct ValSignature[]\",\n        \"name\": \"_sigs\",\n        \"type\": \"tuple[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_amounts\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"address[]\",\n        \"name\": \"_destinations\",\n        \"type\": \"address[]\"\n      },\n      {\n        \"internalType\": \"uint256[]\",\n        \"name\": \"_fees\",\n        \"type\": \"uint256[]\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_batchNonce\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"internalType\": \"address\",\n        \"name\": \"_tokenContract\",\n        \"type\": \"address\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_batchTimeout\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"submitBatch\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"components\": [\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"validators\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"powers\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"valsetNonce\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"rewardAmount\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"rewardToken\",\n            \"type\": \"address\"\n          }\n        ],\n        \"internalType\": \"struct ValsetArgs\",\n        \"name\": \"_currentValset\",\n        \"type\": \"tuple\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"uint8\",\n            \"name\": \"v\",\n            \"type\": \"uint8\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"r\",\n            \"type\": \"bytes32\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"s\",\n            \"type\": \"bytes32\"\n          }\n        ],\n        \"internalType\": \"struct ValSignature[]\",\n        \"name\": \"_sigs\",\n        \"type\": \"tuple[]\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"transferAmounts\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"transferTokenContracts\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"feeAmounts\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"feeTokenContracts\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"logicContractAddress\",\n            \"type\": \"address\"\n          },\n          {\n            \"internalType\": \"bytes\",\n            \"name\": \"payload\",\n            \"type\": \"bytes\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"timeOut\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"invalidationId\",\n            \"type\": \"bytes32\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"invalidationNonce\",\n            \"type\": \"uint256\"\n          }\n        ],\n        \"internalType\": \"struct LogicCallArgs\",\n        \"name\": \"_args\",\n        \"type\": \"tuple\"\n      }\n    ],\n    \"name\": \"submitLogicCall\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"components\": [\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"validators\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"powers\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"valsetNonce\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"rewardAmount\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"rewardToken\",\n            \"type\": \"address\"\n          }\n        ],\n        \"internalType\": \"struct ValsetArgs\",\n        \"name\": \"_currentValset\",\n        \"type\": \"tuple\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"uint8\",\n            \"name\": \"v\",\n            \"type\": \"uint8\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"r\",\n            \"type\": \"bytes32\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"s\",\n            \"type\": \"bytes32\"\n          }\n        ],\n        \"internalType\": \"struct ValSignature[]\",\n        \"name\": \"_sigs\",\n        \"type\": \"tuple[]\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_theHash\",\n        \"type\": \"bytes32\"\n      },\n      {\n        \"internalType\": \"uint256\",\n        \"name\": \"_powerThreshold\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"testCheckValidatorSignatures\",\n    \"outputs\": [],\n    \"stateMutability\": \"pure\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"components\": [\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"validators\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"powers\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"valsetNonce\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"rewardAmount\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"rewardToken\",\n            \"type\": \"address\"\n          }\n        ],\n        \"internalType\": \"struct ValsetArgs\",\n        \"name\": \"_valsetArgs\",\n        \"type\": \"tuple\"\n      },\n      {\n        \"internalType\": \"bytes32\",\n        \"name\": \"_gravityId\",\n        \"type\": \"bytes32\"\n      }\n    ],\n    \"name\": \"testMakeCheckpoint\",\n    \"outputs\": [],\n    \"stateMutability\": \"pure\",\n    \"type\": \"function\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"components\": [\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"validators\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"powers\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"valsetNonce\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"rewardAmount\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"rewardToken\",\n            \"type\": \"address\"\n          }\n        ],\n        \"internalType\": \"struct ValsetArgs\",\n        \"name\": \"_newValset\",\n        \"type\": \"tuple\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"address[]\",\n            \"name\": \"validators\",\n            \"type\": \"address[]\"\n          },\n          {\n            \"internalType\": \"uint256[]\",\n            \"name\": \"powers\",\n            \"type\": \"uint256[]\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"valsetNonce\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"uint256\",\n            \"name\": \"rewardAmount\",\n            \"type\": \"uint256\"\n          },\n          {\n            \"internalType\": \"address\",\n            \"name\": \"rewardToken\",\n            \"type\": \"address\"\n          }\n        ],\n        \"internalType\": \"struct ValsetArgs\",\n        \"name\": \"_currentValset\",\n        \"type\": \"tuple\"\n      },\n      {\n        \"components\": [\n          {\n            \"internalType\": \"uint8\",\n            \"name\": \"v\",\n            \"type\": \"uint8\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"r\",\n            \"type\": \"bytes32\"\n          },\n          {\n            \"internalType\": \"bytes32\",\n            \"name\": \"s\",\n            \"type\": \"bytes32\"\n          }\n        ],\n        \"internalType\": \"struct ValSignature[]\",\n        \"name\": \"_sigs\",\n        \"type\": \"tuple[]\"\n      }\n    ],\n    \"name\": \"updateValset\",\n    \"outputs\": [],\n    \"stateMutability\": \"nonpayable\",\n    \"type\": \"function\"\n  }\n]\n") . expect ("invalid abi")
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
