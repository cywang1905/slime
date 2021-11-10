package controller

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"slime.io/slime/framework/model/metric"
	"slime.io/slime/framework/pkg/metric/trigger"
	metricUser "slime.io/slime/framework/pkg/metric/user"
	"sync"
	"time"
)

type Controller struct {
	WatchersMap     map[schema.GroupVersionKind]metric.UserKind
	TimersMap       map[time.Duration]metric.UserKind
	WatcherUsersMap map[schema.GroupVersionKind][]metric.User // get user list to call user.TriggerEventHandler
	TimerUsersMap   map[time.Duration][]metric.User           // get user list to call user.TriggerEventHandler
	//TriggerEventChan <-chan metric.TriggerEvent
	Config   *metric.ControllerConfig
	Trigger  metric.Trigger
	StopChan chan struct{}
	sync.RWMutex
}

func NewController(config *metric.ControllerConfig) metric.Controller {
	log := log.WithField("reporter", "Controller").WithField("function", "NewController")
	log.Infof("init new metric controller")
	c := &Controller{
		WatchersMap:     make(map[schema.GroupVersionKind]metric.UserKind),
		TimersMap:       make(map[time.Duration]metric.UserKind),
		WatcherUsersMap: make(map[schema.GroupVersionKind][]metric.User),
		TimerUsersMap:   make(map[time.Duration][]metric.User),
		Config:          config,
	}
	c.Trigger = trigger.NewTrigger(c)
	return c
}

func (c *Controller) GetConfig() *metric.ControllerConfig {
	return c.Config
}

func (c *Controller) AddUsers(users []metric.User) error {
	log := log.WithField("reporter", "Controller").WithField("function", "AddUsers")
	// check userkind and stop related watchers
	if err := c.Trigger.Stop(); err != nil {
		return errors.New(fmt.Sprintf("stop trigger failed: %v", err))
	}

	// update related usersmap
	c.Lock()
	defer c.Unlock()

	for _, user := range users {
		switch user.GetUserKind() {
		case metric.UserKind_GVKEvent:
			config, ok := user.GetInterestConfig().(metricUser.PrometheusUserInterestConfig)
			if !ok {
				return errors.New("wrong InterestConfig type of user")
			}
			for _, gvk := range config.InterestKinds {
				if _, ok = c.WatchersMap[gvk]; !ok {
					log.Infof("add watcher key %s to metric controller", gvk.String())
					// update WatchersMap
					c.WatchersMap[gvk] = metric.UserKind_GVKEvent
					// update WatcherUsersMap
					c.WatcherUsersMap[gvk] = []metric.User{user}
				} else {
					log.Infof("update watcher key %s in metric controller", gvk.String())
					// only update WatcherUsersMap, WatchersMap[gvk] already exists
					c.WatcherUsersMap[gvk] = append(c.WatcherUsersMap[gvk], user)
				}
			}

		case metric.UserKind_Timer:
			config, ok := user.GetInterestConfig().(metricUser.TimerUserInterestConfig)
			if !ok {
				return errors.New("wrong InterestConfig type of user")
			}
			t := config.InterestInterval
			if _, ok := c.TimersMap[t]; !ok {
				log.Infof("add new timer key %s to metric controller", t.String())
				// update TimersMap
				c.TimersMap[t] = metric.UserKind_GVKEvent
				// update TimerUsersMap
				c.TimerUsersMap[t] = []metric.User{user}
			} else {
				log.Infof("update timer key %s in metric controller", t.String())
				// only update TimerUsersMap, TimersMap[t] already exists
				c.TimerUsersMap[t] = append(c.TimerUsersMap[t], user)
			}

		}
	}

	// handle trigger event
	c.handleTriggerEvent()

	// start trigger
	if err := c.Trigger.Start(c.WatchersMap); err != nil {
		return err
	}
	if err := c.Trigger.Start(c.TimersMap); err != nil {
		return err
	}
	return nil
}

