package user

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	bootconfig "slime.io/slime/framework/apis/config/v1alpha1"
	"slime.io/slime/framework/bootstrap"
	"slime.io/slime/framework/model/metric"
	"strings"
	"sync"
	"time"
)

type TimerUser struct {
	Name                    string
	UserKind                metric.UserKind
	Source                  metric.Source
	InterestConfig          TimerUserInterestConfig
	InterestMeta            cmap.ConcurrentMap
	NeedUpdateMetricHandler func(event metric.TriggerEvent) map[string]*bootconfig.Prometheus_Source_Handler
	UserEventChan           chan metric.UserEvent
	Env                     *bootstrap.Environment
	sync.RWMutex
}

type TimerUserInterestConfig struct {
	//InterestKinds           []schema.GroupVersionKind
	InterestInterval time.Duration
}

func (u *TimerUser) GetName() string {
	return u.Name
}

func (u *TimerUser) GetUserKind() metric.UserKind {
	return u.UserKind
}

func (u *TimerUser) GetSource() metric.Source {
	return u.Source
}

func (u *TimerUser) GetInterestConfig() interface{} {
	return u.InterestConfig
}

func (u *TimerUser) GetUserEventChan() <-chan metric.UserEvent {
	return u.UserEventChan
}

// AddInterestMeta InterestMetaAdd InterestAdd add resource to interest
func (u *TimerUser) AddInterestMeta(meta interface{}) error {
	//log := log.WithField("reporter", "TimerUser").WithField("function", "AddInterestMeta")
	nn, ok := meta.(types.NamespacedName)
	if !ok {
		return errors.New("interest meta add error: wrong meta type")
	}
	u.InterestMeta.Set(nn.Namespace+"/"+nn.Name, true)
	//log.Debugf("TimerUser %s added interest meta %s successfully", u.Name, nn.String())
	return nil
}

// RemoveInterestMeta InterestMetaRemove InterestRemove remove from interest, the resource will be deleted by k8s
func (u *TimerUser) RemoveInterestMeta(meta interface{}) error {
	//log := log.WithField("reporter", "TimerUser").WithField("function", "RemoveInterestMeta")
	nn, ok := meta.(types.NamespacedName)
	if !ok {
		return errors.New("interest meta remove error: wrong meta type")
	}
	u.InterestMeta.Pop(nn.Namespace + "/" + nn.Name)
	//log.Debugf("TimerUser %s removed interest meta %s successfully", u.Name, nn.String())
	return nil
}

func (u *TimerUser) TriggerEventHandler(te metric.TriggerEvent) error {
	log := log.WithField("reporter", "TimerUser").WithField("function", "TriggerEventHandler")
	if te.TriggerEventKind != metric.TriggerEventKind_Timer {
		return errors.New("handle trigger event error: wrong trigger event kind")
	}

	t, ok := te.Meta.(time.Duration)
	if !ok {
		return errors.New("handle trigger event error: wrong timer trigger event meta type")
	}
	if t != u.InterestConfig.InterestInterval {
		return errors.New("handle trigger event error: wrong timer trigger event meta value")
	}
	handlers := u.NeedUpdateMetricHandler(te)
	if handlers == nil {
		return nil
	}

	log.Debugf("TimerUser %s is handlering timer event", u.Name)
	//if !u.gvkCheck(mi.GVK) {
	//	return errors.New("handle trigger event error: wrong trigger event meta info")
	//}
	u.RLock()
	defer u.RUnlock()
	for k := range u.InterestMeta.Items() {
		if index := strings.Index(k, "/"); index == -1 || index == len(k)-1 {
			continue
		} else {
			ns := k[:index]
			name := k[index+1:]
			nn := types.NamespacedName{Namespace: ns, Name: name}
			material, err := u.Source.GetMetrics(handlers, nn)
			if err != nil {
				return err
			}
			u.UserEventChan <- metric.UserEvent{
				Meta:     nn,
				Material: material,
			}
		}
	}
	log.Debugf("TimerUser %s handled timer event successfully", u.Name)
	return nil
}
