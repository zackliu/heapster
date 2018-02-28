#!/bin/bash

service rsyslog restart
service mdm restart

ldd /usr/lib/x86_64-linux-gnu/libifx.so
ldd /usr/lib/x86_64-linux-gnu/libIfxMetrics.so
/heapster $@