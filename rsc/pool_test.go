package rsc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllPool(t *testing.T) {
	poolList := GetPoolList(MockConn())
	asserts := assert.New(t)

	asserts.NotNil(poolList)
	asserts.Equal(2, poolList.Size())
	for it := poolList.Iterator(); it.Next(); {
		pool := it.Value().(*Pool)
		if pool.Name == "perfpool1132" {
			ValidatePool1132(t, pool)
			break
		}
	}
}

func TestGetPoolByName(t *testing.T) {
	pool := GetPoolByName(MockConn(), "perfpool1132")
	ValidatePool1132(t, pool)
}

func TestGetPoolByName_NotFound(t *testing.T) {
	pool := GetPoolByName(MockConn(), "notFound")
	assert.Nil(t, pool)
}

func TestGetPoolById(t *testing.T) {
	pool := GetPoolById(MockConn(), "pool_1")
	ValidatePool1132(t, pool)
}

func ValidatePool1132(t *testing.T, pool *Pool) {
	asserts := assert.New(t)
	asserts.Equal("pool_1", pool.Id)
	asserts.Equal("", pool.Description)
	asserts.Equal(uint64(1365799600128), pool.SizeFree)
	asserts.Equal(uint64(1374657970176), pool.SizeTotal)
	asserts.Equal(uint64(8858370048), pool.SizeUsed)
}

func TestPool_GetLunList(t *testing.T) {
	pool := &Pool{Rsc: Rsc{conn: MockConn(), type_: "pool", Id: "pool_1"}}
	lunList := pool.GetLunList()
	assert.Equal(t, 2, lunList.Size())

	lun := lunList.Iterator().Value().(*Lun)
	assert.NotEqual(t, 0, lun.SizeTotal)
	assert.Equal(t, "pool_1", lun.Pool.Id)
}
