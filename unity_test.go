package gounity

import (
	"github.com/jealous/gounity/rsc"
	"github.com/stretchr/testify/assert"
	"testing"
)

func MockUnity() Unity {
	return rsc.NewUnityByConn(
		rsc.NewMockConnection("10.244.223.66", "admin", "Password123!"))
}

func TestUnity_GetPoolList(t *testing.T) {
	poolList := MockUnity().GetPoolList()
	assert.Equal(t, 2, poolList.Size())
}

func TestUnity_systemProperty(t *testing.T) {
	unity := MockUnity().(*rsc.Unity)
	Update(unity)
	assert.Equal(t, "FNM00150600267", unity.SerialNumber)
}
