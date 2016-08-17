package lib

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type index struct {
	Url      string                 `json:"url"`
	Body     map[string]interface{} `json:"body"`
	Response string                 `json:"response"`
}

type indices struct {
	Indices []index `json:"indices"`
}

func currentFolder() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func testFolder() string {
	return filepath.Join(currentFolder(), "..", "mocks")
}

func mockFolder(rsc string) string {
	return filepath.Join(testFolder(), rsc)
}

func indexFilename(rsc string) string {
	return filepath.Join(mockFolder(rsc), "index.json")
}

func rscFromUrl(url string) string {
	tokens := strings.Split(url, "/")
	if len(tokens) < 5 {
		log.WithField("url", url).Error(
			"cannot find resource type from url.")
		return ""
	}
	return strings.Split(tokens[5], "?")[0]
}

func removeIp(url string) string {
	return "/" + strings.Join(strings.Split(url, "/")[3:], "/")
}

func getRespFilename(url string, body string) string {
	ret := ""
	filename := indexFilename(rscFromUrl(url))
	strIndices, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("filename", filename).Panic(
			"failed to read index file.")
	} else {
		decoder := json.NewDecoder(strings.NewReader(string(strIndices)))
		decoder.UseNumber()

		var idc indices
		if err = decoder.Decode(&idc); err != nil {
			log.WithField("filename", filename).Panic(
				"failed to parse index file.")
		}
		url = unescape(removeIp(url))
		for i := range idc.Indices {
			idx := idc.Indices[i]
			if strings.EqualFold(idx.Url, url) {
				bodyBytes, err := json.Marshal(idx.Body)
				if err != nil {
					log.WithField("body", idx.Body).Error("failed to marshall index body.")
				}
				definedBody := string(bodyBytes)
				if definedBody == "null" {
					definedBody = ""
				}
				if strings.EqualFold(definedBody, body) {
					ret = idx.Response
					break
				} else {
					log.WithField("definedBody", definedBody).WithField("body", body).Info("body not equal.")
				}
			}
		}
	}
	if ret == "" {
		log.WithField("url", unescape(url)).WithField("body", body).Panic(
			"cannot find response for url.")
	}
	return ret
}

func getMockResp(url string, body string) (string, error) {
	filename := getRespFilename(url, body)
	filename = filepath.Join(mockFolder(rscFromUrl(url)), filename)
	log.WithField("filename", filename).Info("read mock file.")
	resp, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithError(err).WithField("file", filename).Error(
			"failed to read resp file.")
	}
	return string(resp), err
}

func getStatusCode(resp string) int {
	var code int = 200
	re := regexp.MustCompile("\"httpStatusCode\": (\\d+)")
	if match := re.FindSubmatch([]byte(resp)); match != nil {
		if i, err := strconv.Atoi(string(match[1])); err != nil {
			log.WithError(err).Error("failed to convert str to int.")
		} else {
			code = i
		}
	}
	return code
}

func unescape(u string) string {
	ret, err := url.QueryUnescape(u)
	if err != nil {
		log.WithError(err).Error("failed to unescape url.")
	}
	return ret
}

func GbToByte(gb uint64) uint64 {
	return gb * 1024 * 1024 * 1024
}

func cvtBodyType(body interface{}) interface{} {
	var ret interface{}
	switch t := body.(type) {
	default:
		ret = t
	case Rscer:
		ret = map[string]string{"id": t.GetId()}
	case RscLister:
		arr := make([]interface{}, t.Size())
		for it := t.Iterator(); it.Next(); {
			v := it.Value()
			arr[it.NextIndex()-1] = cvtBodyType(v)
		}
		ret = arr
	case map[string]interface{}:
		strMap := make(map[string]interface{})
		for k, v := range t {
			strMap[string(k)] = cvtBodyType(v)
		}
		ret = strMap
	case []interface{}:
		arr := make([]interface{}, len(t))
		for k, v := range t {
			arr[k] = cvtBodyType(v)
		}
		ret = arr
	}
	return ret
}

func makeBody(body interface{}) string {
	bytes, err := json.Marshal(cvtBodyType(body).(map[string]interface{}))
	if err != nil {
		log.WithError(err).Error("failed to marshal json body.")
	}
	return string(bytes)
}
