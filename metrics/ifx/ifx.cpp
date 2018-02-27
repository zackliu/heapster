#include "ifx.hpp"
#include <IfxMetrics.h>

int AttachThreadToIfxDelegate() {
    return AttachThreadToIfx();
}

int DetachThreadFromIfxDelegate() {
    return DetachThreadFromIfx();
}

unsigned int IfxStartupDelegate() {
    return IfxStartup();
}

void IfxCleanupDelegate() {
    return IfxCleanup();
}

long CreateIfxMeasureMetricDelegate(
    void** hMetric,
    const char* monitoringAccount,
    const char* metricNamespace,
    const char* metricName,
    unsigned long countDimension,
    const char** listDimensionNames,
    int addDefaultDimension
) {
    return CreateIfxMeasureMetric(hMetric, monitoringAccount, metricNamespace, metricName, countDimension, listDimensionNames, addDefaultDimension);
}

long SetIfxMeasureMetricDelegate(
    void* hMetric,
    long long rawData,
    unsigned long countDimension,
    const char** listDimensionValues
) {
    return SetIfxMeasureMetric(hMetric, rawData, countDimension, listDimensionValues);
}

// The API to set the measure metric; use the explicit UTC timestamp.
long SetIfxMeasureMetricWithTimestampDelegate(
    void* hMetric,
    unsigned long long timestampUtc,
    long long rawData,
    unsigned long countDimension,
    const char** listDimensionValues
) {
    return SetIfxMeasureMetricWithTimestamp(hMetric, timestampUtc, rawData, countDimension, listDimensionValues);
}