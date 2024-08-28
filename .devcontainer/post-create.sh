#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Fix DNS resolution issues in the KIND cluster
(rc=$(sed 's/^nameserver.*/nameserver 8.8.8.8/' /etc/resolv.conf)
     echo "$rc" | sudo tee /etc/resolv.conf > /dev/null)

git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.1
echo ". $HOME/.asdf/asdf.sh" >> ~/.bashrc
echo ". $HOME/.asdf/completions/asdf.bash" >> ~/.bashrc

source $HOME/.asdf/asdf.sh

cp $SCRIPT_DIR/.tool-versions $HOME/.tool-versions
cat "$HOME/.tool-versions" | awk '{print $1}' | xargs -I {} asdf plugin add {} || true
asdf install
go install github.com/google/ko@v0.16.0
asdf reshim golang

echo 'source <(kubectl completion bash)' >> ~/.bashrc
echo 'alias k=kubectl' >> ~/.bashrc
echo 'complete -o default -F __start_kubectl k' >> ~/.bashrc