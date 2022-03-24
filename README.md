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

- Make /start working, so that when a new user joins they receive a welcome or a help message
- It's a little bit weird that when you're suscribed, the bot doesn't tell you anything except that you're going to receive the Gospel at some point
- Backup DynamoDB table
- Setup usage CloudWatch alarms for lambda usage
- Setup usage CloudWatch alarms for SQS usage
- Add schedule hours feature

### Optional

- Explore the possibility of making the bot private
