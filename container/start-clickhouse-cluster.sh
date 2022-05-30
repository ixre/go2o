#/usr/bin/env sh
podman='podman'
if ! type podman > /dev/null; then podman='docker';fi;

podman network remove -f clickhouse_net
podman network create clickhouse_net

$podman run -d --name zoo1 -p 2181:2181 -e ALLOW_ANONYMOUS_LOGIN=yes bitnami/zookeeper:latest

$podman run -d --name clickhouse --ulimit nofile=262144:262144 \
    --hostname server0 \
    --network clickhouse_net \
    -p 8123:8123 -p 9000:9000 \
    -v $(pwd)/clickhouse/data:/var/lib/clickhouse \
    -v $(pwd)/clickhouse/config.xml:/etc/clickhouse-server/config.xml \
    -v $(pwd)/clickhouse/users.xml:/etc/clickhouse-server/users.xml \
    -v $(pwd)/clickhouse$i/metrika.xml:/etc/clickhouse-server/metrika.xml \
    --add-host server0:192.168.122.1 \
    --add-host server1:192.168.122.1 \
    --add-host server2:192.168.122.1 \
    --add-host server3:192.168.122.1 \
    clickhouse/clickhouse-server

for i in `seq 1`;do

$podman run -d --name clickhouse$i --ulimit nofile=262144:262144 \
    --hostname server$i \
    --network clickhouse_net \
	-p 900$i:9000 \
	-v $(pwd)/clickhouse$i/data:/var/lib/clickhouse \
	-v $(pwd)/clickhouse$i/config.xml:/etc/clickhouse-server/config.xml \
	-v $(pwd)/clickhouse$i/users.xml:/etc/clickhouse-server/users.xml \
    -v $(pwd)/clickhouse$i/metrika.xml:/etc/clickhouse-server/metrika.xml \
    --add-host server0:192.168.122.1 \
    --add-host server1:192.168.122.1 \
    --add-host server2:192.168.122.1 \
    --add-host server3:192.168.122.1 \
	clickhouse/clickhouse-server
done
