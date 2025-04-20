# This script removes the default route via management (internet access) interface from internal routers in nssa and stub areas
# This ensures connectivity to external domains, which would have been announced through LSA Type 5

# remove default route for IRs in AS65001 Area3 (stub)
docker exec -it clab-frr01-r131 ip route del default

# remove default route for IRs in AS65001 Area1 (nssa)
docker exec -it clab-frr01-r111 ip route del default