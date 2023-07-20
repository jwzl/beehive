package context

import (
	gocontext "context"
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
	Send(module string, message interface{})
	Receive(module string) (interface{}, error)
	// group broadcast
	SendToGroup(moduleType string, message interface{})
}

// coreContext is global context object
type coreContext struct {
	moduleContext  ModuleContext
	messageContext MessageContext
	//for golang native context.
	ctx            gocontext.Context
	cancel         gocontext.CancelFunc
}
