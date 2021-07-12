#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set +e
set -o pipefail

MAIN_HTTP=""
source ../dapp-test-common.sh

gID=""
gResp=""

glAddr=""
gameAddr1=""
gameAddr2=""
gameAddr3=""
bwExecAddr=""

init() {
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ipara=$ispara"

    if [ "$ispara" == true ]; then
        bwExecAddr=$(curl -ksd '{"method":"DplatformOS.ConvertExectoAddr","params":[{"execname":"user.p.para.blackwhite"}]}' ${MAIN_HTTP} | jq -r ".result")
    else
        bwExecAddr=$(curl -ksd '{"method":"DplatformOS.ConvertExectoAddr","params":[{"execname":"blackwhite"}]}' ${MAIN_HTTP} | jq -r ".result")
    fi
    echo "bwExecAddr=$bwExecAddr"
}

dplatformos_NewAccount() {
    label=$1
    req='{"method":"DplatformOS.NewAccount","params":[{"label":"'"$label"'"}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result.acc.addr|length > 0)' "$FUNCNAME" ".result.acc.addr"
    glAddr=$RETURN_RESP
}

dplatformos_SendTransaction() {
    rawTx=$1
    addr=$2
    #
    req='{"method":"DplatformOS.SignRawTx","params":[{"addr":"'"$addr"'","txHex":"'"$rawTx"'","expire":"120s","fee":10000000,"index":0}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not)' "DplatformOS.SignRawTx" ".result"
    signTx=$RETURN_RESP

    req='{"method":"DplatformOS.SendTransaction","params":[{"data":"'"$signTx"'"}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not)' "$FUNCNAME" ".result"

    gResp=$RETURN_RESP
    #
    dplatformos_QueryTx "$RETURN_RESP" "${MAIN_HTTP}"
}

blackwhite_BlackwhiteCreateTx() {
    #
    addr=$1
    req='{"method":"blackwhite.BlackwhiteCreateTx","params":[{"PlayAmount":100000000,"PlayerCount":3,"GameName":"hello","Timeout":600,"Fee":1000000}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not)' "$FUNCNAME" ".result"
    #
    dplatformos_SendTransaction "$RETURN_RESP" "${addr}"
    gID="${gResp}"
}

blackwhite_BlackwhitePlayTx() {
    addr=$1
    round1=$2
    round2=$3
    round3=$4
    req='{"method":"blackwhite.BlackwhitePlayTx","params":[{"gameID":"'"$gID"'","amount":100000000,"Fee":1000000,"hashValues":["'"$round1"'","'"$round2"'","'"$round3"'"]}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not)' "$FUNCNAME" ".result"

    #
    dplatformos_SendTransaction "$RETURN_RESP" "${addr}"
}

blackwhite_BlackwhiteShowTx() {
    addr=$1
    sec=$2
    req='{"method":"blackwhite.BlackwhiteShowTx","params":[{"gameID":"'"$gID"'","secret":"'"$sec"'","Fee":1000000}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not)' "$FUNCNAME" ".result"
    dplatformos_SendTransaction "$RETURN_RESP" "${addr}"
}

blackwhite_BlackwhiteTimeoutDoneTx() {
    gameID=$1
    req='{"method":"blackwhite.BlackwhiteTimeoutDoneTx","params":[{"gameID":"'"$gameID"'","Fee":1000000}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not)' "$FUNCNAME"
}

blackwhite_GetBlackwhiteRoundInfo() {
    gameID=$1
    req='{"method":"DplatformOS.Query","params":[{"execer":"blackwhite","funcName":"GetBlackwhiteRoundInfo","payload":{"gameID":"'"$gameID"'"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result.round | [has("gameID", "status", "playAmount", "playerCount", "curPlayerCount", "loop", "curShowCount", "timeout"),true] | unique | length == 1)' "$FUNCNAME"
}

blackwhite_GetBlackwhiteByStatusAndAddr() {
    addr=$1
    req='{"method":"DplatformOS.Query","params":[{"execer":"blackwhite","funcName":"GetBlackwhiteByStatusAndAddr","payload":{"status":5,"address":"'"$addr"'","count":1,"direction":0,"index":-1}}]}'
    resok='(.error|not) and (.result.round[0].createAddr == "'"$addr"'") and (.result.round[0].status == 5) and (.result.round[0] | [has("gameID", "status", "playAmount", "playerCount", "curPlayerCount", "loop", "curShowCount", "timeout", "winner"),true] | unique | length == 1)'
    dplatformos_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
}

