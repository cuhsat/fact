// FKey implementation details.
package fkey

import (
	"github.com/cuhsat/fact/internal/fmount"
	"github.com/cuhsat/fact/internal/sys"
)

func RecoveryIds(img string) (ids []string, err error) {
	loi, err := fmount.LoSetupAttach(img)

	if err != nil {
		return
	}

	lops, err := fmount.Parts(loi)

	if err != nil {
		return
	}

	for _, lop := range lops {
		dev := fmount.Dev(lop)

		id, err := fmount.DislockerInfo(dev)

		if err != nil {
			sys.Error(err)
			continue
		}

		ids = append(ids, id...)

		err = fmount.LoSetupDetach(dev)

		if err != nil {
			sys.Error(err)
		}
	}

	return
}
