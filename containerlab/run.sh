#!/bin/bash

deployment=$1

if [ -z "$deployment" ] ; then
  clab deploy --topo frr01.clab.yml --reconfigure
  chmod +x ./scripts/*.sh
  ./scripts/pc-interfaces.sh
  ./scripts/remove-default-route.sh
else
  if [ $deployment = "local" ] ; then
    clab deploy --topo local-frr01.clab.yml --reconfigure
    chmod +x ./scripts/*.sh
    ./scripts/pc-interfaces.sh
    ./scripts/remove-default-route.sh
  fi
fi
