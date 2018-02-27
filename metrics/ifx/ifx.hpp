#pragma once

#ifdef __cplusplus
extern "C" {
#endif

    int AttachThreadToIfxDelegate();

    int DetachThreadFromIfxDelegate();

    unsigned int IfxStartupDelegate();

    void IfxCleanupDelegate();

    long CreateIfxMeasureMetricDelegate(
        void** hMetric,
        const char* monitoringAccount,
        const char* metricNamespace,
        const char* metricName,
        unsigned long countDimension,
        const char** listDimensionNames,
        int addDefaultDimension
    );

    long SetIfxMeasureMetricDelegate(
        void* hMetric,
        long long rawData,
        unsigned long countDimension,
        const char** listDimensionValues
    );

// The API to set the measure metric; use the explicit UTC timestamp.
    long SetIfxMeasureMetricWithTimestampDelegate(
        void* hMetric,
        unsigned long long timestampUtc,
        long long rawData,
        unsigned long countDimension,
        const char** listDimensionValues
    );

#ifdef __cplusplus
}
#endif