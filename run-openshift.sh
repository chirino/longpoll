#!/bin/bash
set -e 
source .env
oc login $OC_LOGIN_ARGS
oc project $OC_PROJECT

# docker buildx bake --push

cat > ./deploy/overlays/ocp/kustomization.yaml << EOF
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - ../../base
patches:
  - target:
      kind: Ingress
      name: longpoll
    patch: |-
      - op: replace
        path: /spec/rules/0/host
        value: longpoll${OC_DOMAIN_SUFFIX}
  - target:
      kind: Ingress
      name: envoy
    patch: |-
      - op: replace
        path: /spec/rules/0/host
        value: envoy${OC_DOMAIN_SUFFIX}
EOF
kubectl apply -k ./deploy/overlays/ocp

echo
echo "Test the long polling on the server using:"
echo
echo "    go run ./client 'http://envoy${OC_DOMAIN_SUFFIX}/api?interval=65'"
echo
