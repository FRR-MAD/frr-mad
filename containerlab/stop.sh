#!/bin/bash

deployment=$1

if [ -z "$deployment" ] ; then
  clab destroy --topo frr01.clab.yml
else
  clab destroy --topo local-frr01.clab.yml 
fi
