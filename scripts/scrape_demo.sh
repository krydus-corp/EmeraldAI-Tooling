#!/bin/bash

# This script demos the scraping tooling. The demo will run a search for various boat images, clean any corrupt images, convert the images to JPEG, and resize the images. 
# The output of this demo is a folder (~/Desktop/emld-demo/images) containing the scraped images.

PAGES=50
# SEARCH="boat;aerial boat;military boat;industrial boat;military boat aerial;commercial boat;cargo boat;yacht;fishing boat;boat from drone"
SEARCH="turkey;bird;birds;turkey flying;turkey dinner;animal;hipo;cat;dog;horse;wild turkey;turkey in the wild;turkey imaergy;turkey photo"
EMLD_CLI_BIN=$PWD/build/bin/emld-cli
OUTPATH=~/Desktop/emld-demo/images
LOGPATH=~/Desktop/emld-demo/emld.log
SIZE=200
WORKERS=30

STARTING_COUNT=`ls $OUTPATH | wc -l | xargs`

# Split seach terms
IFS=';' # space is set as delimiter
read -ra TERMS <<< "$SEARCH" # str is read into an array as tokens separated by IFS

# Create outpaths
mkdir -p $OUTPATH
touch $LOGPATH

# Kick off search routines
for term in "${TERMS[@]}"; do
    echo "Kicking off search process for term=$term pages=$PAGES";
    ($EMLD_CLI_BIN fetch $term -t images --stream --server "http://127.0.0.1" --pages $PAGES | jq -r '.img_src_b64' | \
    $EMLD_CLI_BIN download --stream -b -w $WORKERS -p $OUTPATH | \
    $EMLD_CLI_BIN image --stream --replace --silent -f "image/jpeg" -x $SIZE -w $WORKERS) &> $LOGPATH &
done

# Trap ctrl-c and call ctrl_c()
trap cleanup INT

function cleanup() {
        echo ""
        printf "Cleaning up."

        # Softkill then hardkill
        pkill emld-cli
        for i in {1..5}; do
            sleep 1
            printf "."
        done
        pkill -9 emld-cli

        # Retrieve images processed
        end_count=`ls $OUTPATH | wc -l | xargs`
        count=$(($end_count-$STARTING_COUNT))

        printf "Processed $count images; exiting...\n"

        exit
}

# Time for workers to spin up
sleep 60

echo 'looping...'
# Wait for ctrl-c or completetion
while true
do
    # modified=`find $OUTPATH -mtime -60s | wc -l`
    # if [ $modified -eq 0 ]; then
    #     echo 'cleanup...'
    #     cleanup
    # fi

    now_count=`ls $OUTPATH | wc -l | xargs`
    count=$(($now_count-$STARTING_COUNT))

    printf "Processed $count images\n"

    sleep 5
done