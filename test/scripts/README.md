
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



helm repo add traefik https://traefik.github.io/charts
helm repo update
helm install traefik traefik/traefik -f traefik-values.yaml -n traefik --create-namespace --version 37.3.0
-- not needed: kubectl apply -f https://github.com/kubernetes-sigs/gateway-api/releases/download/v1.3.0/standard-install.yaml
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

### Prepare certs for the application 

 Download Cert and token from an existing pod:
```shell
sudo mkdir /var/run/secrets/kubernetes.io/serviceaccount/ -p
sudo chwon codespace /var/run/secrets/kubernetes.io/serviceaccount/

kubectl cp traefik-7ddc96d-28flq:/var/run/secrets/kubernetes.io/serviceaccount/..data/ca.crt ./ca.crt -n traefik

kubectl cp traefik-7ddc96d-28flq:/var/run/secrets/kubernetes.io/serviceaccount/..data/token ./token

cp * /var/run/secrets/kubernetes.io/serviceaccount/
```
### Run the application (in background mode)


```shell
go build cmd/server/main.go 
export KUBERNETES_SERVICE_PORT=$(docker port kind-control-plane|sed -e "s/.*://")
export KUBERNETES_SERVICE_HOST=127.0.0.1
export PORT=8080
./main &
```


### Test 

curl "http://127.0.0.1:8080/api/routes?namespace=traefik"

curl "http://127.0.0.1:8080/api/route/traefik/weight-route"|jq
