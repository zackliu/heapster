package geneva

import (
	"k8s.io/heapster/metrics/core"
)

type GenevaSink struct {
}

// type Metric struct {
// 	metricInternal *MeasureMetric
// 	account        string
// 	namespace      string
// 	metricName     string
// 	dimensionKey   []string
// }

// func (this *Metric) newMetric(account string, namespace string, metricName string, dimensionKey []string) {
// 	this.account = account
// 	this.namespace = namespace
// 	this.metricName = metricName
// 	this.dimensionKey = dimensionKey
// 	this.metricInternal = &MeasureMetric{}
// 	this.metricInternal.NewMeasureMetric(this.account, this.namespace, this.metricName, this.dimensionKey)
// }

// func (this *Metric) sendMetric(value int64, dimensionValue []string) {
// 	if this.metricInternal != nil {
// 		this.metricInternal.LogValue(value, dimensionValue)
// 	}
// }

func (this *GenevaSink) Name() string {
	return "Geneva Sink"
}

func (this *GenevaSink) Stop() {

}

func SendToGeneva(batch *core.DataBatch) {
	metric := &MeasureMetric{}
	metric.NewMeasureMetric("SignalRShoeboxTest", "k8stest", "test", []string{"ResourceId", "InstanceId"})
	metric.LogValue(22, []string{"/asdfasd/asdfasdf", "dfsfh46814sdaf"})

}

func (this *GenevaSink) ExportData(batch *core.DataBatch) {
	IfxStartup()
	SendToGeneva(batch)
	IfxCleanup()
}

func NewGenevaSink() *GenevaSink {
	return &GenevaSink{}
}
