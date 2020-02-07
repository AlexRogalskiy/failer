# Failer

Simulate failure, in specific ways.

## Run locally

```bash
go run main.go
```

```bash
$ go run main.go -h
  -addr string
    	address to serve on (default ":8080")
  -distribution string
    	failure distribution, must be one of: random, contiguous, evenly (default "random")
  -success-rate int
    	server succcess rate percentage [0,100] (default 100)
```

### Failure Distribution

Three options for failure distribution:

```bash
go run main.go --distribution=random

# given success rate 80%, succeed 80 times, then fail 20 times
go run main.go --distribution=contiguous

# given success rate 80%, succeed 4 times, then fail 1 time
go run main.go --distribution=evenly
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
