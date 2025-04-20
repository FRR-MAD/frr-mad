#!/bin/sh


docker run --privileged \
  -v $(pwd)/../src:/app/src \
  -v $(pwd)/../local:/app/local \
  -v $(pwd)/../protobufSource:/app/protobufSource \
  -v $(pwd)/../local/frr.conf:/etc/frr/frr.conf \
  -v $(pwd)/../local/daemons:/etc/frr/daemons \
  -v $(pwd)/../local/vtysh.conf:/etc/frr/vtysh.conf \
  frr-854-dev
