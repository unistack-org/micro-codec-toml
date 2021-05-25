// Package toml provides a toml codec
package toml

import (
	"bytes"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/unistack-org/micro/v3/codec"
	rutil "github.com/unistack-org/micro/v3/util/reflect"
)

type tomlCodec struct{}

const (
	flattenTag = "flatten"
)

func (c *tomlCodec) Marshal(v interface{}) ([]byte, error) {
	switch m := v.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	buf := bytes.NewBuffer(nil)
	defer buf.Reset()

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	err := toml.NewEncoder(buf).Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *tomlCodec) Unmarshal(b []byte, v interface{}) error {
	if len(b) == 0 || v == nil {
		return nil
	}

	if m, ok := v.(*codec.Frame); ok {
		m.Data = b
		return nil
	}

	if nv, nerr := rutil.StructFieldByTag(v, codec.DefaultTagName, flattenTag); nerr == nil {
		v = nv
	}

	return toml.Unmarshal(b, v)
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

func NewCodec() codec.Codec {
	return &tomlCodec{}
}
