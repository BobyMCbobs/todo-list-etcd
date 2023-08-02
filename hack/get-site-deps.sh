#!/bin/bash

cd $(git rev-parse --show-toplevel)

OUTPUT_FOLDER=kodata/web/scripts/ext

rm -rf "${OUTPUT_FOLDER}"
mkdir -p "${OUTPUT_FOLDER}"

SCRIPTS=(
    https://cdn.jsdelivr.net/npm/axios@1.3.4/dist/axios.min.js
    https://cdn.jsdelivr.net/npm/dayjs@1.11.7/dayjs.min.js
    https://cdn.jsdelivr.net/npm/dayjs@1.11.7/plugin/timezone.js
    https://cdn.jsdelivr.net/npm/dayjs@1.11.7/plugin/utc.js
)

cd "${OUTPUT_FOLDER}"
for SCRIPT in ${SCRIPTS[*]}; do
    curl -O -L -s "${SCRIPT}"
done
