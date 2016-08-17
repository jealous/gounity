package lib

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

var HEADERS map[string]string = map[string]string{
	"Accept":            "application/json",
	"Content-Type":      "application/json",
	"Accept_Language":   "en_US",
	"X-EMC-REST-CLIENT": "true",
	"User-Agent":        "gounity",
}

type Connection struct {
	ip       string
	username string
	password string
	client   *http.Client
	fields   map[string]string
	useMock  bool
	csrf     string
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	// if redirected , use the header of the first request
	log.WithField("request", req).Debug("request redirected.")
	req.Header = via[0].Header
	return nil
}

func transport() *http.Transport {
	return &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
}

func cookieJar() *cookiejar.Jar {
	options := cookiejar.Options{PublicSuffixList: publicsuffix.List}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Error(err)
	}
	return jar
}

func NewConnection(ip, username, password string) *Connection {
	c := &Connection{ip, username, password, nil, make(map[string]string), false, ""}
	return c.init()
}

func NewMockConnection(ip, username, password string) *Connection {
	c := NewConnection(ip, username, password)
	c.useMock = true
	return c
}

func (conn *Connection) init() *Connection {
	conn.client = &http.Client{
		Transport: transport(), Jar: cookieJar(), CheckRedirect: redirectPolicyFunc}
	conn.fields["type"] = ""
	return conn
}

func (conn *Connection) request(url, body, method string) (*http.Response, error) {
	if conn.useMock {
		body, err := getMockResp(url, body)
		code := getStatusCode(body)
		mockResp := &http.Response{Body: nopCloser{strings.NewReader(body)}, StatusCode: code}
		return mockResp, err
	}

	req, err := conn.newRequest(url, body, method)
	if err != nil {
		return nil, err
	}

	var resp *http.Response
	resp, err = conn.do(req)
	if resp.StatusCode == 401 && method != "GET" {
		req, err := conn.newRequest(url, body, method)
		if err != nil {
			return nil, err
		}
		resp, err = conn.retryWithCsrfToken(req)
	}
	return resp, err
}

func extractRestError(resp *http.Response) error {
	var err error = nil
	if resp.StatusCode >= 400 {
		err = &RestError{}
		bytes, _ := getErrorFromResp(getRespBody(resp))
		unmarshalErr := json.Unmarshal(bytes, err)
		if unmarshalErr != nil {
			log.Error("failed to unmartial response to a RestError.")
			err = unmarshalErr
		} else {
			log.WithError(err).Warn("REST error met.")
		}
	}
	return err
}

func (conn *Connection) newRequest(url, body, method string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		log.WithError(err).Error("create request error.")
		return nil, err
	}
	req.SetBasicAuth(conn.username, conn.password)
	for k, v := range HEADERS {
		req.Header.Add(k, v)
	}
	return req, err
}

func (conn *Connection) retryWithCsrfToken(req *http.Request) (*http.Response, error) {
	var (
		resp *http.Response
		err  error
	)

	log.Info("token invalid, try to get a new token.")
	pathUser := conn.getAllUrl("user", "")
	resp, err = conn.request(pathUser.String(), "", "GET")
	if err != nil {
		log.WithError(err).Error("failed to get csrf-token.")
	} else {
		conn.updateCsrf(resp)
		resp, err = conn.do(req)
	}
	return resp, err
}

func (conn *Connection) updateCsrf(resp *http.Response) {
	newToken := resp.Header.Get("Emc-Csrf-Token")
	if conn.csrf != newToken {
		conn.csrf = newToken
		log.WithField("csrf-token", conn.csrf).Info("update csrf token.")
	}
}

func (conn *Connection) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("EMC-CSRF-TOKEn", conn.csrf)
	resp, err := conn.client.Do(req)
	log.WithField("request", req).Debug("send request.")
	if err != nil {
		log.WithError(err).Error("http request error.")
		return nil, err
	}
	log.WithField("response", resp).Debug("got response.")
	return resp, err
}

func getRespBody(resp *http.Response) string {
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("failed to read response body.")
	}
	respBody := string(bytes)
	log.WithField("body", respBody).Debug(resp)
	return respBody
}

func (conn *Connection) Get(url string) (string, error) {
	resp, err := conn.request(url, "", "GET")
	var body string
	if err != nil {
		body = ""
	} else {
		conn.updateCsrf(resp)
		body = getRespBody(resp)
	}
	return body, err
}

func (conn *Connection) Post(url string, body string) (string, error) {
	resp, err := conn.request(url, body, "POST")
	if restErr := extractRestError(resp); restErr != nil {
		err = restErr
	}
	return getRespBody(resp), err
}

