package context

import (
	"time"

	"github.com/jwzl/beehive/pkg/core/model"
)

//ModuleContext is interface for context module management
type ModuleContext interface {
	AddModule(module string)
	AddModuleGroup(module, group string)
	Cleanup(module string)
}

//MessageContext is interface for message syncing
type MessageContext interface {
	// async mode
	//send message which module by module name.
	Send(module string, message *model.Message)
	//recieve the message to the module by module name.
	Receive(module string) (*model.Message, error)

	// sync mode
	SendSync(module string, message *model.Message, timeout time.Duration) (*model.Message, error)
	SendResp(message *model.Message)

	// group broadcast
	SendToGroup(group string, message *model.Message)
	SendToGroupSync(group string, message *model.Message, timeout time.Duration) error
}
