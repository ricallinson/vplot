# vplot

## Developer Setup

    cd $GOPATH
    git clone git@github.com:ricallinson/vplotter-driver.git ./src/github.com/ricallinson/vplotter-driver
    go get github.com/tarm/goserial

## Install

    go install

## Use

    vplot -l
    vplot <serial-port> <plotter-file>
