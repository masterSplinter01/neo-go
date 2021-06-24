package testdata

import (
	"github.com/nspcc-dev/neo-go/pkg/rpc/client"
	"github.com/nspcc-dev/neo-go/pkg/smartcontract"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

var contractHash = util.Uint160{0xa, 0x73, 0x77, 0x10, 0x82, 0xe8, 0xa7, 0xe4, 0xe2, 0x9e, 0x9f, 0x6b, 0xe5, 0x72, 0xfb, 0x49, 0x45, 0x55, 0x27, 0xf9}

// Client is a wrapper over RPC-client mirroring methods of smartcontract.
type Client client.Client

// Transfer returns one value.
func (c *Client) Transfer(to util.Uint160, amount int64) (bool, error) {
	args := make([]smartcontract.Parameter, 2)
	args[0] = smartcontract.Parameter{Type: smartcontract.Hash160Type, Value: to.BytesBE()}
	args[1] = smartcontract.Parameter{Type: smartcontract.IntegerType, Value: amount}

	result, err := (*client.Client)(c).InvokeFunction(contractHash, "transfer", args, nil)
	if err != nil {
		return false, err
	}

	err = client.GetInvocationError(result)
	if err != nil {
		return false, err
	}

	return client.TopBoolFromStack(result.Stack)
}
