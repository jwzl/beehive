package context

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/klog"
)

//constants for channel context
const (
	ChannelSizeDefault = 1024

	MessageTimeoutDefault = 30 * time.Second

	TickerTimeoutDefault = 20 * time.Millisecond
)

// ChannelContext is object for Context channel
type ChannelContext struct {
	//channels is for module channel.
	channels     map[string]chan interface{}
	chsLock      sync.RWMutex
	//typeChannels is for module group.
	typeChannels map[string]map[string]chan interface{}
	typeChsLock  sync.RWMutex
}

// NewChannelContext creates and returns object of new channel context
// TODO: Singleton
func NewChannelContext() *ChannelContext {
	channelMap := make(map[string]chan interface{})
	moduleChannels := make(map[string]map[string]chan interface{})

	return &ChannelContext{
		channels:     channelMap,
		typeChannels: moduleChannels,
	}
}

// Cleanup close modules
func (ctx *ChannelContext) Cleanup(module string) {
	if channel := ctx.getChannel(module); channel != nil {
		ctx.delChannel(module)
		// decrease probable exception of channel closing
		time.Sleep(20 * time.Millisecond)
		close(channel)
	}
}

// Send send msg to a module. Todo: do not stuck
func (ctx *ChannelContext) Send(module string, message interface{}) {
	// avoid exception because of channel colsing
	// TODO: need reconstruction
	defer func() {
		if exception := recover(); exception != nil {
			klog.Warningf("Recover when send message, exception: %+v", exception)
		}
	}()

	if channel := ctx.getChannel(module); channel != nil {
		channel <- message
		return
	}
	klog.Warningf("Get bad module name :%s when send message, do nothing", module)
}

// Receive msg from channel of module
func (ctx *ChannelContext) Receive(module string) (interface{}, error) {
	if channel := ctx.getChannel(module); channel != nil {
		content := <-channel
		return content, nil
	}

	klog.Warningf("Failed to get channel for module:%s when receive message", module)
	return nil, fmt.Errorf("failed to get channel for module(%s)", module)
}

// SendToGroup send msg to modules. Todo: do not stuck
func (ctx *ChannelContext) SendToGroup(moduleType string, message interface{}) {
	// avoid exception because of channel closing
	// TODO: need reconstruction
	defer func() {
		if exception := recover(); exception != nil {
			klog.Warningf("Recover when sendToGroup message, exception: %+v", exception)
		}
	}()

	send := func(ch chan interface{}) {
		select {
		case ch <- message:
		default:
			klog.Warningf("the message channel is full, message: %+v", message)
			select {
			case ch <- message:
			}
		}
	}
	if channelList := ctx.getTypeChannel(moduleType); channelList != nil {
		for _, channel := range channelList {
			go send(channel)
		}
		return
	}
	klog.Warningf("Get bad module type:%s when sendToGroup message, do nothing", moduleType)
}

// New Channel
func (ctx *ChannelContext) newChannel() chan interface{} {
	channel := make(chan interface{}, ChannelSizeDefault)
	return channel
}

// getChannel return chan
func (ctx *ChannelContext) getChannel(module string) chan interface{} {
	ctx.chsLock.RLock()
	defer ctx.chsLock.RUnlock()

	if _, exist := ctx.channels[module]; exist {
		return ctx.channels[module]
	}

	klog.Warningf("Failed to get channel, type:%s", module)
	return nil
}

// addChannel return chan
func (ctx *ChannelContext) addChannel(module string, moduleCh chan interface{}) {
	ctx.chsLock.Lock()
	defer ctx.chsLock.Unlock()

	ctx.channels[module] = moduleCh
}

// deleteChannel by module name
func (ctx *ChannelContext) delChannel(module string) {
	// delete module channel from channels map
	ctx.chsLock.Lock()
	_, exist := ctx.channels[module]
	if !exist {
		klog.Warningf("Failed to get channel, module:%s", module)
		return
	}
	delete(ctx.channels, module)

	ctx.chsLock.Unlock()

	// delete module channel from typechannels map
	ctx.typeChsLock.Lock()
	for _, moduleMap := range ctx.typeChannels {
		if _, exist := moduleMap[module]; exist {
			delete(moduleMap, module)
			break
		}
	}
	ctx.typeChsLock.Unlock()
}

// getTypeChannel return chan
func (ctx *ChannelContext) getTypeChannel(moduleType string) map[string]chan interface{} {
	ctx.typeChsLock.RLock()
	defer ctx.typeChsLock.RUnlock()

	if _, exist := ctx.typeChannels[moduleType]; exist {
		return ctx.typeChannels[moduleType]
	}

	klog.Warningf("Failed to get type channel, type:%s", moduleType)
	return nil
}

func (ctx *ChannelContext) getModuleByChannel(ch chan interface{}) string {
	ctx.chsLock.RLock()
	defer ctx.chsLock.RUnlock()

	for module, channel := range ctx.channels {
		if channel == ch {
			return module
		}
	}

	klog.Warning("Failed to get module by channel")
	return ""
}

// addTypeChannel put modules into moduleType map
func (ctx *ChannelContext) addTypeChannel(module, group string, moduleCh chan interface{}) {
	ctx.typeChsLock.Lock()
	defer ctx.typeChsLock.Unlock()

	if _, exist := ctx.typeChannels[group]; !exist {
		ctx.typeChannels[group] = make(map[string]chan interface{})
	}
	ctx.typeChannels[group][module] = moduleCh
}

//AddModule adds module into module context
func (ctx *ChannelContext) AddModule(module string) {
	channel := ctx.newChannel()
	ctx.addChannel(module, channel)
}

//AddModuleGroup adds modules into module context group
func (ctx *ChannelContext) AddModuleGroup(module, group string) {
	if channel := ctx.getChannel(module); channel != nil {
		ctx.addTypeChannel(module, group, channel)
		return
	}
	klog.Warningf("Get bad module name %s when addmodulegroup", module)
}
