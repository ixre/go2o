# Go2o Docker Image

```
docker run -d --name go2o \
    -p 1427:1427 -p 1428:1428 \
    -v $(pwd)/data:/data \
    -e "GO2O_NATS_ADDR:172.17.0.1:4222"
    --restart always \
    docker-base.56x.net:5020/go2o
```
