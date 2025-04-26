#!/bin/bash
# ================== #
#       AS65001      #
# ================== #
# setup test PC in Area0 - PC101
docker exec -d clab-frr01-pc101 ip link set eth1 up
docker exec -d clab-frr01-pc101 ip addr add 10.0.0.100/23 dev eth1
# docker exec -d clab-frr01-pc01 ip route add 10.0.0.0/23 via 10.0.0.1 eth1
docker exec -d clab-frr01-pc101 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc101 ip route add 1.1.1.1/24 via 10.0.0.1 eth1
docker exec -d clab-frr01-pc101 ip route add default via 10.0.0.1 dev eth1

# setup test PC in Area0 - PC102
docker exec -d clab-frr01-pc102 ip link set eth1 up
docker exec -d clab-frr01-pc102 ip addr add 10.0.2.100/24 dev eth1
docker exec -d clab-frr01-pc102 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc102 ip route add default via 10.0.2.1 dev eth1

# setup test PC in Area1 - PC111
docker exec -d clab-frr01-pc111 ip link set eth1 up
docker exec -d clab-frr01-pc111 ip addr add 10.1.0.100/24 dev eth1
docker exec -d clab-frr01-pc111 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc111 ip route add default via 10.1.0.11 dev eth1

# setup test PC in Area1 - PC112
docker exec -d clab-frr01-pc112 ip link set eth1 up
docker exec -d clab-frr01-pc112 ip addr add 10.1.1.100/24 dev eth1
docker exec -d clab-frr01-pc112 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc112 ip route add default via 10.1.1.11 dev eth1

# setup test PC in Area2 -  PC121
docker exec -d clab-frr01-pc121 ip link set eth1 up
docker exec -d clab-frr01-pc121 ip addr add 10.2.0.100/24 dev eth1
docker exec -d clab-frr01-pc121 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc121 ip route add 21.21.21.21/24 via 10.2.0.21 eth1
docker exec -d clab-frr01-pc121 ip route add default via 10.2.0.21 dev eth1

# setup test PC in Area3 -  PC131
docker exec -d clab-frr01-pc131 ip link set eth1 up
docker exec -d clab-frr01-pc131 ip addr add 10.3.0.100/24 dev eth1
docker exec -d clab-frr01-pc131 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc131 ip route add 31.31.31.31/24 via 10.3.0.31 eth1
docker exec -d clab-frr01-pc131 ip route add default via 10.3.0.31 dev eth1

# setup test PC outside OSPF -  PC191
docker exec -d clab-frr01-pc191 ip link set eth1 up
docker exec -d clab-frr01-pc191 ip addr add 192.168.1.100/24 dev eth1
docker exec -d clab-frr01-pc191 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc191 ip route add default via 192.168.1.91 dev eth1

# setup test PC in Customer4 network -  CPC104
docker exec -d clab-frr01-cpc104 ip link set eth1 up
docker exec -d clab-frr01-cpc104 ip addr add 192.168.4.100/24 dev eth1

# setup test PC in Customer5 network -  CPC105
docker exec -d clab-frr01-cpc105 ip link set eth1 up
docker exec -d clab-frr01-cpc105 ip addr add 192.168.5.100/24 dev eth1
docker exec -d clab-frr01-cpc105 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-cpc105 ip route add default via 192.168.5.5 dev eth1

# setup test PC in Customer6 network -  CPC106
docker exec -d clab-frr01-cpc106 ip link set eth1 up
docker exec -d clab-frr01-cpc106 ip addr add 192.168.6.100/24 dev eth1

# setup test PC in Customer7 network -  CPC107
docker exec -d clab-frr01-cpc107 ip link set eth1 up
docker exec -d clab-frr01-cpc107 ip addr add 192.168.7.100/24 dev eth1

# setup test PC in Customer8 network -  CPC108
docker exec -d clab-frr01-cpc108 ip link set eth1 up
docker exec -d clab-frr01-cpc108 ip addr add 192.168.8.100/24 dev eth1

# setup test PC in Customer9 network -  CPC109
docker exec -d clab-frr01-cpc109 ip link set eth1 up
docker exec -d clab-frr01-cpc109 ip addr add 192.168.9.100/24 dev eth1

# ================== #
#       AS65002      #
# ================== #

# setup test PC in Area0 -  PC201
docker exec -d clab-frr01-pc201 ip link set eth1 up
docker exec -d clab-frr01-pc201 ip addr add 10.20.0.100/24 dev eth1
docker exec -d clab-frr01-pc201 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc201 ip route add 65.0.2.2/24 via 10.20.0.1 eth1
docker exec -d clab-frr01-pc201 ip route add default via 10.20.0.2 dev eth1

# setup test PC in Area0 -  PC203
docker exec -d clab-frr01-pc203 ip link set eth1 up
docker exec -d clab-frr01-pc203 ip addr add 10.20.3.100/24 dev eth1
docker exec -d clab-frr01-pc203 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc203 ip route add default via 10.20.3.3 dev eth1

# setup test PC in Area0 -  PC204
docker exec -d clab-frr01-pc204 ip link set eth1 up
docker exec -d clab-frr01-pc204 ip addr add 10.20.4.100/24 dev eth1
docker exec -d clab-frr01-pc204 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc204 ip route add default via 10.20.4.4 dev eth1

# ================== #
#       AS65003      #
# ================== #

# setup test PC in Area1 -  PC301
docker exec -d clab-frr01-pc301 ip link set eth1 up
docker exec -d clab-frr01-pc301 ip addr add 10.30.0.100/24 dev eth1
docker exec -d clab-frr01-pc301 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc301 ip route add 65.0.3.2/24 via 10.30.0.2 eth1
docker exec -d clab-frr01-pc301 ip route add default via 10.30.0.2 dev eth1

# setup test PC outside OSPF -  PC393
docker exec -d clab-frr01-pc393 ip addr add 192.168.33.100/24 dev eth1
docker exec -d clab-frr01-pc393 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc393 ip route add default via 192.168.33.91 dev eth1

# setup test PC outside OSPF -  PC394
docker exec -d clab-frr01-pc394 ip addr add 192.168.34.100/24 dev eth1
docker exec -d clab-frr01-pc394 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc394 ip route add default via 192.168.34.91 dev eth1

# setup test PC outside OSPF -  PC392
docker exec -d clab-frr01-pc392 ip addr add 192.168.32.100/24 dev eth1
docker exec -d clab-frr01-pc392 ip route del default via 172.20.20.1 dev eth0
docker exec -d clab-frr01-pc392 ip route add default via 192.168.32.1 dev eth1
