//go:build !windows

// Volume Shadow Copy functions (Error stub).
package ffind

import "errors"

func ShadowCopy(drive string) (dir string, err error) {
	return "", errors.ErrUnsupported
}
