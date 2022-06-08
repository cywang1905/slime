/*
* @Author: yangdihang
* @Date: 2020/6/2
 */

package controllers

import (
	"regexp"
	"sync"

	"slime.io/slime/framework/model"
	modmodel "slime.io/slime/modules/lazyload/model"
)

var log = modmodel.ModuleLog.WithField(model.LogFieldKeyPkg, "controllers")

type Diff struct {
	Deleted []string
	Added   []string
}

type LabelItem struct {
	Name  string
	Value string
}

type NsSvcCache struct {
	Data map[string]map[string]struct{}
	sync.RWMutex
}

type LabelSvcCache struct {
	Data map[LabelItem]map[string]struct{}
	sync.RWMutex
}

type PortProtocolCache struct {
	Data map[int32]map[Protocol]int
	sync.RWMutex
}

type domainAliasRule struct {
	pattern   string
	templates []string
	re        *regexp.Regexp
}

type PortProtocol string

const (
	HTTP    PortProtocol = "http"
	HTTP2   PortProtocol = "http2"
	HTTPS   PortProtocol = "https"
	TCP     PortProtocol = "tcp"
	TLS     PortProtocol = "tls"
	GRPC    PortProtocol = "grpc"
	GRPCWeb PortProtocol = "grpc-web"
	Mongo   PortProtocol = "mongo"
	MySQL   PortProtocol = "mysql"
	Redis   PortProtocol = "redis"
	Dubbo   PortProtocol = "dubbo"
	Unknown PortProtocol = "unknown"
)

// Protocol is the protocol associated with the port
type Protocol int

const (
	ProtocolUnknown = iota
	ProtocolTCP
	ProtocolHTTP
	ProtocolDubbo
)
