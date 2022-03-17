#!/bin/bash

apt-get update
apt-get install -y \
  sudo \
  vim \
  git \
  firmware-iwlwifi \
  ufw \
  tmux \
  cowsay \
  fortune-mod

modprobe -r iwlwifi ; modprobe iwlwifi