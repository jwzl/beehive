package context

import (
	"sync"
	"time"
	gocontext "context"

	"k8s.io/klog/v2"
	"github.com/jwzl/beehive/pkg/core/model"
)

//define channel type
const (
	MsgCtxTypeChannel = "channel"
)

// coreContext is global context object
type coreContext struct {
	moduleContext  ModuleContext
	messageContext MessageContext
	//for golang native context.
	ctx            gocontext.Context
	cancel         gocontext.CancelFunc
}

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
func Send(module string, message *model.Message) {
	context.messageContext.Send(module, message)
}

// Receive the message
// module : local module name
func Receive(module string) (*model.Message, error) {
	message, err := context.messageContext.Receive(module)
	if err == nil {
		return message, nil
	}
	return message, err
}

// SendSync sends message in sync mode
// module: the destination of the message
// timeout: if <= 0 using default value(30s)
func SendSync(module string,
	message *model.Message, timeout time.Duration) (*model.Message, error) {
	return context.messageContext.SendSync(module, message, timeout)
}

// SendResp sends response
// please get resp message using model.NewRespByMessage
func SendResp(resp *model.Message) {
	context.messageContext.SendResp(resp)
}

// SendToGroup broadcasts the message to all of group members
func SendToGroup(group string, message *model.Message) {
	context.messageContext.SendToGroup(group, message)
}
// SendToGroupSync broadcasts the message to all of group members in sync mode
func SendToGroupSync(group string, message *model.Message, timeout time.Duration) error {
	return context.messageContext.SendToGroupSync(group, message, timeout)
}
