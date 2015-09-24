# RFC03 - Jobs

## Schedule new Jobs

1. via http
-  http post --json

2. via ctl (implicit http)
- zeitd --publish <alias>.job
- zeitd --publish  # publish all from ~/.zeitd/jobs

## Scheduler

1. At bootstrap
- Purge internal jobCache if not empty
- Look up on etcd for all registered jobs and schedule them in jobCache
 (wait-index=1)
- keep track of the last createdIndex 


2. New job received
- watch for new jobs since create-index=last
```
etcdctl watch --recursive --after-index=135 /zeitd/jobs
```
- Job to be schedule on etcd and schedule it in memory




