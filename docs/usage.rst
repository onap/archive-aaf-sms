.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0
.. Copyright 2018 Intel Corporation, Inc

Typical Usage Scenario
======================

.. code-block:: guess

    ## Create a Domain
    ## This is where all your secrets will be stored
    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X POST \
        -d '{
                "name": "mysecretdomain"
            }'
        https://sms:10443/v1/sms/domain

    ## Add a new Secret
    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X POST \
        -d '{
                "name": "mysecret",
                "values": {
                    "name": "rah",
                    "age": 35,
                    "password": "mypassword"
                }
            }'
        https://sms:10443/v1/sms/domain/<domaincurltestdomain/secret


    ## List all Secrets under a Domain
    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X GET \
        https://sms:10443/v1/sms/domain/curltestdomain/secret

    ## Get a Secret in a Domain
    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X GET \
        https://sms:10443/v1/sms/domain/curltestdomain/secret/curltestsecret1

    ## Delete a Secret in specified Domain
    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X DELETE \
        https://sms:10443/v1/sms/domain/curltestdomain/secret/curltestsecret1

    ## Delete a Domain
    ## This will delete all the secrets in that Domain
    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X DELETE \
        https://sms:10443/v1/sms/domain/curltestdomain

.. end
