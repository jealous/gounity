package rsc

import (
	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLunList(t *testing.T) {
	lunList := GetLunList(MockConn())

	assert.Equal(t, 2, lunList.Size())

	for it := lunList.Iterator(); it.Next(); {
		lun := it.Value().(*Lun)
		if lun.Name == "lun1" {
			VerifyLun1(t, lun)
			break
		}
	}
}

func TestUpdatePoolProperty(t *testing.T) {
	lun := GetLunById(MockConn(), "sv_2")
	assert.Equal(t, "gounity", lun.Name)
	assert.Equal(t, 2, len(*lun.HostAccess))

	hostAccess := (*lun.HostAccess)[0]
	assert.Equal(t, SNAPSHOT, hostAccess.AccessMask)
	assert.Equal(t, "Host_5", hostAccess.Host.Id)
}

func VerifyLun1(t *testing.T, lun *Lun) {
	asserts := assert.New(t)

	asserts.Equal("lun1", lun.Name)
	asserts.Equal("sv_1", lun.Id)
	asserts.Equal("", lun.Description)
	asserts.Equal(uint64(5368709120), lun.SizeTotal)
	asserts.Equal(uint64(0), lun.SizeAllocated)
	asserts.Equal(uint64(0), lun.SizeUsed)
	asserts.True(lun.IsThinEnabled)
	asserts.Equal("60:06:01:60:15:E0:3A:00:6C:CC:AA:57:FE:07:BC:D3", lun.Wwn)
	asserts.Equal(uint64(3489660928), lun.MetaDataSize)
	asserts.Equal(uint64(2684354560), lun.MetadataSizeAllocated)
	asserts.Equal("60:06:01:60:15:E0:3A:00:CF:2E:61:29:07:10:46:83", lun.SnapWwn)
	asserts.Equal(uint64(0), lun.SnapsSize)
	asserts.Equal(uint64(0), lun.SnapsSizeAllocated)
	asserts.Equal(uint32(2), lun.SnapCount)
	asserts.Equal("pool_1", lun.Pool.Id)
	asserts.Equal("The LUN is operating normally. No action is required.",
		lun.Health.Description())
	asserts.Equal("sv_1", lun.StorageResource.Id)
}

func TestGetLunByName(t *testing.T) {
	lun := GetLunByName(MockConn(), "lun1")
	VerifyLun1(t, lun)
}

func TestCreateLun(t *testing.T) {
	conn := MockConn()
	pool := GetPoolByName(conn, "perfpool1132")
	lun, err := CreateLun(conn, pool, "gounity", 5)
	assert.Equal(t, "gounity", lun.Name)
	assert.Nil(t, err)
}

func TestCreateLun_NameUsed(t *testing.T) {
	conn := MockConn()
	pool := GetPoolByName(conn, "perfpool1132")
	lun, err := CreateLun(conn, pool, "openstack_dummy_lun", 5)
	assert.Nil(t, lun)
	assert.Contains(t, err.Error(), "LUN name has already been reserved")
}

func TestDeleteLunById(t *testing.T) {
	conn := MockConn()
	err := DeleteLunById(conn, "sv_4")
	assert.Nil(t, err)
}

func TestDeleteLunNotFound(t *testing.T) {
	conn := MockConn()
	err := DeleteLunById(conn, "sv_5")
	assert.Contains(t, err.Error(), "resource does not exist")
}

func TestLun_AttachHost(t *testing.T) {
	conn := MockConn()
	host := GetHostByName(conn, "gohost")
	lun := GetLunByName(conn, "gounity")
	hostLun, err := lun.AttachHost(host)
	assert.Equal(t, uint16(3), hostLun.Hlu)
	assert.Nil(t, err)
}

func TestLun_DetachHost(t *testing.T) {
	conn := MockConn()
	host := GetHostByName(conn, "gohost")
	lun := GetLunByName(conn, "gounity")
	err := lun.DetachHost(host)
	assert.Nil(t, err)
}

func TestLun_DetachAllHosts(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	conn := MockConn()
	lun := GetLunByName(conn, "golun")
	err := lun.DetachAllHosts()
	assert.Nil(t, err)
}

func TestGetHostLUNByLun(t *testing.T) {
	lun := GetLunByName(MockConn(), "golun")
	hostLunList := lun.GetHostLUN()
	assert.Equal(t, 1, hostLunList.Size())
	assert.Equal(t, uint16(1), hostLunList.Iterator().Value().(*HostLUN).Hlu)
}
