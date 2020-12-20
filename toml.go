// Package toml provides a toml codec
package toml

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/unistack-org/micro/v3/codec"
)

type tomlCodec struct{}

func (c *tomlCodec) Marshal(b interface{}) ([]byte, error) {
	switch m := b.(type) {
	case nil:
		return nil, nil
	case *codec.Frame:
		return m.Data, nil
	}

	buf := bytes.NewBuffer(nil)
	defer buf.Reset()
	err := toml.NewEncoder(buf).Encode(b)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (c *tomlCodec) Unmarshal(b []byte, v interface{}) error {
	if b == nil {
		return nil
	}
	switch m := v.(type) {
	case nil:
		return nil
	case *codec.Frame:
		m.Data = b
		return nil
	}

	return toml.Unmarshal(b, v)
}

func (c *tomlCodec) ReadHeader(conn io.Reader, m *codec.Message, t codec.MessageType) error {
	return nil
}

func (c *tomlCodec) ReadBody(conn io.Reader, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		buf, err := ioutil.ReadAll(conn)
		if err != nil {
			return err
		}
		m.Data = buf
		return nil
	}

	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		return err
	}

	return toml.Unmarshal(buf, b)
}

func (c *tomlCodec) Write(conn io.Writer, m *codec.Message, b interface{}) error {
	switch m := b.(type) {
	case nil:
		return nil
	case *codec.Frame:
		_, err := conn.Write(m.Data)
		return err
	}

	return toml.NewEncoder(conn).Encode(b)
}

func (c *tomlCodec) String() string {
	return "toml"
}

func NewCodec() codec.Codec {
	return &tomlCodec{}
}
