// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pkg/apis/microservice/v1alpha1/envoy_plugin.proto

package v1alpha1

import (
	fmt "fmt"
	proto "github.com/gogo/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// `WorkloadSelector` specifies the criteria used to determine if the
// `Gateway`, `Sidecar`, or `EnvoyFilter` or `ServiceEntry`
// configuration can be applied to a proxy. The matching criteria
// includes the metadata associated with a proxy, workload instance
// info such as labels attached to the pod/VM, or any other info that
// the proxy provides to Istio during the initial handshake. If
// multiple conditions are specified, all conditions need to match in
// order for the workload instance to be selected. Currently, only
// label based selection mechanism is supported.
type WorkloadSelector struct {
	// One or more labels that indicate a specific set of pods/VMs
	// on which the configuration should be applied. The scope of
	// label search is restricted to the configuration namespace in which the
	// the resource is present.
	Labels               map[string]string `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *WorkloadSelector) Reset()         { *m = WorkloadSelector{} }
func (m *WorkloadSelector) String() string { return proto.CompactTextString(m) }
func (*WorkloadSelector) ProtoMessage()    {}
func (*WorkloadSelector) Descriptor() ([]byte, []int) {
	return fileDescriptor_35868063e6636962, []int{0}
}
func (m *WorkloadSelector) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WorkloadSelector.Unmarshal(m, b)
}
func (m *WorkloadSelector) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WorkloadSelector.Marshal(b, m, deterministic)
}
func (m *WorkloadSelector) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WorkloadSelector.Merge(m, src)
}
func (m *WorkloadSelector) XXX_Size() int {
	return xxx_messageInfo_WorkloadSelector.Size(m)
}
func (m *WorkloadSelector) XXX_DiscardUnknown() {
	xxx_messageInfo_WorkloadSelector.DiscardUnknown(m)
}

var xxx_messageInfo_WorkloadSelector proto.InternalMessageInfo

func (m *WorkloadSelector) GetLabels() map[string]string {
	if m != nil {
		return m.Labels
	}
	return nil
}

type EnvoyPlugin struct {
	WorkloadSelector *WorkloadSelector `protobuf:"bytes,9,opt,name=workload_selector,json=workloadSelector,proto3" json:"workload_selector,omitempty"`
	// route level plugin
	Route []string `protobuf:"bytes,1,rep,name=route,proto3" json:"route,omitempty"`
	// host level plugin
	Host []string `protobuf:"bytes,2,rep,name=host,proto3" json:"host,omitempty"`
	// service level plugin
	Service []string  `protobuf:"bytes,3,rep,name=service,proto3" json:"service,omitempty"`
	Plugins []*Plugin `protobuf:"bytes,4,rep,name=plugins,proto3" json:"plugins,omitempty"`
	// which gateway should use this plugin setting
	Gateway []string `protobuf:"bytes,5,rep,name=gateway,proto3" json:"gateway,omitempty"`
	// which user should use this plugin setting
	User []string `protobuf:"bytes,6,rep,name=user,proto3" json:"user,omitempty"`
	// Deprecated
	IsGroupSetting bool `protobuf:"varint,7,opt,name=isGroupSetting,proto3" json:"isGroupSetting,omitempty"`
	// listener level
	Listener             []*EnvoyPlugin_Listener `protobuf:"bytes,8,rep,name=listener,proto3" json:"listener,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *EnvoyPlugin) Reset()         { *m = EnvoyPlugin{} }
func (m *EnvoyPlugin) String() string { return proto.CompactTextString(m) }
func (*EnvoyPlugin) ProtoMessage()    {}
func (*EnvoyPlugin) Descriptor() ([]byte, []int) {
	return fileDescriptor_35868063e6636962, []int{1}
}
func (m *EnvoyPlugin) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnvoyPlugin.Unmarshal(m, b)
}
func (m *EnvoyPlugin) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnvoyPlugin.Marshal(b, m, deterministic)
}
func (m *EnvoyPlugin) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnvoyPlugin.Merge(m, src)
}
func (m *EnvoyPlugin) XXX_Size() int {
	return xxx_messageInfo_EnvoyPlugin.Size(m)
}
func (m *EnvoyPlugin) XXX_DiscardUnknown() {
	xxx_messageInfo_EnvoyPlugin.DiscardUnknown(m)
}

var xxx_messageInfo_EnvoyPlugin proto.InternalMessageInfo

func (m *EnvoyPlugin) GetWorkloadSelector() *WorkloadSelector {
	if m != nil {
		return m.WorkloadSelector
	}
	return nil
}

