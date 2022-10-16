#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DB_FILE="$DIR/realm.cypher"

echo -e "\e[33m### \e[32m🌍 Initializing Neo4j database from\e[0m $DB_FILE"
"$DIR"/tools/cypher-shell < "$DB_FILE"