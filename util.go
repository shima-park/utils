package n_utils

import (
	"bytes"
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var guidre *regexp.Regexp
var numre *regexp.Regexp
var strre *regexp.Regexp

func init() {
	guidre, _ = regexp.Compile("[^0-9A-Z]")
	numre, _ = regexp.Compile("[^0-9]")
	strre, _ = regexp.Compile("[^1-9a-np-z]")
}

func Guid(l int, check func(string) bool, step ...int) string {
	h := md5.New()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Fprintf(h, "%d.%d.%f", step, r.Int63(), r.Float64())
	str := base32.StdEncoding.EncodeToString(h.Sum(nil))
	str = guidre.ReplaceAllString(str, "")

	n := len(str)
	if n < l {
		str = str + Guid(l-n, func(_ string) bool { return true }, 5)
	} else {
		str = str[0:l]
	}

	if check(str) {
		return str
	} else {
		return Guid(l, check, step[0]+1)
	}
}

func GrandNum(l int) string {
	h := md5.New()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Fprintf(h, "%d.%d.%f", r.Intn(100), r.Int63(), r.Float64())
	str := hex.EncodeToString(h.Sum(nil))
	str = numre.ReplaceAllString(str, "")

	n := len(str)
	if n < l {
		str = str + GrandNum(l-n)
	} else {
		str = str[0:l]
	}

	return str
}

func GrandStr(l int) string {
	h := md5.New()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	fmt.Fprintf(h, "%d.%d.%f", r.Intn(100), r.Int63(), r.Float64())
	str := hex.EncodeToString(h.Sum(nil))

	n := len(str)
	if n < l {
		str = str + GrandStr(l-n)
	} else {
		str = str[0:l]
	}

	return str
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))

}

