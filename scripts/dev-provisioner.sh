# this script can help someone setting up wapikit from scratch from source on a ubuntu machine

#!/bin/sh

set -e

sudo apt-get -y update
sudo apt-get install -y make
sudo apt install -y wget
sudo apt-get -y update

# install golang
wget https://go.dev/dl/go1.22.3.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# install n and node
curl -L https://bit.ly/n-install | bash
n i lts

git clone https://github.com/wapikit/wapikit

# this will ensure that all the other dependencies are installed and the binary is built with static files
make dist
