package native

import (
	"github.com/nspcc-dev/neo-go/pkg/core/interop"
	"github.com/nspcc-dev/neo-go/pkg/core/native/nativenames"
	"github.com/nspcc-dev/neo-go/pkg/crypto/hash"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/callflag"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract/manifest"
	"github.com/nspcc-dev/neo-go/pkg/vm/stackitem"
)

// Crypto represents CryptoLib contract.
type Crypto struct {
	interop.ContractMD
}

const cryptoContractID = -19

func newCrypto() *Crypto {
	c := &Crypto{ContractMD: *interop.NewContractMD(nativenames.CryptoLib, cryptoContractID)}

	desc := newDescriptor("sha256", smartcontract.ByteArrayType,
		manifest.NewParameter("data", smartcontract.ByteArrayType))
	md := newMethodAndPrice(c.sha256, 1<<15, callflag.NoneFlag)
	c.AddMethod(md, desc)

	desc = newDescriptor("ripemd160", smartcontract.ByteArrayType,
		manifest.NewParameter("data", smartcontract.ByteArrayType))
	md = newMethodAndPrice(c.ripemd160, 1<<15, callflag.NoneFlag)
	c.AddMethod(md, desc)

	return c
}

func (c *Crypto) sha256(_ *interop.Context, args []stackitem.Item) stackitem.Item {
	bs, err := args[0].TryBytes()
	if err != nil {
		panic(err)
	}
	return stackitem.NewByteArray(hash.Sha256(bs).BytesBE())
}

func (c *Crypto) ripemd160(_ *interop.Context, args []stackitem.Item) stackitem.Item {
	bs, err := args[0].TryBytes()
	if err != nil {
		panic(err)
	}
	return stackitem.NewByteArray(hash.RipeMD160(bs).BytesBE())
}

// Metadata implements Contract interface.
func (c *Crypto) Metadata() *interop.ContractMD {
	return &c.ContractMD
}

// Initialize implements Contract interface.
func (c *Crypto) Initialize(ic *interop.Context) error {
	return nil
}

// OnPersist implements Contract interface.
func (c *Crypto) OnPersist(ic *interop.Context) error {
	return nil
}

// PostPersist implements Contract interface.
func (c *Crypto) PostPersist(ic *interop.Context) error {
	return nil
}
