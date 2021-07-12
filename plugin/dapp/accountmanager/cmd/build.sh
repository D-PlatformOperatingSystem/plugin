#!/bin/sh
#   ci
strpwd=$(pwd)
strcmd=${strpwd##*dapp/}
strapp=${strcmd%/cmd*}

OUT_DIR="${1}/$strapp"
#FLAG=$2
echo "${OUT_DIR}"
# mkdir -p "${OUT_DIR}"
# cp ./build/* "${OUT_DIR}"
