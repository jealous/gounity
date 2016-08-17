package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func mockConn() *Connection {
	return NewMockConnection("10.244.223.61", "osadmin", "Password123!")
}

func systemType() (*Type, error) {
	url := "https://1.1.1.1/api/types/system?compact=True"
	resp, _ := getMockResp(url, "")
	ret := &Type{}
	err := updateInstFromResp(resp, ret)
	return ret, err
}

func TestNewTypeFromResp(t *testing.T) {
	systemType, err := systemType()

	asserts := assert.New(t)
	asserts.Equal(systemType.Name, "system")
	asserts.Contains(systemType.Description, "Information about")
	asserts.Contains(systemType.Documentation, "system.html")
	asserts.Equal(len(systemType.Attributes), 14)
	asserts.Nil(err)
}

func TestType_AllFieldString(t *testing.T) {
	systemType, _ := systemType()
	assert.Contains(t, systemType.AllFieldString(), "id,health,name,model,")
}

func TestEmptyRscList(t *testing.T) {
	rscList := RscList{}
	assert.Equal(t, 0, rscList.Size())
}

func TestRsc_JsonId(t *testing.T) {
	rsc := Rsc{Id: "abc", type_: "pool"}
	expected := `{"id": "abc"}`
	assert.Equal(t, expected, rsc.JsonId())
}
