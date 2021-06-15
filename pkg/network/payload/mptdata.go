package payload

import (
	"errors"

	"github.com/nspcc-dev/neo-go/pkg/io"
)

// MPTData represents the set of serialized MPT nodes.
type MPTData struct {
	Nodes [][]byte
}

// EncodeBinary implements io.Serializable.
func (d *MPTData) EncodeBinary(w *io.BinWriter) {
	w.WriteArray(d.Nodes)
}

// DecodeBinary implements io.Serializable.
func (d *MPTData) DecodeBinary(r *io.BinReader) {
	r.ReadArray(&d.Nodes)
	if len(d.Nodes) == 0 {
		r.Err = errors.New("empty MPT nodes list")
	}
}
