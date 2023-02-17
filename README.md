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
$ escli index settings my-index-000001 number_of_replicas 2
```
`escli` should be make your elasticsearch experience more powerful.

## Release Note 
Release Note is [Here](RELEASENOTE.md)

>**WARNING**
>If you used escli 0.0.4, you have to change your config file. configuration field `elasticsearch_url` is changed `url`.

[AS-IS]
```bash
- profile: localhost
  elasticsearch_url: http://localhost:9200
  aws_region: ap-northeast-2
```

[TO-BE]
```bash
- profile: localhost
  url: http://localhost:9200
  aws_region: ap-northeast-2
```

## Installation
### Required
- [elasticsearch](https://elastic.co) version 6.0 or higher
- [opensearch](https://opensearch.org) version 1.0 or higher

### Install escli binary 

#### On MAC
```bash
$ brew tap devopsartfactory/devopsart
$ brew update
$ brew install escli 
``` 

#### On Linux
```bash
$ curl -Lo escli https://escli.s3.ap-northeast-2.amazonaws.com/escli/releases/latest/escli-linux-amd64
$ sudo install escli /usr/bin
```

### initialize escli
configuration of escli is stored at `~/.escli/config.yaml` file.
for the first time, there is no configuration so you have to initialize configuration with `escli profiles add`
```bash
? Your Profile Name :  localhost
? Your ElasticSearch or OpenSearch URL :  https://localhost:9200
? Select your product (elasticsearch or opensearch) :  elasticsearch
? Your AWS Default Region (If you don't use AWS, type blank) :
? Your HTTP Username (If you don't use http basic authentication, type blank) :  elastic
? Your HTTP Password (If you don't use http basic authentication, type blank) :  ********************
? Your certificateFingerPrint (If you don't use certificate finger print, type blank) :  67c2d588a7a6a50e773d0cc91f83cab7a6c11a19e199f3f075b9b6d873a5992a
- profile: localhost
  url: https://localhost:9200
  product: elasticsearch
  http_username: elastic
  http_password: qZzEp0Hc112zYx=Z+xQb
  certificate_finger_print: 67c2d588a7a6a50e773d0cc91f83cab7a6c11a19e199f3f075b9b6d873a5992a

? Are you sure to add profile to configuration file?  yes
Adding profile to configuration file is successfully in /Users/alden/.escli/config.yaml
```

## How to use

### configuration
| field | description                                                                             |
| ------|-----------------------------------------------------------------------------------------|
| profile | name of profile                                                                         |
| url | url of target system                                                                    |
| product | product of target system. elasticsearch or opensearch         (default : elasticsearch) |
| aws_region | aws region that you use                                                                 |
| http_username | http username of target system. It is needed if you use basic http authentication.      |
| http_password | http password of target system.                                                         |
| certificate_fingerprint | certificate fingerprint of target system                                                |
#### configuration example
```bash
- profile: dev-access-log
  url: https://dev-access-log.ap-northeast-2.es.amazonaws.com
  product: opensearch
- profile: prod-access-log
  url: https://prod-access-log.ap-northeast-2.es.amazonaws.com
  product: opensearch
  http_username: elastic
  http_password: abcdefg
- profile: localhost
  url: https://localhost:9200
  product: elasticsearch
  http_username: elastic
  http_password: qZzEp0Hc112zYx=Z+xQb
  certificate_finger_print: 67c2d588a7a6a50e773d0cc91f83cab7a6c11a19e199f3f075b9b6d873a5992a
```

### common options
* `--profile` : you can specify profile from configuration file. if you don't specify `--profile` option, escli use first profile of configuration file.
* `--config` : you can specify configuration file. if you don't specify `--config` option, escli use `~/.escli/config.yaml` configuration file.

### `profiles` command
you can add one more elasticsearch clusters to your configuration file by `profiles` command. and then you can use profile with `--profile` option

#### command list
| command     | description                                               |
| ----------- | --------------------------------------------------------- |
| profiles list | shows profiles.        |
| profiles add | add profile to configuration file  |
| profiles remove | remove profile from configuration file    |

#### examples

```bash
$ escli profiles add
? Your Profile Name :  localhost
? Your ElasticSearch or OpenSearch URL :  https://localhost:9200
? Select your product (elasticsearch or opensearch) :  elasticsearch
? Your AWS Default Region (If you don't use AWS, type blank) :
? Your HTTP Username (If you don't use http basic authentication, type blank) :  elastic
? Your HTTP Password (If you don't use http basic authentication, type blank) :  ********************
? Your certificateFingerPrint (If you don't use certificate finger print, type blank) :  67c2d588a7a6a50e773d0cc91f83cab7a6c11a19e199f3f075b9b6d873a5992a
- profile: localhost
  url: https://localhost:9200
  product: elasticsearch
  http_username: elastic
  http_password: qZzEp0Hc112zYx=Z+xQb
  certificate_finger_print: 67c2d588a7a6a50e773d0cc91f83cab7a6c11a19e199f3f075b9b6d873a5992a

? Are you sure to add profile to configuration file?  yes
Adding profile to configuration file is successfully in /Users/alden/.escli/config.yaml
```

```bash
$ escli profiles list
Profile                  : localhost
URL                      : https://localhost:9200
Product                  : elasticsearch
HTTP Username            : elastic
HTTP Password            : ************
Certificate Finger Print : 67c2d588a7a6a50e773d0cc91f83cab7a6c11a19e199f3f075b9b6d873a5992a
```

```bash
$ escli cat health --profile log-es
```


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

`snapshot` command doesn't support OpenSearch.

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
| command        | description |
|----------------| ---------------- |
| index settings | get or set index settings |
| index delete   | delete index |
| index create   | create index |
| index stats    | show statistics of index |


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

```bash
$ escli index stats access_log-2023.02.13 1 --profile=localhost
time      	index               	        total shards	   successful shards	       failed shards	       indexing rate	indexing latency (ms)	          query rate	  query latency (ms)	          fetch rate	  fetch latency (ms)
16:08:13  	access_log-2023.02.13	                  12	                  12	                   0	                3182	                0.13	                   0	                0.00	                   0	                0.00
16:08:14  	access_log-2023.02.13	                  12	                  12	                   0	                2348	                0.09	                   0	                0.00	                   0	                0.00
16:08:15  	access_log-2023.02.13	                  12	                  12	                   0	                2466	                0.12	                   0	                0.00	                   0	                0.00
16:08:16  	access_log-2023.02.13	                  12	                  12	                   0	                   0	                0.00	                   0	                0.00	                   0	                0.00
16:08:17  	access_log-2023.02.13	                  12	                  12	                   0	                6046	                0.14	                   0	                0.00	                   0	                0.00
16:08:18  	access_log-2023.02.13	                  12	                  12	                   0	                5056	                0.17	                   0	                0.00	                   0	                0.00
16:08:19  	access_log-2023.02.13	                  12	                  12	                   0	                1286	                0.10	                   0	                0.00	                   0	                0.00
```

### `cluster` command

#### command list
| command | description |
| ------- | ----------- |
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

### `stats` command

#### examples

```bash
$ escli stats 1 --profile=localhost
time      	        total shards	   successful shards	       failed shards	       indexing rate	indexing latency (ms)	          query rate	  query latency (ms)	          fetch rate	  fetch latency (ms)
16:10:32  	                 204	                 204	                   0	               10591	                0.47	                   0	                0.00	                   0	                0.00
16:10:33  	                 204	                 204	                   0	                1099	                8.16	                   0	                0.00	                   0	                0.00
16:10:34  	                 204	                 204	                   0	                3267	                0.13	                   0	                0.00	                   0	                0.00
16:10:35  	                 204	                 204	                   0	                1869	                0.13	                   0	                0.00	                   0	                0.00
```

## Autocompletion
* zsh
```bash
$ echo "source <(escli completion zsh)" >> ~/.zshrc
$ source  ~/.zshrc
```

* bash
```bash
$ echo "source <(escli completion bash)" >> ~/.bash_rc or ~/.bash_profile
$ source  ~/.bashrc
```