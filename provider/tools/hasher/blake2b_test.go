package hasher_test

import (
	"testing"

	"github.com/klever-io/klever-go-sdk/provider/tools/hasher"
	"github.com/stretchr/testify/assert"
)

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
