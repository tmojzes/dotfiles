#!/usr/bin/env bash

current_folder=${pwd}
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
root_folder=$(cd ${script_dir} && git rev-parse --show-toplevel && cd ${current_folder})

yay -S --needed - <${root_folder}/programs/.local/packages.list
