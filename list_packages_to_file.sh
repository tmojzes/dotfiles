#!/bin/bash
sudo dpkg-query -f '${binary:Package}\n' -W > packages_list.txt
sudo snap list > packages_list_snap.txt
