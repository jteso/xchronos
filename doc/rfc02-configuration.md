#Configuration

```
version = "0.1"
job_store = "etcd" 

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