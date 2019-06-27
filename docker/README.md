# Go2o Docker Image

```
docker run -d --name go2o \
    -p 1427:1427 -p 1428:1428 \
    -v $(pwd)/data:/data \
    -e "GO2O_KAFKA_ADDR:172.17.0.1:9092"
    --restart always \
    docker-base.to2.net:5020/go2o
```
