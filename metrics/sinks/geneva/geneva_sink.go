package geneva

import (
	"net/url"
	"runtime"
	"strings"

	"github.com/golang/glog"
	"k8s.io/heapster/metrics/core"
)

type GenevaConfig struct {
	account   string
	namespace string
}

func getConfig(uri *url.URL) GenevaConfig {
	cfg := GenevaConfig{
		account:   "SignalRShoeboxTest",
		namespace: "k8stest",
	}

	opts := uri.Query()
	if len(opts["account"]) >= 1 {
		cfg.account = opts["account"][0]
	}

	if len(opts["namespace"]) >= 1 {
		cfg.namespace = opts["namespace"][0]
	}

	return cfg
}

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
	glog.Info(metricSetName)
	metricSlice := strings.Split(metricSetName, "/")
	var metricSetKey []string
	var metricSetVal []string
	for _, key := range metricSlice {
		stringSlice := strings.SplitN(key, ":", 2)
		if len(stringSlice) == 2 {
			metricSetKey = append(metricSetKey, stringSlice[0])
			metricSetVal = append(metricSetVal, stringSlice[1])
		}
	}
	length := len(metricSetKey)
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

func getCustomLabel(labels string, labelName string) string {
	labelsSlice := strings.Split(labels, ",")
	for _, label := range labelsSlice {
		labelSlice := strings.SplitN(label, ":", 2)
		if labelSlice[0] == labelName {
			return labelSlice[1]
		}
	}
	return ""
}

func (this *GenevaSink) dealNodeMetric(nodeName string, metricSet *core.MetricSet) {
	if value, ok := metricSet.MetricValues["cpu/node_utilization"]; ok {
		this.metricsDefinition["NodeCpuPercentage"].LogValue("NodeCpuPercentage", getValue(value, 100), []string{nodeName})
	}
	if value, ok := metricSet.MetricValues["cpu/usage_rate"]; ok {
		this.metricsDefinition["NodeCpuUsage"].LogValue("NodeCpuUsage", getValue(value, 1), []string{nodeName})
	}
	if value, ok := metricSet.MetricValues["memory/working_set"]; ok {
		this.metricsDefinition["NodeMemoryUsage"].LogValue("NodeMemoryUsage", getValue(value, 1), []string{nodeName})
	}
	if value, ok := metricSet.MetricValues["memory/node_utilization"]; ok {
		this.metricsDefinition["NodeMemoryPercentage"].LogValue("NodeMemoryPercentage", getValue(value, 100), []string{nodeName})
	}
	if value, ok := metricSet.MetricValues["network/rx_rate"]; ok {
		this.metricsDefinition["NodeNetworkIn"].LogValue("NodeNetworkIn", getValue(value, 1), []string{nodeName})
	}
	if value, ok := metricSet.MetricValues["network/tx_rate"]; ok {
		this.metricsDefinition["NodeNetworkOut"].LogValue("NodeNetworkOut", getValue(value, 1), []string{nodeName})
	}
}

func (this *GenevaSink) dealPodMetric(podName string, metricSet *core.MetricSet) {
	resourceKubeId := getCustomLabel(metricSet.Labels["labels"], "resourceKubeId")
	if resourceKubeId == "" {
		resourceKubeId = "unknown"
	}

	if value, ok := metricSet.MetricValues["network/rx_rate"]; ok {
		this.metricsDefinition["PodNetworkIn"].LogValue("PodNetworkIn", getValue(value, 1), []string{resourceKubeId, podName, "total"})
	}
	if value, ok := metricSet.MetricValues["network/tx_rate"]; ok {
		this.metricsDefinition["PodNetworkOut"].LogValue("PodNetworkOut", getValue(value, 1), []string{resourceKubeId, podName, "total"})
	}
}

