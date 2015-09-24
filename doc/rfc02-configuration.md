# RFC02 - Configuration

1. Xchronosfile
```
version = "0.1"
etcd_nodes = ["10.1.1.1", "10.1.1.2"] 

```

2. <alias>.batch

```
job "outbound_email_marketing" {
    trigger {
        cron = "* * * * 0/30"
        max_executions = -1
    }
    exec = "~/bin/outbound-email.sh"
    on_error {
        max_retries = "3"
        abort_on_failure = true
    }
}
```
