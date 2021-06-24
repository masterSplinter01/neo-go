package blockchainer

import (
	"github.com/nspcc-dev/neo-go/pkg/core/mpt"
	"github.com/nspcc-dev/neo-go/pkg/core/state"
	"github.com/nspcc-dev/neo-go/pkg/crypto/keys"
	"github.com/nspcc-dev/neo-go/pkg/util"
)

// StateRoot represents local state root module.
type StateRoot interface {
	AddStateRoot(root *state.MPTRoot) error
	Collapse(depth int)
	CurrentLocalStateRoot() util.Uint256
	CurrentValidatedHeight() uint32
	GetStateProof(root util.Uint256, key []byte) ([][]byte, error)
	GetStateRoot(height uint32) (*state.MPTRoot, error)
	GetStateValidators(height uint32) keys.PublicKeys
	RestoreMPTNode(path []byte, node mpt.Node) error
	SetUpdateValidatorsCallback(func(uint32, keys.PublicKeys))
	Traverse(root util.Uint256, f func(nodeBytes []byte), stop func(node []byte) bool) error
	UpdateStateValidators(height uint32, pubs keys.PublicKeys)
}
