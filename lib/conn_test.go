package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testConn() *Connection {
	return NewMockConnection("10.244.223.61", "osadmin", "Password123!")
}

func unityConn() *Connection {
	return NewConnection("10.244.223.61", "osadmin", "Password123!")
}

func TestConnection_getAllUrl(t *testing.T) {
	u := unescape(testConn().getAllUrl("system", "").String())
	assert.Equal(t, "https://10.244.223.61/api/types/system/instances?"+
		"compact=true&fields=id,health,name,model,serialNumber,"+
		"internalModel,platform,macAddress,isEULAAccepted,isUpgradeComplete,"+
		"isAutoFailbackEnabled,currentPower,avgPower,supportedUpgradeModels", u)
}

func TestConnection_getAllUrlWithFilter(t *testing.T) {
	u := unescape(testConn().getAllUrl("system", `name eq "abc"`).String())
	assert.Equal(t, "https://10.244.223.61/api/types/system/instances?"+
		"compact=true&fields=id,health,name,model,serialNumber,"+
		"internalModel,platform,macAddress,isEULAAccepted,isUpgradeComplete,"+
		"isAutoFailbackEnabled,currentPower,avgPower,supportedUpgradeModels"+
		`&filter=name eq "abc"`, u)
}

func TestConnection_getInstUrl(t *testing.T) {
	u := unescape(testConn().getInstUrl("system", "0").String())
	assert.Equal(t, "https://10.244.223.61/api/instances/system/0?"+
		"compact=true&fields=id,health,name,model,serialNumber,internalModel,"+
		"platform,macAddress,isEULAAccepted,isUpgradeComplete,"+
		"isAutoFailbackEnabled,currentPower,avgPower,supportedUpgradeModels", u)
}

func TestConnection_getTypeUrl(t *testing.T) {
	assert.Equal(t, testConn().getTypeUrl("system").String(),
		"https://10.244.223.61/api/types/system?compact=true")
}

func TestConnection_GetActionUrl(t *testing.T) {
	assert.Equal(t, testConn().getActionUrl("storageResource", "createLun").String(),
		"https://10.244.223.61/api/types/storageResource/action/createLun?compact=true")
}
