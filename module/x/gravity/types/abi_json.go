package types

// The go-ethereum ABI encoder *only* encodes function calls and then it only encodes
// function calls for which you provide an ABI json just like you would get out of the
// solidity compiler with your compiled contract.
// You are supposed to compile your contract, use abigen to generate an ABI , import
// this generated go module and then use for that for all testing and development.
// This abstraction layer is more trouble than it's worth, because we don't want to
// encode a function call at all, but instead we want to emulate a Solidity encode operation
// which has no equal available from go-ethereum.
//
// In order to work around this absurd series of problems we have to manually write the below
// 'function specification' that will encode the same arguments into a function call. We can then
// truncate the first several bytes where the call name is encoded to finally get the equal of the

const (
	// BatchTxCheckpointABIJSON checks the ETH ABI for compatability of the OutgoingBatchTx message
	BatchTxCheckpointABIJSON = `[{
		"inputs": [
			{
			  "components": [
				{
				  "internalType": "address[]",
				  "name": "validators",
				  "type": "address[]"
				},
				{
				  "internalType": "uint256[]",
				  "name": "powers",
				  "type": "uint256[]"
				},
				{
				  "internalType": "uint256",
				  "name": "valsetNonce",
				  "type": "uint256"
				},
				{
				  "internalType": "uint256",
				  "name": "rewardAmount",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "rewardToken",
				  "type": "address"
				}
			  ],
			  "internalType": "struct ValsetArgs",
			  "name": "_currentValset",
			  "type": "tuple"
			},
			{
			  "components": [
				{
				  "internalType": "uint8",
				  "name": "v",
				  "type": "uint8"
				},
				{
				  "internalType": "bytes32",
				  "name": "r",
				  "type": "bytes32"
				},
				{
				  "internalType": "bytes32",
				  "name": "s",
				  "type": "bytes32"
				}
			  ],
			  "internalType": "struct Signature[]",
			  "name": "_sigs",
			  "type": "tuple[]"
			},
			{
			  "internalType": "uint256[]",
			  "name": "_amounts",
			  "type": "uint256[]"
			},
			{
			  "internalType": "address[]",
			  "name": "_destinations",
			  "type": "address[]"
			},
			{
			  "internalType": "uint256[]",
			  "name": "_fees",
			  "type": "uint256[]"
			},
			{
			  "internalType": "uint256",
			  "name": "_batchNonce",
			  "type": "uint256"
			},
			{
			  "internalType": "address",
			  "name": "_tokenContract",
			  "type": "address"
			},
			{
			  "internalType": "uint256",
			  "name": "_batchTimeout",
			  "type": "uint256"
			}
		  ],
		  "name": "submitBatch",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
	}]`

	// SignerSetTxCheckpointABIJSON checks the ETH ABI for compatability of the signer set update message
	SignerSetTxCheckpointABIJSON = `[{
		"inputs": [
			{
			  "components": [
				{
				  "internalType": "address[]",
				  "name": "validators",
				  "type": "address[]"
				},
				{
				  "internalType": "uint256[]",
				  "name": "powers",
				  "type": "uint256[]"
				},
				{
				  "internalType": "uint256",
				  "name": "valsetNonce",
				  "type": "uint256"
				},
				{
				  "internalType": "uint256",
				  "name": "rewardAmount",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "rewardToken",
				  "type": "address"
				}
			  ],
			  "internalType": "struct ValsetArgs",
			  "name": "_newValset",
			  "type": "tuple"
			},
			{
			  "components": [
				{
				  "internalType": "address[]",
				  "name": "validators",
				  "type": "address[]"
				},
				{
				  "internalType": "uint256[]",
				  "name": "powers",
				  "type": "uint256[]"
				},
				{
				  "internalType": "uint256",
				  "name": "valsetNonce",
				  "type": "uint256"
				},
				{
				  "internalType": "uint256",
				  "name": "rewardAmount",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "rewardToken",
				  "type": "address"
				}
			  ],
			  "internalType": "struct ValsetArgs",
			  "name": "_currentValset",
			  "type": "tuple"
			},
			{
			  "components": [
				{
				  "internalType": "uint8",
				  "name": "v",
				  "type": "uint8"
				},
				{
				  "internalType": "bytes32",
				  "name": "r",
				  "type": "bytes32"
				},
				{
				  "internalType": "bytes32",
				  "name": "s",
				  "type": "bytes32"
				}
			  ],
			  "internalType": "struct Signature[]",
			  "name": "_sigs",
			  "type": "tuple[]"
			}
		  ],
		  "name": "updateValset",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
		
	}]`

	// ContractCallTxABIJSON checks the ETH ABI for compatability of the logic call message
	ContractCallTxABIJSON = `[{
		"inputs": [
			{
			  "components": [
				{
				  "internalType": "address[]",
				  "name": "validators",
				  "type": "address[]"
				},
				{
				  "internalType": "uint256[]",
				  "name": "powers",
				  "type": "uint256[]"
				},
				{
				  "internalType": "uint256",
				  "name": "valsetNonce",
				  "type": "uint256"
				},
				{
				  "internalType": "uint256",
				  "name": "rewardAmount",
				  "type": "uint256"
				},
				{
				  "internalType": "address",
				  "name": "rewardToken",
				  "type": "address"
				}
			  ],
			  "internalType": "struct ValsetArgs",
			  "name": "_currentValset",
			  "type": "tuple"
			},
			{
			  "components": [
				{
				  "internalType": "uint8",
				  "name": "v",
				  "type": "uint8"
				},
				{
				  "internalType": "bytes32",
				  "name": "r",
				  "type": "bytes32"
				},
				{
				  "internalType": "bytes32",
				  "name": "s",
				  "type": "bytes32"
				}
			  ],
			  "internalType": "struct Signature[]",
			  "name": "_sigs",
			  "type": "tuple[]"
			},
			{
			  "components": [
				{
				  "internalType": "uint256[]",
				  "name": "transferAmounts",
				  "type": "uint256[]"
				},
				{
				  "internalType": "address[]",
				  "name": "transferTokenContracts",
				  "type": "address[]"
				},
				{
				  "internalType": "uint256[]",
				  "name": "feeAmounts",
				  "type": "uint256[]"
				},
				{
				  "internalType": "address[]",
				  "name": "feeTokenContracts",
				  "type": "address[]"
				},
				{
				  "internalType": "address",
				  "name": "logicContractAddress",
				  "type": "address"
				},
				{
				  "internalType": "bytes",
				  "name": "payload",
				  "type": "bytes"
				},
				{
				  "internalType": "uint256",
				  "name": "timeOut",
				  "type": "uint256"
				},
				{
				  "internalType": "bytes32",
				  "name": "invalidationId",
				  "type": "bytes32"
				},
				{
				  "internalType": "uint256",
				  "name": "invalidationNonce",
				  "type": "uint256"
				}
			  ],
			  "internalType": "struct LogicCallArgs",
			  "name": "_args",
			  "type": "tuple"
			}
		  ],
		  "name": "submitLogicCall",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
    }]`
)
