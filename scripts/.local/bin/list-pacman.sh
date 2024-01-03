#!/bin/bash

all_packages=($(pacman -Qqe))
base_devel_packages=($(pacman -Sg base-devel | cut -f 2 -d " "))

# Removes duplicates from other array
for i in "${base_devel_packages[@]}"; do
	all_packages=(${all_packages[@]//*$i*/})
done

all_packages=("${all_packages[@]/"$del"/}")

printf "%s\n" "${all_packages[@]}" >~/.pacman.list

cat ~/.pacman.list
