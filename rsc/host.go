package rsc

import (
	"fmt"
	"regexp"
	"strings"
)

type HostTypeEnum int

const (
	Unknown HostTypeEnum = iota
	HostManual
	Subnet
	NetGroup
	RPA
	HostAuto
	VNXSanCopy HostTypeEnum = 255
)

type Host struct {
	Rsc
	Name                string
	Health              Health
	Description         string
	OsType              string
	Type                HostTypeEnum
	HostUUID            string
	HostPushedUUID      string
	HostPolledUUID      string
	FcHostInitiators    []HostInitiator
	IscsiHostInitiators []HostInitiator
	HostIPPorts         []HostIPPort
	StorageResources    []StorageResource
	HostLUNs            []HostLUN
	Tenant              Tenant
}

type HostList struct {
	RscList
}

func (hl *HostList) initRsc() Rscer {
	return &Host{Rsc: Rsc{conn: hl.conn, type_: "host"}}
}

type HostListCtor struct {
	conn *Connection
}

func (host *Host) AddInitiator(uid string) (*HostInitiator, error) {
	return CreateInitiator(host.conn, host, uid)
}

func (host *Host) Attach(lun *Lun) (*HostLUN, error) {
	hostLun, err := lun.AttachHost(host)
	if err == nil {
		Update(host)
	}
	return hostLun, err
}

func (host *Host) Detach(lun *Lun) error {
	err := lun.DetachHost(host)
	if err == nil {
		Update(host)
	}
	return err
}

func (hl *HostListCtor) initList(filter string) RscLister {
	ret := &HostList{RscList: RscList{type_: "host", conn: hl.conn}}
	ret.filter = filter
	return ret
}

func GetHostByName(conn *Connection, name string) *Host {
	host := getRscByName(name, &HostListCtor{conn})
	if host == nil {
		return nil
	} else {
		return host.(*Host)
	}
}

func GetHostById(conn *Connection, id string) *Host {
	host := &Host{Rsc: Rsc{conn: conn, type_: "host", Id: id}}
	Update(host)
	return host
}

func GetHostList(conn *Connection) *HostList {
	return getRscList(&HostListCtor{conn}).(*HostList)
}

func CreateHost(conn *Connection, name string) (*Host, error) {
	body := map[string]interface{}{
		"type": HostManual,
		"name": name,
	}
	if id, err := conn.PostInstance("host", makeBody(body)); err != nil {
		return nil, err
	} else {
		return GetHostById(conn, id), nil
	}
}

func (host *Host) Delete() error {
	return DeleteHostById(host.conn, host.Id)
}

func DeleteHostById(conn *Connection, id string) error {
	_, err := conn.DeleteRscInst("host", id, "")
	return err
}

type HostInitiatorType int

const (
	UNKNOWN HostInitiatorType = iota
	FC
	ISCSI
)

type HostIPPort struct {
	Rsc
	Name           string
	Type           HostPortTypeEnum
	Address        string
	Netmask        string
	V6PrefixLength uint32
	IsIgnored      bool
	Host           Host
}

type HostLUNTypeEnum int

const (
	HostLUNType_UNKNOWN HostLUNTypeEnum = iota
	HostLUNType_LUN
	HostLUNType_Snap
)

type HostLUN struct {
	Rsc
	Host Host
	Type HostLUNTypeEnum
	Hlu  uint16
	Lun  Lun
	// todo: snap
	//Snap Snap
	IsReadOnly    bool
	IsDefaultSnap bool
}

type HostLUNList struct {
	RscList
}

func (hll *HostLUNList) initRsc() Rscer {
	return &HostLUN{Rsc: Rsc{conn: hll.conn, type_: "hostLUN"}}
}

type HostLunListCtor struct {
	conn *Connection
}

func (hll *HostLunListCtor) initList(filter string) RscLister {
	ret := &HostLUNList{RscList: RscList{type_: "hostLUN", conn: hll.conn}}
	ret.filter = filter
	return ret
}

func GetHostLUN(conn *Connection, host *Host, lun *Lun) *HostLUN {
	ctor := &HostLunListCtor{conn}
	list := ctor.initList(fmt.Sprintf(`host eq "%v" and lun eq "%v"`, host.Id, lun.Id))
	UpdateList(list)
	hi := list.Iterator().Value()
	if hi == nil {
		return nil
	} else {
		return hi.(*HostLUN)
	}
}

type Tenant struct {
	Rsc
	Name  string
	Uuid  string
	Vlans []uint32
	Host  []Host
}

type HostPortTypeEnum int

const (
	HostPortType_IPv4 HostPortTypeEnum = iota
	HostPortType_IPv6
	HostPortType_NetworkName
)

type HostInitiator struct {
	Rsc
	Health      Health
	Type        HostInitiatorType
	InitiatorId string
	IsIgnored   bool
	NodeWWN     string
	PortWWN     string
	ParentHost  Host
	Paths       []HostInitiatorPath
	IsBound     bool
}

type HostInitiatorList struct {
	RscList
}

func (hil *HostInitiatorList) initRsc() Rscer {
	return &HostInitiator{Rsc: Rsc{conn: hil.conn, type_: "hostInitiator"}}
}

