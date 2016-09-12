package rsc

import (
	log "github.com/Sirupsen/logrus"
	"time"
)

type SnapStateEnum int

const (
	SNAP_STATE_READY SnapStateEnum = iota + 2
	SNAP_STATE_FAULTED
	SNAP_STATE_OFFLINE SnapStateEnum = iota + 4
	SNAP_STATE_INVALID
	SNAP_STATE_INITIALIZING
	SNAP_STATE_DESTROYING
)

type Snap struct {
	Rsc
	Name             string
	Description      string
	StorageResource  *StorageResource
	Lun              Lun
	SnapGroup        *Snap
	ParentSnap       *Snap
	CreationTime     *time.Time
	ExpirationTime   *time.Time
	LastRefreshTime  *time.Time
	LastWritableTime *time.Time
	IsSystemSnap     bool
	IsModifiable     bool
	AttachedWWN      string
	IsReadOnly       bool
	IsModified       bool
	IsAutoDelete     bool
	State            SnapStateEnum
	Size             uint64
	HostAccess       *[]SnapHostAccess
}

type SnapAccessLevelEnum int

const (
	SNAP_ACCESS_READ_ONLY SnapAccessLevelEnum = iota
	SNAP_ACCESS_READ_WRITE
)

type SnapHostAccess struct {
	Host          *Host
	AllowedAccess SnapAccessLevelEnum
}

type SnapList struct {
	RscList
}

func (sl *SnapList) initRsc() Rscer {
	return &Snap{Rsc: Rsc{conn: sl.conn, type_: "snap"}}
}

type SnapListCtor struct {
	conn *Connection
}

func (pi *SnapListCtor) initList(filter string) RscLister {
	ret := &SnapList{RscList: RscList{type_: "snap", conn: pi.conn}}
	ret.filter = filter
	return ret
}

func GetSnapByName(conn *Connection, name string) *Snap {
	snap := getRscByName(name, &SnapListCtor{conn})
	if snap == nil {
		return nil
	} else {
		return snap.(*Snap)
	}
}

func (snap *Snap) Delete() error {
	return DeleteSnapById(snap.conn, snap.Id)
}

func GetSnapById(conn *Connection, id string) *Snap {
	snap := &Snap{Rsc: Rsc{conn: conn, type_: "snap", Id: id}}
	Update(snap)
	return snap
}

func GetSnapList(conn *Connection) *SnapList {
	return getRscList(&SnapListCtor{conn}).(*SnapList)
}

func CreateLunSnap(conn *Connection, lun *Lun, name string) (*Snap, error) {
	param := map[string]interface{}{
		// the storage resource is is the same as the lun id,
		// so we could safely use lun here.
		"storageResource": lun,
		"name":            name,
		"description":     "created by gounity",
		"isAutoDelete":    true,
		"isReadOnly":      false,
	}
	if id, err := conn.PostInstance("snap", makeBody(param)); err != nil {
		log.WithError(err).WithField("name", name).Error("create snap failed.")
		return nil, err
	} else {
		log.WithField("name", name).Info("create snap success.")
		return GetSnapById(conn, id), nil
	}
}

func DeleteSnapById(conn *Connection, id string) error {
	log.SetLevel(log.DebugLevel)
	log.WithField("id", id).Info("delete snap.")
	_, err := conn.DeleteRscInst("snap", id, "")
	return err
}