func (this *GenevaSink) dealContainerMetric(podName string, containerName string, metricSet *core.MetricSet) {
	resourceKubeId := getCustomLabel(metricSet.Labels["labels"], "resourceKubeId")
	if resourceKubeId == "" {
		resourceKubeId = "unknown"
	}

	if value, ok := metricSet.MetricValues["cpu/usage_rate"]; ok {
		this.metricsDefinition["PodCpuUsage"].LogValue("PodCpuUsage", getValue(value, 1), []string{resourceKubeId, podName, containerName})
	}
	if value, ok := metricSet.MetricValues["memory/working_set"]; ok {
		this.metricsDefinition["PodMemory"].LogValue("PodMemory", getValue(value, 1), []string{resourceKubeId, podName, containerName})
	}

	for _, labeledMetric := range metricSet.LabeledMetrics {
		if labeledMetric.Name == "disk/io_read_bytes_rate" {
			this.metricsDefinition["PodDiskRead"].LogValue("PodDiskRead", getValue(labeledMetric.MetricValue, 1), []string{resourceKubeId, podName, containerName})
		}
		if labeledMetric.Name == "disk/io_write_bytes_rate" {
			this.metricsDefinition["PodDiskWrite"].LogValue("PodDiskWrite", getValue(labeledMetric.MetricValue, 1), []string{resourceKubeId, podName, containerName})
		}
	}
}

func (this *GenevaSink) SendToGeneva(batch *core.DataBatch) {
	for metricSetName, metricSet := range batch.MetricSets {
		resourceType, node, pod, container := switchMetricSet(metricSetName)
		if resourceType == "node" {
			this.dealNodeMetric(node, metricSet)
		} else if resourceType == "pod" {
			this.dealPodMetric(pod, metricSet)
		} else if resourceType == "container" {
			this.dealContainerMetric(pod, container, metricSet)
		}
	}
}

func (this *GenevaSink) InitMetrics(config GenevaConfig) {
	account := config.account
	namespace := config.namespace

	this.metricsDefinition["NodeCpuPercentage"] = NewMeasureMetric(account, namespace, "NodeCpuPercentage", []string{"NodeName"})
	this.metricsDefinition["NodeCpuUsage"] = NewMeasureMetric(account, namespace, "NodeCpuUsage", []string{"NodeName"})
	this.metricsDefinition["NodeMemoryUsage"] = NewMeasureMetric(account, namespace, "NodeMemoryUsage", []string{"NodeName"})
	this.metricsDefinition["NodeMemoryPercentage"] = NewMeasureMetric(account, namespace, "NodeMemoryPercentage", []string{"NodeName"})
	this.metricsDefinition["NodeNetworkIn"] = NewMeasureMetric(account, namespace, "NodeNetworkIn", []string{"NodeName"})
	this.metricsDefinition["NodeNetworkOut"] = NewMeasureMetric(account, namespace, "NodeNetworkOut", []string{"NodeName"})
	this.metricsDefinition["PodCpuUsage"] = NewMeasureMetric(account, namespace, "PodCpuUsage", []string{"resourceKubeId", "PodName", "ContainerName"})
	this.metricsDefinition["PodMemory"] = NewMeasureMetric(account, namespace, "PodMemory", []string{"resourceKubeId", "PodName", "ContainerName"})
	this.metricsDefinition["PodNetworkIn"] = NewMeasureMetric(account, namespace, "PodNetworkIn", []string{"resourceKubeId", "PodName", "ContainerName"})
	this.metricsDefinition["PodNetworkOut"] = NewMeasureMetric(account, namespace, "PodNetworkOut", []string{"resourceKubeId", "PodName", "ContainerName"})
	this.metricsDefinition["PodDiskWrite"] = NewMeasureMetric(account, namespace, "PodDiskWrite", []string{"resourceKubeId", "PodName", "ContainerName"})
	this.metricsDefinition["PodDiskRead"] = NewMeasureMetric(account, namespace, "PodDiskRead", []string{"resourceKubeId", "PodName", "ContainerName"})
}

func (this *GenevaSink) ExportData(batch *core.DataBatch) {
	runtime.LockOSThread()
	AttachThreadToIfx()
	this.SendToGeneva(batch)
	DetachThreadFromIfx()
}

func NewGenevaSink(uri *url.URL) *GenevaSink {
	genevaSink := &GenevaSink{make(map[string]*MeasureMetric)}
	IfxStartup()
	config := getConfig(uri)
	genevaSink.InitMetrics(config)
	return genevaSink
}
