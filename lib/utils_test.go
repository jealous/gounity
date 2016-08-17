package lib

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCurrentFolder(t *testing.T) {
	folder := currentFolder()
	assert.Contains(t, folder, "gounity")
	assert.Contains(t, folder, "lib")
}

func TestTestFolder(t *testing.T) {
	folder := testFolder()
	assert.Contains(t, folder, "gounity")
	assert.Contains(t, folder, "mocks")
}

func TestMockFolder(t *testing.T) {
	folder := mockFolder("system")
	assert.Contains(t, folder, "gounity")
	assert.Contains(t, folder, "system")
}

func TestIndexFilename(t *testing.T) {
	file := indexFilename("system")
	assert.Contains(t, file, "system")
	assert.Contains(t, file, "index.json")
}

func TestSrcFromUrl_types(t *testing.T) {
	url := "https://1.1.1.1/api/types/system/instances"
	assert.Equal(t, rscFromUrl(url), "system")
}

func TestSrcFromUrl_instances(t *testing.T) {
	url := "https://1.1.1.1/api/instances/system/0"
	assert.Equal(t, rscFromUrl(url), "system")
}

func TestSrcFromUrl_variable(t *testing.T) {
	url := "https://1.1.1.1/api/types/system?compact=True"
	assert.Equal(t, rscFromUrl(url), "system")
}

func TestSrcFromUrl_meta(t *testing.T) {
	url := "https://1.1.1.1/api/types/system"
	assert.Equal(t, rscFromUrl(url), "system")
}

func TestSrcFromUrl_invalid(t *testing.T) {
	url := "https://1.1.1.1/apidocs"
	assert.Equal(t, rscFromUrl(url), "")
}

func TestRemoveIp(t *testing.T) {
	url := "https://1.1.1.1/apidocs/abc"
	assert.Equal(t, removeIp(url), "/apidocs/abc")
}

func TestGetRespFilename(t *testing.T) {
	url := "https://1.1.1.1/api/types/system?compact=True"
	assert.Equal(t, getRespFilename(url, ""), "type.json")
}

func TestGetResp(t *testing.T) {
	url := "https://1.1.1.1/api/types/system?compact=True"
	resp, _ := getMockResp(url, "")
	assert.Contains(t, resp,
		"Unique identifier of the system instance.")
}

func TestMakeBody_plainMap(t *testing.T) {
	outerBody := make(map[string]interface{})
	innerBody := make(map[string]interface{})
	pool := make(map[string]interface{})
	pool["id"] = "pool_1"
	innerBody["pool"] = pool
	innerBody["size"] = 1234567
	outerBody["param"] = innerBody
	outerBody["name"] = "test"
	bs := makeBody(outerBody)

	expected := `{"name":"test","param":{"pool":{"id":"pool_1"},"size":1234567}}`
	assert.Equal(t, expected, string(bs))
}

func TestMakeBody_rscWithId(t *testing.T) {
	body := make(map[string]interface{})
	body["pool"] = &Pool{Rsc: Rsc{Id: "pool_1"}}

	expected := `{"pool":{"id":"pool_1"}}`
	assert.Equal(t, expected, string(makeBody(body)))
}

func TestMakeBody_rscLister(t *testing.T) {
	body := make(map[string]interface{})
	arr := []Rscer{&Pool{Rsc: Rsc{Id: "pool_1"}}, &Pool{Rsc: Rsc{Id: "pool_2"}}}
	body["pools"] = &PoolList{RscList: RscList{list_: arr}}

	expected := `{"pools":[{"id":"pool_1"},{"id":"pool_2"}]}`
	assert.Equal(t, expected, string(makeBody(body)))
}

func TestMakeBody_nestedRscer(t *testing.T) {
	host := GetHostById(mockConn(), "Host_6")
	hostAccess := []interface{}{
		map[string]interface{}{"host": host, "mask": PRODUCTION},
	}
	lunParam := map[string]interface{}{
		"param": map[string]interface{}{
			"hostAccess": hostAccess,
		},
	}
	assert.Contains(t, string(makeBody(lunParam)), `{"id":"Host_6"}`)
}
