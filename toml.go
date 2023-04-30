// Package toml provides a toml codec
package toml // import "go.unistack.org/micro-codec-toml/v4"

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
	pb "go.unistack.org/micro-proto/v4/codec"
	"go.unistack.org/micro/v4/codec"
	rutil "go.unistack.org/micro/v4/util/reflect"
)

type tomlCodec struct {
	opts codec.Options
}

var _ codec.Codec = &tomlCodec{}

const (
	flattenTag = "flatten"
)

func (c *tomlCodec) Marshal(v interface{}, opts ...codec.Option) ([]byte, error) {
	if v == nil {
		return nil, nil
	}

	options := c.opts
	for _, o := range opts {
		o(&options)
	}

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		return m.Data, nil
	case *pb.Frame:
		return m.Data, nil
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

	if nv, nerr := rutil.StructFieldByTag(v, options.TagName, flattenTag); nerr == nil {
		v = nv
	}

	switch m := v.(type) {
	case *codec.Frame:
		m.Data = d
		return nil
	case *pb.Frame:
		m.Data = d
		return nil
	}

	return toml.Unmarshal(d, v)
}

func (c *tomlCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *tomlCodec) ReadBody(conn io.Reader, v interface{}) error {
	if v == nil {
		return nil
	}

	buf, err := io.ReadAll(conn)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return nil
	}
	return c.Unmarshal(buf, v)
}

func (c *tomlCodec) Write(conn io.Writer, m *codec.Message, v interface{}) error {
	if v == nil {
		return nil
	}
	buf, err := c.Marshal(v)
	if err != nil {
		return err
	} else if len(buf) == 0 {
		return codec.ErrInvalidMessage
	}

	_, err = conn.Write(buf)
	return err
}

func (c *tomlCodec) String() string {
	return "toml"
}

func NewCodec(opts ...codec.Option) *tomlCodec {
	return &tomlCodec{opts: codec.NewOptions(opts...)}
}
