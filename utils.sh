
de='docker exec -it kafka-offset-tester'

start_kafka() {
  docker run -d --name kafka-offset-tester -p 2181:2181 -p 9092:9092 kafka-offset-tester
}

create_topic() {
  $de bin/kafka-topics.sh --create --topic test-2 --zookeeper localhost:2181 --partitions 1 --replication-factor 1
  $de bin/kafka-configs.sh --zookeeper localhost:2181 --alter --entity-type topics --entity-name test-2 --add-config retention.ms=10000
  $de bin/kafka-configs.sh --zookeeper localhost:2181 --alter --entity-type topics --entity-name test-2 --add-config segment.ms=10000
}

describe_topic() {
  $de bin/kafka-configs.sh --describe --entity-type topics --entity-name test-2 --zookeeper localhost:2181
}

describe_consumers() {
  $de bin/kafka-consumer-groups.sh --describe --bootstrap-server localhost:9092 --group group-test
}

tail_logs() {
  $de tail -f logs/server.log
}
