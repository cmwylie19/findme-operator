# Minikube Environment
```
minikube start --cpus=4 --memory=8192
```

# Install Istio
```
# Create istio-system namespace
kubectl create ns istio-system

# Label default namespace to allow sidecars
kubectl label ns default istio-injection=enabled

# Verify Label
kubectl get ns default -ojsonpath='{.metadata.labels.istio-injection}'  

# Install Istio base chart which contains cluster-wide resources use by the istio control plane
helm install istio-base helm-charts/istio -n istio-system

# Verify resource were applied to the cluster
kubectl get crd | grep istio.io

# Install Istio discovery chart which deploys `istiod` service
helm install istiod helm-charts/istiod -n istio-system

# Verify installation
kubectl get svc istiod -n istio-system

# Install Istio ingress gateway
helm install istio-ingress helm-charts/istio-ingress -n istio-system 

# Verify installation
kubectl get svc -n istio-system istio-ingressgateway
```
# OLM
```
operator-sdk olm install/uninstall

make bundle bundle-build bundle-push 

operator-sdk run bundle docker.io/cmwylie19/findme-operator-bundle:v0.0.1 

operator-sdk run bundle-upgrade docker.io/cmwylie19/findme-operator-bundle:v0.0.1 
```

k logs deploy/operator-controller-manager -n operator-system -c manager 

k logs deploy/operator-controller-manager -n operator-system -c manager | grep controller-runtime.manager.controller.findme

make docker-build docker-push


https://github.com/spotahome/redis-operator/blob/master/operator/redisfailover/service/generator_test.go

### RBAC
/config/rbac/role.yaml
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;

### Istio
```
istioctl install --set profile=demo -y
kubectl label namespace default istio-injection=enabled

```
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
```

```
istioctl install --set profile=demo -y; k apply -f samples/addons; k apply -f samples/addonsl k label ns istio-system istio-injection=enabled;k apply -f 
heredoc> apiVersion: networking.istio.io/v1alpha3
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
```