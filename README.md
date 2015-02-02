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
    vplot -m ./fixtures/vplotter.mock ./fixtures/square.vplot

## Setup a Raspberry Pi

Kickoff with an update;

    sudo apt-get update

Install Screen;

    sudo apt-get install screen

Install Go;

    wget http://dave.cheney.net/paste/go1.4.linux-arm~multiarch-armv6-1.tar.gz
    sudo tar -C /usr/bin -xzf go1.4.linux-arm~multiarch-armv6-1.tar.gz
    export PATH=$PATH:/usr/local/go/bin

Install Git;

    sudo apt-get install git-core

Install Go Workspace Manager;

    git clone https://github.com/ricallinson/gwm.git ~/.gwm

Add the following line to the end of `.bashrc`;

    source ~/.gwm/gwm.sh

Then active your changes;

    source ~/.bashrc

Create the Go root;

    mkdir ~/Go
    cd ~/Go
    gwm use .

Get the vplot source and build it;

    go get github.com/ricallinson/vplot

Create a directory to upload the plot files into;

    mkdir ~/plots

## Plot a file

Send a plot to the Raspberry Pi;

    scp <file>.plot pi@192.168.0.189:~/plots

Now connect to the Raspberry Pi;

    ssh 192.168.0.189 -l pi

To keep the plot going after the `ssh` session use `screen`. Once in a screen you can use `ctrl + A + D` to exit and keep it running. Use `screen -l` to list the screens. Use `screen -r <screen>` to resume screen.

    screen
    ~/Go/bin/vplot /dev/ttyACM0 ~/plots/plot.vplot ~/Go/src/github.com/ricallinson/vplot/config/wall.cfg

Once the plot has started you can exit the screen any time by pressing `ctrl + A + D`.
