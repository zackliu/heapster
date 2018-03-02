package geneva

import (
	"strings"

	"k8s.io/heapster/metrics/core"
)

type GenevaSink struct {
	metricsDefinition map[string]*MeasureMetric
}

func (this *GenevaSink) Name() string {
	return "Geneva Sink"
}

func (this *GenevaSink) Stop() {
	IfxCleanup()
}

// Return type, node, pod, container.
func switchMetricSet(metricSetName string) (string, string, string, string) {
	metricSlice := strings.Split(metricSetName, "/")
	var metricSetKey []string
	var metricSetVal []string
	for _, key := range metricSlice {
		metricSetKey = append(metricSetKey, strings.SplitN(key, ":", 2)[0])
		metricSetVal = append(metricSetVal, strings.SplitN(key, ":", 2)[1])
	}
	length := len(metricSlice)
	// We only consider three situation:
	//   Only node: node:aks-agentpool-39637307-0
	//   Pod: namespace:default/pod:signalr-5bc9678c77-ndzbt
	//   Container: namespace:default/pod:signalr-5bc9678c77-ndzbt/container:signalr
	if length == 1 && metricSetKey[0] == "node" {
		return "node", metricSetVal[0], "", ""
	} else if length == 2 && metricSetKey[0] == "namespace" && metricSetKey[1] == "pod" {
		return "pod", "", metricSetVal[1], ""
	} else if length == 3 && metricSetKey[0] == "namespace" && metricSetKey[1] == "pod" && metricSetKey[2] == "container" {
		return "container", "", metricSetVal[1], metricSetVal[2]
	} else {
		return "others", "", "", ""
	}
}

func getValue(value core.MetricValue, magnification int64) int64 {
	if core.ValueInt64 == value.ValueType {
		return value.IntValue * magnification
	} else if core.ValueFloat == value.ValueType {
		return int64(value.FloatValue * float32(magnification))
	} else {
		return 0
	}
}

func (this *GenevaSink) dealNodeMetric(nodeName string, metricSet *core.MetricSet) {
	if value, ok := metricSet.MetricValues["cpu/node_utilization"]; ok {
		this.metricsDefinition["NodeCpuPercentage"].LogValue(getValue(value, 100), []string{nodeName})
	} else if value, ok := metricSet.MetricValues["cpu/usage_rate"]; ok {
		this.metricsDefinition["NodeCpuUsage"].LogValue(getValue(value, 1), []string{nodeName})
	} else if value, ok := metricSet.MetricValues["memory/working_set"]; ok {
		this.metricsDefinition["NodeMemoryUsage"].LogValue(getValue(value, 1), []string{nodeName})
	} else if value, ok := metricSet.MetricValues["memory/node_utilization"]; ok {
		this.metricsDefinition["NodeMemoryPercentage"].LogValue(getValue(value, 100), []string{nodeName})
	} else if value, ok := metricSet.MetricValues["network/rx_rate"]; ok {
		this.metricsDefinition["NodeNetworkIn"].LogValue(getValue(value, 1), []string{nodeName})
	} else if value, ok := metricSet.MetricValues["network/tx_rate"]; ok {
		this.metricsDefinition["NodeNetworkOut"].LogValue(getValue(value, 1), []string{nodeName})
	}
}

func (this *GenevaSink) SendToGeneva(batch *core.DataBatch) {
	for metricSetName, metricSet := range batch.MetricSets {
		resourceType, node, pod, container := switchMetricSet(metricSetName)
		if resourceType == "node" {
			this.dealNodeMetric(node, metricSet)
		} else if resourceType == "pod" {

		} else if resourceType == "container" {

		}

	}

	metric := NewMeasureMetric("SignalRShoeboxTest", "k8stest", "test", []string{"ResourceId", "InstanceId"})
	metric.LogValue(22, []string{"/asdfasd/asdfasdf", "dfsfh46814sdaf"})

}

func (this *GenevaSink) InitMetrics() {
	this.metricsDefinition["NodeCpuPercentage"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "NodeCpuPercentage", []string{"NodeName"})
	this.metricsDefinition["NodeCpuUsage"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "NodeCpuUsage", []string{"NodeName"})
	this.metricsDefinition["NodeMemoryUsage"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "NodeMemoryUsage", []string{"NodeName"})
	this.metricsDefinition["NodeMemoryPercentage"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "NodeMemoryPercentage", []string{"NodeName"})
	this.metricsDefinition["NodeNetworkIn"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "NodeNetworkIn", []string{"NodeName"})
	this.metricsDefinition["NodeNetworkOut"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "NodeNetworkOut", []string{"NodeName"})
	this.metricsDefinition["PodCpuUsage"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "PodCpuUsage", []string{"ResourceId", "PodName", "ContainerName"})
	this.metricsDefinition["PodMemory"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "PodMemory", []string{"ResourceId", "PodName", "ContainerName"})
	this.metricsDefinition["PodNetworkIn"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "PodNetworkIn", []string{"ResourceId", "PodName", "ContainerName"})
	this.metricsDefinition["PodNetworkOut"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "PodNetworkOut", []string{"ResourceId", "PodName", "ContainerName"})
	this.metricsDefinition["PodDiskWrite"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "PodDiskWrite", []string{"ResourceId", "PodName", "ContainerName"})
	this.metricsDefinition["PodDiskRead"] = NewMeasureMetric("SignalRShoeboxTest", "k8stest", "PodDiskRead", []string{"ResourceId", "PodName", "ContainerName"})
}

func (this *GenevaSink) ExportData(batch *core.DataBatch) {
	this.SendToGeneva(batch)
}

func NewGenevaSink() *GenevaSink {
	genevaSink := &GenevaSink{make(map[string]*MeasureMetric)}
	IfxStartup()
	genevaSink.InitMetrics()
	return genevaSink
}
