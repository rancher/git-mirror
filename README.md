# git-mirror

Mirror git repositories.

## Architecture

```mermaid
flowchart LR
  rancher

  git-clone  -. create .-> persistent-volume
  git-porter -. read   .-> persistent-volume
  git-scrub  -. delete .-> persistent-volume
  git-mirror -. update .-> persistent-volume

  rancher -. read .-> git-porter

  subgraph chart
    git-scrub
    git-clone
    git-porter
    git-mirror
    persistent-volume
  end
```
