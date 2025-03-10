# VOLUME
docker volume create clickhouse_vol

# NETWORK
docker network create app_net

# SUPERSET
docker run -d --net=app_net -p 80:8088 -e "SUPERSET_SECRET_KEY=diplomClickhouse" --name superset apache/superset
docker exec -it superset superset fab create-admin \
              --username admin \
              --firstname Superset \
              --lastname Admin \
              --email admin@superset.com \
              --password aboba322
docker exec -it superset superset db upgrade
docker exec -it superset superset init

# CLICKHOUSE
docker run -d \
    --name clickhouseDB \
    --net=app_net \
    -p 8123:8123 \
    -v clickhouse_vol:/var/lib/clickhouse \
    clickhouse/clickhouse-server

# SUPERSET-CLICKHOUSE
docker exec superset pip install clickhouse-connect
docker restart superset
////// clickhousedb://clickhouseDB/default

# CLICKHOUSE Query
CREATE TABLE app.app(
  uid UInt32
) engine=MergeTree() ORDER BY uid;