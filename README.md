# git-mirror

Mirror git repositories.

## Architecture

```mermaid
flowchart LR
  rancher

  git-scrub -. clean .-> persistent-volume
  git-clone -. write .-> persistent-volume
  git-porter -. read .-> persistent-volume
  git-mirror -. write .-> persistent-volume

  rancher -. read .-> git-porter

  subgraph chart
    git-scrub
    git-clone
    git-porter
    git-mirror
    persistent-volume
  end
```

