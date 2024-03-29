package context

import (
	"sync"
	"time"

	"github.com/jwzl/beehive/pkg/core/model"
)

// UnixSocketContext unixsocket struct
type UnixSocketContext struct {
	filename string
	bufsize  int
	handler  func(string) string
}

var (
	// singleton
	usContext *UnixSocketContext
	usOnce    sync.Once
)

//GetUnixSocketContext defines and returns unix socket context object
func GetUnixSocketContext() *UnixSocketContext {
	usOnce.Do(func() {
		usContext = &UnixSocketContext{}
	})
	return usContext
}

//AddModule adds module to context
func (ctx *UnixSocketContext) AddModule(module string) {

}

//AddModuleGroup adds module to module context group
func (ctx *UnixSocketContext) AddModuleGroup(module, group string) {

}

//Cleanup cleans up module
func (ctx *UnixSocketContext) Cleanup(module string) {

}

// Send async mode
func (ctx *UnixSocketContext) Send(module string, content *model.Message) {

}

//Receive the message
//local module name
func (ctx *UnixSocketContext) Receive(module string) *model.Message {
	return nil
}

func (ctx *UnixSocketContext) SendSync(module string, message *model.Message, timeout time.Duration) (*model.Message, error){
	return nil, nil
}

func (ctx *UnixSocketContext) SendResp(message *model.Message) {

}
// SendToGroup send msg to modules. Todo: do not stuck
func (ctx *UnixSocketContext) SendToGroup(moduleType string, message *model.Message) {

}

func (ctx *UnixSocketContext) SendToGroupSync(moduleType string, message *model.Message, timeout time.Duration) error {
	return nil
}
