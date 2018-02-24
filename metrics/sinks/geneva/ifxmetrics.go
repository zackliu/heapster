package geneva

/*
#cgo CFLAGS: -I/usr/include/azurepal/IfxMetrics
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu
#include <stdlib.h>
#include <IfxMetrics.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type MeasureMetric struct {
	hMetric        unsafe.Pointer
	account        C.CString
	namespace      C.CString
	metricName     C.CString
	countDimension C.uint
}

func (this *MeasureMetric) NewMeasureMetric(account string, namespace string, metricName string, dimensionKey []string) int {
	this.hMetric = unsafe.Pointer()
	this.account = C.CString(account)
	this.namespace = C.CString(namespace)
	this.metricName = C.CString(metricName)
	this.countDimension = C.uint(len(dimensionKey))

	cDimensionKey := make([]*C.char, len(dimensionKey))
	for i := range dimensionKey {
		cDimensionKey[i] = C.CString(goString[i])
	}
	result := C.CreateIfxMeasureMetric(&this.hMetric, this.account, this.namespace, this.metricName, this.countDimension, (**C.char)(unsafe.Pointer(&cDimensionKey[0])))
	return int(result)
}

func (this *MeasureMetric) LogValue(value int64, dimensionValue []string) error {
	if len(dimensionValue) != this.countDimension {
		return fmt.Errorf("The length of dimensionValue is different from dimensionKey", err)
	}

	cDimensionKey := make([]*C.char, len(dimensionValue))
	for i := range dimensionValue {
		cDimensionKey[i] = C.CString(dimensionValue[i])
		defer C.free(unsafe.Pointer(cDimensionKey[i]))
	}

	result := C.SetIfxMeasureMetric(this.hMetric, C.int64(value), this.countDimension, (**C.char)(unsafe.Pointer(&cDimensionKey[0])))
	if int(result) == 0 {
		return nil
	}
	return fmt.Errorf("Send metrics error: %v", result)
}

func IfxStartup() {
	C.IfxStartup()
}

func IfxCleanup() {
	C.IfxCleanup()
}

func AttachThreadToIfx() {
	C.AttachThreadToIfx()
}

func DetachThreadFromIfx() {
	C.DetachThreadFromIfx()
}
