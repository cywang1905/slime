package source

import (
	"context"
	"errors"
	"fmt"
	prometheus "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/types"
	bootconfig "slime.io/slime/framework/apis/config/v1alpha1"
	"slime.io/slime/framework/model/metric"
	"strings"
	"time"
)

type PrometheusSource struct {
	Name       string
	SourceKind metric.SourceKind
	Api        prometheus.API
}

func (ps *PrometheusSource) GetName() string {
	return ps.Name
}

func (ps *PrometheusSource) GetKind() metric.SourceKind {
	return ps.SourceKind
}

func (ps *PrometheusSource) GetMetrics(handlers interface{}, meta interface{}) (interface{}, error) {
	log := log.WithField("reporter", "PrometheusSource").WithField("function", "GetMetrics")
	psHandlers, ok := handlers.(map[string]*bootconfig.Prometheus_Source_Handler)
	if !ok {
		return nil, errors.New("get metrics error: wrong handler type")
	}
	nn, ok := meta.(types.NamespacedName)
	if !ok {
		return nil, errors.New("get metrics error: wrong meta type")
	}
	material := make(map[string]model.Value)

	for k, v := range psHandlers {
		if v.Query == "" {
			continue
		}
		// TODO: Use more accurate replacements
		query := strings.ReplaceAll(v.Query, "$namespace", nn.Namespace)
		query = strings.ReplaceAll(query, "$source_app", nn.Name)

		qv, w, e := ps.Api.Query(context.Background(), query, time.Now())
		if e != nil {
			log.Debugf("failed get metric from prometheus, handler: %s, error: %+v", v.Query, e)
			return nil, errors.New(fmt.Sprintf("failed get metric from prometheus, %+v", e))
		} else if w != nil {
			log.Debugf("failed get metric from prometheus, handler: %s, warning: %s", v.Query, strings.Join(w, ";"))
			return nil, errors.New(fmt.Sprintf("%s, failed get metric from prometheus", strings.Join(w, ";")))
		}
		log.Debugf("successfully get metric from prometheus, handler: %s", query)
		material[k] = qv
	}
	return material, nil
}
