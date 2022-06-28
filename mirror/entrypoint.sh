#!/bin/sh
set -eu

TMP=$(mcookie)
TMP="/tmp/${TMP}"

mkdir -p "${TMP}"

cd "${TMP}"

git clone "${GIT_REPOSITORY_URL}" "${GIT_REPOSITORY_NAME}"

cd "${GIT_REPOSITORY_NAME}"

git remote add local "file:///var/git/${GIT_REPOSITORY_NAME}.git"

git push local "${GIT_REPOSITORY_BRANCH:-master}"

rm -rf "${TMP}"

