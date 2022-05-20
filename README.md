# git-mirror

```mermaid
flowchart LR
    rancher-server --> git-mirror
    serve --> rancher-catalog-stats
    rancher-catalog-stats --> metrics
    subgraph git-mirror
    mirror((git-mirror)) --> serve((git-serve))
    logrotate((git-logrotate)) --> serve 
    analyze((git-analyze)) --> serve
    end
    subgraph metrics
    influxdb --> grafana
    end
```

Mirror a set of Git repositories. This can be useful for a variety of reasons:

1. Bypass geo-restrictions
2. MITM client traffic for forensic purposes

# Deployment Options

## Rancher

A template is provided for running on Rancher. Add the github.com/rancher/eio-charts repository as a Rancher Catalog to use it.
