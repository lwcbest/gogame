package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func RequestJSON(method, urlstr string, headers map[string]string, bodyinterface interface{}) ([]byte, error) {
	var bodys io.Reader
	if bodyinterface != nil {
		switch bodyinterface.(type) {
		case string:
			bodys = bytes.NewBufferString(bodyinterface.(string))
		case io.Reader:
			bodys = bodyinterface.(io.Reader)
		default:
			by, _ := json.Marshal(bodyinterface)
			bodys = bytes.NewBuffer(by)
		}
	}

	req, err := http.NewRequest(method, urlstr, bodys)
	if err != nil {
		return nil, err
	}

	for hk, hv := range headers {
		req.Header.Set(hk, hv)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	return b, err
}
