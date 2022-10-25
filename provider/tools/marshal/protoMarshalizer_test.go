package marshal_test

import (
	"fmt"
	"testing"

	"github.com/klever-io/klever-go-sdk/models/proto"
	"github.com/klever-io/klever-go-sdk/provider/tools/marshal"
	"github.com/stretchr/testify/assert"
)

var protoMarsh = marshal.ProtoMarshalizer{}

func recovedMarshal(obj interface{}) (buf []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			if panicError, ok := p.(error); ok {
				err = panicError
			} else {
				err = fmt.Errorf("%#v", p)
			}
			buf = nil
		}
	}()
	buf, err = protoMarsh.Marshal(obj)
	return
}

func TestProtoMarshalizer_MarshalWrongObj(t *testing.T) {

	obj := "klever"
	encNode, err := recovedMarshal(obj)
	assert.Nil(t, encNode)
	assert.NotNil(t, err)
}

func TestProtoMarshalizer_NilMarshalizer(t *testing.T) {
	var m *marshal.ProtoMarshalizer
	assert.True(t, m.IsInterfaceNil())
}

func TestProtoMarshalizer_MarshalNilObj(t *testing.T) {
	encNode, err := recovedMarshal(nil)
	assert.Nil(t, encNode)
	assert.Contains(t, err.Error(), "can not serialize the object")
}

func TestProtoMarshalizer_ShouldWork(t *testing.T) {
	m := marshal.NewProtoMarshalizer()

	tx := &proto.Transaction{
		RawData: &proto.Transaction_Raw{
			Nonce:  10,
			Sender: []byte("Klever"),
		},
		Block: 1234,
	}
	data, err := m.Marshal(tx)
	assert.Nil(t, err)

	tx2 := &proto.Transaction{}

	err = m.Unmarshal(tx2, data)
	assert.Nil(t, err)
	assert.Equal(t, tx.RawData.Sender, tx2.RawData.Sender)
	assert.Equal(t, tx.Block, tx2.Block)
}
