#!/bin/bash
#
# -------------------------------------------------------------------------
#   Copyright 2018 Intel Corporation, Inc
#
#   Licensed under the Apache License, Version 2.0 (the "License");
#   you may not use this file except in compliance with the License.
#   You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
#   Unless required by applicable law or agreed to in writing, software
#   distributed under the License is distributed on an "AS IS" BASIS,
#   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
#   See the License for the specific language governing permissions and
#   limitations under the License.
#
# -------------------------------------------------------------------------

SS=0
if [ -e /opt/config/nexus_docker_repo.txt ]
then
	NEXUS_DOCKER_REPO=$(cat /opt/config/nexus_docker_repo.txt)
else
	NEXUS_DOCKER_REPO=nexus3.onap.org:10001
fi
echo "Using ${NEXUS_DOCKER_REPO} for docker Repo"

SMS_IMG=${NEXUS_DOCKER_REPO}/onap/aaf/sms:latest
QUO_IMG=${NEXUS_DOCKER_REPO}/onap/aaf/smsquorumclient:latest
VAU_IMG=library/vault:0.10.0
CON_IMG=library/consul:1.0.6
WORK_DIR=${PWD}

if [ "$1" = "start" ]; then

# Create Volume for mapping war file and tomcat
docker volume create sms-service;
docker volume create sms-consul;
docker volume create sms-quorum;

# Create a network for all the containers to run in.
docker network create sms-net;

# Create Consul Container
docker create --rm --name sms-consul --network sms-net \
--hostname sms-consul -p "8500:8500" \
-v sms-consul:/consul/data \
${CON_IMG} \
consul agent -server -client 0.0.0.0 -bootstrap-expect=1 -config-file /consul/config/config.json;

# Copy the configuration for Consul
docker cp consul.json sms-consul:/consul/config/config.json;

# Start the consul container
docker start sms-consul;

#Wait for Consul to start
sleep 10

# Create Vault Container
docker create --rm --name sms-vault --network sms-net \
--hostname sms-vault -p "8200:8200" \
-e SKIP_SETCAP=true \
${VAU_IMG} \
vault server -config /vault/config/config.json;

docker cp vault.json sms-vault:/vault/config/config.json;
docker start sms-vault;

# Start SMS
# Matching hostname with cert name
docker create --rm --name aaf-sms.onap --network sms-net \
--hostname aaf-sms.onap -p "10443:10443" \
-v sms-service:/sms/auth \
${SMS_IMG};

docker cp smsconfig.json aaf-sms.onap:/sms/smsconfig.json
docker start aaf-sms.onap

# Start 3 Quorum Clients
for i in {0..2}
do
	docker create --rm --name sms-quorum-$i --network sms-net \
	--hostname sms-quorum-$i \
	-v sms-quorum:/quorumclient/auth \
	${QUO_IMG};

	docker cp quorumconfig.json sms-quorum-$i:/quorumclient/config.json
	docker start sms-quorum-$i
done

# Connect service to host bridge network so that its port can be seen.
docker network connect bridge aaf-sms.onap;
SS=1;
fi

# Shutdown and clean up.
if [ "$1" = "stop" ]; then
docker stop sms-vault sms-consul aaf-sms.onap;
for i in {0..2}; do
docker stop sms-quorum-$i
done
docker network rm sms-net;
sleep 5;
docker volume rm sms-service;
docker volume rm sms-consul;
docker volume rm sms-quorum;
SS=1
fi

if [ $SS = 0 ]; then
	echo "Please type ${0} start or ${0} stop"
fi