func (m *EnvoyPlugin) GetRoute() []string {
	if m != nil {
		return m.Route
	}
	return nil
}

func (m *EnvoyPlugin) GetHost() []string {
	if m != nil {
		return m.Host
	}
	return nil
}

func (m *EnvoyPlugin) GetService() []string {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *EnvoyPlugin) GetPlugins() []*Plugin {
	if m != nil {
		return m.Plugins
	}
	return nil
}

func (m *EnvoyPlugin) GetGateway() []string {
	if m != nil {
		return m.Gateway
	}
	return nil
}

func (m *EnvoyPlugin) GetUser() []string {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *EnvoyPlugin) GetIsGroupSetting() bool {
	if m != nil {
		return m.IsGroupSetting
	}
	return false
}

func (m *EnvoyPlugin) GetListener() []*EnvoyPlugin_Listener {
	if m != nil {
		return m.Listener
	}
	return nil
}

type EnvoyPlugin_Listener struct {
	Port                 uint32   `protobuf:"varint,1,opt,name=port,proto3" json:"port,omitempty"`
	Outbound             bool     `protobuf:"varint,2,opt,name=outbound,proto3" json:"outbound,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnvoyPlugin_Listener) Reset()         { *m = EnvoyPlugin_Listener{} }
func (m *EnvoyPlugin_Listener) String() string { return proto.CompactTextString(m) }
func (*EnvoyPlugin_Listener) ProtoMessage()    {}
func (*EnvoyPlugin_Listener) Descriptor() ([]byte, []int) {
	return fileDescriptor_35868063e6636962, []int{1, 0}
}
func (m *EnvoyPlugin_Listener) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnvoyPlugin_Listener.Unmarshal(m, b)
}
func (m *EnvoyPlugin_Listener) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnvoyPlugin_Listener.Marshal(b, m, deterministic)
}
func (m *EnvoyPlugin_Listener) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnvoyPlugin_Listener.Merge(m, src)
}
func (m *EnvoyPlugin_Listener) XXX_Size() int {
	return xxx_messageInfo_EnvoyPlugin_Listener.Size(m)
}
func (m *EnvoyPlugin_Listener) XXX_DiscardUnknown() {
	xxx_messageInfo_EnvoyPlugin_Listener.DiscardUnknown(m)
}

var xxx_messageInfo_EnvoyPlugin_Listener proto.InternalMessageInfo

func (m *EnvoyPlugin_Listener) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

func (m *EnvoyPlugin_Listener) GetOutbound() bool {
	if m != nil {
		return m.Outbound
	}
	return false
}

func init() {
	proto.RegisterType((*WorkloadSelector)(nil), "slime.microservice.v1alpha1.WorkloadSelector")
	proto.RegisterMapType((map[string]string)(nil), "slime.microservice.v1alpha1.WorkloadSelector.LabelsEntry")
	proto.RegisterType((*EnvoyPlugin)(nil), "slime.microservice.v1alpha1.EnvoyPlugin")
	proto.RegisterType((*EnvoyPlugin_Listener)(nil), "slime.microservice.v1alpha1.EnvoyPlugin.Listener")
}

func init() {
	proto.RegisterFile("pkg/apis/microservice/v1alpha1/envoy_plugin.proto", fileDescriptor_35868063e6636962)
}

var fileDescriptor_35868063e6636962 = []byte{
	// 412 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0xc1, 0x6a, 0xdb, 0x40,
	0x10, 0x45, 0x76, 0x6c, 0xcb, 0x63, 0x5a, 0xdc, 0x25, 0x87, 0xc5, 0xbd, 0x88, 0x14, 0x8a, 0x2e,
	0x91, 0x50, 0x72, 0x69, 0x02, 0xbd, 0x14, 0x42, 0x2f, 0x29, 0xb4, 0x9b, 0x43, 0x21, 0x17, 0xb3,
	0x76, 0x06, 0x65, 0xf1, 0x5a, 0x2b, 0x76, 0x57, 0x32, 0xfa, 0x9b, 0xfe, 0x50, 0xff, 0xa9, 0x68,
	0x57, 0x32, 0xc6, 0x07, 0x95, 0xdc, 0xe6, 0x3d, 0x76, 0xde, 0x7b, 0x33, 0x3b, 0x90, 0x95, 0xbb,
	0x3c, 0xe5, 0xa5, 0x30, 0xe9, 0x5e, 0x6c, 0xb5, 0x32, 0xa8, 0x6b, 0xb1, 0xc5, 0xb4, 0xce, 0xb8,
	0x2c, 0x5f, 0x79, 0x96, 0x62, 0x51, 0xab, 0x66, 0x5d, 0xca, 0x2a, 0x17, 0x45, 0x52, 0x6a, 0x65,
	0x15, 0xf9, 0x68, 0xa4, 0xd8, 0x63, 0x72, 0xfa, 0x3e, 0xe9, 0xdf, 0xaf, 0x6e, 0xff, 0xa3, 0xe7,
	0x95, 0xd6, 0x7b, 0x5e, 0xf0, 0x1c, 0xb5, 0x57, 0xbc, 0xfa, 0x13, 0xc0, 0xf2, 0xb7, 0xd2, 0x3b,
	0xa9, 0xf8, 0xcb, 0x13, 0x4a, 0xdc, 0x5a, 0xa5, 0xc9, 0x2f, 0x98, 0x4a, 0xbe, 0x41, 0x69, 0x68,
	0x10, 0x8d, 0xe3, 0xc5, 0xcd, 0x5d, 0x32, 0xe0, 0x9b, 0x9c, 0xb7, 0x27, 0x8f, 0xae, 0xf7, 0xa1,
	0xb0, 0xba, 0x61, 0x9d, 0xd0, 0xea, 0x0e, 0x16, 0x27, 0x34, 0x59, 0xc2, 0x78, 0x87, 0x0d, 0x0d,
	0xa2, 0x20, 0x9e, 0xb3, 0xb6, 0x24, 0x97, 0x30, 0xa9, 0xb9, 0xac, 0x90, 0x8e, 0x1c, 0xe7, 0xc1,
	0xfd, 0xe8, 0x4b, 0x70, 0xf5, 0x77, 0x0c, 0x8b, 0x87, 0x76, 0x17, 0x3f, 0xdd, 0x00, 0xe4, 0x19,
	0x3e, 0x1c, 0x3a, 0xcb, 0xb5, 0xe9, 0x3c, 0xe9, 0x3c, 0x0a, 0xe2, 0xc5, 0xcd, 0xf5, 0x9b, 0x82,
	0xb2, 0xe5, 0xe1, 0x7c, 0xf2, 0x4b, 0x98, 0x68, 0x55, 0x59, 0x74, 0x83, 0xcf, 0x99, 0x07, 0x84,
	0xc0, 0xc5, 0xab, 0x32, 0x96, 0x8e, 0x1c, 0xe9, 0x6a, 0x42, 0x61, 0xd6, 0x19, 0xd0, 0xb1, 0xa3,
	0x7b, 0x48, 0xbe, 0xc2, 0xcc, 0xaf, 0xda, 0xd0, 0x0b, 0xb7, 0xbe, 0x4f, 0x83, 0xa9, 0xfc, 0x54,
	0xac, 0xef, 0x69, 0x85, 0x73, 0x6e, 0xf1, 0xc0, 0x1b, 0x3a, 0xf1, 0xc2, 0x1d, 0x6c, 0x63, 0x54,
	0x06, 0x35, 0x9d, 0xfa, 0x18, 0x6d, 0x4d, 0x3e, 0xc3, 0x7b, 0x61, 0xbe, 0x6b, 0x55, 0x95, 0x4f,
	0x68, 0xad, 0x28, 0x72, 0x3a, 0x8b, 0x82, 0x38, 0x64, 0x67, 0x2c, 0xf9, 0x01, 0xa1, 0x14, 0xc6,
	0x62, 0x81, 0x9a, 0x86, 0x2e, 0x55, 0x36, 0x98, 0xea, 0x64, 0xe1, 0xc9, 0x63, 0xd7, 0xc8, 0x8e,
	0x12, 0xab, 0x7b, 0x08, 0x7b, 0xb6, 0x8d, 0x55, 0x2a, 0x6d, 0xdd, 0x67, 0xbe, 0x63, 0xae, 0x26,
	0x2b, 0x08, 0x55, 0x65, 0x37, 0xaa, 0x2a, 0x5e, 0xdc, 0x87, 0x86, 0xec, 0x88, 0xbf, 0xa5, 0xcf,
	0xd7, 0xde, 0x59, 0xa8, 0xd4, 0x15, 0xe9, 0xf0, 0xe1, 0x6e, 0xa6, 0xee, 0x54, 0x6f, 0xff, 0x05,
	0x00, 0x00, 0xff, 0xff, 0x35, 0x69, 0xa0, 0x91, 0x31, 0x03, 0x00, 0x00,
}
