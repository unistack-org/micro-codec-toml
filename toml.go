// Package toml provides a toml codec
package toml

import (
	"bytes"

	"github.com/BurntSushi/toml"
	pb "go.unistack.org/micro-proto/v3/codec"
	"go.unistack.org/micro/v3/codec"
	rutil "go.unistack.org/micro/v3/util/reflect"
)

type tomlCodec struct {
	opts codec.Options
}

var _ codec.Codec = &tomlCodec{}

func (c *tomlCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
	case codec.RawMessage:
		return []byte(m), nil
	case *codec.RawMessage:
		return []byte(*m), nil
	}

	buf := bytes.NewBuffer(nil)
	err := toml.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (c *tomlCodec) Unmarshal(d []byte, v interface{}, opts ...codec.Option) error {
	if v == nil || len(d) == 0 {
		return nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if options.Flatten {
		if nv, nerr := rutil.StructFieldByTag(v, options.TagName, "flatten"); nerr == nil {
			v = nv
		}
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	case *codec.RawMessage:
		*m = append((*m)[0:0], d...)
		return nil
	case codec.RawMessage:
		copy(m, d)
		return nil
	}

	return toml.Unmarshal(d, v)
}

func (c *tomlCodec) String() string {
	return "toml"
}

func NewCodec(opts ...codec.Option) *tomlCodec {
	return &tomlCodec{opts: codec.NewOptions(opts...)}
}
