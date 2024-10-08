package types

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"cosmossdk.io/errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	signaturePrefix = "\x19Ethereum Signed Message:\n32"
)

// NewEthereumSignature creates a new signuature over a given byte array
func NewEthereumSignature(hash []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.Wrap(ErrInvalid, "did not pass in private key")
	}
	protectedHash := crypto.Keccak256Hash(append([]byte(signaturePrefix), hash...))
	return crypto.Sign(protectedHash.Bytes(), privateKey)
}

// decodeSignature was duplicated from go-ethereum with slight modifications
func decodeSignature(sig []byte) (r, s *big.Int, v byte) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	if sig[64] == 27 || sig[64] == 28 {
		v = sig[64] - 27
	} else {
		v = sig[64]
	}
	return r, s, v
}

// ValidateEthereumSignature takes a message, an associated signature and public key and
// returns an error if the signature isn't valid
func ValidateEthereumSignature(hash []byte, signature []byte, ethAddress common.Address) error {

	/// signature to public key: invalid signature length: invalid
	/// signature not matching: invalid: invalid
	if len(signature) < 65 {
		return errors.Wrapf(ErrInvalid, "signature too short signature %x", signature)
	}

	// Copy to avoid mutating signature slice by accident
	var sigCopy = make([]byte, len(signature))
	copy(sigCopy, signature)

	r, s, v := decodeSignature(sigCopy)
	if !crypto.ValidateSignatureValues(v, r, s, true) {
		return errors.Wrap(ErrInvalid, "Signature values failed validation")
	}

	// To verify signature
	// - use crypto.SigToPub to get the public key
	// - use crypto.PubkeyToAddress to get the address
	// - compare this to the address given.

	// for backwards compatibility reasons  the V value of an Ethereum sig is presented
	// as 27 or 28, internally though it should be a 0-3 value due to changed formats.
	// It seems that go-ethereum expects this to be done before sigs actually reach it's
	// internal validation functions. In order to comply with this requirement we check
	// the sig an dif it's in standard format we correct it. If it's in go-ethereum's expected
	// format already we make no changes.
	//
	// We could attempt to break or otherwise exit early on obviously invalid values for this
	// byte, but that's a task best left to go-ethereum
	if sigCopy[64] == 27 || sigCopy[64] == 28 {
		sigCopy[64] -= 27
	}

	hash = append([]uint8(signaturePrefix), hash...)

	pubkey, err := crypto.SigToPub(crypto.Keccak256Hash(hash).Bytes(), sigCopy)
	if err != nil {
		return errors.Wrapf(err, "signature to public key sig %x hash %x", sigCopy, hash)
	}

	if addr := crypto.PubkeyToAddress(*pubkey); addr != ethAddress {
		return errors.Wrapf(ErrInvalid, "signature not matching addr %x sig %x hash %x", addr, signature, hash)
	}

	return nil
}
