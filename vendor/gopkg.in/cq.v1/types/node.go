package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Conn interface {
	BaseURL() string
}

type Node struct {
	LabelURI   string                 `json:"labels"`
	SelfURI    string                 `json:"self"`
	Properties map[string]CypherValue `json:"data"`
}

func (n *Node) Scan(value interface{}) error {
	if value == nil {
		return ErrScanOnNil
	}

	switch value.(type) {
	case map[string]CypherValue:
		cv := value.(map[string]CypherValue)
		var ok = false
		var inner CypherValue
		inner, ok = cv["data"]
		if ok != true {
			break
		}
		properties, ok := inner.Val.(map[string]CypherValue)
		if ok {
			n.Properties = properties
		} else {
			n.Properties = map[string]CypherValue{}

			properties := inner.Val.(map[string]string)

			for k, v := range properties {
				n.Properties[k] = CypherValue{CypherString, v}
			}
		}

		inner, ok = cv["self"]
		if ok != true {
			break
		}
		n.SelfURI = inner.Val.(string)
		inner, ok = cv["labels"]
		if ok != true {
			break
		}
		n.LabelURI = inner.Val.(string)
		return nil
	case []byte:
		err := json.Unmarshal(value.([]byte), &n)
		return err
	}
	return errors.New(fmt.Sprintf("cq: invalid Scan value for %T: %T", n, value))
}

func (n *Node) Labels(baseURL string) ([]string, error) {
	base, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}
	labelURL, err := url.Parse(n.LabelURI)
	if err != nil {
		return nil, err
	}
	labelURL.Scheme = base.Scheme
	labelURL.User = base.User
	req, err := http.NewRequest("GET", n.LabelURI, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println(labelURL.User)
	pass, _ := labelURL.User.Password()
	user := labelURL.User.Username()
	req.SetBasicAuth(user, pass)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	ret := []string{}
	err = json.NewDecoder(resp.Body).Decode(&ret)
	return ret, err
}
