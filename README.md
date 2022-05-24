# git-mirror

Mirror git repositories.

## Architecture

```mermaid
flowchart LR
  git
  shared-storage[(efs)]
  nginx
  rancher
  etl %% extract, transform, load
  grafana

  git -. write .-> shared-storage
  nginx -. read .-> shared-storage
  rancher -. read .-> nginx
  nginx -. fluentd protocol .-> etl
  archive --> shared-storage
  enrich -. influx .-> grafana

  subgraph chart %% currently git-mirror
    git
    nginx
  end

  %% log aggregation
  subgraph etl %% currently fluentd configmap
    fluentbit
    archive
    enrich

    fluentbit --> archive
    fluentbit --> enrich
  end
```

