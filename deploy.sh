#!/bin/bash

set -x

MDM_ACCOUNT_PROD="MicrosoftSignalRServiceShoebox"

deleteConfig(){
    kubectl delete configmap mdm-mdsd-conf -n kube-system
    kubectl delete secret secret-conf -n kube-system
    kubectl delete secret genevaregistry -n kube-system
}

deployCommon(){
    #Create Secret for cert.pem and key.pem
    kubectl create secret generic secret-conf --from-file=Secret -n kube-system

    #Create ConfigMap of mdm
    kubectl create configmap mdm-mdsd-conf --from-file=MdmMdsdConf -n kube-system

    #Create docker-registry secret
    kubectl create secret docker-registry genevaregistry --docker-server=geneva.azurecr.io --docker-username=geneva --docker-password= --docker-email=chenyl@microsoft.com -n kube-system

    #Deploy heapster
    kubectl replace -f deploy/kube-config/geneva/heapster.yaml
}


deployINT(){
    deleteConfig

    template=`cat deploy/kube-config/geneva/heapster.template.yaml`
    printf "ACCOUNT=SignalRShoeboxTest\ncat << EOF\n$template\nEOF" | bash > deploy/kube-config/geneva/heapster.yaml
    
    deployCommon

    rm deploy/kube-config/geneva/heapster.yaml
}

deployPROD(){
    deleteConfig

    template=`cat deploy/kube-config/geneva/heapster.template.yaml`
    printf "ACCOUNT=${MDM_ACCOUNT_PROD}${2}\ncat << EOF\n$template\nEOF" | bash > deploy/kube-config/geneva/heapster.yaml

    #Change mdm account according to location
    sed -i "s/MDM_ACCOUNT=.*/MDM_ACCOUNT=\"${MDM_ACCOUNT_PROD}${2}\"/g" MdmMdsdConf_PROD/mdm

    deployCommon

    sed -i "s/MDM_ACCOUNT=.*/MDM_ACCOUNT=\"SignalRShoeboxTest\"/g" MdmMdsdConf_PROD/mdm

    rm deploy/kube-config/geneva/heapster.yaml
}



if [[ -z "$1" ||  -n "$1" && "$1" == "INT" ]]; then
    deployINT
elif [[ -n "$1" && "$1" == "PROD" && -n "$2" ]]; then
    deployPROD $1 $2
else
    echo "Need endpoint and location"
    exit 1
fi
