#!/bin/bash

deployment=$1

if [ -z "$deployment" ] ; then
  chmod +x ./scripts/*.sh
  ./scripts/custom-bridges.sh
  clab deploy --topo frr01.clab.yml --reconfigure
  ./scripts/pc-interfaces.sh
  ./scripts/remove-default-route.sh
else
  if [ $deployment = "local" ] ; then
    chmod +x ./scripts/*.sh
    ./scripts/custom-bridges.sh
    clab deploy --topo local-frr01.clab.yml --reconfigure
    ./scripts/pc-interfaces.sh
    ./scripts/remove-default-route.sh
  fi
fi
