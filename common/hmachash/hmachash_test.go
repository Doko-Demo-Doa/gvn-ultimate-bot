package hmachash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHmacHash(t *testing.T) {
	hm := NewHMAC("sample_secret")
	hm2 := NewHMAC("sample_secret2")

	t.Run("Hash string", func(t *testing.T) {
		txt := "test"
		hashed := hm.Hash(txt)
		hashed2 := hm2.Hash(txt)

		assert.NotEqual(t, hashed2, hashed)
		assert.NotEqual(t, txt, hashed)
		assert.NotEqual(t, txt, hashed2)
		assert.GreaterOrEqual(t, len(hashed), 32)
	})
}