blackwhite_GetBlackwhiteloopResult() {
    gameID=$1
    req='{"method":"DplatformOS.Query","params":[{"execer":"blackwhite","funcName":"GetBlackwhiteloopResult","payload":{"gameID":"'"$gameID"'","loopSeq":0}}]}'
    resok='(.error|not) and (.result.gameID == "'"$gameID"'") and (.result.results|length >= 1)'
    dplatformos_Http "$req" ${MAIN_HTTP} "$resok" "$FUNCNAME"
}

function run_testcases() {
    #
    sect1="123"
    black1="6vm6gJ2wvEIxC8Yc6r/N6lIU5OZk633YMnIfwcZBD0o="
    black2="6FXx5aeDSCaq1UrhLO8u0H31Hl8TpvzxuHrgGo9WeFk="
    white0="DrNPzA68XiGimZE/igx70kTPJxnIJnVf8NCGnb7XoYU="
    white1="SB5Pnf6Umf2Wba0dqyNOezq5FEqTd22WPVYAhSA6Lxs="

    #
    dplatformos_NewAccount "label188"
    gameAddr1="${glAddr}"
    dplatformos_NewAccount "label288"
    gameAddr2="${glAddr}"
    dplatformos_NewAccount "label388"
    gameAddr3="${glAddr}"

    #
    origAddr="12oupcayRT7LvaC4qW4avxsTE7U41cKSio"

    dplatformos_GetAccounts "${MAIN_HTTP}"

    #
    M_HTTP=${MAIN_HTTP//8901/28803}
    dplatformos_SendToAddress "${origAddr}" "${gameAddr1}" 1000000000 "${M_HTTP}"
    dplatformos_SendToAddress "${origAddr}" "${gameAddr2}" 1000000000 "${M_HTTP}"
    dplatformos_SendToAddress "${origAddr}" "${gameAddr3}" 1000000000 "${M_HTTP}"

    #
    dplatformos_SendToAddress "${origAddr}" "${gameAddr1}" 1000000000 "${MAIN_HTTP}"
    dplatformos_SendToAddress "${origAddr}" "${gameAddr2}" 1000000000 "${MAIN_HTTP}"
    dplatformos_SendToAddress "${origAddr}" "${gameAddr3}" 1000000000 "${MAIN_HTTP}"

    #
    dplatformos_SendToAddress "${gameAddr1}" "${bwExecAddr}" 500000000 "${MAIN_HTTP}"
    dplatformos_SendToAddress "${gameAddr2}" "${bwExecAddr}" 500000000 "${MAIN_HTTP}"
    dplatformos_SendToAddress "${gameAddr3}" "${bwExecAddr}" 500000000 "${MAIN_HTTP}"

    blackwhite_BlackwhiteCreateTx "${gameAddr1}"

    blackwhite_BlackwhitePlayTx "${gameAddr1}" "${white0}" "${white1}" "${black2}"
    blackwhite_BlackwhitePlayTx "${gameAddr2}" "${white0}" "${black1}" "${black2}"
    blackwhite_BlackwhitePlayTx "${gameAddr3}" "${white0}" "${black1}" "${black2}"

    blackwhite_BlackwhiteShowTx "${gameAddr1}" "${sect1}"
    blackwhite_BlackwhiteShowTx "${gameAddr2}" "${sect1}"
    blackwhite_BlackwhiteShowTx "${gameAddr3}" "${sect1}"

    blackwhite_BlackwhiteTimeoutDoneTx "$gID"
    #
    blackwhite_GetBlackwhiteRoundInfo "$gID"
    blackwhite_GetBlackwhiteByStatusAndAddr "${gameAddr1}"
    blackwhite_GetBlackwhiteloopResult "$gID"
}

function main() {
    dplatformos_RpcTestBegin blackwhite
    MAIN_HTTP="$1"
    echo "main_ip=$MAIN_HTTP"

    init
    run_testcases
    dplatformos_RpcTestRst blackwhite "$CASE_ERR"
}

dplatformos_debug_function main "$1"
