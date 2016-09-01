package rsc

type Unity struct {
	Rsc
	Name                  string
	Model                 string
	Health                Health
	SerialNumber          string
	Platform              string
	MacAddress            string
	IsEULAAccepted        bool
	IsUpgradeComplete     bool
	IsAutoFailbackEnabled bool
	CurrentPower          int
	AvgPower              int
}

func NewUnity(ip, username, password string) *Unity {
	conn := NewConnection(ip, username, password)
	return NewUnityByConn(conn)
}

func NewUnityByConn(conn *Connection) *Unity {
	return &Unity{Rsc: Rsc{conn: conn, type_: "system", Id: "0"}}
}

func (unity *Unity) Ip() string {
	return unity.conn.ip
}

func (unity *Unity) Username() string {
	return unity.GetConn().username
}

func (unity *Unity) GetPoolList() *PoolList {
	return GetPoolList(unity.GetConn())
}

func (unity *Unity) GetPoolByName(name string) *Pool {
	return GetPoolByName(unity.GetConn(), name)
}

func (unity *Unity) GetPoolById(id string) *Pool {
	return GetPoolById(unity.GetConn(), id)
}

func (unity *Unity) GetLunList() *LunList {
	return GetLunList(unity.GetConn())
}

func (unity *Unity) GetLunById(id string) *Lun {
	return GetLunById(unity.GetConn(), id)
}

func (unity *Unity) GetLunByName(name string) *Lun {
	return GetLunByName(unity.GetConn(), name)
}

func (unity *Unity) GetHostList() *HostList {
	return GetHostList(unity.GetConn())
}

func (unity *Unity) GetHostById(id string) *Host {
	return GetHostById(unity.GetConn(), id)
}

func (unity *Unity) GetHostByName(name string) *Host {
	return GetHostByName(unity.GetConn(), name)
}

func (unity *Unity) CreateHost(name string) (*Host, error) {
	return CreateHost(unity.GetConn(), name)
}
