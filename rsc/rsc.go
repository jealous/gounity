package rsc

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"strings"
)

type Rscer interface {
	GetType() string
	GetId() string
	GetConn() *Connection
	JsonId() string
}

type RscLister interface {
	getRscType() string
	appendRsc(rsc Rscer) Rscer
	GetConn() *Connection
	GetFilter() string
	Iterator() *RscListIterator
	Size() int

	// following method need to be implemented by individual resource
	initRsc() Rscer
}

type RscListCtor interface {
	initList(filter string) RscLister
}

func Update(r Rscer) Rscer {
	if resp, err := r.GetConn().GetRscInst(r.GetType(), r.GetId()); err != nil {
		log.WithError(err).Error("failed to get resource instance.")
	} else {
		if err = updateInstFromResp(resp, r); err != nil {
			log.WithError(err).Error("failed to update resource.")
		}
	}
	return r
}

func getElemFromResp(input, elem string) ([]byte, error) {
	data := []byte(input)
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(data, &objMap)
	if err != nil {
		log.WithError(err).WithField("input", input).Error("failed to unmarshal response.")
		return nil, err
	}
	return *objMap[elem], nil
}

func getContentFromResp(input string) ([]byte, error) {
	return getElemFromResp(input, "content")
}

type idJson struct {
	Id string
}

func getIdFromPostResp(resp string) (string, error) {
	content, err := getContentFromResp(resp)
	if err != nil {
		return "", err
	}
	data := []byte(content)
	var id idJson
	err = json.Unmarshal(data, &id)
	if err != nil {
		log.WithError(err).WithField("input", content).Error("failed to unmarshal content.")
		return "", err
	}
	return id.Id, nil
}

func getErrorFromResp(input string) ([]byte, error) {
	return getElemFromResp(input, "error")
}

func updateInstFromResp(input string, rsc interface{}) error {
	if content, err := getContentFromResp(input); err != nil {
		return err
	} else {
		if err = json.Unmarshal(content, rsc); err != nil {
			return err
		}
	}
	return nil
}

func UpdateList(lister RscLister) RscLister {
	resp, err := lister.GetConn().GetRscAll(
		lister.getRscType(), lister.GetFilter())
	if err == nil {
		updateListerFromResp(resp, lister)
	}
	return lister
}

func updateListerFromResp(input string, lister RscLister) error {
	data := []byte(input)
	var objMap map[string]*json.RawMessage
	err := json.Unmarshal(data, &objMap)
	if err != nil {
		log.WithError(err).Error("failed to unmarshal response.")
		return err
	}

	var entries []json.RawMessage
	_ = json.Unmarshal(*objMap["entries"], &entries)

	for i := range entries {
		rsc := lister.appendRsc(lister.initRsc())
		err = updateInstFromResp(string(entries[i]), rsc)
	}
	return err
}

type Rsc struct {
	conn  *Connection
	type_ string
	Id    string
}

func (r *Rsc) GetType() string {
	return r.type_
}

func (r *Rsc) GetId() string {
	return r.Id
}

func (r *Rsc) GetConn() *Connection {
	return r.conn
}

func (r *Rsc) JsonId() string {
	return fmt.Sprintf(`{"id": "%v"}`, r.Id)
}

type RscList struct {
	conn   *Connection
	type_  string
	list_  []Rscer
	filter string
}

func (rl *RscList) getRscType() string {
	return rl.type_
}

func (rl *RscList) GetConn() *Connection {
	return rl.conn
}

func (rl *RscList) Iterator() *RscListIterator {
	return &RscListIterator{items: rl.list_, current: 0}
}

func (rl *RscList) Size() int {
	return len(rl.list_)
}

func (rl *RscList) appendRsc(rsc Rscer) Rscer {
	rl.list_ = append(rl.list_, rsc)
	return rsc
}

func (rl *RscList) GetFilter() string {
	return rl.filter
}

type RscListIterator struct {
	items   []Rscer
	current int
}

func (it *RscListIterator) Next() bool {
	return it.current < len(it.items)
}

func (it *RscListIterator) Value() Rscer {
	if it.current >= len(it.items) {
		return nil
	}
	ret := it.items[it.current]
	it.current++
	return ret
}

func (it *RscListIterator) NextIndex() int {
	return it.current
}

type Type struct {
	Name          string
	Description   string
	Documentation string
	Attributes    []Attributes
}

type Attributes struct {
	Name         string
	Type         string
	Description  string
	DisplayValue string
}

func (t *Type) AllFieldString() string {
	attributes := make([]string, len(t.Attributes))
	for i, attr := range t.Attributes {
		attributes[i] = attr.Name
	}
	return strings.Join(attributes, ",")
}

type Health struct {
	Value          int
	DescriptionIds []string
	Descriptions   []string
}

func getRscByName(name string, ctor RscListCtor) Rscer {
	list := ctor.initList(fmt.Sprintf(`name eq "%v"`, name))
	UpdateList(list)
	return list.Iterator().Value()
}

func getRscList(ctor RscListCtor) RscLister {
	rscList := ctor.initList("")
	UpdateList(rscList)
	return rscList
}

type ErrorMessage struct {
	EnUs string `json:"en-US"`
}

type RestError struct {
	ErrorCode      uint64
	HttpStatusCode int
	Messages       []ErrorMessage
	Created        string
}

func (err *RestError) Error() string {
	if len(err.Messages) != 0 {
		msg := err.Messages[0]
		return err.formattedMsg(msg.EnUs)
	}
	return err.formattedMsg("error message not available.")
}

func (err *RestError) formattedMsg(msg string) string {
	return fmt.Sprintf("[%v] %v %v", err.ErrorCode, err.Created, msg)
}
