# escli
`escli` is a command-line tool for managing elasticsearch cluster. If you want to set number_of_replicas of index, you should do with curl like below example.
```bash
$ curl -X PUT "localhost:9200/my-index-000001/_settings?pretty" -H 'Content-Type: application/json' -d'
{
  "index" : {
    "number_of_replicas" : 2
  }
}
'
```
But with escli, you type below command.
```bash
$ escli index setting my-index-000001 number_of_replicas 2
```
`escli` should be make your elasticsearch experience more powerful.

## Installation
### Required
- [elasticsearch](https://elastic.co) version 6.0 or higher

### Install escli binary 

#### On MAC
```bash
$ brew tap devopsartfactory/devopsart
$ brew update
$ brew install escli 
``` 

#### On Linux
```bash
$ curl -Lo escli https://escli.s3.ap-northeast-2.amazonaws.com/releases/latest/escli-linux-amd64
$ sudo install escli /usr/bin
```

### initialize escli
configuration of escli is stored at `~/.escli/config.yaml` file.
for the first time, there is no configuration so you have to initialize configuration with `escli init`
```bash
$ escli init
? Your ElasticSearch URL :  http://elasticsearch.domain.com:9200
? Your AWS Default Region :  ap-northeast-2
elasticsearchurl: http://elasticsearch.domain.com:9200
awsregion: ap-northeast-2

? Are you sure to generate configuration file?  y
New configuration file is successfully generated in /Users/benx/.escli/config.yaml

```

## How to use

### `cat` command

#### command list
| command     | description                                               |
| ----------- | --------------------------------------------------------- |
| cat health  | shows health of cluster. it calls `_cat/health` API       |
| cat indices | shows information of indices. it calls `_cat/indices` API |
| cat nodes   | shows information of nodes. it calls `_cat/nodes` API     |
| cat shards  | shows information of shards. it calls `_cat/shards` API   |

#### available options
| option        | description                                                      |
| ------------- | ---------------------------------------------------------------- |
| troubled-only | shows objects with trouble. such as yellow status or red status  |
| sort-by       | set sort key.                                                    |

#### examples

```bash
$ escli cat indices

index                                                   health  status  pri     rep     store.size
.kibana_1                                               green   open    1       1          918.1kb
.kibana_task_manager                                    green   open    1       1           26.5kb
.monitoring-es-6-2021.01.04                             green   open    1       1           21.2gb
.monitoring-es-6-2021.01.05                             green   open    1       1           21.7gb
```

```bash
$ escli cat indices --troubled-only
index                                             health  status  pri     rep     store.size
application-log-2020.11.01                        red     open    40      0          202.4gb
application-log-2020.11.02                        red     open    40      0          289.8gb
application-log-2020.11.03                        red     open    40      0          199.9gb
```

```bash
$ escli cat indices --sorted-by store.size:desc
index                                        health  status  pri     rep     store.size
application-log-2021.01.07                   green   open    100     1            3.3tb
application-log-2021.01.04                   green   open    100     1              3tb
application-log-2021.01.06                   green   open    100     1            2.9tb
```

### `snapshot` command

#### command list
| command     | description                                               |
| ----------- | --------------------------------------------------------- |
| snapshot list        | shows information of repositories and snapshots. it calls `_cat/snapshot` API. |
| snapshot archive     | change storage class of snapshots to S3 glacier. it works on AWS only.  |
| snapshot restore     | change storage class of snapshots to S3 standard and restore snapshot.    |
| snapshot create        | create snapshot of indices.    |

#### available options
| option        | description                                                      |
| ------------- | ---------------------------------------------------------------- |
| force | do not ask continue. it will be used for automated batch job. (only `archive`, `restore` command) |
| with-repo | shows snapshots of specified repo (only `list` command) |
| repo-only | shows only information of repos (only `list` command) |

#### examples

```bash
$ escli snapshot list --repo-only
Repository ID : log-archive
Repository ID : log-archive-standard-ia
Repository ID : log-archive-standard
```

```bash
$ escli snapshot create prod-snapshot snapshots-2021-01-01 result-prod-2021-01-01
snapshots-2021-01-01 is created
```

```bash
$ escli snapshot archive send-mail-result-prod-snapshot snapshots-2020-12-31 --region us-east-1
bucket name : result-prod
base path : elasticsearch-snapshot-standard
Downloaded /tmp/index-456 81207 bytes
index name : result-prod-2020-12-31
elasticsearch-snapshot-standard/indices/z8bqmUmAQxy8tuwSsmFEKg/0/__-urzTmmuR8K6s6kpLryZ5g
? Change Storage Class to GLACIER 
```

- If you use `--force` option to `snapshot arvhice` command, escli doesn't ask you to continue. It makes all archiving job.

### `index` command

#### command list
| command | description |
| index settings | get or set index settings |

#### examples

```bash
$ escli index settings prod-2021-01-01 
{
  "prod-2021-01-01" : {
    "settings" : {
      "index" : {
        "creation_date" : "1609906432373",
        "number_of_shards" : "5",
        "number_of_replicas" : "2",
        "uuid" : "ha6Y6uiCSfOV_syHJwFCqA",
        "version" : {
          "created" : "6080099"
        },
        "provided_name" : "send-mail-result-prod-2021-01-01"
      }
    }
  }
}
```

```bash
$ escli index settings prod-2021-01-12 number_of_replicas                                                                                                                                                                                               ok  3s 
{
  "prod-2021-01-12" : {
    "settings" : {
      "index" : {
        "number_of_replicas" : "1"
      }
    }
  }
}
```

```bash
$ escli index settings send-mail-result-prod-2021-01-12 number_of_replicas 2                                                                                                                                                                                                 ok 
{
  "acknowledged" : true
}
```

### `cluster` command

#### command list
| command | description |
| cluster settings | get or set index settings |

#### examples

```bash
$ escli cluster settings
{
  "persistent" : {
    "cluster" : {
      "routing" : {
        "allocation" : {
......
}
```

```bash
$ escli cluster settings persistent indices.recovery.max_bytes_per_sec 50mb
```

### `diag` command

#### examples

```bash
$ escli diag
check cluster status...........................[green] 😎
check yellow status indices....................[0] 😎
check red status indices.......................[0] 😎
check number of master nodes...................[3]
check maximum disk used percent of nodes.......[36]
```