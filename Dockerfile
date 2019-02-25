FROM spotify/kafka

ENV ADVERTISED_HOST=localhost
ENV ADVERTISED_PORT=9092
EXPOSE 9092
EXPOSE 2181

ADD server.properties /opt/kafka_2.11-0.10.1.0/config/server.properties
WORKDIR /opt/kafka_2.11-0.10.1.0

