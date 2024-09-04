#!/bin/bash


echo "Install the postgres client ..."
sudo apt-get update 
export DEBIAN_FRONTEND=noninteractive 
sudo apt-get -y install --no-install-recommends postgresql-client
psql --version



echo "Setup git ..."
git config --global --add safe.directory /workspace
git config --list

echo "Install sql-migrate command line tool ..."
go install github.com/rubenv/sql-migrate/...@latest
sql-migrate --version

echo "Bootstrap the project ..."
make bootstrap


echo "Ready!"