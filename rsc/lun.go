package rsc

import (
	log "github.com/Sirupsen/logrus"
)

type Lun struct {
	Rsc
	Name                  string
	Health                Health
	Description           string
	SizeTotal             uint64
	SizeUsed              uint64
	SizeAllocated         uint64
	IsThinEnabled         bool
	Wwn                   string
	MetaDataSize          uint64
	MetadataSizeAllocated uint64
	SnapWwn               string
	SnapsSize             uint64
	SnapsSizeAllocated    uint64
	SnapCount             uint32
	Pool                  Pool
	HostAccess            []BlockHostAccess
}

type LunList struct {
	RscList
}

func (ll *LunList) initRsc() Rscer {
	return &Lun{Rsc: Rsc{conn: ll.conn, type_: "lun"}}
}

type LunListCtor struct {
	conn *Connection
}

func (pi *LunListCtor) initList(filter string) RscLister {
	ret := &LunList{RscList: RscList{type_: "lun", conn: pi.conn}}
	ret.filter = filter
	return ret
}

func GetLunByName(conn *Connection, name string) *Lun {
	lun := getRscByName(name, &LunListCtor{conn})
	if lun == nil {
		return nil
	} else {
		return lun.(*Lun)
	}
}

func GetLunById(conn *Connection, id string) *Lun {
	lun := &Lun{Rsc: Rsc{conn: conn, type_: "lun", Id: id}}
	Update(lun)
	return lun
}

func GetLunList(conn *Connection) *LunList {
	return getRscList(&LunListCtor{conn}).(*LunList)
}

func CreateLun(conn *Connection, pool *Pool, name string, sizeGb uint64) (*Lun, error) {
	lunParam := map[string]interface{}{
		"pool": pool,
		"size": GbToByte(sizeGb),
	}
	body := map[string]interface{}{
		"lunParameters": lunParam,
		"name":          name,
	}
	resp, err := StorageRscDo(conn, "createLun", body)
	if err != nil {
		log.WithError(err).Error("failed to create lun.")
		return nil, err
	} else {
		id := getStorageRscId(resp)
		return GetLunById(conn, id), nil
	}
}

type HostLUNAccessEnum int

const (
	NO_ACCESS HostLUNAccessEnum = iota
	PRODUCTION
	SNAPSHOT
	BOTH
	MIXED
)

type BlockHostAccess struct {
	Host       Host
	AccessMask HostLUNAccessEnum
}

func (lun *Lun) AttachHost(host *Host) (*HostLUN, error) {
	hostAccess := []interface{}{
		map[string]interface{}{"host": host, "accessMask": PRODUCTION},
	}
	for _, access := range lun.HostAccess {
		item := map[string]interface{}{"host": &access.Host, "accessMask": access.AccessMask}
		hostAccess = append(hostAccess, item)
	}
	lunParam := *assembleLunParameter(&hostAccess)
	if _, err := StorageRscInstDo(lun.conn, lun.Id, "modifyLun", lunParam); err != nil {
		log.WithError(err).Error("attach host failed.")
		return nil, err
	}
	hostLun := GetHostLUN(lun.conn, host, lun)
	return hostLun, nil
}

func assembleLunParameter(hostAccess *[]interface{}) *map[string]interface{} {
	return &map[string]interface{}{
		"lunParameters": map[string]interface{}{
			"hostAccess": *hostAccess,
		},
	}
}

func (lun *Lun) DetachHost(host *Host) error {
	hostAccess := make([]interface{}, 0)
	for _, access := range lun.HostAccess {
		if access.Host.Id == host.Id {
			continue
		}
		item := map[string]interface{}{"host": &access.Host, "accessMask": access.AccessMask}
		hostAccess = append(hostAccess, item)
	}
	lunParam := *assembleLunParameter(&hostAccess)
	if _, err := StorageRscInstDo(lun.conn, lun.Id, "modifyLun", lunParam); err != nil {
		log.WithError(err).Error("detach host failed.")
		return err
	}
	return nil
}

func (lun *Lun) Delete() error {
	return DeleteLunById(lun.conn, lun.Id)
}

func DeleteLunById(conn *Connection, id string) error {
	_, err := conn.DeleteRscInst("storageResource", id, "")
	return err
}
