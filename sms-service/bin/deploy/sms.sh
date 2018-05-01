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

SMS_IMG=${NEXUS_DOCKER_REPO}/onap/music/cassandra_music:latest
QUO_IMG=${NEXUS_DOCKER_REPO}/onap/music/music:latest
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

# Start Consul 
docker create --rm --name sms-consul --network sms-net \
-p "8500:8500" \
${CON_IMG};
docker cp consul.hcl /consul/config
docker start sms-consul

# Start Vault
docker create --rm --name sms-vault --network sms-net \
-p "8200:8200" \
${VAU_IMG};
docker cp vault.hcl /vault/config
docker start sms-vault

# Start SMS
docker create --rm --name sms-service --network sms-net \
-p "10443:10443" \
${SMS_IMG};
docker cp smsconfig.json /sms/
docker start sms-service

# Start 3 SMS Quorum Clients
docker create --rm --name sms-quorum-0 --network sms-net \
-v sms-quorum:/quorumclient/auth \
${QUO_IMG};
docker cp quorumconfig.json /quorumclient
docker start sms-quorum-0

docker create --rm --name sms-quorum-1 --network sms-net \
-v sms-quorum:/quorumclient/auth \
${QUO_IMG};
docker cp quorumconfig.json /quorumclient
docker start sms-quorum-1

docker create --rm --name sms-quorum-2 --network sms-net \
-v sms-quorum:/quorumclient/auth \
${QUO_IMG};
docker cp quorumconfig.json /quorumclient
docker start sms-quorum-2

# Connect tomcat to host bridge network so that its port can be seen.
docker network connect bridge sms-service;
SS=1;
fi

# Shutdown and clean up.
if [ "$1" = "stop" ]; then
docker stop sms-vault sms-consul sms-service;
docker stop sms-quorum-0 sms-quorum-1 sms-quorum-2;
docker rm sms-vault sms-consul sms-service 
docker rm sms-quorum-0 sms-quorum-1 sms-quorum-2
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