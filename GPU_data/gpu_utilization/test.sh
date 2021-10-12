#!/bin/sh

ls *.tar.gz > list.txt

for TAR in `cat list.txt`
do
        tar zxf $TAR
done

rm -rf list.txt