# magnifibot
Magnifibot is an app that sends you the Gospel and lectures of the day

## Get the lectures

This is a simple POST:

```
curl -XPOST 'https://www.archimadrid.org/index.php/oracion-y-liturgia/index.php?option=com_archimadrid&format=ajax&task=leer_lecturas' \
-H 'Accept: application/json, text/javascript, */*; q=0.01' \
-H 'Content-Type: application/x-www-form-urlencoded; charset=UTF-8' \
-d 'dia=2022-01-10'
```

## To Do

### Required

- Tests
- Elasticache for storing the gospel. The key would be the actual day, the value would be the response we get from Archimadrid, before being parsed.

### Optional

- We could also monitor Redis Cache hits and misses with [OpenTelemetry](https://blog.uptrace.dev/posts/opentelemetry-metrics-cache-stats/).
