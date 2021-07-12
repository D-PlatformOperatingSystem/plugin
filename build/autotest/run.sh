#!/usr/bin/env bash

set -e
set -o pipefail
#set -o verbose
#set -o xtrace

# os: ubuntu18.04 x64

sedfix=""
if [ "$(uname)" == "Darwin" ]; then
    sedfix=".bak"
fi

## get dplatformos path
DPOS_PATH=$(go list -f "{{.Dir}}" github.com/D-PlatformOperatingSystem/dplatformos)

function build_auto_test() {

    trap "rm -f ../autotest/main.go" INT TERM EXIT
    local AutoTestMain="${DPOS_PATH}/cmd/autotest/main.go"
    cp "${AutoTestMain}" ./
    sed -i $sedfix '/^package/a import _ \"github.com\/D-PlatformOperatingSystem\/plugin\/plugin\"' main.go
    go build -v -i -o autotest
}

function copyAutoTestConfig() {

    declare -a DplatformOSAutoTestDirs=("${DPOS_PATH}/system" "../../plugin")
    echo "#copy auto test config to path \"$1\""
    local AutoTestConfigFile="$1/autotest.toml"

    #pre config auto test
    {

        echo 'cliCmd="./dplatformos-cli"'
        echo "checkTimeout=60"
    } >"${AutoTestConfigFile}"

    #copy all the dapp test case config file
    for rootDir in "${DplatformOSAutoTestDirs[@]}"; do

        if [[ ! -d ${rootDir} ]]; then
            continue
        fi

        testDirArr=$(find "${rootDir}" -type d -name autotest)

        for autotest in ${testDirArr}; do

            dapp=$(basename "$(dirname "${autotest}")")
            dappConfig=${autotest}/${dapp}.toml

            #make sure dapp have auto test config
            if [ -e "${dappConfig}" ]; then

                cp "${dappConfig}" "$1"/

                #add dapp test case config
                {
                    echo "[[TestCaseFile]]"
                    echo "dapp=\"$dapp\""
                    echo "filename=\"$dapp.toml\""
                } >>"${AutoTestConfigFile}"

            fi

        done
    done
}

function copyDplatformOS() {

    echo "# copy dplatformos bin to path \"$1\", make sure build dplatformos"
    cp ../dplatformos ../dplatformos-cli ../dplatformos.toml "$1"
    cp "${DPOS_PATH}"/cmd/dplatformos/dplatformos.test.toml "$1"
}

function copyAll() {

    dir="$1"
    #check dir exist
    if [[ ! -d ${dir} ]]; then
        mkdir "${dir}"
    fi
    cp autotest "${dir}"
    copyAutoTestConfig "${dir}"
    copyDplatformOS "${dir}"
    echo "# all copy have done!"
}

function main() {

    if [[ $1 == "build" ]]; then #build autotest
        build_auto_test
    else
        dir="$1"
        echo "$dir"
        rm -rf ../autotest/"$dir" && mkdir "$dir"
        cp -r "$DPOS_PATH"/build/autotest/"$dir"/* ./"$dir"/ && copyAll "$dir"
        chmod -R 755 "$dir" && cd "$dir" && ./autotest.sh "${@:2}" && cd ../
    fi
}

main "$@"