func Sha1(str string) string {
	h := sha1.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func BeInt(i interface{}) (number int64, err error) {
	switch i.(type) {
	case int, int8, int16, int32, int64:
		switch i.(type) {
		case int, int8, int16, int32, int64:
			number = reflect.ValueOf(i).Int()
			return
		}
	case uint, uint8, uint16, uint32, uint64:
		switch i.(type) {
		case uint, uint8, uint16, uint32, uint64:
			number = reflect.ValueOf(i).Int()
			return
		}
	case string:
		number, err = strconv.ParseInt(i.(string), 10, 0)
		return
	case float64, float32:
		number = int64(reflect.ValueOf(i).Float())
		return
	}
	return
}

func Be_int(i interface{}) int64 {
	switch i.(type) {
	case int, int8, int16, int32, int64:
		switch i.(type) {
		case int, int8, int16, int32, int64:
			return reflect.ValueOf(i).Int()
		}
	case uint, uint8, uint16, uint32, uint64:
		switch i.(type) {
		case uint, uint8, uint16, uint32, uint64:
			return reflect.ValueOf(i).Int()
		}
	case string:
		d, _ := strconv.ParseInt(i.(string), 10, 0)
		return d
	case float64, float32:
		return int64(reflect.ValueOf(i).Float())
	case bool:
		println("bbool")
	}
	return 0
}

func Be_string(i interface{}) string {
	switch i.(type) {
	case string:
		return i.(string)
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(i).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatInt(reflect.ValueOf(i).Int(), 10)
	case float64, float32:
		return strconv.FormatFloat(reflect.ValueOf(i).Float(), 'f', 2, 64)
	}
	return ""
}

func Be_byte(i interface{}) (b []byte) {
	switch v := i.(type) {
	case string:
		return []byte(v)
	case interface{}:
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		enc.Encode(v)
		return buf.Bytes()
	}
	return
}

func Str_to_float(i string, flag int) (float64, error) {
	if i != "" {
		return strconv.ParseFloat(i, flag)
	}
	return 0, nil
}

func GetDobjInfo(object interface{}) (key, fileds string) {
	object_type := reflect.TypeOf(object).Elem()
	num_filed := object_type.NumField()
	if num_filed > 0 {
		key = strings.ToLower(object_type.Field(0).Name)
	} else {
		return
	}

	for i := 0; i < num_filed; i++ {
		object_filed := object_type.Field(i)
		filed_name := strings.ToLower(object_filed.Name)

		field_tag := object_filed.Tag
		if string(field_tag) != "" {
			if field_tag.Get("skip") == "true" {
				continue
			}

			if field_tag.Get("primary_key") == "true" {
				key = filed_name
			}
		}

		fileds += "`" + filed_name + "`, "
	}

	key = "`" + key + "`"
	fileds = strings.TrimRight(fileds, ", ")

	return
}

func ConvertStructToMap(object interface{}) (obj_M map[string]interface{}) {
	obj_M = make(map[string]interface{})
	object_type := reflect.TypeOf(object).Elem()
	value := reflect.ValueOf(object).Elem()
	num_filed := object_type.NumField()
	if num_filed < 1 {
		return
	}

	for i := 0; i < num_filed; i++ {
		object_filed := object_type.Field(i)
		filed_name := strings.ToLower(object_filed.Name)

		field_tag := object_filed.Tag
		json_fieds := field_tag.Get("json")
		if (json_fieds) != "" {
			json_fied_arr := strings.Split(json_fieds, ",")
			if len(json_fied_arr) > 1 {
				if json_fied_arr[0] != "-" {
					filed_name = json_fied_arr[0]
				}
			}
		}

		obj_M[filed_name] = value.Field(i).Interface()
	}

	return
}

func FilterMap(obj_M map[string]interface{}, fileds []string) {
	filter_M := make(map[string]bool)
	for index, _ := range fileds {
		filter_M[fileds[index]] = true
	}

	for filed_name, _ := range obj_M {
		if filter_M[filed_name] == false {
			delete(obj_M, filed_name)
		}
	}

	return
}

func DiffStrings(dst_strs []string, src_strs []string) (diff_strs []string) {
	diff_strs = make([]string, 0)
	for _, dst_str := range dst_strs {
		i := -1
		for j, src_str := range src_strs {
			if dst_str == src_str {
				i = j
				break
			}
		}
		if i == -1 {
			diff_strs = append(diff_strs, dst_str)
		}
	}

	return
}

func InterStrings(dst_strs []string, src_strs []string) (inter_strs []string) {
	inter_strs = make([]string, 0)
	for _, dst_str := range dst_strs {
		i := -1
		for j, src_str := range src_strs {
			if dst_str == src_str {
				i = j
				break
			}
		}
		if i != -1 {
			inter_strs = append(inter_strs, dst_str)
		}
	}

	return
}

func GetDobjInfoForJoin(object interface{}) (key, fileds, rel_fileds string) {
	object_type := reflect.TypeOf(object).Elem()
	num_filed := object_type.NumField()
	if num_filed > 0 {
		key = strings.ToLower(object_type.Field(0).Name)
	} else {
		return
	}

	table_name := GetDobjTableName(object)
	for i := 0; i < num_filed; i++ {
		object_filed := object_type.Field(i)
		filed_name := strings.ToLower(object_filed.Name)

		field_tag := object_filed.Tag
		if string(field_tag) != "" {
			if field_tag.Get("skip") == "true" {
				continue
			}

			if field_tag.Get("primary_key") == "true" {
				key = filed_name
			}
		}

		fileds += "`" + filed_name + "`, "
		rel_fileds += table_name + "." + filed_name + " as `" + filed_name + "`, "
	}

	key = "`" + key + "`"
	fileds = strings.TrimRight(fileds, ", ")
	rel_fileds = strings.TrimRight(rel_fileds, ", ")

	return
}

func GetDobjTableName(object interface{}) (table_name string) {
	object_val := reflect.ValueOf(object)
	method := object_val.MethodByName("TableName")
	if method.IsValid() == false {
		object_type := reflect.TypeOf(object).Elem()
		table_name = strings.ToLower(object_type.Name())
	} else {
		res := method.Call([]reflect.Value{})
		table_name = res[0].String()
	}
	return
}

func RandString(n int) string {
	const alphanum = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols := big.NewInt(int64(len(alphanum)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(n)), nil)
	r, err := crand.Int(crand.Reader, states)
	if err != nil {
		log.Println(err)
	}
	var bytes = make([]byte, n)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)
	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = alphanum[symbol.Int64()]
	}
	return string(bytes)
}

