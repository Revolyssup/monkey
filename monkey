#!bin/sh
go build .
if [ "$2" = "" ];then
        ./monkeyrun $1  #running a particular file
    else ./monkeyrun
fi
rm monkeyrun

