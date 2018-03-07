#!/bin/bash

set -x

#Necessary files
CERT_FILE="/etc/mdm/cert.pem"
KEY_FILE="/etc/mdm/key.pem"
MDM_FILE="/etc/default/mdm"

# Temp path in k8s
MDM_FILE_TEMP="/tmp/geneva/mdm"

# Check all the necessary files
function check_file {
    if [[ -f ${CERT_FILE} && -f ${KEY_FILE} && -f ${MDM_FILE} ]]; then
        echo "All necessary files exist"
        return 0
    else
        echo "Some necessary files not exist"
        return 1
    fi
}

if cp "${MDM_FILE_TEMP}" "${MDM_FILE}"; then
    if ! check_file; then
        echo "ERROR: Necessary files missing"
        exit 1
    fi
else
    echo "ERROR: Copy failed"
    exit 1
fi

if [ -z "$ACCOUNT" ]; then
    ACCOUNT="SignalRShoeboxTest"
fi

service rsyslog restart
service mdm restart

ldd /usr/lib/x86_64-linux-gnu/libifx.so
ldd /usr/lib/x86_64-linux-gnu/libIfxMetrics.so
/heapster --source=kubernetes --sink="geneva:?account=${ACCOUNT}&namespace=systemLoad"