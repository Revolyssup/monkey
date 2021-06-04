go build .
if [ "$2" = "" ];then
        ./monkey $1  #running a particular file
    else ./monkey
fi
rm monkey

