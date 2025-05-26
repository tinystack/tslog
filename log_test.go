package tslog

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLevel(t *testing.T) {
	assert.Equal(t, ParseLevel("debug"), DebugLevel)
	assert.Equal(t, ParseLevel("info"), InfoLevel)
	assert.Equal(t, ParseLevel("warn"), WarnLevel)
	assert.Equal(t, ParseLevel("error"), ErrorLevel)
}

func TestWithLevel(t *testing.T) {
	opt := &Options{
		lvl: DebugLevel,
	}
	WithLevel(ErrorLevel)(opt)
	assert.Equal(t, ErrorLevel, opt.lvl)
}

func TestWithWriter(t *testing.T) {
	opt := &Options{
		w: []io.Writer{os.Stdout},
	}
	WithWriter(os.Stderr, os.Stdout)(opt)
	assert.Equal(t, []io.Writer{os.Stderr, os.Stdout}, opt.w)
}

func TestWithCaller(t *testing.T) {
	opt := &Options{
		caller: true,
	}
	WithCaller(false)(opt)
	assert.Equal(t, false, opt.caller)
}

func TestWithEncoder(t *testing.T) {
	opt := &Options{
		encoder: EncoderConsole,
	}
	WithEncoder(EncoderJSON)(opt)
	assert.Equal(t, EncoderJSON, opt.encoder)
}
