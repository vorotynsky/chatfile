#!/bin/sh

MODELS=""

for file in $(ls ./models -p | grep -v /); do 
    MODELS="$MODELS --load-as-model $file=./models/$file"
done


for file in $(find ./lib -name "*.chatfile" -type f); do
    new_file="${file%.*}.ai.go"

    chatfile --seed=42 --temperature=0.0 $MODELS $file > $new_file
done
