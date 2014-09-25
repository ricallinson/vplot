# vplot

## Developer Setup

    cd $GOPATH
    git clone git@github.com:ricallinson/vplot.git ./src/github.com/ricallinson/vplot
    go get github.com/tarm/goserial

## Install

    cd ./src/github.com/ricallinson/vplot
    go install

## Use

    vplot -l
    vplot /dev/tty.1234 ./fixtures/square.vplot
