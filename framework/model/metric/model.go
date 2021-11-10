package metric

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
)

type UserKind int32

const (
	UserKind_GVKEvent UserKind = 1
	UserKind_Timer    UserKind = 2
	//UserKind_LogEvent UserKind = 3
)

type SourceKind int32

const (
	SourceKind_Prometheus SourceKind = 1
	//SourceKind_AccessLog  SourceKind = 2
)

type TriggerKind int32

const (
	TriggerEventKind_Watcher TriggerKind = 1
	TriggerEventKind_Timer   TriggerKind = 2
	//TriggerEventKind_Logger  TriggerKind = 3
)

type User interface {
	GetName() string
	GetUserKind() UserKind
	GetSource() Source
	GetInterestConfig() interface{}
	GetUserEventChan() <-chan UserEvent
	AddInterestMeta(meta interface{}) error
	RemoveInterestMeta(meta interface{}) error
	TriggerEventHandler(te TriggerEvent) error
}

type Source interface {
	GetName() string
	GetKind() SourceKind
	GetMetrics(hs interface{}, meta interface{}) (interface{}, error)
}

type Controller interface {
	AddUsers(users []User) error
	RemoveUsers(users []User) error
	Stop()
	GetConfig() *ControllerConfig
	//GetTriggersMap() map[UserKind][]interface{}
}

type Trigger interface {
	GetTriggerEventChan() chan TriggerEvent
	Stop() error
	Start(triggersMap interface{}) error
}

type ControllerConfig struct {
	DynamicClient dynamic.Interface
}

type UserEvent struct {
	Meta     interface{}
	Material interface{}
}

type TriggerEvent struct {
	TriggerEventKind TriggerKind
	Meta             interface{}
}

type MetaInfo struct {
	GVK schema.GroupVersionKind
	NN  types.NamespacedName
}
