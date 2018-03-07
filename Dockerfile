##############################################
# Build heapster and event
##############################################
FROM ubuntu:xenial

RUN apt-get update && apt-get install -y wget apt-utils apt-transport-https g++-4.8-multilib make git

RUN wget https://dl.google.com/go/go1.10.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.10.linux-amd64.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

RUN rm /usr/bin/gcc && \
    ln -s /usr/bin/gcc-4.8 /usr/bin/gcc && \
    ln -s /usr/bin/g++-4.8 /usr/bin/g++

RUN repoPkg=azure-repoclient-https-noauth-public-xenial_1.0.2-47_amd64.deb && \
	wget --no-check-certificate https://apt-mo.trafficmanager.net/repos/azurecore/pool/main/a/azure-repoclient-https-noauth-public-xenial/$repoPkg && \
	dpkg -i $repoPkg && \
	apt-get update && \
    apt-get install -y metricsext && \
    apt-get install -y libazurepal && \
    apt-get install -y libazurepal-dev

ENV SOURCE ${GOPATH}/src/k8s.io/heapster
RUN mkdir -p ${SOURCE}
Add . ${SOURCE}

WORKDIR ${SOURCE}/metrics/ifx
RUN bash ./build.sh
WORKDIR ${SOURCE}

RUN make

##############################################
# Build final image
##############################################
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

COPY --from=0 /go/src/k8s.io/heapster/metrics/ifx/libifx.so /usr/lib/x86_64-linux-gnu/libifx.so

COPY --from=0 /go/src/k8s.io/heapster/heapster /

COPY run.sh /run.sh

RUN chmod +x /run.sh

ENV LD_LIBRARY_PATH="/usr/lib/x86_64-linux-gnu"

# COPY ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/run.sh"]
