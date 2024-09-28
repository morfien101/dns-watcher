# dns-watcher

Watches for changes in dns and triggers defined actions.

This is useful if you have something that is sensitive to DNS changes but is unable to fix the problem itself.

Originally created as a add on for nginx to watch upstream servers that change regularly (AWS API Gateway).
Nginx does the DNS lookup for upstreams at config read and never again. For API Gateway, within a few minutes these will have changed.
You can reload the service to read new records.

# Usage

Download the correct binary for your system arch. Then create a yaml file with your records and actions. Then run it.

You can get an example of the config using the `-example` flag.

```yaml
sleep: 60
log_level: info
watchers:
  - records: 
    - a.example.com
    command: /bin/bash
    args:
      - echo
      - example.com has changed!
  - records: 
    - b.example.com
    command: /bin/bash
    args:
      - systemctl
      - restart
      - nginx
```

Tests are done sequentially and then the sleep happens.
You can specify multiple records and trigger a single action to stop the checks from flapping the service, if a service restart is what you are doing.

Failed tests or actions will be logged but will not stop the watcher. Therefore it is important that you watch the logs for failures.