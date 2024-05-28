// Evt implementation tests.
package evt

import (
	"testing"

	"github.com/cuhsat/fact/internal/fact"
)

func TestTool(t *testing.T) {
	t.Run("", func(t *testing.T) {
		asm, err := fact.EzTools("EvtxECmd.dll")

		if err != nil {
			t.Fatal(err)
		}

		if len(asm) == 0 {
			t.Fatal("tool not found")
		}
	})
}
