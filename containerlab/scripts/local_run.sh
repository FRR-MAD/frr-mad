#!/bin/bash

clab deploy --topo local.frr01.clab.yml --reconfigure
chmod +x ./*.sh
./pc-interfaces.sh
./remove-default-route.sh
