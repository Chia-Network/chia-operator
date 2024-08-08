#!/usr/bin/env bash

mkdir -p certs
cd certs
chia-tools certs generate -o ./
kubectl create secret generic ${SECRET_NAME} --namespace ${NAMESPACE} --from-file=ca/chia_ca.crt --from-file=ca/chia_ca.key --from-file=ca/private_ca.crt --from-file=ca/private_ca.key
