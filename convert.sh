#!/bin/bash

for file in `ls *.png`
do
    newfile=`echo $file | sed 's/png/jpg/'`
    convert $file -quality 90 $newfile
    rm $file
done
