#!/bin/bash
sed -i 's/frr854-testing/registry.gitlab.ost.ch:45023\/ins-stud\/sa-ba\/ba-fs25-frr-monitoring-analytics\/frr-deployment\/frr854:latest/' frr01.clab.yml
sed -i 's/ffr854-monitoring/registry.gitlab.ost.ch:45023\/ins-stud\/sa-ba\/ba-fs25-frr-monitoring-analytics\/frr-deployment\/frr854-moni:latest/' frr01.clab.yml
