#!/bin/bash


echo "Install the postgres client ..."
sudo apt-get update 
export DEBIAN_FRONTEND=noninteractive 
sudo apt-get -y install --no-install-recommends postgresql-client
psql --version

# echo "Install NATS CLI ..."
# LATEST=$(curl -s https://api.github.com/repos/nats-io/natscli/releases/latest | grep browser_download_url | grep linux-amd64.zip | cut -d '"' -f 4) \
#     && wget -O /tmp/natscli.zip $LATEST \
#     && unzip /tmp/natscli.zip -d /tmp \
#     && sudo mv /tmp/nats-*-linux-amd64/nats /usr/local/bin/nats \
#     && rm -r /tmp/natscli.zip /tmp/nats-*-linux-amd64
# nats --version


echo "Setup git ..."
git config --global --add safe.directory /workspace
git config --list



echo "Bootstrap the project ..."
make bootstrap

echo "Create a default nats context to connect to local nats server ..."
nats context save local --server nats://s3cr3t@localhost --description 'Local Host' --select 

echo "Ready!"