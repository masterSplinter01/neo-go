package crypto

import (
	"github.com/nspcc-dev/neo-go/pkg/interop"
	"github.com/nspcc-dev/neo-go/pkg/interop/contract"
)

// Hash represents GAS contract hash.
const Hash = "\x63\x4f\x72\x0f\xf1\x39\xb9\xe0\xce\x73\x1c\xb7\x74\x41\x18\x35\x13\x34\xcb\x08"

// Sha256 computes SHA256 hash of b. It uses `Neo.Crypto.SHA256` syscall.
func Sha256(b []byte) interop.Hash256 {
	return contract.Call(interop.Hash160(Hash), "sha256", contract.NoneFlag, b).(interop.Hash256)
}

// Ripemd160 computes RIPEMD160 hash of b. It uses `Neo.Crypto.RIPEMD160` syscall.
func Ripemd160(b []byte) interop.Hash160 {
	return contract.Call(interop.Hash160(Hash), "ripemd160", contract.NoneFlag, b).(interop.Hash160)
}
