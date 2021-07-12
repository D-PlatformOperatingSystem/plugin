#!/usr/bin/env bash
# shellcheck disable=SC2128
# shellcheck source=/dev/null

MAIN_HTTP=""
oracle_addPublisher_unsignedTx="0a066d616e61676512410a3f0a146f7261636c652d7075626c6973682d6576656e741222313271796f6361794e46374c7636433971573461767873324537553431664b5366761a0361646420a08d0630e6b685d696ee9394163a223151344e687572654a784b4e4266373164323642394a336642516f5163666d657a32"
oracle_addPublisher_unsignedTx_para="0a12757365722e702e706172612e6d616e61676512410a3f0a146f7261636c652d7075626c6973682d6576656e741222313271796f6361794e46374c7636433971573461767873324537553431664b5366761a0361646420a08d0630a186de8894c9aa864d3a22314469484633317577783977356a6a733571514269474a6b4e686e71656564763157"
oracle_publisher_key="4257D8692EF7FE13C68B65D6A52F03933DB2FA5CE8FAF210B5B8B80C721CED01"
eventId=""
txhash=""

source ../dapp-test-common.sh

oracle_AddPublisher() {
    echo "=============== # Add publisher ==============="
    ispara=$(echo '"'"${MAIN_HTTP}"'"' | jq '.|contains("8901")')
    echo "ispara=$ispara"
    if [ "$ispara" == true ]; then
        dplatformos_SignAndSendTx "${oracle_addPublisher_unsignedTx_para}" "${oracle_publisher_key}" "${MAIN_HTTP}"
    else
        dplatformos_SignAndSendTx "${oracle_addPublisher_unsignedTx}" "${oracle_publisher_key}" "${MAIN_HTTP}"
    fi
}

oracle_publish_transaction() {
    req='{"method":"DplatformOS.CreateTransaction","params":[{"execer":"oracle","actionName":"EventPublish","payload":{"type":"football", "subType":"Premier League","time":1747814996,"content":"test","introduction":"test"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatformos_SignAndSendTx "$RETURN_RESP" "${oracle_publisher_key}" "${MAIN_HTTP}"
    eventId="${txhash}"
    echo "eventId $eventId"
}

oracle_prePublishResult_transaction() {
    event_id=$1
    req='{"method":"DplatformOS.CreateTransaction","params":[{"execer":"oracle","actionName":"ResultPrePublish","payload":{"eventID":"'"$event_id"'", "source":"sina sport","result":"0:1"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatformos_SignAndSendTx "$RETURN_RESP" "${oracle_publisher_key}" "${MAIN_HTTP}"
}

oracle_eventAbort_transaction() {
    event_id=$1
    req='{"method":"DplatformOS.CreateTransaction","params":[{"execer":"oracle","actionName":"EventAbort","payload":{"eventID":"'"$event_id"'"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatformos_SignAndSendTx "$RETURN_RESP" "${oracle_publisher_key}" "${MAIN_HTTP}"
}

oracle_resultAbort_transaction() {
    event_id=$1
    req='{"method":"DplatformOS.CreateTransaction","params":[{"execer":"oracle","actionName":"ResultAbort","payload":{"eventID":"'"$event_id"'"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatformos_SignAndSendTx "$RETURN_RESP" "${oracle_publisher_key}" "${MAIN_HTTP}"
}

oracle_publishResult_transaction() {
    event_id=$1
    req='{"method":"DplatformOS.CreateTransaction","params":[{"execer":"oracle","actionName":"ResultPublish","payload":{"eventID":"'"$event_id"'", "source":"sina sport","result":"1:1"}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result != null)' "$FUNCNAME" ".result"
    dplatformos_SignAndSendTx "$RETURN_RESP" "${oracle_publisher_key}" "${MAIN_HTTP}"
}

oracle_QueryOraclesByID() {
    event_id=$1
    req='{"method":"DplatformOS.Query", "params":[{"execer":"oracle","funcName":"QueryOraclesByIDs","payload":{"eventID":["'"$event_id"'"]}}]}'
    dplatformos_Http "$req" ${MAIN_HTTP} '(.error|not) and (.result.status[0] | [has("eventID", "status", "type", "subType", "source"),true] | unique | length == 1)' "$FUNCNAME"
}

function run_test() {
    #
    oracle_AddPublisher
    #
    oracle_publish_transaction
    #
    oracle_prePublishResult_transaction "$eventId"
    #
    oracle_publishResult_transaction "$eventId"
    #   ID
    dplatformos_BlockWait 2 "${MAIN_HTTP}"
    oracle_QueryOraclesByID "$eventId"

    #
    oracle_publish_transaction
    #
    oracle_eventAbort_transaction "$eventId"
    #   ID
    dplatformos_BlockWait 2 "${MAIN_HTTP}"
    oracle_QueryOraclesByID "$eventId"

    #
    oracle_publish_transaction
    #
    oracle_prePublishResult_transaction "$eventId"
    #
    oracle_resultAbort_transaction "$eventId"
    #   ID
    dplatformos_BlockWait 2 "${MAIN_HTTP}"
    oracle_QueryOraclesByID "$eventId"

}

function main() {
    dplatformos_RpcTestBegin oracle
    MAIN_HTTP="$1"
    echo "main_ip=$MAIN_HTTP"

    run_test
    dplatformos_RpcTestRst oracle "$CASE_ERR"
}

dplatformos_debug_function main "$1"
