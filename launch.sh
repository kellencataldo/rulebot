go build cmd/main.go
retVal=$?
if [ $retVal -ne 0 ]; then
    exit
fi

exit