func (conn *Connection) Delete(url, body string) (string, error) {
	resp, err := conn.request(url, body, "DELETE")
	if restErr := extractRestError(resp); restErr != nil {
		err = restErr
	}
	return getRespBody(resp), err
}

func (conn *Connection) getFields(rsc string) string {
	var ret string = ""
	if val, ok := conn.fields[rsc]; ok {
		ret = val
	} else {
		resp, err := conn.GetRscMeta(rsc)
		if err != nil {
			log.WithError(err).WithField("rsc", rsc).Error(
				"failed to get fields.")
		} else {
			t := &Type{}
			err := updateInstFromResp(resp, t)
			if err == nil {
				ret = t.AllFieldString()
				conn.fields[rsc] = ret
			}
		}
	}
	return ret
}

func (conn *Connection) GetRscAll(rsc, filter string) (string, error) {
	return conn.Get(conn.getAllUrl(rsc, filter).String())
}

func (conn *Connection) GetRscInst(rsc, id string) (string, error) {
	return conn.Get(conn.getInstUrl(rsc, id).String())
}

func (conn *Connection) GetRscMeta(rsc string) (string, error) {
	return conn.Get(conn.getTypeUrl(rsc).String())
}

func (conn *Connection) DeleteRscInst(rsc, id, body string) (string, error) {
	return conn.Delete(conn.deleteInstUrl(rsc, id).String(), body)
}

func (conn *Connection) PostInstance(rsc, body string) (string, error) {
	_url := conn.postInstUrl(rsc)
	log.WithField("url", _url).WithField("body", body).Info("post new instance.")
	if resp, err := conn.Post(_url.String(), body); err != nil {
		log.WithError(err).Error("failed to create new instance.")
		return "", err
	} else {
		var id string
		if id, err = getIdFromPostResp(resp); err != nil {
			return "", err
		}
		return id, nil
	}
}

func (conn *Connection) baseUrl() string {
	return fmt.Sprintf("https://%v/api", conn.ip)
}

func (conn *Connection) assembleGetUrl(rsc, raw string) *url.URL {
	return conn.assembleGetUrlWithFilter(rsc, raw, "")
}

func (conn *Connection) assembleGetUrlWithFilter(rsc, raw, filter string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.WithError(err).WithField("url", raw).Error("parsing url error")
		return nil
	}
	q := u.Query()
	q.Add("compact", "true")
	fields := conn.getFields(rsc)
	if fields != "" {
		q.Add("fields", fields)
	}
	if filter != "" {
		q.Add("filter", filter)
	}
	u.RawQuery = q.Encode()
	return u
}

func (conn *Connection) assemblePostUrl(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		log.WithError(err).WithField("url", raw).Error("parsing url error")
		return nil
	}
	q := u.Query()
	q.Add("compact", "true")
	u.RawQuery = q.Encode()
	return u
}

func (conn *Connection) assembleDeleteUrl(raw string) *url.URL {
	return conn.assemblePostUrl(raw)
}

func (conn *Connection) getAllUrl(rsc, filter string) *url.URL {
	rawUrl := fmt.Sprintf("%v/types/%v/instances", conn.baseUrl(), rsc)
	return conn.assembleGetUrlWithFilter(rsc, rawUrl, filter)
}

func (conn *Connection) getInstUrl(rsc string, id string) *url.URL {
	rawUrl := fmt.Sprintf("%v/instances/%v/%v", conn.baseUrl(), rsc, id)
	return conn.assembleGetUrl(rsc, rawUrl)
}

func (conn *Connection) deleteInstUrl(rsc string, id string) *url.URL {
	rawUrl := fmt.Sprintf("%v/instances/%v/%v", conn.baseUrl(), rsc, id)
	return conn.assembleDeleteUrl(rawUrl)
}

func (conn *Connection) postInstUrl(rsc string) *url.URL {
	rawUrl := fmt.Sprintf("%v/types/%v/instances", conn.baseUrl(), rsc)
	return conn.assemblePostUrl(rawUrl)
}

func (conn *Connection) getTypeUrl(rsc string) *url.URL {
	rawUrl := fmt.Sprintf("%v/types/%v", conn.baseUrl(), rsc)
	return conn.assembleGetUrl("type", rawUrl)
}

func (conn *Connection) getActionUrl(rsc, action string) *url.URL {
	rawUrl := fmt.Sprintf("%v/types/%v/action/%v", conn.baseUrl(), rsc, action)
	return conn.assemblePostUrl(rawUrl)
}

func (conn *Connection) getInstActionUrl(rsc, inst, action string) *url.URL {
	rawUrl := fmt.Sprintf("%v/instances/%v/%v/action/%v",
		conn.baseUrl(), rsc, inst, action)
	return conn.assemblePostUrl(rawUrl)
}