func Contain(obj interface{}, target interface{}) (bExist bool, err error) {
	targetValue := reflect.ValueOf(target)
	targetKind := reflect.TypeOf(target).Kind()
	switch targetKind {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				bExist = true
				return
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			bExist = true
			return
		}

	default:
		err = errors.New("target type error")
		return
	}

	return
}

func ParseDefaultDura(dura_s, default_dura_s string) (dura time.Duration) {
	dura, err := time.ParseDuration(dura_s)
	if err == nil {
		return
	} else {
		fmt.Printf("[ERROR] Failed to parse duration from config, err: %v\n", dura)
	}

	dura, err = time.ParseDuration(default_dura_s)
	if err == nil {
		return
	} else {
		fmt.Printf("[ERROR] Failed to parse duration from api, err: %v\n", dura)
	}

	return
}

func ParseTimestamp(ts string) (hour, min int) {
	if ts == "" {
		return
	}

	t_arr := strings.Split(ts, ":")
	hour = int(Be_int(t_arr[0]))
	if len(t_arr) > 1 {
		min = int(Be_int(t_arr[1]))
	}

	return
}

func CronTask(hour, min int, t_f func(int, int)) {
	now := time.Now()
	year, month, day := now.Date()

	spec := time.Date(year, month, day, hour, min, 0, 0, now.Location())
	dura := spec.Sub(now)
	if dura.Nanoseconds() < 1 {
		go t_f(hour, min)
		spec := spec.AddDate(0, 0, 1)
		dura = spec.Sub(now)
	}

	tch := time.After(dura)
	go func() {
		<-tch
		go t_f(hour, min)
		dura, _ = time.ParseDuration("24h")
		tick := time.Tick(dura)
		go func() {
			for c := range tick {
				log.Println(c)
				go t_f(hour, min)
			}
		}()
	}()
}

//  目前周期为一天
// 构造时间戳，以及判断是否还未到发生时刻
func ConstructCronTs(hour, min int) (
	spec time.Time,
	last_spec time.Time,
	ttl_dura time.Duration,
	has_yet bool) {

	has_yet = true
	now := time.Now()
	year, month, day := now.Date()
	spec = time.Date(year, month, day, hour, min, 0, 0, now.Location())
	dura := spec.Sub(now)
	buf_dura, _ := time.ParseDuration("-5m")
	if dura.Minutes() < 4.999 || dura.Nanoseconds() < 0 {
		ttl_dura = spec.AddDate(0, 0, 1).Add(buf_dura).Sub(now)
	} else {
		has_yet = false
		return
	}
	last_spec = spec.AddDate(0, 0, -1)

	return
}

func ConstructCronTtlTs(hour, min int) (ttl_dura time.Duration) {
	now := time.Now()
	year, month, day := now.Date()
	spec := time.Date(year, month, day, hour, min, 0, 0, now.Location())
	dura := spec.Sub(now)
	buf_dura, _ := time.ParseDuration("-5m")
	if dura.Minutes() < 4.999 || dura.Nanoseconds() < 0 {
		ttl_dura = spec.AddDate(0, 0, 1).Add(buf_dura).Sub(now)
	} else {
		ttl_dura = buf_dura
	}

	return
}

func Unmarshal(src []byte, dst interface{}) (err error) {
	decode := json.NewDecoder(bytes.NewBuffer(src))
	decode.UseNumber()
	if err = decode.Decode(dst); err != nil {
		return
	}

	return
}
