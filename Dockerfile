FROM ubuntu:xenial

RUN apt-get update && apt-get install -y wget apt-utils apt-transport-https

RUN repoPkg=azure-repoclient-https-noauth-public-xenial_1.0.2-47_amd64.deb && \
	wget --no-check-certificate https://apt-mo.trafficmanager.net/repos/azurecore/pool/main/a/azure-repoclient-https-noauth-public-xenial/$repoPkg && \
	dpkg -i $repoPkg && \
	apt-get update && \
    apt-get install -y rsyslog && \
    apt-get install -y metricsext && \
    apt-get install -y libazurepal

RUN apt-get install -y elfutils

COPY mdm /etc/default/mdm

COPY cert.pem /etc/mdm/cert.pem

COPY key.pem /etc/mdm/key.pem

COPY ./metrics/ifx/libifx.so /usr/lib/x86_64-linux-gnu/libifx.so

COPY heapster /

COPY run.sh /run.sh

ENV LD_LIBRARY_PATH="/usr/lib/x86_64-linux-gnu"

# COPY ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/run.sh"]