func (c *Controller) RemoveUsers(users []metric.User) error {
	log := log.WithField("reporter", "Controller").WithField("function", "RemoveUsers")
	// check userkind and stop related watchers
	if err := c.Trigger.Stop(); err != nil {
		return errors.New(fmt.Sprintf("stop trigger failed: %v", err))
	}

	// update related usersmap
	c.Lock()
	defer c.Unlock()

	for _, user := range users {
		switch user.GetUserKind() {
		case metric.UserKind_GVKEvent:
			config, ok := user.GetInterestConfig().(metricUser.PrometheusUserInterestConfig)
			if !ok {
				return errors.New("wrong InterestConfig type of user")
			}
			for _, gvk := range config.InterestKinds {
				if us, ok := c.WatcherUsersMap[gvk]; !ok {
					return errors.New("controller does not contain key " + gvk.String())
				} else {
					handle := false
					for index, u := range us {
						if u.GetName() == user.GetName() {
							if len(us) == 1 {
								// the last user of this gvk, delete key
								delete(c.WatchersMap, gvk)
								delete(c.WatcherUsersMap, gvk)
							} else {
								// only update WatcherUsersMap, WatchersMap[gvk] still has other users
								c.WatcherUsersMap[gvk] = append(us[:index], us[index+1:]...)
							}
							handle = true
							log.Debugf("successfuly remove key %s of user %s from metric controller", gvk.String(), user.GetName())
							break
						}
					}
					if !handle {
						return errors.New(fmt.Sprintf("failure remove key %s of user %s from metric controller", gvk.String(), user.GetName()))
					}
				}
			}

		case metric.UserKind_Timer:
			config, ok := user.GetInterestConfig().(metricUser.TimerUserInterestConfig)
			if !ok {
				return errors.New("wrong InterestConfig type of user")
			}
			t := config.InterestInterval
			if us, ok := c.TimerUsersMap[t]; !ok {
				return errors.New("controller does not contain key " + t.String())
			} else {
				handle := false
				for index, u := range us {
					if u.GetName() == user.GetName() {
						if len(us) == 1 {
							// the last user of this duration, delete key
							delete(c.TimersMap, t)
							delete(c.TimerUsersMap, t)
						} else {
							// only update WatcherUsersMap, WatchersMap[gvk] still has other users
							c.TimerUsersMap[t] = append(us[:index], us[index+1:]...)
						}
						handle = true
						log.Debugf("successfuly remove key %s of user %s from metric controller", t.String(), user.GetName())
						break
					}
				}
				if !handle {
					return errors.New(fmt.Sprintf("failure remove key %s of user %s from metric controller", t.String(), user.GetName()))
				}
			}
		}
	}

	// handle trigger event
	c.handleTriggerEvent()

	// start trigger
	if err := c.Trigger.Start(c.WatchersMap); err != nil {
		return err
	}
	if err := c.Trigger.Start(c.TimersMap); err != nil {
		return err
	}
	return nil

}

func (c *Controller) handleTriggerEvent() {
	log := log.WithField("reporter", "Controller").WithField("function", "handleTriggerEvent")
	// consume trigger event
	go func() {
		for {
			select {
			case <-c.StopChan:
				log.Infof("stopped consuming trigger event")
				return
			case te, ok := <-c.Trigger.GetTriggerEventChan():
				log.Debugf("got an trigger event")
				if !ok {
					log.Warningf("trigger event channel closed, break process loop")
					return
				}
				switch te.TriggerEventKind {
				case metric.TriggerEventKind_Watcher:
					mi, ok := te.Meta.(metric.MetaInfo)
					if !ok {
						log.Errorf("wrong trigger event meta type")
						continue
					}
					for _, u := range c.WatcherUsersMap[mi.GVK] {
						if err := u.TriggerEventHandler(te); err != nil {
							log.Errorf("user %s handled trigger event error: %v", u.GetName(), err)
							continue
						}
						//log.Debugf("controller sent an trigger event to user %s successfully", u.GetName())
					}

				case metric.TriggerEventKind_Timer:
					t, ok := te.Meta.(time.Duration)
					if !ok {
						log.Errorf("wrong trigger event meta type")
						continue
					}
					for _, u := range c.TimerUsersMap[t] {
						if err := u.TriggerEventHandler(te); err != nil {
							log.Errorf("user %s handled trigger event error: %v", u.GetName(), err)
							continue
						}
						//log.Debugf("controller sent an trigger event to user %s", u.GetName())
					}
				}
			}
		}
	}()
}

// todo exit gracefully
// useless? controller lifecycle is same with main function
func (c *Controller) Stop() {
	c.StopChan <- struct{}{}
}
