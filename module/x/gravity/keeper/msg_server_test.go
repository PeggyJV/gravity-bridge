package keeper

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"testing"

	types1 "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v5/x/gravity/types"
)

var (
	nonexistentOrcAddr, _ = sdk.AccAddressFromBech32("cosmos13txzft28sfwqlg38vkkparzaxyzewhws5ucqhe")
)

func TestMsgServer_SubmitEthereumSignature(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr2, orcAddr2)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr3, orcAddr3)
	}

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	// setup for GetOutgoingTx
	signerSetTx := gk.CreateSignerSetTx(ctx)

	// setup for ValidateEthereumSignature
	gravityId := gk.getGravityID(ctx)
	checkpoint := signerSetTx.GetCheckpoint([]byte(gravityId))
	signature, err := types.NewEthereumSignature(checkpoint, ethPrivKey)
	require.NoError(t, err)

	signerSetTxConfirmation := &types.SignerSetTxConfirmation{
		SignerSetNonce: signerSetTx.Nonce,
		EthereumSigner: ethAddr1.Hex(),
		Signature:      signature,
	}

	confirmation, err := types.PackConfirmation(signerSetTxConfirmation)
	require.NoError(t, err)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSubmitEthereumTxConfirmation{
		Confirmation: confirmation,
		Signer:       orcAddr1.String(),
	}

	_, err = msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// Test error scenarios for SubmitEthereumTxConfirmation

	t.Run("Invalid confirmation type", func(t *testing.T) {
		invalidConfirmation := &types1.Any{
			TypeUrl: "invalid-type",
			Value:   []byte("invalid confirmation"),
		}
		msg := &types.MsgSubmitEthereumTxConfirmation{
			Confirmation: invalidConfirmation,
			Signer:       orcAddr1.String(),
		}

		_, err = msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed unpacking protobuf message from Any")
	})

	t.Run("Non-existent validator", func(t *testing.T) {
		msg := &types.MsgSubmitEthereumTxConfirmation{
			Confirmation: confirmation,
			Signer:       nonexistentOrcAddr.String(),
		}

		_, err = msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not orchestrator or validator")
	})

	t.Run("Invalid Ethereum signature", func(t *testing.T) {
		invalidSignature, _ := types.NewEthereumSignature(checkpoint, ethPrivKey)
		invalidSignature[0] ^= 0xFF // Flip some bits to make it invalid

		invalidConfirmation := &types.SignerSetTxConfirmation{
			SignerSetNonce: signerSetTx.Nonce,
			EthereumSigner: ethAddr1.Hex(),
			Signature:      invalidSignature,
		}

		packedInvalidConfirmation, _ := types.PackConfirmation(invalidConfirmation)

		msg := &types.MsgSubmitEthereumTxConfirmation{
			Confirmation: packedInvalidConfirmation,
			Signer:       orcAddr1.String(),
		}

		_, err = msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signature verification failed")
	})

	t.Run("Duplicate confirmation", func(t *testing.T) {
		// First submission should succeed
		msg := &types.MsgSubmitEthereumTxConfirmation{
			Confirmation: confirmation,
			Signer:       orcAddr1.String(),
		}

		// Second submission of the same confirmation should fail
		_, err = msgServer.SubmitEthereumTxConfirmation(sdk.WrapSDKContext(ctx), msg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signature duplicate: invalid")
	})
}

