# escli
`escli` is a command-line tool for managing elasticsearch cluster. 

## Installation
### Required
- [elasticsearch](https://elastic.co) version 6.0 or higher

### Install escli binary (TODO)
```bash
$ brew tap DevopsArtFactory/escli
$ brew update
$ brew install escli 
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
| snapshot take        | take snapshot of indices.    |

#### available options
| option        | description                                                      |
| ------------- | ---------------------------------------------------------------- |
| force | do not ask continue. it will be used for automated batch job.  |

#### examples

```bash
$ escli snapshot list
```

### `index` command

#### command list
| command     | description                                               |
| ----------- | --------------------------------------------------------- |
| index settings       | get or set settings of index.  |

#### available options
| option        | description                                                      |
| ------------- | ---------------------------------------------------------------- |
| force | do not ask continue. it will be used for automated batch job.  |

#### examples

