package rsc

import "fmt"

type Pool struct {
	Rsc
	Name        string
	Health      Health
	Description string
	SizeFree    uint64
	SizeTotal   uint64
	SizeUsed    uint64
}

type PoolList struct {
	RscList
}

func (pl *PoolList) initRsc() Rscer {
	return &Pool{Rsc: Rsc{conn: pl.conn, type_: "pool"}}
}

type PoolListCtor struct {
	conn *Connection
}

func (pi *PoolListCtor) initList(filter string) RscLister {
	ret := &PoolList{RscList: RscList{type_: "pool", conn: pi.conn}}
	ret.filter = filter
	return ret
}

func GetPoolByName(conn *Connection, name string) *Pool {
	pool := getRscByName(name, &PoolListCtor{conn})
	if pool == nil {
		return nil
	} else {
		return pool.(*Pool)
	}
}

func GetPoolById(conn *Connection, id string) *Pool {
	pool := &Pool{Rsc: Rsc{conn: conn, type_: "pool", Id: id}}
	Update(pool)
	return pool
}

func GetPoolList(conn *Connection) *PoolList {
	return getRscList(&PoolListCtor{conn}).(*PoolList)
}

func (pool *Pool) GetLunList() *LunList {
	lunList := (&LunListCtor{pool.GetConn()}).initList(
		fmt.Sprintf(`pool eq "%v"`, pool.Id))
	UpdateList(lunList)
	return lunList.(*LunList)
}

func (pool *Pool) CreateLun(name string, sizeGb uint32) (*Lun, error) {
	return CreateLun(pool.conn, pool, name, uint64(sizeGb))
}
