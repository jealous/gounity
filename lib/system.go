package lib

type System struct {
	ip       string
	username string
	password string
}

func NewSystem(ip, username, password string) *System {
	ret := &System{}
	ret.ip = ip
	ret.username = username
	ret.password = password
	return ret
}

func (sys *System) Ip() string {
	return sys.ip
}

func (sys *System) Username() string {
	return sys.username
}