func TestMsgServer_SendToEthereum(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)

		testDenom    = "stake"
		testContract = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")

		balance = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10000),
		}
		amount = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(1000),
		}
		fee = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10),
		}
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
	}

	{ // add balance to bank
		err = env.AddBalanceToBank(ctx, orcAddr1, sdk.Coins{balance})
		require.NoError(t, err)
	}

	// create denom in keeper
	gk.setCosmosOriginatedDenomToERC20(ctx, testDenom, testContract)

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEthereum{
		Sender:            orcAddr1.String(),
		EthereumRecipient: ethAddr1.String(),
		Amount:            amount,
		BridgeFee:         fee,
	}

	_, err = msgServer.SendToEthereum(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_CancelSendToEthereum(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)

		testDenom    = "stake"
		testContract = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")

		balance = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10000),
		}
		amount = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(1000),
		}
		fee = sdk.Coin{
			Denom:  testDenom,
			Amount: sdk.NewInt(10),
		}
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
	}

	{ // add balance to bank
		err = env.AddBalanceToBank(ctx, orcAddr1, sdk.Coins{balance})
		require.NoError(t, err)
	}

	// create denom in keeper
	gk.setCosmosOriginatedDenomToERC20(ctx, testDenom, testContract)

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEthereum{
		Sender:            orcAddr1.String(),
		EthereumRecipient: ethAddr1.String(),
		Amount:            amount,
		BridgeFee:         fee,
	}

	response, err := msgServer.SendToEthereum(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	cancelMsg := &types.MsgCancelSendToEthereum{
		Id:     response.Id,
		Sender: orcAddr1.String(),
	}

	_, err = msgServer.CancelSendToEthereum(sdk.WrapSDKContext(ctx), cancelMsg)
	require.NoError(t, err)

	// Test error cases for CancelSendToEthereum

	t.Run("Invalid ID", func(t *testing.T) {
		// Test case: Invalid ID
		invalidIDMsg := &types.MsgCancelSendToEthereum{
			Id:     999999, // Non-existent ID
			Sender: orcAddr1.String(),
		}
		_, err = msgServer.CancelSendToEthereum(sdk.WrapSDKContext(ctx), invalidIDMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not found")
	})

	t.Run("Sender is not the original sender", func(t *testing.T) {
		// Create a new send to ethereum message
		msg := &types.MsgSendToEthereum{
			Sender:            orcAddr1.String(),
			EthereumRecipient: ethAddr1.String(),
			Amount:            amount,
			BridgeFee:         fee,
		}
		response, err = msgServer.SendToEthereum(sdk.WrapSDKContext(ctx), msg)
		require.NoError(t, err)

		wrongSenderMsg := &types.MsgCancelSendToEthereum{
			Id:     response.Id,
			Sender: orcAddr2.String(), // Different sender
		}
		_, err = msgServer.CancelSendToEthereum(sdk.WrapSDKContext(ctx), wrongSenderMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "can't cancel a message you didn't send")
	})
}

func TestMsgServer_SubmitEthereumEvent(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)

		testContract = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr2, orcAddr2)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr3, orcAddr3)
	}

	// setup for GetValidatorEthereumAddress
	gk.setValidatorEthereumAddress(ctx, valAddr1, ethAddr1)

	sendToCosmosEvent := &types.SendToCosmosEvent{
		EventNonce:     1,
		TokenContract:  testContract.Hex(),
		Amount:         sdk.NewInt(1000),
		EthereumSender: ethAddr1.String(),
		CosmosReceiver: orcAddr1.String(),
		EthereumHeight: 200,
	}

	event, err := types.PackEvent(sendToCosmosEvent)
	require.NoError(t, err)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSubmitEthereumEvent{
		Event:  event,
		Signer: orcAddr1.String(),
	}

	_, err = msgServer.SubmitEthereumEvent(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// Test error cases for SubmitEthereumEvent
	t.Run("Invalid signer address", func(t *testing.T) {
		invalidMsg := &types.MsgSubmitEthereumEvent{
			Event:  event,
			Signer: "invalid_address",
		}
		_, err := msgServer.SubmitEthereumEvent(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "signer address: invalid")
	})

	t.Run("Non-existent signer", func(t *testing.T) {
		invalidMsg := &types.MsgSubmitEthereumEvent{
			Event:  event,
			Signer: nonexistentOrcAddr.String(), // Using a different orchestrator address
		}
		res, err := msgServer.SubmitEthereumEvent(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Nil(t, res)
		require.Error(t, err)
		require.Contains(t, err.Error(), "not orchestrator or validator")
	})

	t.Run("Non-contiguous event nonce", func(t *testing.T) {
		invalidEvent, err := types.PackEvent(&types.ContractCallExecutedEvent{
			EventNonce: 10,
		})
		require.NoError(t, err)

		invalidMsg := &types.MsgSubmitEthereumEvent{
			Event:  invalidEvent,
			Signer: orcAddr1.String(),
		}
		_, err = msgServer.SubmitEthereumEvent(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "non contiguous event nonce")
	})
}

func TestMsgServer_SetDelegateKeys(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env         = CreateTestEnv(t)
		ctx         = env.Context
		gk          = env.GravityKeeper
		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey)
	)

	// setup for getSignerValidator
	gk.StakingKeeper = NewStakingKeeperMock(valAddr1)

	// Set the sequence to 1 because the antehandler will do this in the full
	// chain.
	acc := env.AccountKeeper.NewAccountWithAddress(ctx, orcAddr1)
	acc.SetSequence(1)
	env.AccountKeeper.SetAccount(ctx, acc)

	msgServer := NewMsgServerImpl(gk)

	ethMsg := types.DelegateKeysSignMsg{
		ValidatorAddress: valAddr1.String(),
		Nonce:            0,
	}
	signMsgBz := env.Marshaler.MustMarshal(&ethMsg)
	hash := ethCrypto.Keccak256Hash(signMsgBz).Bytes()

	sig, err := types.NewEthereumSignature(hash, ethPrivKey)
	require.NoError(t, err)

	msg := &types.MsgDelegateKeys{
		ValidatorAddress:    valAddr1.String(),
		OrchestratorAddress: orcAddr1.String(),
		EthereumAddress:     ethAddr1.String(),
		EthSignature:        sig,
	}

	_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	// Test error cases for SetDelegateKeys
	t.Run("Invalid validator address", func(t *testing.T) {
		invalidMsg := &types.MsgDelegateKeys{
			ValidatorAddress:    "invalid_address",
			OrchestratorAddress: orcAddr1.String(),
			EthereumAddress:     ethAddr1.String(),
			EthSignature:        sig,
		}
		_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid validator address")
	})

	t.Run("Invalid orchestrator address", func(t *testing.T) {
		invalidMsg := &types.MsgDelegateKeys{
			ValidatorAddress:    valAddr1.String(),
			OrchestratorAddress: "invalid_address",
			EthereumAddress:     ethAddr1.String(),
			EthSignature:        sig,
		}
		_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid orchestrator address")
	})

	t.Run("Ethereum address already in use", func(t *testing.T) {

		invalidMsg := &types.MsgDelegateKeys{
			ValidatorAddress:    valAddr1.String(),
			OrchestratorAddress: orcAddr1.String(),
			EthereumAddress:     ethAddr1.String(),
			EthSignature:        sig,
		}
		_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("ethereum address %s in use", ethAddr1))
	})

	t.Run("Orchestrator address already in use", func(t *testing.T) {
		invalidMsg := &types.MsgDelegateKeys{
			ValidatorAddress:    valAddr1.String(),
			OrchestratorAddress: orcAddr1.String(),
			EthereumAddress:     "anything",
			EthSignature:        sig,
		}
		_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), fmt.Sprintf("orchestrator address %s in use", orcAddr1))
	})

	t.Run("Invalid ethereum signature", func(t *testing.T) {
		invalidSig := []byte("invalid_signature")
		invalidMsg := &types.MsgDelegateKeys{
			ValidatorAddress:    valAddr1.String(),
			OrchestratorAddress: orcAddr2.String(),
			EthereumAddress:     "anything",
			EthSignature:        invalidSig,
		}
		_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "failed to validate delegate keys signature for Ethereum address")
	})

	t.Run("Validator not found", func(t *testing.T) {
		nonExistentValAddr := sdk.ValAddress(bytes.Repeat([]byte{1}, 20))
		invalidMsg := &types.MsgDelegateKeys{
			ValidatorAddress:    nonExistentValAddr.String(),
			OrchestratorAddress: orcAddr1.String(),
			EthereumAddress:     ethAddr1.String(),
			EthSignature:        sig,
		}
		_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), invalidMsg)
		require.Error(t, err)
		require.Contains(t, err.Error(), "validator does not exist")
	})
}

