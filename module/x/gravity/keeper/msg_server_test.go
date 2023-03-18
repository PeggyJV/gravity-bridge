package keeper

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"

	"github.com/peggyjv/gravity-bridge/module/v3/x/gravity/types"
)

func TestMsgServer_SubmitEVMSignature(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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

	// setup for GetValidatorEVMAddress
	gk.setValidatorEVMAddress(ctx, valAddr1, ethAddr1)

	// setup for GetOutgoingTx
	signerSetTx := gk.CreateSignerSetTx(ctx, types.EthereumChainID)

	// setup for ValidateEVMSignature
	gravityId := gk.getGravityID(ctx, types.EthereumChainID)
	checkpoint := signerSetTx.GetCheckpoint([]byte(gravityId))
	signature, err := types.NewEVMSignature(checkpoint, ethPrivKey)
	require.NoError(t, err)

	signerSetTxConfirmation := &types.SignerSetTxConfirmation{
		SignerSetNonce: signerSetTx.Nonce,
		EVMSigner:      ethAddr1.Hex(),
		Signature:      signature,
		ChainId:        types.EthereumChainID,
	}

	confirmation, err := types.PackConfirmation(signerSetTxConfirmation)
	require.NoError(t, err)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSubmitEVMTxConfirmation{
		Confirmation: confirmation,
		Signer:       orcAddr1.String(),
		ChainId:      types.EthereumChainID,
	}

	_, err = msgServer.SubmitEVMTxConfirmation(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_SendToEVM(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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
	gk.setCosmosOriginatedDenomToERC20(ctx, types.EthereumChainID, testDenom, testContract)

	// setup for GetValidatorEVMAddress
	gk.setValidatorEVMAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEVM{
		Sender:       orcAddr1.String(),
		EVMRecipient: ethAddr1.String(),
		Amount:       amount,
		BridgeFee:    fee,
		ChainId:      types.EthereumChainID,
	}

	_, err = msgServer.SendToEVM(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_CancelSendToEVM(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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
	gk.setCosmosOriginatedDenomToERC20(ctx, types.EthereumChainID, testDenom, testContract)

	// setup for GetValidatorEVMAddress
	gk.setValidatorEVMAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEVM{
		Sender:       orcAddr1.String(),
		EVMRecipient: ethAddr1.String(),
		Amount:       amount,
		BridgeFee:    fee,
		ChainId:      types.EthereumChainID,
	}

	response, err := msgServer.SendToEVM(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	cancelMsg := &types.MsgCancelSendToEVM{
		Id:      response.Id,
		Sender:  orcAddr1.String(),
		ChainId: types.EthereumChainID,
	}

	_, err = msgServer.CancelSendToEVM(sdk.WrapSDKContext(ctx), cancelMsg)
	require.NoError(t, err)
}

func TestMsgServer_RequestBatchTx(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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
	gk.setCosmosOriginatedDenomToERC20(ctx, types.EthereumChainID, testDenom, testContract)

	// setup for GetValidatorEVMAddress
	gk.setValidatorEVMAddress(ctx, valAddr1, ethAddr1)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSendToEVM{
		Sender:       orcAddr1.String(),
		EVMRecipient: ethAddr1.String(),
		Amount:       amount,
		BridgeFee:    fee,
		ChainId:      types.EthereumChainID,
	}

	_, err = msgServer.SendToEVM(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)

	requestMsg := &types.MsgRequestBatchTx{
		Signer:  orcAddr1.String(),
		Denom:   testDenom,
		ChainId: types.EthereumChainID,
	}

	_, err = msgServer.RequestBatchTx(sdk.WrapSDKContext(ctx), requestMsg)
	require.NoError(t, err)
}

func TestMsgServer_RequestEmptyBatchTx(t *testing.T) {
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

		testDenom    = "stake"
		testContract = common.HexToAddress("0x429881672B9AE42b8EbA0E26cD9C73711b891Ca5")
	)

	{ // setup for getSignerValidator
		gk.StakingKeeper = NewStakingKeeperMock(valAddr1, valAddr2, valAddr3)
		gk.SetOrchestratorValidatorAddress(ctx, valAddr1, orcAddr1)
	}

	// create denom in keeper
	gk.setCosmosOriginatedDenomToERC20(ctx, types.EthereumChainID, testDenom, testContract)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgRequestBatchTx{
		Signer:  orcAddr1.String(),
		Denom:   testDenom,
		ChainId: types.EthereumChainID,
	}

	_, err := msgServer.RequestBatchTx(sdk.WrapSDKContext(ctx), msg)

	require.Error(t, err)
}

func TestMsgServer_SubmitEVMEvent(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env = CreateTestEnv(t)
		ctx = env.Context
		gk  = env.GravityKeeper

		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)

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

	// setup for GetValidatorEVMAddress
	gk.setValidatorEVMAddress(ctx, valAddr1, ethAddr1)

	sendToCosmosEvent := &types.SendToCosmosEvent{
		EventNonce:     1,
		TokenContract:  testContract.Hex(),
		Amount:         sdk.NewInt(1000),
		EVMSender:      ethAddr1.String(),
		CosmosReceiver: orcAddr1.String(),
		EVMHeight:      200,
		ChainId:        types.EthereumChainID,
	}

	event, err := types.PackEvent(sendToCosmosEvent)
	require.NoError(t, err)

	msgServer := NewMsgServerImpl(gk)

	msg := &types.MsgSubmitEVMEvent{
		Event:   event,
		Signer:  orcAddr1.String(),
		ChainId: types.EthereumChainID,
	}

	_, err = msgServer.SubmitEVMEvent(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_SetDelegateKeys(t *testing.T) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	require.NoError(t, err)

	var (
		env         = CreateTestEnv(t)
		ctx         = env.Context
		gk          = env.GravityKeeper
		orcAddr1, _ = sdk.AccAddressFromBech32("cosmos1dg55rtevlfxh46w88yjpdd08sqhh5cc3xhkcej")
		valAddr1    = sdk.ValAddress(orcAddr1)
		ethAddr1    = crypto.PubkeyToAddress(ethPrivKey.PublicKey)
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
	hash := crypto.Keccak256Hash(signMsgBz).Bytes()

	sig, err := types.NewEVMSignature(hash, ethPrivKey)
	require.NoError(t, err)

	msg := &types.MsgDelegateKeys{
		ValidatorAddress:    valAddr1.String(),
		OrchestratorAddress: orcAddr1.String(),
		EVMAddress:          ethAddr1.String(),
		EVMSignature:        sig,
	}

	_, err = msgServer.SetDelegateKeys(sdk.WrapSDKContext(ctx), msg)
	require.NoError(t, err)
}

func TestMsgServer_SubmitEVMHeightVote(t *testing.T) {
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

	msg := &types.MsgEVMHeightVote{
		EvmHeight: 5,
		Signer:    orcAddr1.String(),
		ChainId:   types.EthereumChainID,
	}

	_, err := msgServer.SubmitEVMHeightVote(sdk.WrapSDKContext(ctx), msg)

	require.NoError(t, err)
	require.Equal(t, gk.GetEVMHeightVote(ctx, types.EthereumChainID, valAddr1).EVMHeight, uint64(5))
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

	privKey, err := crypto.ToECDSA(privKeyBz)
	require.NoError(t, err)
	require.NotNil(t, privKey)

	require.True(t, bytes.Equal(privKeyBz, crypto.FromECDSA(privKey)))
	require.Equal(t, privKeyHexStr, hexutil.Encode(crypto.FromECDSA(privKey)))

	publicKey := privKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	require.True(t, ok)

	address := crypto.PubkeyToAddress(*publicKeyECDSA)
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
	hash := crypto.Keccak256Hash(signMsgBz).Bytes()

	sig, err := types.NewEVMSignature(hash, privKey)
	sig[64] += 27 // change the V value
	require.NoError(t, err)

	err = types.ValidateEVMSignature(hash, sig, address)
	require.NoError(t, err)

	// replace gorcSig with what the following command produces:
	// $ gorc sign-delegate-keys <your-eth-key-name> cosmosvaloper1dmly9yyhd5lyhyl8qhs7wtcd4xt7gyxlesgvmc 0
	gorcSig := "0xbda7037e448ca07ac91f5f386b72df37b6bbacf102b2c8f5acb58b5e053d68d96875ce9e442433bea55ac083230f492670ca2c07a8303c332dca06b1c0758c661b"
	require.Equal(t, hexutil.Encode(sig), gorcSig)
}
