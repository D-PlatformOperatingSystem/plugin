#!/usr/bin/env bash

RAW_TX_HASH=""
LAST_BLOCK_HASH=""
LAST_BLOCK_HEIGHT=0
CASE_ERR=""
RETURN_RESP=""

#color
RED='\033[1;31m'
GRE='\033[1;32m'
NOC='\033[0m'

echo_rst() {
    if [ "$2" -eq 0 ]; then
        echo -e "${GRE}$1 ok${NOC}"
    elif [ "$2" -eq 2 ]; then
        echo -e "${GRE}$1 not support${NOC}"
    else
        echo -e "${RED}$1 fail${NOC}"
        echo -e "${RED}$3 ${NOC}"
        CASE_ERR="err"
        echo $CASE_ERR
    fi
}

dplatformos_Http() {
    #  echo "#$4 request: request="$1" MAIN_HTTP="$2" js="$3" FUNCNAME="$4" response="$5""
    local body
    body=$(curl -ksd "$1" "$2")
    RETURN_RESP=$(echo "$body" | jq -r "$5")
    echo "#response: $body" "$RETURN_RESP"
    ok=$(echo "$body" | jq -r "$3")
    [ "$ok" == true ]
    rst=$?
    echo_rst "$4" "$rst" "$body"
}

dplatformos_SignAndSendTxWait() {
    # txHex="$1" priKey="$2" MAIN_HTTP="$3" FUNCNAME="$4"
    req='{"method":"DplatformOS.DecodeRawTransaction","params":[{"txHex":"'"$1"'"}]}'
    dplatformos_Http "$req" "$3" '(.result.txs[0].execer != "") and (.result.txs[0].execer != null)' "$4"
    dplatformos_SignAndSendTx "$1" "$2" "$3"
    dplatformos_BlockWait 1 "$3"
}

dplatformos_BlockWait() {
    local MAIN_HTTP=$2
    local req='"method":"DplatformOS.GetLastHeader","params":[]'

    cur_height=$(curl -ksd "{$req}" "${MAIN_HTTP}" | jq ".result.height")
    expect=$((cur_height + ${1}))

    local count=0
    while true; do
        new_height=$(curl -ksd "{$req}" "${MAIN_HTTP}" | jq ".result.height")
        if [ "${new_height}" -ge "${expect}" ]; then
            break
        fi
        count=$((count + 1))
        sleep 1
    done
    echo "wait new block $count/10 s, cur height=$expect,old=$cur_height"
}

dplatformos_QueryTx() {
    local MAIN_HTTP=$2
    dplatformos_BlockWait 1 "$MAIN_HTTP"
    local txhash="$1"
    local req='"method":"DplatformOS.QueryTransaction","params":[{"hash":"'"$txhash"'"}]'

    local times=10
    while true; do
        ret=$(curl -ksd "{$req}" "${MAIN_HTTP}" | jq -r ".result.tx.hash")
        if [ "${ret}" != "${1}" ]; then
            dplatformos_BlockWait 1 "$MAIN_HTTP"
            times=$((times - 1))
            if [ $times -le 0 ]; then
                echo "====query tx=$1 failed"
                curl -ksd "{$req}" "${MAIN_HTTP}"
                exit 1
            fi
        else
            RAW_TX_HASH=$txhash
            echo "====query tx=$RAW_TX_HASH success"
            break
        fi
    done
}

dplatformos_SendTx() {
    local signedTx=$1
    local MAIN_HTTP=$2

    req='"method":"DplatformOS.SendTransaction","params":[{"token":"DOM","data":"'"$signedTx"'"}]'
    resp=$(curl -ksd "{$req}" "${MAIN_HTTP}")
    err=$(jq '(.error)' <<<"$resp")
    txhash=$(jq -r ".result" <<<"$resp")

    if [ "$err" == null ]; then
        dplatformos_QueryTx "$txhash" "$MAIN_HTTP"
    else
        echo "send tx error:$err"
    fi
}

dplatformos_SendToAddress() {
    local from="$1"
    local to="$2"
    local amount=$3
    local MAIN_HTTP=$4

    local req='"method":"DplatformOS.SendToAddress", "params":[{"from":"'"$from"'","to":"'"$to"'", "amount":'"$amount"', "note":"test\n"}]'
    resp=$(curl -ksd "{$req}" "${MAIN_HTTP}")
    ok=$(jq '(.error|not) and (.result.hash|length==66)' <<<"$resp")

    [ "$ok" == true ]

    hash=$(jq -r ".result.hash" <<<"$resp")
    echo "hash"
    dplatformos_QueryTx "$hash" "$MAIN_HTTP"
}

