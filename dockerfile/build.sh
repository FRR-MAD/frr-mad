#!/bin/sh
cd ..
docker build -t frr-854-dev -f dockerfile/frr-dev.dockerfile .
