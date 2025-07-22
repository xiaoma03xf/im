#!/bin/bash

# 进入 kafka 容器，创建 topic
docker exec -it kafka /bin/bash -c "/opt/kafka/bin/kafka-topics.sh --create --topic topic-edge-01 --bootstrap-server localhost:9092"
echo "Topic topic-edge-01 创建完成"
