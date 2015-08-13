Key: Leader election
/xchronos/var/scheduler/election value=<node_id> (TTL:heartbeat)

Dir: Job executors
/xchronos/etc/executors/<node_id> (TTL:heartbeat)

Dir: Jobs
/xchronos/etc/jobs/<job_id>



# Xchronos

It is a distributed and fault-tolerant scheduler that runs on top of a number of job stores (etcd/consul), that can be used for task orchestration.


## Features

- Support for ISO8061 
- Stats 
- Job history
- Docker support
- Configurable retry policy
- Fault tolerance
- Integration with golang apps

## Misfire Instructions

MISFIRE_INSTRUCTION_FIRE_NOW - execute as soon misfire has been identified
MISFIRE_INSTRUCTION_IGNORE_MISFIRE_POLICY - ignore
MISFIRE_INSTRUCTION_RESCHEDULE_NEXT_WITH_EXISTING_COUNT - it will honor the # total of fires
MISFIRE_INSTRUCTION_RESCHEDULE_NEXT_WITH_REMAINING_COUNT - it will ignore the misfires, not honoring the # total of fires
MISFIRE_INSTRUCTION_RESCHEDULE_NOW_WITH_EXISTING_REPEAT_COUNT
MISFIRE_INSTRUCTION_RESCHEDULE_NOW_WITH_EXISTING_REPEAT_COUNT

## Triggers

- simple:
```
name — the name that identifies the trigger;
startDelay — delay (in milliseconds) between scheduler startup and first job’s execution;
repeatInterval — timeout (in milliseconds) between consecutive job’s executions;
repeatCount — trigger will fire job execution (1 + repeatCount) times and stop after that (specify 0 here to have one-shot job or -1 to repeat job executions indefinitely);
```
- cron:
```
name — the name that identifies the trigger;
startDelay — delay (in milliseconds) between scheduler startup and first job’s execution;
cronExpression — cron expression
```

- custom:
```
triggerClass — your class which implements Trigger interface;
```

Example:
```
>> statement_generation.service

[Unit]
Description=My Apache Frontend
After=docker.service
Requires=docker.service

[Service]
TimeoutStartSec=0
ExecStartPre=-/usr/bin/docker kill apache1
ExecStartPre=-/usr/bin/docker rm apache1
ExecStartPre=/usr/bin/docker pull coreos/apache
ExecStart=/usr/bin/docker run -rm --name apache1 -p 80:80 coreos/apache /usr/sbin/apache2ctl -D FOREGROUND
ExecStop=/usr/bin/docker stop apache1

[X-Fleet]
MachineMetadata="region=us-east-1" "diskType=SSD"

[X-Chronos]
JobStore=etcd (*)

TriggerStartDelay=10000 
TriggerRepeatInterval=2000 
TriggerRepeatCount=-1 (*)
TriggerCron= * * 1 * * *

MisfirePolicy=MISFIRE_INSTRUCTION_RESCHEDULE_NEXT_WITH_REMAINING_COUNT  
MaxAttempts=3
TimeBeetweenAttempts=10000
Availability=2
```




