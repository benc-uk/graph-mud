#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo -e "\e[33m### \e[32m🌍 Running Neo4j cypher script\e[0m $1"
"$DIR"/tools/cypher-shell < "$1"