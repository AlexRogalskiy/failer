# Failer

Simulate failure, in specific ways.

## Run locally

```bash
go run main.go
```

## Docker build

```bash
docker build . -t buoyantio/failer:latest
```

## Deploy to Kubernetes

```bash
cat k8s.yml | kubectl apply -f -
```

### Observe failure in Kubernetes

```bash
kubectl -n failer get po --watch
```
