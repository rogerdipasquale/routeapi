
### Install local kind cluster

```shell 
# For AMD64 / x86_64
[ $(uname -m) = x86_64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.32.0/kind-linux-amd64
# For ARM64
[ $(uname -m) = aarch64 ] && curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.32.0/kind-linux-arm64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

kind create cluster
```
### Install traefik and API gateway routes

kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml

helm repo add traefik https://traefik.github.io/charts
helm update
helm install traefik traefik/traefik -f values.yaml -n traefik --create-namespace --version 37.3.0

### Install assets
Assets are taken from:
https://github.com/rogerdipasquale/k8s-gateway/tree/main/traefik/application
and 
https://github.com/rogerdipasquale/k8s-gateway/tree/main/traefik/gateway

Those files are downloaded here as applications and routes.yaml
RBAC can be taken from (../../manifests/rbac.yaml)[../../manifests/rbac.yaml]

Run: 

```
kubectl apply -f applications.yaml
kubectl apply -f routes.yaml
kubectl apply -f ../../manifests/rbac.yaml
```

### Run the application (in background mode)

go run cmd/server/main.go &

### Test 

curl "http://127.0.0.1:8080/api/routes?namespace=traefik"

curl "http://127.0.0.1:8080/api/route/traefik/weight-route"|jq
