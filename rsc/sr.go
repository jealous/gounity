package rsc

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/url"
)

func StorageRscDo(conn *Connection, action string, body map[string]interface{}) (string, error) {
	_url := conn.getActionUrl("storageResource", action)
	return _storageRscDo(conn, _url, body)
}

func StorageRscInstDo(conn *Connection, id, action string, body map[string]interface{}) (string, error) {
	_url := conn.getInstActionUrl("storageResource", id, action)
	return _storageRscDo(conn, _url, body)
}

func _storageRscDo(conn *Connection, url *url.URL, body map[string]interface{}) (string, error) {
	strUrl := url.String()
	strBody := makeBody(body)
	log.WithField("url", strUrl).WithField("body", strBody).Info(
		"do storage resource action.")
	return conn.Post(strUrl, strBody)
}

type StorageResourceIdResp struct {
	StorageResource Rsc
}

func getStorageRscId(resp string) string {
	content, err := getContentFromResp(resp)
	if err != nil {
		return ""
	}
	srResp := &StorageResourceIdResp{}
	err = json.Unmarshal(content, srResp)
	if err != nil {
		log.WithError(err).Error("failed to parse response body to id.")
		return ""
	} else {
		return srResp.StorageResource.Id
	}
}
