/*
Package crypto provides an interface to cryptographic syscalls.
*/
package crypto

import "github.com/nspcc-dev/neo-go/pkg/interop"

// ECDsaSecp256r1Verify checks that sig is correct msg's signature for a given pub
// (serialized public key). It uses `Neo.Crypto.VerifyWithECDsaSecp256r1` syscall.
func ECDsaSecp256r1Verify(msg []byte, pub interop.PublicKey, sig interop.Signature) bool {
	return false
}

// ECDsaSecp256k1Verify checks that sig is correct msg's signature for a given pub
// (serialized public key). It uses `Neo.Crypto.VerifyWithECDsaSecp256k1` syscall.
func ECDsaSecp256k1Verify(msg []byte, pub interop.PublicKey, sig interop.Signature) bool {
	return false
}

// ECDSASecp256r1CheckMultisig checks multiple ECDSA signatures at once. It uses
// `Neo.Crypto.CheckMultisigWithECDsaSecp256r1` syscall.
func ECDSASecp256r1CheckMultisig(msg []byte, pubs []interop.PublicKey, sigs []interop.Signature) bool {
	return false
}

// ECDSASecp256k1CheckMultisig checks multiple ECDSA signatures at once. It uses
// `Neo.Crypto.CheckMultisigWithECDsaSecp256k1` syscall.
func ECDSASecp256k1CheckMultisig(msg []byte, pubs []interop.PublicKey, sigs []interop.Signature) bool {
	return false
}