func TestMsgServer_SubmitEthereumHeightVote(t *testing.T) {
	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)

		orcAddr2, _ = sdk.AccAddressFromBech32("cosmos164knshrzuuurf05qxf3q5ewpfnwzl4gj4m4dfy")
		valAddr2    = sdk.ValAddress(orcAddr2)

		orcAddr3, _ = sdk.AccAddressFromBech32("cosmos193fw83ynn76328pty4yl7473vg9x86alq2cft7")
		valAddr3    = sdk.ValAddress(orcAddr3)
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr2, orcAddr2)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr3, orcAddr3)
	}

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgEthereumHeightVote{
		EthereumHeight: 5,
		Signer:         orcAddr1.String(),
	}

	_, err := msgServer.SubmitEthereumHeightVote(sdk.WrapSDKContext(ctx), msg)

	require.NoError(t, err)
	require.Equal(t, gk.GetEthereumHeightVote(ctx, valAddr1).EthereumHeight, uint64(5))
}

func TestEthVerify(t *testing.T) {
	// Replace privKeyHexStr and addrHexStr with your own private key and address
	// HEX values.
	privKeyHexStr := "0xee63225c8a0928168d362147cd19859de6459e972ffcf9294a69382b4ad99720"
	addrHexStr := "0xA093773C30Ad5c3e83B20E66CB4e6136aEa098B7"

	// ==========================================================================
	// setup
	// ==========================================================================
	privKeyBz, err := hexutil.Decode(privKeyHexStr)
	require.NoError(t, err)

	privKey, err := ethCrypto.ToECDSA(privKeyBz)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	require.True(t, bytes.Equal(privKeyBz, ethCrypto.FromECDSA(privKey)))
	require.Equal(t, privKeyHexStr, hexutil.Encode(ethCrypto.FromECDSA(privKey)))

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := ethCrypto.PubkeyToAddress(*publicKeyECDSA)
	require.Equal(t, addrHexStr, address.Hex())

	// ==========================================================================
	// signature verification
	// ==========================================================================
	cdc := MakeTestMarshaler()

	valAddr := "cosmosvaloper16k7rf90uvt4tgslqh280wvdzxp5q9ah6nxxupc"
	signMsgBz, err := cdc.Marshal(&types.DelegateKeysSignMsg{
		ValidatorAddress: valAddr,
		Nonce:            0,
	})

	require.NoError(t, err)

	fmt.Println("MESSAGE BYTES TO SIGN:", hexutil.Encode(signMsgBz))
	hash := ethCrypto.Keccak256Hash(signMsgBz).Bytes()

	sig, err := types.NewEthereumSignature(hash, privKey)
	sig[64] += 27 // change the V value
	require.NoError(t, err)

	err = types.ValidateEthereumSignature(hash, sig, address)
	require.NoError(t, err)

	// replace gorcSig with what the following command produces:
	// $ gorc sign-delegate-keys <your-eth-key-name> cosmosvaloper1dmly9yyhd5lyhyl8qhs7wtcd4xt7gyxlesgvmc 0
	gorcSig := "0xbda7037e448ca07ac91f5f386b72df37b6bbacf102b2c8f5acb58b5e053d68d96875ce9e442433bea55ac083230f492670ca2c07a8303c332dca06b1c0758c661b"
	require.Equal(t, hexutil.Encode(sig), gorcSig)
}