dplatformos_ImportPrivkey() {
    local pri="$1"
    local acc="$2"
    local label="$3"
    local MAIN_HTTP=$4

    local req='"method":"DplatformOS.ImportPrivkey", "params":[{"privkey":"'"$pri"'", "label":"'"$label"'"}]'
    resp=$(curl -ksd "{$req}" "$MAIN_HTTP")
    #ok=$(jq '(((.error|not) and (.result.label=="'"$label"'") and (.result.acc.addr == "'"$acc"'")) or (.error=="ErrPrivkeyExist"))' <<<"$resp")
    ok=$(jq '(((.error|not) and (.result.label=="'"$label"'") and (.result.acc.addr == "'"$acc"'")) or (.error=="ErrPrivkeyExist") or (.error=="ErrLabelHasUsed"))' <<<"$resp")

    [ "$ok" == true ]
}

dplatformos_SignAndSendTx() {
    local txHex="$1"
    local priKey="$2"
    local MAIN_HTTP=$3
    local expire="120s"
    if [ -n "$4" ]; then
        expire=$4
    fi

    local req='"method":"DplatformOS.SignRawTx","params":[{"privkey":"'"$priKey"'","txHex":"'"$txHex"'","expire":"'"$expire"'"}]'
    signedTx=$(curl -ksd "{$req}" "${MAIN_HTTP}" | jq -r ".result")

    if [ "$signedTx" != null ]; then
        dplatformos_SendTx "$signedTx" "${MAIN_HTTP}"
    else
        echo "signedTx null error"
    fi
}

dplatformos_QueryBalance() {
    local addr=$1
    local MAIN_HTTP=$2
    req='"method":"DplatformOS.GetAllExecBalance","params":[{"addr":"'"${addr}"'"}]'
    #echo "#request: $req"
    resp=$(curl -ksd "{$req}" "${MAIN_HTTP}")
    echo "#response: $resp"
    ok=$(jq '(.error|not) and (.result != "")' <<<"$resp")
    [ "$ok" == true ]

    echo "$resp" | jq -r ".result"
}

dplatformos_QueryExecBalance() {
    local addr=$1
    local exec=$2
    local MAIN_HTTP=$3

    req='{"method":"DplatformOS.GetBalance", "params":[{"addresses" : ["'"${addr}"'"], "execer" : "'"${exec}"'"}]}'
    resp=$(curl -ksd "$req" "${MAIN_HTTP}")
    echo "#response: $resp"
    ok=$(jq '(.error|not) and (.result[0] | [has("balance", "frozen"), true] | unique | length == 1)' <<<"$resp")
    [ "$ok" == true ]
}

dplatformos_GetAccounts() {
    local MAIN_HTTP=$1
    resp=$(curl -ksd '{"jsonrpc":"2.0","id":2,"method":"DplatformOS.GetAccounts","params":[{}]}' -H 'content-type:text/plain;' "${MAIN_HTTP}")
    echo "$resp"
}

dplatformos_LastBlockhash() {
    local MAIN_HTTP=$1
    result=$(curl -ksd '{"method":"DplatformOS.GetLastHeader","params":[{}]}' -H 'content-type:text/plain;' "${MAIN_HTTP}" | jq -r ".result.hash")
    LAST_BLOCK_HASH=$result
    echo -e "######\\n  last blockhash is $LAST_BLOCK_HASH  \\n######"
}

dplatformos_LastBlockHeight() {
    local MAIN_HTTP=$1
    result=$(curl -ksd '{"method":"DplatformOS.GetLastHeader","params":[{}]}' -H 'content-type:text/plain;' "${MAIN_HTTP}" | jq -r ".result.height")
    LAST_BLOCK_HEIGHT=$result
    echo -e "######\\n  last blockheight is $LAST_BLOCK_HEIGHT \\n######"
}

dplatformos_applyCoins() {
    echo "dplatformos_getMainChainCoins"
    if [ "$#" -lt 3 ]; then
        echo "dplatformos_getMainCoins wrong params"
        exit 1
    fi
    local targetAddr=$1
    local count=$2
    local ip=$3
    if [ "$count" -gt 15000000000 ]; then
        echo "dplatformos_getMainCoins wrong coins count,should less than 150 00000000"
        exit 1
    fi

    local poolAddr="1PcGKYYoLn1PLLJJodc1UpgWGeFAQasAkx"
    dplatformos_SendToAddress "${poolAddr}" "${targetAddr}" "$count" "${ip}"

}

dplatformos_RpcTestBegin() {
    echo -e "${GRE}====== $1 Rpc Test Begin ===========${NOC}"
}

dplatformos_RpcTestRst() {
    if [ -n "$2" ]; then
        echo -e "${RED}====== $1 Rpc Test Fail ===========${NOC}"
        exit 1
    else
        echo -e "${GRE}====== $1 Rpc Test Pass ===========${NOC}"
    fi
}

dplatformos_debug_function() {
    set -x
    eval "$@"
    set +x
}
