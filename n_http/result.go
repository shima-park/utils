package n_http

import (
	"encoding/json"
	"errors"
	"fmt"

	"git.wdwd.com/nova/n_utils"
	log "github.com/Sirupsen/logrus"
)

type Result interface {
	CheckStatus() error
	GetJson(v interface{}) error
	GetInt(key string) (int64, error)
	GetString(key string) (string, error)
	GetContent() (string, error)
	GetRawData() error
	SetRawData([]byte)
}

type DefaultResult struct {
	Status   string      `json:"status"`
	Msg      string      `json:"msg"`
	Code     int64       `json:"code"`
	Data     interface{} `json:"data"`
	raw_data []byte      `json:"-,omitempty"`
}

func (this *DefaultResult) CheckStatus() (err error) {
	switch {
	case this.Status == "success":
		return
	case this.Code == 0:
		return
	case this.Status == "error":
		log.Error(fmt.Sprintf("Code: %d, Msg: %s", this.Code, this.Msg))
		err = errors.New(this.Msg)
		return
	}
	return
}

func (this *DefaultResult) get(key string, v interface{}) (err error) {
	if err = this.CheckStatus(); err != nil {
		log.Warn(err)
		return
	}

	switch v.(type) {
	case *string:
		data_map := this.Data.(map[string]interface{})
		var s *string = v.(*string)
		*s = n_utils.Be_string(data_map[key])
		v = s
		return
	case *int64:
		data_map := this.Data.(map[string]interface{})
		var i *int64 = v.(*int64)
		*i = n_utils.Be_int(data_map[key])
		v = i
		return
	default:
		this.Data = v
		return json.Unmarshal(this.GetRawData(), this)
	}
	return
}

func (this *DefaultResult) GetString(key string) (val string, err error) {
	err = this.get(key, &val)
	return
}

func (this *DefaultResult) GetInt(key string) (i int64, err error) {
	err = this.get(key, &i)
	return
}

func (this *DefaultResult) GetJson(v interface{}) (err error) {
	err = this.get("", v)
	return
}

func (this *DefaultResult) GetRawData() (b []byte) {
	return this.raw_data
}

func (this *DefaultResult) GetContext() (c string, err error) {
	if err = this.CheckStatus(); err != nil {
		log.Warn(err)
		return
	}

	defer func() {
		if f := recover(); f != nil {
			err = errors.New(fmt.Sprintf("%v", f))
		}
	}()

	if this.Data != nil {
		c = this.Data.(string)
	}
	return
}
