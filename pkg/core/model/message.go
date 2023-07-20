package model

import(
	"time"
	"reflect"

	"github.com/satori/go.uuid"
)

const (
	ResponseOperation      = "response"
)

type MessageHeader struct {
	// the message id
	ID string `json:"ID,omitempty"`
	// the response message parentid must be same with message received
	// please use NewRespByMessage to new response message
	ParentID string `json:"parent_msg_id,omitempty"`
	// the time of creating
	Timestamp int64 `json:"timestamp,omitempty"`
	// specific resource version for the message, if any.
	ResourceVersion string `json:"resourceversion,omitempty"`
	// the flag will be set in sendsync
	Sync bool `json:"sync,omitempty"`
	//tag for other need
	Tag	 string  `json:"tag,omitempty"`
	// message type indicates the context type that delivers the message, such as channel, unixsocket, etc.
	// if the value is empty, the channel context type will be used.
	MessageType string `json:"type,omitempty"`	
}

//MessageRoute contains structure of message
type MessageRoute struct {
	// where the message come from
	Source string `json:"source,omitempty"`
	// where the message come to
	Destination string `json:"target,omitempty"`
	// where the message will broadcasted to
	Group string `json:"group, omitempty"`

	// what's the operation on resource
	Operation string `json:"operation,omitempty"`
	// what's the resource want to operate
	Resource string `json:"resource,omitempty"`
}

// Message struct
type Message struct {
	Header  MessageHeader `json:"header"`
	Router  MessageRoute  `json:"route,omitempty"`
	Content interface{}   `json:"content,omitempty"`
}

//GetID returns message ID
func (msg *Message) GetID() string {
	return msg.Header.ID
}

// GetParentID returns message parent id
func (msg *Message) GetParentID() string {
	return msg.Header.ParentID
}

//BuildRouter sets route and resource operation in message
func (msg *Message) BuildRouter(source, group, target, res, opr string) *Message {
	msg.SetRoute(source, group, target)
	msg.SetResourceOperation(res, opr)
	return msg
}

//SetResourceOperation sets router resource and operation in message
func (msg *Message) SetResourceOperation(res, opr string) *Message {
	msg.Router.Resource = res
	msg.Router.Operation = opr
	return msg
}
//SetRoute sets router source and group in message
func (msg *Message) SetRoute(source, group, target string) *Message {
	msg.Router.Source = source
	msg.Router.Group = group
	msg.Router.Destination = target
	return msg
}

// SetResourceVersion sets resource version in message header
func (msg *Message) SetResourceVersion(resourceVersion string) *Message {
	msg.Header.ResourceVersion = resourceVersion
	return msg
}

//SetTag set message tags
func (msg *Message) SetTag(tag string) *Message {
	msg.Header.Tag = tag
	return msg
}

//GetTag get message tags
func (msg *Message) GetTag() string {
	return msg.Header.Tag
}

// SetType set message context type
func (msg *Message) SetType(msgType string) *Message {
	msg.Header.MessageType = msgType
	return msg
}

// SetDestination set destination
func (msg *Message) SetDestination(dest string) *Message {
	msg.Router.Destination = dest
	return msg
}

//GetTimestamp returns message timestamp
func (msg *Message) GetTimestamp() int64 {
	return msg.Header.Timestamp
}

//GetContent returns message content
func (msg *Message) GetContent() interface{} {
	return msg.Content
}

//GetResource returns message route resource
func (msg *Message) GetResource() string {
	return msg.Router.Resource
}

//GetOperation returns message route operation string
func (msg *Message) GetOperation() string {
	return msg.Router.Operation
}

//GetSource returns message route source string
func (msg *Message) GetSource() string {
	return msg.Router.Source
}

//GetGroup returns message route group
func (msg *Message) GetGroup() string {
	return msg.Router.Group
}

//GetTarget returns message route Target
func (msg *Message) GetTarget() string {
	return msg.Router.Destination
}

// GetDestination get destination
func (msg *Message) GetDestination() string {
	return msg.Router.Destination
}

// IsEmpty is empty
func (msg *Message) IsEmpty() bool {
	return reflect.DeepEqual(msg, &Message{})
}

// IsSync : msg.Header.Sync will be set in sendsync
func (msg *Message) IsSync() bool {
	return msg.Header.Sync
}

//UpdateID returns message object updating its ID
func (msg *Message) UpdateID() *Message {
	msg.Header.ID = uuid.NewV4().String()
	return msg
}

// BuildHeader builds message header. You can also use for updating message header
func (msg *Message) BuildHeader(ID, parentID string, timestamp int64) *Message {
	msg.Header.ID = ID
	msg.Header.ParentID = parentID
	msg.Header.Timestamp = timestamp
	return msg
}

//FillBody fills message  content that you want to send
func (msg *Message) FillBody(content interface{}) *Message {
	msg.Content = content
	return msg
}

// GetType get message context type
func (msg *Message) GetType() string {
	return msg.Header.MessageType
}

// NewRawMessage returns a new raw message:
// model.NewRawMessage().BuildHeader().BuildRouter().FillBody()
func NewRawMessage() *Message {
	return &Message{}
}

// NewMessage returns a new basic message:
// model.NewMessage().BuildRouter().FillBody()
func NewMessage(parentID string) *Message {
	msg := &Message{}
	msg.Header.ID = uuid.NewV4().String()
	msg.Header.ParentID = parentID
	msg.Header.Timestamp = time.Now().UnixNano() / 1e6
	return msg
}

// NewRespByMessage returns a new response message by a message received
func (msg *Message) NewRespByMessage(message *Message, content interface{}) *Message {
	return NewMessage(message.GetID()).SetRoute(message.GetDestination(),
		message.GetGroup(), message.GetSource()).
		SetResourceOperation(message.GetResource(), ResponseOperation).
		SetType(message.GetType()).
		FillBody(content)
}
