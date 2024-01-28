#!/usr/bin/env bash

current_folder=${pwd}
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
root_folder=$(cd ${script_dir} && git rev-parse --show-toplevel && cd ${current_folder})

sudo pacman -S --needed git base-devel &&
	git clone https://aur.archlinux.org/yay-bin.git ${root_folder}/yay-bin &&
	cd ${root_folder}/yay-bin &&
	makepkg -si &&
	cd ${root_folder} &&
	rm -rf yay-bin &&
	cd ${current_folder}
