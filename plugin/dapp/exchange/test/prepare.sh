#!/usr/bin/env bash

BUILD=$(cd "$(dirname "$0")" && cd ../../../../build && pwd)
echo "$BUILD"

cd "$BUILD" || return

seed=$(./dplatformos-cli seed generate -l 0)
echo "$seed"

./dplatformos-cli seed save -p dpos123456 -s "$seed"
sleep 1
./dplatformos-cli wallet unlock -p dpos123456
sleep 1
./dplatformos-cli account list

## account
#genesis -- 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli account import_key -k 3990969DF92A5914F7B71EEB9A4E58D6E255F32BF042FEA5318FC8B3D50EE6E8 -l genesis

#A -- 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
./dplatformos-cli account import_key -k 0x6da92a632ab7deb67d38c0f6560bcfed28167998f6496db64c258d5e8393a81b -l A

#B -- 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
./dplatformos-cli account import_key -k 0x19c069234f9d3e61135fefbeb7791b149cdf6af536f26bebb310d4cd22c3fee4 -l B

#C -- 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
./dplatformos-cli account import_key -k 0x7a80a1f75d7360c6123c32a78ecf978c1ac55636f87892df38d8b85a9aeff115 -l C

#D -- 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs
./dplatformos-cli account import_key -k 0xcacb1f5d51700aea07fca2246ab43b0917d70405c65edea9b5063d72eb5c6b71 -l D

## config token
./dplatformos-cli send config config_tx -c token-finisher -o add -v 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
sleep 1
./dplatformos-cli config query -k token-finisher

./dplatformos-cli send config config_tx -c token-blacklist -o add -v DOM -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
sleep 1
./dplatformos-cli config query -k token-blacklist

## 10
./dplatformos-cli send token precreate -f 0.001 -i "test ccny" -n "ccny" -a 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs -p 0 -s CCNY -t 1000000000 -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
sleep 1
./dplatformos-cli token precreated
./dplatformos-cli send token finish -s CCNY -f 0.001 -a 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
sleep 1
./dplatformos-cli token created
./dplatformos-cli token balance -a 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs -s CCNY -e token

## transfer dpos
./dplatformos-cli send coins transfer -a 10000 -t 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4 -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli send coins transfer -a 10000 -t 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli send coins transfer -a 10000 -t 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli send coins transfer -a 10000 -t 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli account list

## send dpos to execer，   10000 dpos
./dplatformos-cli send coins send_exec -e exchange -a 10000 -k 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
./dplatformos-cli send coins send_exec -e exchange -a 10000 -k 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
./dplatformos-cli send coins send_exec -e exchange -a 10000 -k 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
./dplatformos-cli send coins send_exec -e exchange -a 10000 -k 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs
echo "account balance in execer"
./dplatformos-cli account balance -e exchange -a 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
./dplatformos-cli account balance -e exchange -a 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
./dplatformos-cli account balance -e exchange -a 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
./dplatformos-cli account balance -e exchange -a 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs

##transfer token，  2  ccny
./dplatformos-cli send token transfer -a 200000000 -s CCNY -t 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4 -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli send token transfer -a 200000000 -s CCNY -t 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli send token transfer -a 200000000 -s CCNY -t 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs
./dplatformos-cli send token transfer -a 200000000 -s CCNY -t 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs -k 1CbEVT9RnM5oZhWMj4fxUrJX94VtRotzvs

## send token to excer
./dplatformos-cli send token send_exec -a 200000000 -e exchange -s CCNY -k 1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4
./dplatformos-cli send token send_exec -a 200000000 -e exchange -s CCNY -k 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR
./dplatformos-cli send token send_exec -a 200000000 -e exchange -s CCNY -k 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k
./dplatformos-cli send token send_exec -a 200000000 -e exchange -s CCNY -k 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs
echo "token balance in execer"
./dplatformos-cli token balance -e exchange -s CCNY -a "1KSBd17H7ZK8iT37aJztFB22XGwsPTdwE4 1JRNjdEqp4LJ5fqycUBm9ayCKSeeskgMKR 1NLHPEcbTWWxxU3dGUZBhayjrCHD3psX7k 1MCftFynyvG2F4ED5mdHYgziDxx6vDrScs"
