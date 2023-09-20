#!/bin/bash

# kind delete clusters longpoll
kind get clusters | grep longpoll || (
  kind create cluster --config ./deploy/kind.yaml
  kubectl cluster-info --context kind-longpoll
  kubectl apply -f ./deploy/kind-ingress.yaml
)
kubectl rollout restart deployment ingress-nginx-controller -n ingress-nginx
kubectl rollout status deployment ingress-nginx-controller -n ingress-nginx --timeout=5m

docker build . -t docker.io/chirino/longpoll:latest
kind load --name longpoll docker-image docker.io/chirino/longpoll:latest
kubectl rollout restart deployment longpoll -n longpoll
kubectl apply -k ./deploy/overlays/kind

echo
echo "Test the long polling on the server using:"
echo
echo "    go run ./client 'http://nvoy.127.0.0.1.nip.io:8080/api?interval=65'"
echo