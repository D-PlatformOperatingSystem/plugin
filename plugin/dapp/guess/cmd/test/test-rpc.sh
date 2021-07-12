#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null
set -e
set -o pipefail

MAIN_HTTP=""
source ../dapp-test-common.sh

MAIN_HTTP=""
guess_admin_addr=12oupcayRT7LvaC4qW4avxsTE7U41cKSio
guess_user1_addr=1NrfEBfdFJUUqgbw5ZbHXhdew6NNQumYhM
guess_user2_addr=17tRkBrccmFiVcLPXgEceRxDzJ2WaDZumN
guess_addr=""
guess_exec=""

eventId=""
txhash=""

guess_game_start() {
    tx=$(curl -ksd '{"method":"DplatformOS.CreateTransaction","params":[{"execer":"guess","actionName":"Start", "payload":{"topic":"WorldCup Final","options":"A:France;B:Claodia","category":"football","maxBetsOneTime":10000000000,"maxBetsNumber":100000000000,"devFeeFactor":5,"devFeeAddr":"1D6RFZNp2rh6QdbcZ1d7RWuBUz61We6SD7","platFeeFactor":5,"platFeeAddr":"1PHtChNt3UcfssR7v7trKSk3WJtAWjKjjX"}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatformos_SignAndSendTxWait "$tx" "4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01" ${MAIN_HTTP} "$FUNCNAME"
    eventId="${txhash}"
}

guess_game_bet() {
    local priv=$1
    local opt=$2
    tx=$(curl -ksd '{"method":"DplatformOS.CreateTransaction","params":[{"execer":"guess","actionName":"Bet", "payload":{"gameID":"'"${eventId}"'","option":"'"${opt}"'", "betsNum":500000000}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatformos_SignAndSendTxWait "$tx" "${priv}" ${MAIN_HTTP} "$FUNCNAME"
}

guess_game_stop() {
    tx=$(curl -ksd '{"method":"DplatformOS.CreateTransaction","params":[{"execer":"guess","actionName":"StopBet", "payload":{"gameID":"'"${eventId}"'"}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatformos_SignAndSendTxWait "$tx" "4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01" ${MAIN_HTTP} "$FUNCNAME"
}

guess_game_publish() {
    tx=$(curl -ksd '{"method":"DplatformOS.CreateTransaction","params":[{"execer":"guess","actionName":"Publish", "payload":{"gameID":"'"${eventId}"'","result":"A"}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatformos_SignAndSendTxWait "$tx" "4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01" ${MAIN_HTTP} "$FUNCNAME"
}

guess_game_abort() {
    tx=$(curl -ksd '{"method":"DplatformOS.CreateTransaction","params":[{"execer":"guess","actionName":"Abort", "payload":{"gameID":"'"${eventId}"'"}}]}' ${MAIN_HTTP} | jq -r ".result")
    dplatformos_SignAndSendTxWait "$tx" "4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01" ${MAIN_HTTP} "$FUNCNAME"
}

guess_QueryGameByID() {
    local event_id=$1
    local status=$2
    local req='{"method":"DplatformOS.Query", "params":[{"execer":"guess","funcName":"QueryGameByID","payload":{"gameID":"'"$event_id"'"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.result|has("game")) and (.result.game.status == '"$status"')' "$FUNCNAME"
}

init() {
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ipara=$ispara"
    if [ "$ispara" == true ]; then
        guess_addr=$(curl -ksd '{"method":"DplatformOS.ConvertExectoAddr","params":[{"execname":"user.p.para.guess"}]}' ${MAIN_HTTP} | jq -r ".result")
        guess_exec="user.p.para.guess"
    else
        guess_addr=$(curl -ksd '{"method":"DplatformOS.ConvertExectoAddr","params":[{"execname":"guess"}]}' ${MAIN_HTTP} | jq -r ".result")
        guess_exec="guess"
    fi
    echo "guess_addr=$guess_addr"

    local main_ip=${MAIN_HTTP//8901/28803}
    dplatformos_ImportPrivkey "0xc889d2958843fc96d4bd3f578173137d37230e580d65e9074545c61e7e9c1932" "1NrfEBfdFJUUqgbw5ZbHXhdew6NNQumYhM" "guess11" "${main_ip}"
    dplatformos_ImportPrivkey "0xf10c79470dc74c229c4ee73b05d14c58322b771a6c749d27824f6a59bb6c2d73" "17tRkBrccmFiVcLPXgEceRxDzJ2WaDZumN" "guess22" "$main_ip"

    local guess1="1NrfEBfdFJUUqgbw5ZbHXhdew6NNQumYhM"
    local guess2="17tRkBrccmFiVcLPXgEceRxDzJ2WaDZumN"

    if [ "$ispara" == false ]; then
        dplatformos_applyCoins "$guess1" 12000000000 "${main_ip}"
        dplatformos_QueryBalance "${guess1}" "$main_ip"

        dplatformos_applyCoins "$guess2" 12000000000 "${main_ip}"
        dplatformos_QueryBalance "${guess2}" "$main_ip"
    else
        dplatformos_applyCoins "$guess1" 1000000000 "${main_ip}"
        dplatformos_QueryBalance "${guess1}" "$main_ip"

        dplatformos_applyCoins "$guess2" 1000000000 "${main_ip}"
        dplatformos_QueryBalance "${guess2}" "$main_ip"
        local para_ip="${MAIN_HTTP}"
        dplatformos_ImportPrivkey "0xc889d2958843fc96d4bd3f578173137d37230e580d65e9074545c61e7e9c1932" "1NrfEBfdFJUUqgbw5ZbHXhdew6NNQumYhM" "guess11" "$para_ip"
        dplatformos_ImportPrivkey "0xf10c79470dc74c229c4ee73b05d14c58322b771a6c749d27824f6a59bb6c2d73" "17tRkBrccmFiVcLPXgEceRxDzJ2WaDZumN" "guess22" "$para_ip"

        dplatformos_applyCoins "$guess1" 12000000000 "${para_ip}"
        dplatformos_QueryBalance "${guess1}" "$para_ip"
        dplatformos_applyCoins "$guess2" 12000000000 "${para_ip}"
        dplatformos_QueryBalance "${guess2}" "$para_ip"
    fi

    dplatformos_SendToAddress "$guess1" "$guess_addr" 10000000000 ${MAIN_HTTP}
    dplatformos_QueryExecBalance "${guess1}" "guess" "$MAIN_HTTP"
    dplatformos_SendToAddress "$guess2" "$guess_addr" 10000000000 ${MAIN_HTTP}
    dplatformos_QueryExecBalance "${guess2}" "guess" "$MAIN_HTTP"

    dplatformos_BlockWait 1 "${MAIN_HTTP}"
}

function run_test() {
    #
    dplatformos_ImportPrivkey "0xc889d2958843fc96d4bd3f578173137d37230e580d65e9074545c61e7e9c1932" "1NrfEBfdFJUUqgbw5ZbHXhdew6NNQumYhM" "user1" "$MAIN_HTTP"
    dplatformos_ImportPrivkey "0xf10c79470dc74c229c4ee73b05d14c58322b771a6c749d27824f6a59bb6c2d73" "17tRkBrccmFiVcLPXgEceRxDzJ2WaDZumN" "user2" "$MAIN_HTTP"
    dplatformos_ImportPrivkey "4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01" "12oupcayRT7LvaC4qW4avxsTE7U41cKSio" "admin" "$MAIN_HTTP"

    dplatformos_QueryBalance "${guess_admin_addr}" "$MAIN_HTTP"
    dplatformos_QueryBalance "${guess_user1_addr}" "$MAIN_HTTP"
    dplatformos_QueryBalance "${guess_user2_addr}" "$MAIN_HTTP"
    dplatformos_QueryExecBalance "${guess_user1_addr}" "${guess_exec}" "$MAIN_HTTP"
    dplatformos_QueryExecBalance "${guess_user2_addr}" "${guess_exec}" "$MAIN_HTTP"

    #  1：start -> bet -> bet -> stop -> publish
    #
    guess_game_start

    #
    guess_QueryGameByID "$eventId" 11

    #  1
    guess_game_bet "0xc889d2958843fc96d4bd3f578173137d37230e580d65e9074545c61e7e9c1932" "A"

    #
    guess_QueryGameByID "$eventId" 12

    #  2
    guess_game_bet "0xf10c79470dc74c229c4ee73b05d14c58322b771a6c749d27824f6a59bb6c2d73" "B"

    #
    guess_QueryGameByID "$eventId" 12

    #
    guess_game_stop

    #
    guess_QueryGameByID "$eventId" 13

    #
    guess_game_publish

    #
    guess_QueryGameByID "$eventId" 15

    #
    dplatformos_QueryExecBalance "${guess_user1_addr}" "${guess_exec}" "$MAIN_HTTP"
    dplatformos_QueryExecBalance "${guess_user2_addr}" "${guess_exec}" "$MAIN_HTTP"

    #  2：start->stop->abort
    guess_game_start

    #
    guess_QueryGameByID "$eventId" 11

    #
    guess_game_stop

    #
    guess_QueryGameByID "$eventId" 13

    #
    guess_game_abort

    #
    guess_QueryGameByID "$eventId" 14

    #  3：start->abort
    guess_game_start

    #
    guess_QueryGameByID "$eventId" 11

    #
    guess_game_abort

    #
    guess_QueryGameByID "$eventId" 14

    #  4：start->bet->abort

    #
    guess_game_start

    #
    guess_QueryGameByID "$eventId" 11

    #  1
    guess_game_bet "0xc889d2958843fc96d4bd3f578173137d37230e580d65e9074545c61e7e9c1932" "A"

    #
    guess_QueryGameByID "$eventId" 12

    #  2
    guess_game_bet "0xf10c79470dc74c229c4ee73b05d14c58322b771a6c749d27824f6a59bb6c2d73" "B"

    #
    guess_QueryGameByID "$eventId" 12

    #
    guess_game_abort

    #
    guess_QueryGameByID "$eventId" 14

    #  5：start->bet->stop->abort
    #
    guess_game_start

    #
    guess_QueryGameByID "$eventId" 11

    #  1
    guess_game_bet "0xc889d2958843fc96d4bd3f578173137d37230e580d65e9074545c61e7e9c1932" "A"

    #
    guess_QueryGameByID "$eventId" 12

    #  2
    guess_game_bet "0xf10c79470dc74c229c4ee73b05d14c58322b771a6c749d27824f6a59bb6c2d73" "B"

    #
    guess_QueryGameByID "$eventId" 12

    #
    guess_game_stop

    #
    guess_QueryGameByID "$eventId" 13

    #
    guess_game_abort

    #
    guess_QueryGameByID "$eventId" 14

    #
    dplatformos_QueryExecBalance "${guess_user1_addr}" "${guess_exec}" "$MAIN_HTTP"
    dplatformos_QueryExecBalance "${guess_user2_addr}" "${guess_exec}" "$MAIN_HTTP"
}

function main() {
    dplatformos_RpcTestBegin guess
    MAIN_HTTP="$1"
    echo "main_ip=$MAIN_HTTP"

    init
    run_test
    dplatformos_RpcTestRst guess "$CASE_ERR"
}

dplatformos_debug_function main "$1"
