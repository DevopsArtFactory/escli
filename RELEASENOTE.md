# Release Note

## 0.0.4 (2021/03/30)
- autocomplete is included.
- `profiles list` command shows all profiles of escli.
- `profiles add` command add profile to configuration file of escli
- `profiles remove [profileName]` command remove profile from configuration file of escli
- `init` command is deprecated.
 
## 0.0.3 (2021/02/05)
- `cat health` command shows all things of `_cat/health` API
- `index delete` command is added. you can delete index that command.
- `update` command is added. you can update escli that command.
- `snapshot list --repo-only` command shows more information of snapshot repository
- bug fix : can't restore and archive snapshot when base path of repository is blank. it is fixed.
- `diag` command shows minimum disk used percent of data node.

## 0.0.2 (2021/01/14)
- support `_cluster/reroute?retry_failed`. it helps reroute unassigned shards. `_cluster/reroute` with `retry_failed` query string reroute unassigned shards automatically.
- add `force` options to `snapshot delete` , `snapshot create`, `index settings`, `cluster settings`
- when initialize escli, `aws_default_region` config is not mandatory. if you use elasticsearch on-premise environments, so you don't have `aws_region`, you type blank that parameter.

## 0.0.1 (2021/01/13)
- first version
  