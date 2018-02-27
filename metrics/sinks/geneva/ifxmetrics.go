package geneva

/*
#cgo CFLAGS: -I../../ifx -I/usr/include
#cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lifx
#include <ifx.hpp>
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type MeasureMetric struct {
	hMetric        unsafe.Pointer
	account        *C.char
	namespace      *C.char
	metricName     *C.char
	countDimension C.ulong
}

func (this *MeasureMetric) NewMeasureMetric(account string, namespace string, metricName string, dimensionKey []string) int {
	this.hMetric = nil
	this.account = C.CString(account)
	this.namespace = C.CString(namespace)
	this.metricName = C.CString(metricName)
	this.countDimension = C.ulong(len(dimensionKey))

	cDimensionKey := make([]*C.char, len(dimensionKey))
	for i := range dimensionKey {
		cDimensionKey[i] = C.CString(dimensionKey[i])
	}
	result := C.CreateIfxMeasureMetricDelegate(&this.hMetric, this.account, this.namespace, this.metricName, this.countDimension, (**C.char)(unsafe.Pointer(&cDimensionKey[0])), 0)
	return int(result)
}

func (this *MeasureMetric) LogValue(value int64, dimensionValue []string) error {
	if C.ulong(len(dimensionValue)) != this.countDimension {
		return fmt.Errorf("The length of dimensionValue is different from dimensionKey")
	}

	cDimensionKey := make([]*C.char, len(dimensionValue))
	for i := range dimensionValue {
		cDimensionKey[i] = C.CString(dimensionValue[i])
		defer C.free(unsafe.Pointer(cDimensionKey[i]))
	}

	result := C.SetIfxMeasureMetricDelegate(this.hMetric, C.longlong(value), this.countDimension, (**C.char)(unsafe.Pointer(&cDimensionKey[0])))
	if int(result) == 0 {
		return nil
	}
	return fmt.Errorf("Send metrics error: %v", result)
}

func IfxStartup() {
	C.IfxStartupDelegate()
}

func IfxCleanup() {
	C.IfxCleanupDelegate()
}

func AttachThreadToIfx() {
	C.AttachThreadToIfxDelegate()
}

func DetachThreadFromIfx() {
	C.DetachThreadFromIfxDelegate()
}
