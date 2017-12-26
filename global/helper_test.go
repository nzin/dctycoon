package global

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMega(t *testing.T) {
	assert.Equal(t, int32(0), ParseMega(""), "ParseMega('0')")
	assert.Equal(t, int32(4), ParseMega("4M"), "ParseMega('4M')")
	assert.Equal(t, int32(4096), ParseMega("4G"), "ParseMega('4G')")
	assert.Equal(t, int32(4096*1024), ParseMega("4T"), "ParseMega('4T')")
	assert.Equal(t, int32(0), ParseMega("M"), "ParseMega('M')")
	assert.Equal(t, int32(4), ParseMega("4MT"), "ParseMega('4MT')")
	assert.Equal(t, int32(2147483647), ParseMega("1234567890G"), "ParseMega('1234567890G')")
	assert.Equal(t, int32(2147483647), ParseMega("123456789012M"), "ParseMega('123456789012M')")
}
