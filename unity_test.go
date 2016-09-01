package gounity

import (
	"testing"
	"github.com/jealous/gounity/rsc"
	"github.com/stretchr/testify/assert"
)

func MockUnity() Unity{
	return rsc.NewUnityByConn(
		rsc.NewMockConnection("10.244.223.66", "admin", "Password123!"))
}

func TestUnity_GetPoolList(t *testing.T) {
	poolList := MockUnity().GetPoolList()
	assert.Equal(t, 2, poolList.Size())
}