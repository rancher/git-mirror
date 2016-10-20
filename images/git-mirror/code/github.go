package main

import (
  "encoding/json"
)

var exampleJsonPush = []byte(`{
  "ref": "refs/heads/test",
  "before": "0000000000000000000000000000000000000000",
  "after": "7abc80b1847e64af2658317ef3307875e4ca8b84",
  "created": true,
  "deleted": false,
  "forced": false,
  "base_ref": "refs/heads/master",
  "compare": "https://github.com/LLParse/rancher-catalog/compare/test",
  "commits": [],
  "head_commit": {
    "id": "7abc80b1847e64af2658317ef3307875e4ca8b84",
    "tree_id": "83647aba0849f06c02451cd5f66cc07aae8c9fa9",
    "distinct": true,
    "message": "Merge pull request #197 from alena1108/114\n\nkubernetes-agent/ingress-controller/kubernetes images update for k8s 1.2.4",
    "timestamp": "2016-09-23T15:55:46-07:00",
    "url": "https://github.com/LLParse/rancher-catalog/commit/7abc80b1847e64af2658317ef3307875e4ca8b84",
    "author": {
      "name": "Denise",
      "email": "denise@rancher.com",
      "username": "deniseschannon"
    },
    "committer": {
      "name": "GitHub",
      "email": "noreply@github.com",
      "username": "web-flow"
    },
    "added": [
      "templates/k8s/11/docker-compose.yml",
      "templates/k8s/11/rancher-compose.yml",
      "templates/kubernetes/11/docker-compose.yml",
      "templates/kubernetes/11/rancher-compose.yml"
    ],
    "removed": [],
    "modified": [
      "templates/kubernetes/config.yml"
    ]
  },
  "repository": {
    "id": 54731015,
    "name": "rancher-catalog",
    "full_name": "LLParse/rancher-catalog",
    "owner": {
      "name": "LLParse",
      "email": "joliver@rancher.com"
    },
    "private": false,
    "html_url": "https://github.com/LLParse/rancher-catalog",
    "description": null,
    "fork": true,
    "url": "https://github.com/LLParse/rancher-catalog",
    "forks_url": "https://api.github.com/repos/LLParse/rancher-catalog/forks",
    "keys_url": "https://api.github.com/repos/LLParse/rancher-catalog/keys{/key_id}",
    "collaborators_url": "https://api.github.com/repos/LLParse/rancher-catalog/collaborators{/collaborator}",
    "teams_url": "https://api.github.com/repos/LLParse/rancher-catalog/teams",
    "hooks_url": "https://api.github.com/repos/LLParse/rancher-catalog/hooks",
    "issue_events_url": "https://api.github.com/repos/LLParse/rancher-catalog/issues/events{/number}",
    "events_url": "https://api.github.com/repos/LLParse/rancher-catalog/events",
    "assignees_url": "https://api.github.com/repos/LLParse/rancher-catalog/assignees{/user}",
    "branches_url": "https://api.github.com/repos/LLParse/rancher-catalog/branches{/branch}",
    "tags_url": "https://api.github.com/repos/LLParse/rancher-catalog/tags",
    "blobs_url": "https://api.github.com/repos/LLParse/rancher-catalog/git/blobs{/sha}",
    "git_tags_url": "https://api.github.com/repos/LLParse/rancher-catalog/git/tags{/sha}",
    "git_refs_url": "https://api.github.com/repos/LLParse/rancher-catalog/git/refs{/sha}",
    "trees_url": "https://api.github.com/repos/LLParse/rancher-catalog/git/trees{/sha}",
    "statuses_url": "https://api.github.com/repos/LLParse/rancher-catalog/statuses/{sha}",
    "languages_url": "https://api.github.com/repos/LLParse/rancher-catalog/languages",
    "stargazers_url": "https://api.github.com/repos/LLParse/rancher-catalog/stargazers",
    "contributors_url": "https://api.github.com/repos/LLParse/rancher-catalog/contributors",
    "subscribers_url": "https://api.github.com/repos/LLParse/rancher-catalog/subscribers",
    "subscription_url": "https://api.github.com/repos/LLParse/rancher-catalog/subscription",
    "commits_url": "https://api.github.com/repos/LLParse/rancher-catalog/commits{/sha}",
    "git_commits_url": "https://api.github.com/repos/LLParse/rancher-catalog/git/commits{/sha}",
    "comments_url": "https://api.github.com/repos/LLParse/rancher-catalog/comments{/number}",
    "issue_comment_url": "https://api.github.com/repos/LLParse/rancher-catalog/issues/comments{/number}",
    "contents_url": "https://api.github.com/repos/LLParse/rancher-catalog/contents/{+path}",
    "compare_url": "https://api.github.com/repos/LLParse/rancher-catalog/compare/{base}...{head}",
    "merges_url": "https://api.github.com/repos/LLParse/rancher-catalog/merges",
    "archive_url": "https://api.github.com/repos/LLParse/rancher-catalog/{archive_format}{/ref}",
    "downloads_url": "https://api.github.com/repos/LLParse/rancher-catalog/downloads",
    "issues_url": "https://api.github.com/repos/LLParse/rancher-catalog/issues{/number}",
    "pulls_url": "https://api.github.com/repos/LLParse/rancher-catalog/pulls{/number}",
    "milestones_url": "https://api.github.com/repos/LLParse/rancher-catalog/milestones{/number}",
    "notifications_url": "https://api.github.com/repos/LLParse/rancher-catalog/notifications{?since,all,participating}",
    "labels_url": "https://api.github.com/repos/LLParse/rancher-catalog/labels{/name}",
    "releases_url": "https://api.github.com/repos/LLParse/rancher-catalog/releases{/id}",
    "deployments_url": "https://api.github.com/repos/LLParse/rancher-catalog/deployments",
    "created_at": 1458923046,
    "updated_at": "2016-03-25T16:24:07Z",
    "pushed_at": 1476410867,
    "git_url": "git://github.com/LLParse/rancher-catalog.git",
    "ssh_url": "git@github.com:LLParse/rancher-catalog.git",
    "clone_url": "https://github.com/LLParse/rancher-catalog.git",
    "svn_url": "https://github.com/LLParse/rancher-catalog",
    "homepage": null,
    "size": 435,
    "stargazers_count": 0,
    "watchers_count": 0,
    "language": "Python",
    "has_issues": false,
    "has_downloads": true,
    "has_wiki": true,
    "has_pages": false,
    "forks_count": 0,
    "mirror_url": null,
    "open_issues_count": 0,
    "forks": 0,
    "open_issues": 0,
    "watchers": 0,
    "default_branch": "master",
    "stargazers": 0,
    "master_branch": "master"
  },
  "pusher": {
    "name": "LLParse",
    "email": "joliver@rancher.com"
  },
  "sender": {
    "login": "LLParse",
    "id": 6145659,
    "avatar_url": "https://avatars.githubusercontent.com/u/6145659?v=3",
    "gravatar_id": "",
    "url": "https://api.github.com/users/LLParse",
    "html_url": "https://github.com/LLParse",
    "followers_url": "https://api.github.com/users/LLParse/followers",
    "following_url": "https://api.github.com/users/LLParse/following{/other_user}",
    "gists_url": "https://api.github.com/users/LLParse/gists{/gist_id}",
    "starred_url": "https://api.github.com/users/LLParse/starred{/owner}{/repo}",
    "subscriptions_url": "https://api.github.com/users/LLParse/subscriptions",
    "organizations_url": "https://api.github.com/users/LLParse/orgs",
    "repos_url": "https://api.github.com/users/LLParse/repos",
    "events_url": "https://api.github.com/users/LLParse/events{/privacy}",
    "received_events_url": "https://api.github.com/users/LLParse/received_events",
    "type": "User",
    "site_admin": false
  }`)

type GHPushEvent struct {
  Repo *GHRepository `json:"repository"`
}

type GHRepository struct {
  Id       int64  `json:"id"`
  Name     string `json:"name"`
  FullName string `json:"full_name"`
}

func ParsePushEvent(data []byte) (*GHPushEvent, error) {
  event := GHPushEvent{}
  err := json.Unmarshal(data, &event)
  return &event, err
}
