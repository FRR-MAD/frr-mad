# create bridge SW101 to connect R101, R102 and PC102
docker network create \
  --driver bridge \
  --opt com.docker.network.bridge.name=sw101 \
  --opt com.docker.network.bridge.enable_ip_masquerade=false \
  sw101-net

docker network create \
  --driver bridge \
  --opt com.docker.network.bridge.name=sw111 \
  --opt com.docker.network.bridge.enable_ip_masquerade=false \
  sw111-net

