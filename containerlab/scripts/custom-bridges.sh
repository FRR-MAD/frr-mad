# create bridge SW101 to connect R101, R102 and PC102
docker network create \
  --driver bridge \
  --opt com.docker.network.bridge.name=sw101 \
  --opt com.docker.network.bridge.enable_ip_masquerade=false \
  --subnet=172.21.2.0/24 \
  sw101-net

docker network create \
  --driver bridge \
  --opt com.docker.network.bridge.name=sw102 \
  --opt com.docker.network.bridge.enable_ip_masquerade=false \
  --subnet=172.21.1.0/24 \
  sw102-net

docker network create \
  --driver bridge \
  --opt com.docker.network.bridge.name=sw111 \
  --opt com.docker.network.bridge.enable_ip_masquerade=false \
  --subnet=172.21.0.0/24 \
  sw111-net
