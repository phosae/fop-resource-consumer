/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/phosae/fop-resource-consumer/common"
)

// ResourceConsumerHandler holds metrics for a resource consumer.
type ResourceConsumerHandler struct {
	metrics     map[string]float64
	metricsLock sync.Mutex
}

// NewResourceConsumerHandler creates and initializes a ResourceConsumerHandler to defaults.
func NewResourceConsumerHandler() *ResourceConsumerHandler {
	return &ResourceConsumerHandler{metrics: map[string]float64{}}
}

func parseCmdToForm(req *http.Request) error {
	cmd := req.Form.Get("cmd")
	parts := strings.Split(cmd, "/")
	if N := len(parts); N <= 2 || N%2 != 0 {
		return fmt.Errorf("invalid cmd: %s", cmd)
	}
	// ?cmd=<any-app>/<any-cmd>/<arg-1>/<argv-1>/.../<arg-n>/<argv-n>
	for i := 2; i < len(parts); i = i + 2 {
		req.Form.Set(parts[i], parts[i+1])
	}
	return nil
}

func (handler *ResourceConsumerHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// handle exposing metrics in Prometheus format (both GET & POST)
	if req.URL.Path == common.MetricsAddress {
		handler.handleMetrics(w)
		return
	}
	if req.Method != "POST" {
		http.Error(w, common.BadRequest, http.StatusBadRequest)
		return
	}

	// parsing POST request data and URL data
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// parsing cmd query in to Form
	if err := parseCmdToForm(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// handle consumeCPU
	if req.Form.Get("type") == common.ConsumeCPUAddress[1:] {
		handler.handleConsumeCPU(w, req.Form)
		return
	}
	// handle consumeMem
	if req.Form.Get("type") == common.ConsumeMemAddress[1:] {
		handler.handleConsumeMem(w, req.Form)
		return
	}
	// handle getCurrentStatus
	if req.Form.Get("type") == common.GetCurrentStatusAddress[1:] {
		handler.handleGetCurrentStatus(w)
		return
	}
	// handle bumpMetric
	if req.Form.Get("type") == common.BumpMetricAddress[1:] {
		handler.handleBumpMetric(w, req.Form)
		return
	}
	http.Error(w, fmt.Sprintf("%s: %s", common.UnknownFunction, req.URL.Path), http.StatusNotFound)
}

func (handler *ResourceConsumerHandler) handleConsumeCPU(w http.ResponseWriter, query url.Values) {
	// getting string data for consumeCPU
	durationSecString := query.Get(common.DurationSecQuery)
	millicoresString := query.Get(common.MillicoresQuery)
	if durationSecString == "" || millicoresString == "" {
		http.Error(w, common.NotGivenFunctionArgument, http.StatusBadRequest)
		return
	}

	// convert data (strings to ints) for consumeCPU
	durationSec, durationSecError := strconv.Atoi(durationSecString)
	millicores, millicoresError := strconv.Atoi(millicoresString)
	if durationSecError != nil || millicoresError != nil {
		http.Error(w, common.IncorrectFunctionArgument, http.StatusBadRequest)
		return
	}

	go ConsumeCPU(millicores, durationSec)
	fmt.Fprintln(w, common.ConsumeCPUAddress[1:])
	fmt.Fprintln(w, millicores, common.MillicoresQuery)
	fmt.Fprintln(w, durationSec, common.DurationSecQuery)
}

func (handler *ResourceConsumerHandler) handleConsumeMem(w http.ResponseWriter, query url.Values) {
	// getting string data for consumeMem
	durationSecString := query.Get(common.DurationSecQuery)
	megabytesString := query.Get(common.MegabytesQuery)
	if durationSecString == "" || megabytesString == "" {
		http.Error(w, common.NotGivenFunctionArgument, http.StatusBadRequest)
		return
	}

	// convert data (strings to ints) for consumeMem
	durationSec, durationSecError := strconv.Atoi(durationSecString)
	megabytes, megabytesError := strconv.Atoi(megabytesString)
	if durationSecError != nil || megabytesError != nil {
		http.Error(w, common.IncorrectFunctionArgument, http.StatusBadRequest)
		return
	}

	go ConsumeMem(megabytes, durationSec)
	fmt.Fprintln(w, common.ConsumeMemAddress[1:])
	fmt.Fprintln(w, megabytes, common.MegabytesQuery)
	fmt.Fprintln(w, durationSec, common.DurationSecQuery)
}

func (handler *ResourceConsumerHandler) handleGetCurrentStatus(w http.ResponseWriter) {
	GetCurrentStatus()
	fmt.Fprintln(w, "Warning: not implemented!")
	fmt.Fprint(w, common.GetCurrentStatusAddress[1:])
}

func (handler *ResourceConsumerHandler) handleMetrics(w http.ResponseWriter) {
	handler.metricsLock.Lock()
	defer handler.metricsLock.Unlock()
	for k, v := range handler.metrics {
		fmt.Fprintf(w, "# HELP %s info message.\n", k)
		fmt.Fprintf(w, "# TYPE %s gauge\n", k)
		fmt.Fprintf(w, "%s %f\n", k, v)
	}
}

func (handler *ResourceConsumerHandler) bumpMetric(metric string, delta float64, duration time.Duration) {
	handler.metricsLock.Lock()
	if _, ok := handler.metrics[metric]; ok {
		handler.metrics[metric] += delta
	} else {
		handler.metrics[metric] = delta
	}
	handler.metricsLock.Unlock()

	time.Sleep(duration)

	handler.metricsLock.Lock()
	handler.metrics[metric] -= delta
	handler.metricsLock.Unlock()
}

func (handler *ResourceConsumerHandler) handleBumpMetric(w http.ResponseWriter, query url.Values) {
	// getting string data for handleBumpMetric
	metric := query.Get(common.MetricNameQuery)
	deltaString := query.Get(common.DeltaQuery)
	durationSecString := query.Get(common.DurationSecQuery)
	if durationSecString == "" || metric == "" || deltaString == "" {
		http.Error(w, common.NotGivenFunctionArgument, http.StatusBadRequest)
		return
	}

	// convert data (strings to ints/floats) for handleBumpMetric
	durationSec, durationSecError := strconv.Atoi(durationSecString)
	delta, deltaError := strconv.ParseFloat(deltaString, 64)
	if durationSecError != nil || deltaError != nil {
		http.Error(w, common.IncorrectFunctionArgument, http.StatusBadRequest)
		return
	}

	go handler.bumpMetric(metric, delta, time.Duration(durationSec)*time.Second)
	fmt.Fprintln(w, common.BumpMetricAddress[1:])
	fmt.Fprintln(w, metric, common.MetricNameQuery)
	fmt.Fprintln(w, delta, common.DeltaQuery)
	fmt.Fprintln(w, durationSec, common.DurationSecQuery)
}
