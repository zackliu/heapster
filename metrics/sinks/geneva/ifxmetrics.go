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
	countDimension C.ulong
}

func NewMeasureMetric(account string, namespace string, metricName string, dimensionKey []string) *MeasureMetric {
	this := &MeasureMetric{}
	this.hMetric = nil
	cAccount := C.CString(account)
	cNamespace := C.CString(namespace)
	cMetricName := C.CString(metricName)
	this.countDimension = C.ulong(len(dimensionKey))

	defer C.free(unsafe.Pointer(cAccount))
	defer C.free(unsafe.Pointer(cNamespace))
	defer C.free(unsafe.Pointer(cMetricName))

	cDimensionKey := make([]*C.char, len(dimensionKey))
	for i := range dimensionKey {
		cDimensionKey[i] = C.CString(dimensionKey[i])
		defer C.free(unsafe.Pointer(cDimensionKey[i]))
	}
	result := C.CreateIfxMeasureMetricDelegate(&this.hMetric, cAccount, cNamespace, cMetricName, this.countDimension, (**C.char)(unsafe.Pointer(&cDimensionKey[0])), 0)
	if int(result) == 0 {
		return this
	}
	return nil
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
