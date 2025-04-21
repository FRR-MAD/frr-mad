#!/bin/bash

if [ -n "$1" ] ; then
  mkdir -p /tmp/containerlab/r101
  chmod 777 /tmp/containerlab/r101
  clab deploy --topo frr01-dev.clab.yml --reconfigure
  chmod +x ./scripts/*.sh
  ./scripts/pc-interfaces.sh
  ./scripts/remove-default-route.sh
else
  clab destroy --topo frr01-dev.clab.yml
fi
