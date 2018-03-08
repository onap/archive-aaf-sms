#!/bin/bash
RED='\033[0;31m'
NC='\033[0m'
URL=$1
for i in `seq 1 2`; 
do
  echo -e "${RED}----------------BEGIN GET STATUS----------------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X GET \
    https://${URL}:10443/v1/sms/status

  echo -e "${RED}----------------BEGIN CREATE SECRET DOMAIN------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X POST \
    -d @test/test_create_domain.json https://${URL}:10443/v1/sms/domain

  echo -e "${RED}----------------BEGIN CREATE SECRET 1-----------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X POST \
    -d @test/test_create_secret1.json https://${URL}:10443/v1/sms/domain/curltestdomain/secret

  echo -e "${RED}----------------BEGIN CREATE SECRET 2-----------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X POST \
    -d @test/test_create_secret2.json https://${URL}:10443/v1/sms/domain/curltestdomain/secret

  echo -e "${RED}----------------BEGIN CREATE SECRET 3-----------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X POST \
    -d @test/test_create_secret3.json https://${URL}:10443/v1/sms/domain/curltestdomain/secret

  echo -e "${RED}----------------BEGIN LIST SECRET---------------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X GET \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret

  echo -e "${RED}----------------BEGIN GET SECRET 1--------------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X GET \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret/curltestsecret1

  echo -e "${RED}----------------BEGIN GET SECRET 2--------------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X GET \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret/curltestsecret2

  echo -e "${RED}----------------BEGIN GET SECRET 3--------------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X GET \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret/curltestsecret3

  echo -e "${RED}----------------BEGIN DELETE SECRET 1-----------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X DELETE \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret/curltestsecret1

  echo -e "${RED}----------------BEGIN DELETE SECRET 2-----------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X DELETE \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret/curltestsecret2

  echo -e "${RED}----------------BEGIN DELETE SECRET 3-----------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X DELETE \
    https://${URL}:10443/v1/sms/domain/curltestdomain/secret/curltestsecret3

  echo -e "${RED}----------------BEGIN DELETE SECRET DOMAIN------${NC}"
  curl -i -w "\n" -H "Accept: application/json" --cacert auth/selfsignedca.pem --cert auth/client.cert --key auth/client.key -X DELETE \
    https://${URL}:10443/v1/sms/domain/curltestdomain
done
