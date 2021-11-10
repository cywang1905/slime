package trigger

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"slime.io/slime/framework/util"
	"time"

	"slime.io/slime/framework/model/metric"
)

type Trigger struct {
	WatchersMap      map[watch.Interface]chan struct{}
	TimersMap        map[time.Duration]chan struct{}
	TriggerEventChan chan metric.TriggerEvent
	Config           *metric.ControllerConfig
}

func NewTrigger(c metric.Controller) metric.Trigger {
	log := log.WithField("reporter", "Trigger").WithField("function", "NewTrigger")
	log.Infof("init new metric trigger")
	return &Trigger{
		//WatchersMap:      make(map[watch.Interface]chan struct{}),
		//TimersMap:        make(map[time.Duration]chan struct{}),
		TriggerEventChan: make(chan metric.TriggerEvent),
		Config:           c.GetConfig(),
	}
}

func (t *Trigger) GetTriggerEventChan() chan metric.TriggerEvent {
	return t.TriggerEventChan
}

func (t *Trigger) Stop() error {
	for _, ch := range t.WatchersMap {
		go func(c chan struct{}) {
			c <- struct{}{}
		}(ch)
	}
	for _, ch := range t.TimersMap {
		go func(c chan struct{}) {
			c <- struct{}{}
		}(ch)
	}
	return nil
}

func (t *Trigger) Start(triggersMap interface{}) error {
	log := log.WithField("reporter", "Trigger").WithField("function", "Start")
	// change to switch triggersMap.(type)
	switch triggersMap.(type) {
	// watchers
	case map[schema.GroupVersionKind]metric.UserKind:
		// init again
		t.WatchersMap = make(map[watch.Interface]chan struct{})

		wm, _ := triggersMap.(map[schema.GroupVersionKind]metric.UserKind)
		for gvk := range wm {
			dc := t.Config.DynamicClient
			gvr, _ := meta.UnsafeGuessKindToResource(gvk)
			lw := &cache.ListWatch{
				ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
					return dc.Resource(gvr).List(options)
				},
				WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
					return dc.Resource(gvr).Watch(options)
				},
			}
			watcher := util.ListWatcher(context.Background(), lw)
			t.WatchersMap[watcher] = make(chan struct{})
			log.Infof("add watcher %s to metric trigger", gvr.String())
		}

	// timers
	case map[time.Duration]metric.UserKind:
		// init again
		t.TimersMap = make(map[time.Duration]chan struct{})

		tm, _ := triggersMap.(map[time.Duration]metric.UserKind)
		for td := range tm {
			t.TimersMap[td] = make(chan struct{})
			log.Infof("add timer %s to metric trigger", td.String())
		}

	default:
		return errors.New("start error: wrong triggersMap kind")
	}

	for wat, channel := range t.WatchersMap {
		go func(w watch.Interface, ch chan struct{}) {
			for {
				select {
				case <-ch:
					log.Debugf("stop a watcher")
					return
				case e, ok := <-w.ResultChan():
					log.Debugf("got watcher event: type %s, kind %s", e.Type, e.Object.GetObjectKind().GroupVersionKind().String())
					if !ok {
						log.Warningf("a result chan of watcher is closed, break process loop")
						return
					}
					object, ok := e.Object.(*unstructured.Unstructured)
					if !ok {
						log.Errorf("invalid type of object in watcher event")
						continue
					}
					mi := metric.MetaInfo{
						GVK: e.Object.GetObjectKind().GroupVersionKind(),
						NN: types.NamespacedName{
							Name:      object.GetName(),
							Namespace: object.GetNamespace(),
						},
					}
					te := metric.TriggerEvent{
						TriggerEventKind: metric.TriggerEventKind_Watcher,
						Meta:             mi,
					}
					t.TriggerEventChan <- te
					//log.Debugf("sent watcher event to controller: type %s, kind %s", e.Type, e.Object.GetObjectKind().GroupVersionKind().String())
				}
			}
		}(wat, channel)
	}

	for timer, channel := range t.TimersMap {
		ticker := time.NewTicker(timer)
		go func(ti time.Duration, ch chan struct{}) {
			for {
				select {
				case <-ch:
					log.Debugf("stop a timer")
					return
				case <-ticker.C:
					log.Debugf("got timer event: duration %s", ti.String())
					te := metric.TriggerEvent{
						TriggerEventKind: metric.TriggerEventKind_Timer,
						Meta:             ti,
					}
					t.TriggerEventChan <- te
					//log.Debugf("sent timer event to controller: duration %s", timer.String())
				}
			}
		}(timer, channel)
	}
	return nil
}