type HostInitiatorListCtor struct {
	conn *Connection
}

func (hil *HostInitiatorListCtor) initList(filter string) RscLister {
	ret := &HostInitiatorList{RscList: RscList{type_: "hostInitiator", conn: hil.conn}}
	ret.filter = filter
	return ret
}

func GetInitiatorByUid(conn *Connection, uid string) *HostInitiator {
	ctor := &HostInitiatorListCtor{conn}
	list := ctor.initList(fmt.Sprintf(`initiatorId eq "%v"`, uid))
	UpdateList(list)
	hi := list.Iterator().Value()
	if hi == nil {
		return nil
	} else {
		return hi.(*HostInitiator)
	}
}

func GetInitiatorById(conn *Connection, id string) *HostInitiator {
	initiator := &HostInitiator{Rsc: Rsc{conn: conn, type_: "hostInitiator", Id: id}}
	Update(initiator)
	return initiator
}

func CreateInitiator(conn *Connection, host *Host, uid string) (*HostInitiator, error) {
	initiator := GetInitiatorByUid(conn, uid)
	if initiator == nil {
		initiatorType := getInitiatorType(uid)
		body := map[string]interface{}{
			"host":              host,
			"initiatorType":     initiatorType,
			"initiatorWWNorIqn": uid,
		}
		if id, err := conn.PostInstance("hostInitiator", makeBody(body)); err != nil {
			return nil, err
		} else {
			initiator = GetInitiatorById(conn, id)
		}
	}
	return initiator, nil
}

func (initiator *HostInitiator) Delete() error {
	return DeleteInitiatorById(initiator.conn, initiator.Id)
}

func DeleteInitiatorById(conn *Connection, id string) error {
	_, err := conn.DeleteRscInst("hostInitiator", id, "")
	return err
}

var wwnRe *regexp.Regexp = regexp.MustCompile("(\\w{2}:){15}\\w{2}")

func getInitiatorType(uid string) HostInitiatorType {
	if wwnRe.MatchString(uid) {
		return FC
	} else if strings.HasPrefix(uid, "iqn.") {
		return ISCSI
	}
	return UNKNOWN
}

type HostInitiatorPathTypeEnum int

const (
	Manual HostInitiatorPathTypeEnum = iota
	ESX_Auto
	Other_Auto
)

type HostInitiatorPath struct {
	Rsc
	RegistrationType HostInitiatorPathTypeEnum
	IsLoggedIn       bool
	HostPushName     string
	SessionIds       []string
	Initiator        HostInitiator
	FcPort           FcPort
	IscsiPortal      IscsiPortal
}

type IpProtocolVersionEnum int

const (
	IPv4 IpProtocolVersionEnum = 4
	IPv6 IpProtocolVersionEnum = 6
)

type IscsiPortal struct {
	Rsc
	EthernetPort      EthernetPort
	IscsiNode         IscsiNode
	IpAddress         string
	Netmask           string
	V6PrefixLength    uint32
	Gateway           string
	VlanId            uint32
	IpProtocolVersion IpProtocolVersionEnum
}

type EthernetPort struct {
	Rsc
	Health           Health
	StorageProcessor StorageProcessor
	Name             string
	PortNumber       uint32
	Bond             bool
	IsLinkUp         bool
	MacAddress       string
	IsRDMACapable    bool
}

type IscsiNode struct {
	Rsc
	Name         string
	EthernetPort EthernetPort
	Alias        string
}

type FcPort struct {
	Rsc
	Health           Health
	SlotNumber       uint32
	Wwn              string
	StorageProcessor StorageProcessor
	NPortId          uint32
	name             string
}

type SPModelNameEnum int

const (
	SP300 SPModelNameEnum = 1 + iota
	SP400
	SP500
	SP600
)

type StorageProcessor struct {
	Rsc
	Health               Health
	IsRescueMode         bool
	Model                string
	slotNumber           uint32
	Name                 string
	EmcPartNumber        string
	EmcSerialNumber      string
	Manufacturer         string
	VendorPartNumber     string
	VendorSerialNumber   string
	SasExpanderVersion   string
	BiosFirmwareRevision string
	PostFirmwareRevision string
	MemorySize           uint32
	ModelName            SPModelNameEnum
}

type StorageResourceTypeEnum int

const (
	SRType_Filesystem StorageResourceTypeEnum = 1 + iota
	SRType_ConsistencyGroup
	SRType_VmwareFs
	SRType_VmwareIscsi
)

const (
	SRType_Lun StorageResourceTypeEnum = 8 + iota
	SRType_VVolDataStoreFs
	SRType_VVolDataStoreIscsi
)

type StorageResource struct {
	Rsc
	Health                Health
	Name                  string
	Description           string
	SizeTotal             uint64
	SizeUsed              uint64
	SizeAllocated         uint64
	MetadataSize          uint64
	MetadataSizeAllocated uint64
	SnapsSizeTotal        uint64
	SnapsSizeAllocated    uint64
	SnapCount             uint32
	Pools                 []Pool
	Luns                  []Lun
}
