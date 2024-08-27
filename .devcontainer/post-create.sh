#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# Fix DNS resolution issues
(rc=$(sed 's/^nameserver.*/nameserver 8.8.8.8/' /etc/resolv.conf)
     echo "$rc" | sudo tee /etc/resolv.conf > /dev/null)

git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.1
echo ". $HOME/.asdf/asdf.sh" >> ~/.bashrc
echo ". $HOME/.asdf/completions/asdf.bash" >> ~/.bashrc

source $HOME/.asdf/asdf.sh

cp "$SCRIPT_DIR/.tool-versions" "$HOME/.tool-versions"
cat "$HOME/.tool-versions" | awk '{print $1}' | xargs -I {} asdf plugin add {} || true
asdf install