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

- Fix too long Telegram Messages by splitting the message into multiple messages
- Backup DynamoDB table
- Add on-demand feature
- Setup usage CloudWatch alarms for lambda usage
- Add schedule hours feature

### Optional

- Explore the possibility of making the bot private
