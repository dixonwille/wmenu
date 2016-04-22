package wmenu

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testInvalid    = newMenuError(ErrInvalid, "1")
	testDuplicate  = newMenuError(ErrDuplicate, "2")
	testTooMany    = newMenuError(ErrTooMany, "")
	testNoResponse = newMenuError(ErrNoResponse, "")
	errNormal      = errors.New("Opps")
)

func TestIsInvalidErr(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsInvalidErr(testInvalid))
	assert.False(IsInvalidErr(testDuplicate))
	assert.False(IsInvalidErr(errNormal))
}

func TestIsDuplicateErr(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsDuplicateErr(testDuplicate))
	assert.False(IsDuplicateErr(testInvalid))
	assert.False(IsDuplicateErr(errNormal))
}

func TestIsTooManyErr(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsTooManyErr(testTooMany))
	assert.False(IsTooManyErr(testDuplicate))
	assert.False(IsTooManyErr(errNormal))
}

func TestIsNoResponseErr(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsNoResponseErr(testNoResponse))
	assert.False(IsNoResponseErr(testDuplicate))
	assert.False(IsNoResponseErr(errNormal))
}
