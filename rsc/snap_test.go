package rsc

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSnapStateEnum(t *testing.T) {
	assert.Equal(t, SnapStateEnum(3), SNAP_STATE_FAULTED)
	assert.Equal(t, SnapStateEnum(9), SNAP_STATE_DESTROYING)
}

func TestSnapAccessLevelEnum(t *testing.T) {
	assert.Equal(t, SnapAccessLevelEnum(0), SNAP_ACCESS_READ_ONLY)
	assert.Equal(t, SnapAccessLevelEnum(1), SNAP_ACCESS_READ_WRITE)
}

func VerifyVolcanSnap0(t *testing.T, snap *Snap) {
	asserts := assert.New(t)
	asserts.Equal("Volcan_snap0", snap.Name)
	asserts.Equal("38654705727", snap.Id)
	asserts.Equal(SNAP_STATE_READY, snap.State)
	asserts.Equal("The initial snap for Volcan.", snap.Description)
	asserts.Equal(2016, snap.CreationTime.Year())
	asserts.Equal(time.September, snap.CreationTime.Month())
	asserts.Equal(12, snap.CreationTime.Day())
	asserts.Equal(9, snap.CreationTime.Hour())
	asserts.Equal(44, snap.CreationTime.Minute())
	asserts.Equal(10, snap.CreationTime.Second())
	asserts.False(snap.IsSystemSnap)
	asserts.False(snap.IsModifiable)
	asserts.False(snap.IsReadOnly)
	asserts.False(snap.IsModified)
	asserts.True(snap.IsAutoDelete)
	asserts.Equal(uint64(5368709120), snap.Size)
	asserts.Equal("sv_22", snap.StorageResource.Id)
	asserts.Equal("sv_22", snap.Lun.Id)
}

func TestGetSnapByName(t *testing.T) {
	snap := GetSnapByName(MockConn(), "Volcan_snap0")
	VerifyVolcanSnap0(t, snap)
}

func TestGetSnapById(t *testing.T) {
	VerifyVolcanSnap0(t, GetSnapById(MockConn(), "38654705727"))
}

func TestGetSnapList(t *testing.T) {
	snaps := GetSnapList(MockConn())
	assert.Equal(t, 2, snaps.Size())
}

func TestCreateLunSnap(t *testing.T) {
	lun := GetLunByName(MockConn(), "Volcan")
	snap, err := CreateLunSnap(MockConn(), lun, "Volcan_snap1")
	Update(snap)
	assert.Equal(t, "Volcan_snap1", snap.Name)
	assert.Equal(t, "sv_22", snap.Lun.Id)
	assert.Nil(t, err)
}

func TestSnap_Delete(t *testing.T) {
	snap := GetSnapById(MockConn(), "38654705730")
	err := snap.Delete()
	assert.Nil(t, err)
}
