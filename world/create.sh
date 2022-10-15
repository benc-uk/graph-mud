#!/bin/bash

echo -e "\e[31m### \e[32m🌍 Initializing Neo4j database\e[0m"
cat world.cypher | ./cypher-shell 