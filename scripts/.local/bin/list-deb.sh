#!/usr/bin/env bash

current_folder=${pwd}
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
root_folder=$(cd ${script_dir} && git rev-parse --show-toplevel && cd ${current_folder})

manually_installed_packages=($(apt-mark showmanual))

printf "%s\n" "${manually_installed_packages[@]}" >${root_folder}/programs/.local/deb_packages.list

cat ${root_folder}/programs/.local/deb_packages.list
