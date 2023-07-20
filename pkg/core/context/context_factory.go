package context

import (
	gocontext "context"
	"sync"

	"k8s.io/klog"
)

//define channel type
const (
	MsgCtxTypeChannel = "channel"
)

var (
	// singleton
	context *coreContext
	once    sync.Once
)

// InitContext init global context instance
func InitContext(contextType string){
	once.Do(func() {
		ctx, cancel := gocontext.WithCancel(gocontext.Background())
		context = &coreContext{
			ctx:    ctx,
			cancel: cancel,
		}
		switch contextType {
		case MsgCtxTypeChannel:
			channelContext := NewChannelContext()
			context.messageContext = channelContext
			context.moduleContext = channelContext
		default:
			klog.Warningf("Do not support context type:%s", contextType)
		}
	})
}

func GetContext() gocontext.Context {
	return context.ctx
}

func Done() <-chan struct{} {
	return context.ctx.Done()
}

// Cancel gocontext function
func Cancel() {
	context.cancel()
}

// AddModule adds module into module context
func AddModule(module string) {
	context.moduleContext.AddModule(module)
}

// AddModuleGroup adds module into module context group
func AddModuleGroup(module, group string) {
	context.moduleContext.AddModuleGroup(module, group)
}

// Cleanup cleans up module
func Cleanup(module string) {
	context.moduleContext.Cleanup(module)
}

// Send the message
func Send(module string, message interface{}) {
	context.messageContext.Send(module, message)
}

// Receive the message
// module : local module name
func Receive(module string) (interface{}, error) {
	message, err := context.messageContext.Receive(module)
	if err == nil {
		return message, nil
	}
	klog.Warning("Receive: failed to receive message")
	return message, err
}

// SendToGroup broadcasts the message to all of group members
func SendToGroup(moduleType string, message interface{}) {
	context.messageContext.SendToGroup(moduleType, message)
}
