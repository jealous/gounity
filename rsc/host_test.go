package rsc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHostRelatedEnumValues(t *testing.T) {
	assert.Equal(t, SPModelNameEnum(3), SP500)
	assert.Equal(t, IpProtocolVersionEnum(6), IPv6)
	assert.Equal(t, StorageResourceTypeEnum(9), SRType_VVolDataStoreFs)
}

func TestGetHostList(t *testing.T) {
	hosts := GetHostList(MockConn())
	assert.Equal(t, 7, hosts.Size())
	for it := hosts.Iterator(); it.Next(); {
		host := it.Value().(*Host)
		if host.Id == "Host_6" {
			VerifyHost6(t, host)
		}
	}
}

func TestGetHostByName(t *testing.T) {
	host := GetHostByName(MockConn(), "gohost")
	VerifyHost6(t, host)
}

func TestGetHostById(t *testing.T) {
	host, err := GetHostById(MockConn(), "Host_6")
	VerifyHost6(t, host)
	assert.Nil(t, err)
}

func VerifyHost6(t *testing.T, host *Host) {
	asserts := assert.New(t)
	asserts.Equal("Host_6", host.Id)
	asserts.Equal("gohost", host.Name)
	asserts.Equal("", host.Description)
	asserts.Equal("Windows Client", host.OsType)
	asserts.Equal(HostManual, host.Type)
	asserts.Equal(1, len(host.IscsiHostInitiators))
	asserts.Equal("HostInitiator_10", host.IscsiHostInitiators[0].Id)
	asserts.Equal(1, len(host.HostIPPorts))
	asserts.Equal("HostNetworkAddress_3", host.HostIPPorts[0].Id)
	asserts.Equal(4, len(host.HostLUNs))
	asserts.Equal("", host.Tenant.Id)
	asserts.Equal(0, len(host.StorageResources))
}

func TestCreateHost(t *testing.T) {
	conn := MockConn()
	host, err := CreateHost(conn, "gohost")
	VerifyHost6(t, host)
	assert.Nil(t, err)
}

func TestDeleteHost(t *testing.T) {
	conn := MockConn()
	err := DeleteHostById(conn, "Host_12")
	assert.Nil(t, err)
}

func TestDeleteHost_notFound(t *testing.T) {
	conn := MockConn()
	err := DeleteHostById(conn, "Host_7")
	assert.Contains(t, err.Error(), "does not exist")
}

func TestGetInitiatorByUid(t *testing.T) {
	uid := "20:00:00:90:FA:53:41:40:10:00:00:90:FA:53:41:40"
	hi := GetInitiatorByUid(MockConn(), uid)
	VerifyHostInitiator2(t, hi)
}

func VerifyHostInitiator2(t *testing.T, hi *HostInitiator) {
	asserts := assert.New(t)
	asserts.Equal("HostInitiator_2", hi.Id)
	asserts.Equal(FC, hi.Type)
	asserts.Equal("20:00:00:90:FA:53:41:40:10:00:00:90:FA:53:41:40", hi.InitiatorId)
	asserts.False(hi.IsIgnored)
	asserts.Equal("20:00:00:90:FA:53:41:40", hi.NodeWWN)
	asserts.Equal("10:00:00:90:FA:53:41:40", hi.PortWWN)
	asserts.False(hi.IsBound)
	asserts.Equal("Host_2", hi.ParentHost.Id)
	asserts.Equal("HostInitiator_2_02:00:00:04", hi.Paths[0].Id)
}

func TestGetInitiatorByUid_notFound(t *testing.T) {
	uid := "20:00:00:90:FA:53:41:40:10:00:00:90:FA:AA:AA:AA"
	hi := GetInitiatorByUid(MockConn(), uid)
	assert.Nil(t, hi)
}

func TestGetInitiatorType_FC(t *testing.T) {
	uid := "20:00:00:90:FA:53:41:40:10:00:00:90:FA:53:41:40"
	assert.Equal(t, FC, getInitiatorType(uid))
}

func TestGetInitiatorType_ISCSI(t *testing.T) {
	uid := "iqn.1993-08.org.debian:01:a4f95ed14d65"
	assert.Equal(t, ISCSI, getInitiatorType(uid))
}

