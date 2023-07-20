package context

import (
	"fmt"
	"testing"

	"github.com/jwzl/beehive/pkg/core/model"
)

func TestSendSync(t *testing.T) {
	InitContext(MsgCtxTypeChannel)
	AddModule("test_dest")
	messsage := &model.Message{}

	go func() {
		Send("test_dest", messsage)
		fmt.Printf("send message %v\n", messsage)
	}()

	msg, err := Receive("test_dest")
	fmt.Printf("receive msg: %v, error: %v\n", msg, err)
}
