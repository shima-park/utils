package n_http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"net/http"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

func Get(api string, params map[string]interface{}, res Result) (err error) {
	return do("GET", api, params, res)
}

func Post(api string, params map[string]interface{}, res Result) (err error) {
	return do("POST", api, params, res)
}

func Put(api string, params map[string]interface{}, res Result) (err error) {
	return do("PUT", api, params, res)
}

func Delete(api string, params map[string]interface{}, res Result) (err error) {
	return do("DELETE", api, params, res)
}

func do(method, api string, params map[string]interface{}, res Result) (err error) {
	return do_http(method, api, params, res)
}

func do_http(method, api_url string, params map[string]interface{}, res Result) (err error) {
	log.Debug("[HTTP]", method, ": ", api_url, "params: ", fmt.Sprintf("%+v", params))

	var body []byte
	r, err := new_request(method, api_url, params)
	if err != nil {
		return
	}
	var c http.Client
	resp, err := c.Do(r)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	switch {
	case err != nil:
		log.Debug(r)
	default:
		if len(body) == 0 {
			err = errors.New("no result return")
		} else {
			log.Info("response: ", string(body))
			res.SetRawData(body)
		}
	}

	return
}

func new_request(method, api string, params map[string]interface{}) (req *http.Request, err error) {
	vals := url.Values{}

	if params != nil {
		for k, v := range params {
			vals.Set(k, param_to_str(v))
		}
	}

	req, err = http.NewRequest(method, api, nil)
	if err != nil {
		return
	}

	use_url_query := method != "POST"
	query_string := vals.Encode()

	if use_url_query {
		req.URL.RawQuery = query_string
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.ContentLength = int64(len(query_string))
		req.Body = &closebuf{bytes.NewBufferString(query_string)}
	}

	return
}

type closebuf struct {
	*bytes.Buffer
}

func (cb *closebuf) Close() error {
	return nil
}

func param_to_str(v interface{}) (v2 string) {
	switch v.(type) {
	case string:
		v2 = v.(string)
	default:
		buf, _ := json.Marshal(v)
		v2 = string(buf)
	}
	return
}
