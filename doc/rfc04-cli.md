# RFC04 - CLI

# Sinopsis
zeitd *command* (*subcommand*) (*params*)

## command-subcommand
run agent (--no-scheduler|--no-executor)
run console

schedule (all|<file>.batch)

status

# Examples
Zeitd cluster of two nodes: Node1 and Node2
```
Node1 => $zeitd run agent --with-config "~/.zeitd/etc/zeitd.conf"
Node2 => $zeitd run agent --with-config "~/.zeitd/etc/zeitd.conf"
```

By default all batch files will be lookup in:
- current directory
- ~/.zeitd/etc

```
Jumpbox => $zeitd schedule all ("~/.zeitd/etc/*.batch")
Jumpbox => $zeitd schedule payment.batch
```

Check the current status:
```
Jumpbox => $zeitd status

Status: Running and Healthy
Nodes
    ----------------------------------------------
   | Hostname | Role                 | Status     |
   | -------- | -------------------- | ---------- |
   | Node1    | Scheduler, Executor  | Running    |
   | Node2    | Executor             | Running    |
    ----------------------------------------------
  
Jobs   

    --------------------------------------------------------------------------
   | Batch    | Job                  | Status              | NextRunAt        |
   | -------- | -------------------- | ------------------- | ---------------- |
   | payment  | cs2files-push        | scheduled           | Mon 12, 08:00:00 |
   | payment  | bpay-push            | Running (23 secs)   | Mon 12, 09:30:00 |
    --------------------------------------------------------------------------
  
```






