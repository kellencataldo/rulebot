#!/bin/bash

for file in `ls *.png`
do

    swap=`echo $file | sed 's/png/jpg/'`

#    basename="${file%.*}"
#    swap="$basename-tmp.jpg"

    echo $file -> $swap

    convert $file -quality 90 $swap
    rm $file
  #  mv $swap $file
done
