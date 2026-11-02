#!/bin/bash

OUTPUT_DIR="./out"
POSSIBLE_FIELDS="industrial mechanical aerospace civil chemical business electrical software data petroleum"


function exitHelp() {
    echo -e "usage:\n\t./collect.sh [sqlite | json] [aggregate | all | $POSSIBLE_FIELDS]"
    echo -e "you can include multiple by spacing"
    exit 1
}

if [ "$#" -lt 2 ]; then
    exitHelp
fi

FEND=
if [ "$1" = "sqlite" ]; then
    FEND="db"
elif [ "$1" = "json" ]; then
    FEND="json"
else
    echo "unknow choice: $1"
    exitHelp
fi

set -euo pipefail

mkdir -p "$OUTPUT_DIR"
find "${OUTPUT_DIR:?}" -mindepth 1 -maxdepth 1 -type f -delete

FIELDS=$2
if [ "$2" = "aggregate" ]; then
    FILE="$OUTPUT_DIR/aggregate.$FEND"
    echo "========== $FILE ==========" # cool
    OUTPUT_FILE="$FILE" COUNT=-1 TYPE="$1" make
    exit 0
elif [ "$2" = "all" ]; then
    FIELDS="$POSSIBLE_FIELDS"
fi



echo "Running for $1: $2"
read -ra fds <<< "$FIELDS"

for i in "${fds[@]}"; do
    FILE="$OUTPUT_DIR/$i.$FEND"
    echo "========== $FILE ==========" # cool
    OUTPUT_FILE="$FILE" MAJOR="$i" COUNT=-1 TYPE="$1" make
done