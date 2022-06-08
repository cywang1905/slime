package controllers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	config "slime.io/slime/framework/apis/config/v1alpha1"
	"slime.io/slime/framework/bootstrap"
	"slime.io/slime/modules/lazyload/api/v1alpha1"
	"strings"
)

var (
	slimeBootGvr = schema.GroupVersionResource{
		Group:    "config.netease.com",
		Version:  "v1alpha1",
		Resource: "slimeboots",
	}
)

type SlimeBoot struct {
	metav1.TypeMeta   `json:",inline,omitempty" yaml:",inline,omitempty"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec              SlimeBootSpec `json:"spec" yaml:"spec"`
}

type SlimeBootSpec struct {
	Modules   []Module  `json:"module" yaml:"module"`
	Image     Image     `json:"image" yaml:"image"`
	Component Component `json:"component" yaml:"component"`
}

type Component struct {
	GlobalSidecar GlobalSidecar `json:"globalSidecar" yaml:"globalSidecar"`
}

type GlobalSidecar struct {
	Enable        bool          `json:"enable" yaml:"enable"`
	Image         Image         `json:"image" yaml:"image"`
	ProbePort     int           `json:"probePort" yaml:"probePort"`
	Resources     Resources     `json:"resources" yaml:"resources"`
	SidecarInject SidecarInject `json:"sidecarInject" yaml:"sidecarInject"`
}

type SidecarInject struct {
	Enable bool              `json:"enable" yaml:"enable"`
	Mode   string            `json:"mode" yaml:"mode"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type Resources struct {
	Limits   Limits   `json:"limits" yaml:"limits"`
	Requests Requests `json:"requests" yaml:"requests"`
}

type Limits struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

type Requests struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

type Image struct {
	PullPolicy string `json:"pullPolicy" yaml:"pullPolicy"`
	Repository string `json:"repository" yaml:"repository"`
	Tag        string `json:"tag" yaml:"tag"`
}

type Module struct {
	Plugin  *config.Plugin     `protobuf:"bytes,1,opt,name=plugin,proto3" json:"plugin,omitempty"`
	Limiter *config.Limiter    `protobuf:"bytes,2,opt,name=limiter,proto3" json:"limiter,omitempty"`
	Global  *config.Global     `protobuf:"bytes,3,opt,name=global,proto3" json:"global,omitempty"`
	Fence   *config.Fence      `protobuf:"bytes,4,opt,name=fence,proto3" json:"fence,omitempty"`
	Metric  *config.Metric     `protobuf:"bytes,6,opt,name=metric,proto3" json:"metric,omitempty"`
	Name    string             `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	Enable  bool               `protobuf:"varint,7,opt,name=enable,proto3" json:"enable,omitempty"`
	General *v1alpha1.Fence    `protobuf:"bytes,8,opt,name=general,proto3" json:"general,omitempty"`
	Bundle  *config.Bundle     `protobuf:"bytes,9,opt,name=bundle,proto3" json:"bundle,omitempty"`
	Mode    config.Config_Mode `protobuf:"varint,10,opt,name=mode,proto3,enum=slime.config.v1alpha1.Config_Mode" json:"mode,omitempty"`
	// like bundle item kind, necessary if not bundle
	Kind string `protobuf:"bytes,11,opt,name=kind,proto3" json:"kind,omitempty"`
}

func updateResources(wormholePort []string, env bootstrap.Environment) {
	dynCli := env.DynamicClient

	// get slimeboot cr name, cr name should equals deployment name
	slimeBootNs := os.Getenv("WATCH_NAMESPACE")
	slimeBootName := strings.Split(os.Getenv("POD_NAME"), "-")[0]

	// Unstructured
	utd, err := dynCli.Resource(slimeBootGvr).Namespace(slimeBootNs).Get(slimeBootName, metav1.GetOptions{}, "")
	if err != nil {
		log.Errorf("get slimeboot [%s/%s] error: %+v", slimeBootNs, slimeBootName, err)
		return
	}
	// Unstructured -> SlimeBoot
	var slimeBoot SlimeBoot
	if err = runtime.DefaultUnstructuredConverter.FromUnstructured(utd.UnstructuredContent(), &slimeBoot); err != nil {
		log.Errorf("convert slimeboot %s/%s to structured error: %v", slimeBootNs, slimeBootName, err)
		return
	}

	for _, module := range slimeBoot.Spec.Modules {
		if module.Kind != "lazyload" {
			continue
		}
		// update wormholePort
		module.General.WormholePort = wormholePort
	}

	// update slimeBoot
	obj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&slimeBoot)
	if err != nil {
		log.Errorf("convert slimeboot %s/%s to unstructured error: %+v", slimeBootNs, slimeBootName, err)
		return
	}
	utd.SetUnstructuredContent(obj)
	utd, err = dynCli.Resource(slimeBootGvr).Namespace(slimeBootNs).Update(utd, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("update slimeboot %s/%s error: %+v", slimeBootNs, slimeBootName, err)
		return
	}
	log.Infof("update slimeboot %s/%s successfully, new wormholePort: %+v", slimeBootNs, slimeBootName, wormholePort)
	return
}
