#!/bin/sh
set -eu

CONFIG_FILE="/etc/git-mirror/config.yaml"

# Check for environment variable otherwise
# attempt to read from configuration file
if [ -z "${MIRROR_DATA_DIR}" ]; then
  mirror_data_dir="$(yq e '.storagedir' ${CONFIG_FILE})"
else
  mirror_data_dir="${MIRROR_DATA_DIR}"
fi

if [ -z "${MIRROR_REPO}" ]; then
  mirror_repo="$(yq e '.repositories[]' ${CONFIG_FILE})"
else
  mirror_repo="$(echo ${MIRROR_REPO} | sed 's/,//g')"
fi

cd "${mirror_data_dir}"

for repository in $mirror_repo; do
  git clone "${repository}" &
done

wait

