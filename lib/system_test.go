package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func getTestSystem() *System {
	return NewSystem("10.0.0.5", "admin", "password")
}

func TestSystem_Ip(t *testing.T) {
	assert.Equal(t, getTestSystem().Ip(), "10.0.0.5")
}

func TestSystem_Username(t *testing.T) {
	assert.Equal(t, getTestSystem().Username(), "admin")
}
