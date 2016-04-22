package wmenu

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testInvalid    = newMenuError(ErrInvalid, "1", 0)
	testDuplicate  = newMenuError(ErrDuplicate, "2", 0)
	testTooMany    = newMenuError(ErrTooMany, "", 0)
	testNoResponse = newMenuError(ErrNoResponse, "", 0)
	errNormal      = errors.New("Opps")
	errMenu        = newMenuError(errors.New("General"), "", 0)
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

func TestIsMenuError(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsMenuErr(testNoResponse))
	assert.True(IsMenuErr(testDuplicate))
	assert.True(IsMenuErr(testInvalid))
	assert.True(IsMenuErr(testTooMany))
	assert.True(IsMenuErr(testNoResponse))
	assert.True(IsMenuErr(errMenu))
	assert.False(IsMenuErr(errNormal))
}
