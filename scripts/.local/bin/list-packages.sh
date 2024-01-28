#!/usr/bin/env bash

current_folder=${pwd}
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
root_folder=$(cd ${script_dir} && git rev-parse --show-toplevel && cd ${current_folder})

all_packages=($(pacman -Qqe))
base_devel_packages=($(pacman -Sg base-devel | cut -f 2 -d " "))

# Removes duplicates from other array
for i in "${base_devel_packages[@]}"; do
	all_packages=(${all_packages[@]//*$i*/})
done

all_packages=("${all_packages[@]/"$del"/}")

printf "%s\n" "${all_packages[@]}" >${root_folder}/programs/.local/packages.list

cat ${root_folder}/programs/.local/packages.list
