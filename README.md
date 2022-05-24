# git-mirror

Mirror git repositories.

## Architecture

```mermaid
flowchart LR
  git-sync
  shared-storage[(efs)]
  git-porter
  rancher
  etl %% extract, transform, load
  grafana

  git-sync -. write .-> shared-storage
  git-porter -. read .-> shared-storage
  rancher -. read .-> git-porter
  git-sync -. fluentd protocol .-> etl
  archive --> shared-storage
  enrich -. influx .-> grafana

  subgraph chart %% currently git-mirror
    git-sync
    git-porter
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

