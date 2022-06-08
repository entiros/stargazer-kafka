# Kafka Stargazer agent

The Kafka Stargazer Agent is used to create services in Starlify matching the topics in your Kafka cluster.

# Setup

Follow the _Install Kafka Stargazer agent_ guide at https://starlify.entiros.se. Take the properties in the final step and use the to configure your Kafka Stargazer agent (see below).

# Configuration

Kafka Stargazer agent can be configured using a configuration file or environment properties.

## Configuration file
**Example**
```yaml
# Starlify configuration
starlify:
  apiKey: "[Starlify Agent API key]"
  agentId: "[Starlify Agent ID]"
  systemId: "[Target system ID]"

# Kafka configuration
kafka:
  host: "[Kafka bootstrap server]"
  oauth:
    token: "[OAUTH token or empty]"
```

## Environment properties

```
STARLIFY_APIKEY=[Starlify Agent API key]
STARLIFY_SYSTEMID=[Starlify Agent ID]
STARLIFY_AGENTID=[Target system ID]

KAFKA_HOST=[Kafka bootstrap server]
KAFKA_OAUTH_TOKEN=[OAUTH token or empty]
```

# Using the Kafka Stargazer agent 
```shell script
$ ./stargazer-kafka /path/to/local/config.yml
```

# Using the Kafka Stargazer agent image

```shell script
docker run \
    --volume=/path/to/local/config.yml:/configs/config.yaml \
    stargazer-kafka:latest
```

## Docker compose using configuration file
```yaml
version: "3"
services:
  stargazer:
    image: stargazer-kafka:latest
    restart: unless-stopped
    volumes:
      - /path/to/local/config.yml:/configs/config.yaml
```

## Docker compose using environment variables
```yaml
version: "3"
services:
  stargazer:
    image: stargazer-kafka:latest
    restart: unless-stopped
    environment:
      - KAFKA_HOST=[Kafka bootstrap server]
      - KAFKA_OAUTH_TOKEN=
      - STARLIFY_APIKEY=[Starlify Agent API key]
      - STARLIFY_AGENTID=[Starlify Agent ID]
      - STARLIFY_SYSTEMID=[Target system ID]
```

# Testing locally

A example docker composer configuration is provided in this repo that starts a Zookeeper, Kafka and Stargazer docker container. Update the environment variables for Starlify in `docker-compose-example.yml` and start using Docker compose:
```shell script
$  docker-compose -f docker-compose-example.yml up
```

**To create a new topic in Kafka**

1. Open your terminal and exec inside the Kafka container
```shell script
$  docker exec -it stargazer-kafka-kafka-1 bash
```

2. In the container, go to the below path
```shell script
$  cd /opt/bitnami/kafka
```

3. Create a topic
```shell script
$  bin/kafka-topics.sh --create --topic my-first-kafka-topic --bootstrap-server localhost:9092 --partitions 1 --replication-factor 1
```