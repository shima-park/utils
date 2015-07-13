package n_code

import (
	"reflect"

	"github.com/ugorji/go/codec"
)

var (
	bh codec.BincHandle
	mh codec.MsgpackHandle
	ch codec.CborHandle
)

const (
	BincHandle    = "bin"
	MsgpackHandle = "msgpack"
	CborHandle    = "cbor"
)

func init() {
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))
}

func BinaryPackEncode(v interface{}, b []byte, handle string) error {
	h := getHandle(handle)
	enc := codec.NewEncoderBytes(&b, h)
	return enc.Encode(v)
}

func BinaryPackDecode(v interface{}, b []byte, handle string) error {
	h := getHandle(handle)
	dec := codec.NewDecoderBytes(b, h)
	return dec.Decode(&v)
}

func getHandle(handle string) codec.Handle {
	var h codec.Handle
	switch handle {
	case "bin":
		h = &bh
	case "msgpack":
		h = &mh
	case "cbor":
		h = &ch
	}
	return h
}
