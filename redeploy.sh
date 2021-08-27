#!/bin/bash

minikube delete;
minikube start --cpus 4 --memory 8192;
istioctl install --set profile=demo -y;
kubectl label ns default istio-injection=enabled;
kubectl apply -f ~/istio-1.11.0/samples/addons;
kubectl apply -f ~/istio-1.11.0/samples/addons;
kubectl apply -f -<<EOF
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: ingress-gateway
spec:
  selector:
    istio: ingressgateway # use istio default controller
  servers:
  - port:
      number: 80
      name: http
      protocol: HTTP
    hosts:
    - "*"
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: main
spec:
  hosts:
  - "*"
  gateways:
  - ingress-gateway
  http:
  - match:
    - uri:
        prefix: /
    route:
    - destination:
        host: findme-resource
        port:
          number: 80
EOF
make docker-build docker-push
kubectl apply -f -<<EOF
apiVersion: application.caseywylie.io/v1alpha1
kind: Findme
metadata:
  name: findme-resource
spec: 
  size: 1
EOF
sleep 10s;
kubectl logs deploy/operator-controller-manager -n operator-system -c manager | grep controller-runtime.manager.controller.findme