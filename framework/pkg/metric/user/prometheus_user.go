package user

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	bootconfig "slime.io/slime/framework/apis/config/v1alpha1"
	"slime.io/slime/framework/model/metric"
	"sync"
)

type PrometheusUser struct {
	Name                    string
	UserKind                metric.UserKind
	Source                  metric.Source
	InterestConfig          PrometheusUserInterestConfig
	InterestMeta            cmap.ConcurrentMap
	NeedUpdateMetricHandler func(event metric.TriggerEvent) map[string]*bootconfig.Prometheus_Source_Handler
	UserEventChan           chan metric.UserEvent
	//Env                     bootstrap.Environment
	sync.RWMutex
}

type PrometheusUserInterestConfig struct {
	InterestKinds []schema.GroupVersionKind
}

func (u *PrometheusUser) GetName() string {
	return u.Name
}

func (u *PrometheusUser) GetUserKind() metric.UserKind {
	return u.UserKind
}

func (u *PrometheusUser) GetSource() metric.Source {
	return u.Source
}

func (u *PrometheusUser) GetInterestConfig() interface{} {
	return u.InterestConfig
}

func (u *PrometheusUser) GetUserEventChan() <-chan metric.UserEvent {
	return u.UserEventChan
}

// AddInterestMeta InterestMetaAdd InterestAdd add resource to interest
func (u *PrometheusUser) AddInterestMeta(meta interface{}) error {
	//log := log.WithField("reporter", "PrometheusUser").WithField("function", "AddInterestMeta")
	nn, ok := meta.(types.NamespacedName)
	if !ok {
		return errors.New("interest meta add error: wrong meta type")
	}
	u.InterestMeta.Set(nn.Namespace+"/"+nn.Name, true)
	//log.Debugf("PrometheusUser %s added interest meta %s successfully", u.Name, nn.String())
	return nil
}

// RemoveInterestMeta InterestMetaRemove InterestRemove remove from interest, the resource will be deleted by k8s
func (u *PrometheusUser) RemoveInterestMeta(meta interface{}) error {
	//log := log.WithField("reporter", "PrometheusUser").WithField("function", "RemoveInterestMeta")
	nn, ok := meta.(types.NamespacedName)
	if !ok {
		return errors.New("interest meta remove error: wrong meta type")
	}
	u.InterestMeta.Pop(nn.Namespace + "/" + nn.Name)
	//log.Debugf("PrometheusUser %s removed interest meta %s successfully", u.Name, nn.String())
	return nil
}

// TriggerEventHandler handle trigger event, send user event to module
func (u *PrometheusUser) TriggerEventHandler(te metric.TriggerEvent) error {
	log := log.WithField("reporter", "PrometheusUser").WithField("function", "TriggerEventHandler")
	if te.TriggerEventKind != metric.TriggerEventKind_Watcher {
		return errors.New("handle trigger event error: wrong trigger event kind")
	}

	mi, ok := te.Meta.(metric.MetaInfo)
	if !ok {
		return errors.New("handle trigger event error: wrong watcher trigger event meta type")
	}
	handlers := u.NeedUpdateMetricHandler(te)
	if handlers == nil {
		return nil
	}

	log.Debugf("PrometheusUser %s is handlering watcher event", u.Name)
	_, metaCheck := u.InterestMeta.Get(mi.NN.Namespace + "/" + mi.NN.Name)
	if !u.gvkCheck(mi.GVK) || !metaCheck {
		log.Debugf("PrometheusUser %s handled watcher event successfully, skip", u.Name)
		return nil
	}
	material, err := u.Source.GetMetrics(handlers, mi.NN)
	if err != nil {
		return err
	}
	u.UserEventChan <- metric.UserEvent{
		Meta:     mi.NN,
		Material: material,
	}
	log.Debugf("PrometheusUser %s handled watcher event successfully, sent", u.Name)

	return nil
}

func (u *PrometheusUser) gvkCheck(gvk schema.GroupVersionKind) bool {
	for _, g := range u.InterestConfig.InterestKinds {
		if g == gvk {
			return true
		}
	}
	return false
}
