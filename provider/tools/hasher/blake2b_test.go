package hasher_test

import (
	"encoding/hex"
	"testing"

	"github.com/klever-io/klever-go-sdk/provider/tools/hasher"
	"github.com/stretchr/testify/assert"
)

func TestBlake2b_NewHasher(t *testing.T) {
	t.Parallel()

	hasher, err := hasher.NewHasher()
	assert.Nil(t, err)

	resEmpty := hasher.Compute("")
	assert.Equal(t, 32, len(resEmpty))
}

func TestBlake2b_ComputeWithDifferentHashSizes(t *testing.T) {
	t.Parallel()

	input := "dummy string"
	sizes := []int{2, 5, 8, 16, 32, 37, 64}
	for _, size := range sizes {
		testComputeOk(t, input, size)
	}
}

func testComputeOk(t *testing.T, input string, size int) {
	hasher := hasher.Blake2b{HashSize: size}
	res := hasher.Compute(input)
	assert.Equal(t, size, len(res))
}

func TestBlake2b_Empty(t *testing.T) {

	hasher := &hasher.Blake2b{HashSize: 64}

	var nilStr string
	resNil := hasher.Compute(nilStr)
	assert.Equal(t, 64, len(resNil))

	resEmpty := hasher.Compute("")
	assert.Equal(t, 64, len(resEmpty))

	assert.Equal(t, resEmpty, resNil)
}

func TestBlake2b_ZeroDefaultSize(t *testing.T) {
	hasher := &hasher.Blake2b{HashSize: 0}
	assert.Equal(t, 32, hasher.Size())

	resEmpty := hasher.Compute("")
	assert.Equal(t, 32, len(resEmpty))

	hash := hasher.EmptyHash()
	assert.Equal(t, "0e5751c026e543b2e8ab2eb06099daa1d1e5df47778f7787faab45cdf12fe3a8", hex.EncodeToString(hash))
}

func TestBlake2b_NilInterface(t *testing.T) {
	var hasher *hasher.Blake2b
	assert.True(t, hasher.IsInterfaceNil())
}