func TestCreateInitiator_existed(t *testing.T) {
	uid := "20:00:00:90:FA:53:41:40:10:00:00:90:FA:53:41:40"
	host := &Host{Rsc: Rsc{Id: "Host_2"}}
	hi, err := CreateInitiator(MockConn(), host, uid)
	assert.Nil(t, err)
	VerifyHostInitiator2(t, hi)
}

func TestCreateInitiator(t *testing.T) {
	uid := "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:66:77:88:99"
	host := &Host{Rsc: Rsc{Id: "Host_13"}}
	hi, err := CreateInitiator(MockConn(), host, uid)
	assert.Nil(t, err)
	assert.Equal(t, uid, hi.InitiatorId)
}

func TestCreateInitiator_hostNotFound(t *testing.T) {
	uid := "AA:BB:CC:DD:EE:FF:00:11:22:33:44:55:88:88:88:88"
	host := &Host{Rsc: Rsc{Id: "Host_999"}}
	hi, err := CreateInitiator(MockConn(), host, uid)
	assert.Nil(t, hi)
	assert.Contains(t, err.Error(), "Cannot find the specified host")
}

func TestDeleteInitiatorById(t *testing.T) {
	err := DeleteInitiatorById(MockConn(), "HostInitiator_49")
	assert.Nil(t, err)
}

func TestGetHostLUN(t *testing.T) {
	hostLun := GetHostLUN(MockConn(), &Host{Rsc: Rsc{Id: "Host_5"}}, &Lun{Rsc: Rsc{Id: "sv_2"}})
	asserts := assert.New(t)
	asserts.Equal("Host_5_sv_2_snap", hostLun.Id)
	asserts.Equal(HostLUNType_Snap, hostLun.Type)
	asserts.Equal(uint16(0), hostLun.Hlu)
	asserts.Equal("sv_2", hostLun.Lun.Id)
	asserts.Equal("Host_5", hostLun.Host.Id)
	asserts.False(hostLun.IsReadOnly)
	asserts.True(hostLun.IsDefaultSnap)
}

func TestGetHostLUN_notFound(t *testing.T) {
	hostLun := GetHostLUN(MockConn(), &Host{Rsc: Rsc{Id: "Host_5"}}, &Lun{Rsc: Rsc{Id: "sv_99"}})
	assert.Nil(t, hostLun)
}

func TestGetHostLUNList(t *testing.T) {
	host := GetHostByName(MockConn(), "gohost")
	hostLunList := host.GetHostLUNList()
	assert.Equal(t, 2, hostLunList.Size())
	for it := hostLunList.Iterator(); it.Next(); {
		hostLUN := it.Value().(*HostLUN)
		assert.Equal(t, "Host_13", hostLUN.Host.Id)
	}
}

func TestGetIscsiPortalById(t *testing.T) {
	portal, err := GetIscsiPortalById(MockConn(), "if_4")
	asserts := assert.New(t)
	asserts.Equal("if_4", portal.Id)
	asserts.Equal(IPv4, portal.IpProtocolVersion)
	asserts.Equal("10.244.213.177", portal.IpAddress)
	asserts.Equal("255.255.255.0", portal.Netmask)
	asserts.Equal("10.244.213.1", portal.Gateway)
	asserts.Equal("spa_eth2", portal.EthernetPort.Id)
	asserts.Nil(err)
}

func TestGetIscsiPortalById_notFound(t *testing.T) {
	portal, err := GetIscsiPortalById(MockConn(), "if_13")
	asserts := assert.New(t)
	asserts.Nil(portal)
	restError := err.(*RestError)
	asserts.Equal(uint64(0x7d13005), restError.ErrorCode)
	asserts.Equal(404, restError.HttpStatusCode)
}

func TestDeleteIscsiPortal(t *testing.T) {
	portal, _ := GetIscsiPortalById(MockConn(), "if_5")
	err := portal.Delete()
	assert.Nil(t, err)
}

func TestCreateIscsiPortal(t *testing.T) {
	conn := MockConn()
	port, _ := GetEthernetPortById(conn, "spa_eth3")
	portal, err := CreateIscsiPortal(conn, port, "10.244.213.179", "255.255.255.0", "10.244.213.1")
	asserts := assert.New(t)
	asserts.Nil(err)
	asserts.Equal("if_5", portal.Id)
	asserts.Equal("10.244.213.179", portal.IpAddress)
}

func TestGetIscsiPortalList(t *testing.T) {
	portals := GetIscsiPortalList(MockConn())
	assert.Equal(t, 2, portals.Size())
}
