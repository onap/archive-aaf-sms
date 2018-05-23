.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0
.. Copyright 2018 Intel Corporation, Inc

Usage Scenario
==============

**Create a Domain**

This is the root where you will store your secrets.

.. code-block:: guess

    curl -H "Accept: application/json" --cacert ca.pem  --cert client.cert --key client.key
        -X POST \
        -d '{
                "name": "mysecretdomain"
            }'
        https://aaf-sms.onap:10443/v1/sms/domain

.. end

---------------

**Add a new Secret**

Store a new secret in your created Domain.
Secrets have a name and a map containing key value pairs.

.. code-block:: guess

    curl -H "Accept: application/json" --cacert ca.pem  --cert client.cert --key client.key
        -X POST \
        -d '{
                "name": "mysecret",
                "values": {
                    "name": "rah",
                    "age": 35,
                    "password": "mypassword"
                }
            }'
        https://aaf-sms.onap:10443/v1/sms/domain/<PREVIOUSLY CREATED DOMAIN NAME>/secret

.. end

---------------

**List all Secret Names in a Domain**

.. code-block:: guess

    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X GET \
        https://aaf-sms.onap:10443/v1/sms/domain/<PREVIOUSLY CREATED DOMAIN NAME>/secret

.. end

---------------

**Get a previously stored Secret from Domain**

.. code-block:: guess

    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X GET \
        https://aaf-sms.onap:10443/v1/sms/domain/<PREVIOUSLY CREATED DOMAIN NAME>/secret/<PREVIOUSLY CREATED SECRET NAME>

.. end

---------------

**Delete a Secret in specified Domain**

.. code-block:: guess

    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X DELETE \
        https://aaf-sms.onap:10443/v1/sms/domain/<PREVIOUSLY CREATED DOMAIN NAME>/secret/<PREVIOUSLY CREATED SECRET NAME>

.. end

---------------

**Delete a Domain**

.. code-block:: guess

    curl -H "Accept: application/json" --cacert ca.pem --cert client.cert --key client.key
        -X DELETE \
        https://aaf-sms.onap:10443/v1/sms/domain/<PREVIOUSLY CREATED DOMAIN NAME>
.. end
