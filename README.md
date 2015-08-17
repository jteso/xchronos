## Run

1. Build the docker image
```
docker build -t xchronos .
```

2. Run, run, run

```
docker run --rm xchronos ./xchronos -etcd-nodes=http://10.1.42.1:4001
```

or whatever ip has been assigned to your docker0 interface
