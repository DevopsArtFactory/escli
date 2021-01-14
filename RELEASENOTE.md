# Release Note

## 0.0.2 (2021/01/14)
- support `_cluster/reroute?retry_failed`. it helps reroute unassigned shards. `_cluster/reroute` with `retry_failed` query string reroute unassigned shards automatically.
- add `force` options to `snapshot delete` , `snapshot create`, `index settings`, `cluster settings`
- when initialize escli, `aws_default_region` config is not mandatory. if you use elasticsearch on-premise environments, so you don't have `aws_region`, you type blank that parameter.

## 0.0.1 (2021/01/13)
- first version
  