#!/bin/zsh
go build -o monkeyrun .
if [ "$2" = "" ];then
        ./monkeyrun $1  #running a particular file
    else ./monkeyrun
fi
rm monkeyrun

