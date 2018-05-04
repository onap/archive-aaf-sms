SMS 1.0.0 API
===============================

.. toctree::
    :maxdepth: 3


Description
~~~~~~~~~~~

This is a service that provides secret management facilities



Contact Information
~~~~~~~~~~~~~~~~~~~



kiran.k.kamineni@intel.com





License
~~~~~~~


`Apache 2.0 <http://www.apache.org/licenses/LICENSE-2.0.html>`_




Base URL
~~~~~~~~

https://aaf.onap.org:10443/v1/sms/

Security
~~~~~~~~


.. _securities_token:

token (API Key)
---------------



**Name:** token

**Located in:** header




DOMAIN
~~~~~~


Operations related to Secret Domains





DELETE ``/domain/{domainName}``
-------------------------------


Summary
+++++++

Deletes a domain by name

Description
+++++++++++

.. raw:: html

    Deletes a domain with provided name

Parameters
++++++++++

.. csv-table::
    :delim: |
    :header: "Name", "Located in", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 15, 10, 10, 10, 20, 30

        domainName | path | Yes | string |  |  | Name of the domain


Request
+++++++


Responses
+++++++++

**204**
^^^^^^^

Successful Deletion


**404**
^^^^^^^

Invalid Path or Path not found






POST ``/domain``
----------------


Summary
+++++++

Add a new domain



Request
+++++++



.. _d_c7bdcff9aff0692da98e588abdbc895b:

Body
^^^^

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        name | No | string |  |  | Name of the secret domain under which all secrets will be stored
        uuid | No | string |  |  | Optional value provided by user. If user does not provide, server will auto generate

.. code-block:: javascript

    {
        "name": "somestring", 
        "uuid": "somestring"
    }

Responses
+++++++++

**201**
^^^^^^^

Successful Creation


Type: :ref:`Domain <d_c7bdcff9aff0692da98e588abdbc895b>`

**Example:**

.. code-block:: javascript

    {
        "name": "somestring", 
        "uuid": "somestring"
    }

**400**
^^^^^^^

Invalid input


**500**
^^^^^^^

Internal Server Error




  
LOGIN
~~~~~


Operations related to username password based authentication





POST ``/login``
---------------


Summary
+++++++

Login with username and password

Description
+++++++++++

.. raw:: html

    Operations related to logging in via username and Password


Request
+++++++



.. _d_8e36d758bad367e4538a291a5dd5355f:

Body
^^^^

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        password | No | string |  |  | 
        username | No | string |  |  | 

.. code-block:: javascript

    {
        "password": "somestring", 
        "username": "somestring"
    }

Responses
+++++++++

**200**
^^^^^^^

Successful Login returns a token


.. _i_bbceffdf8441c1c476ca77c42ad12f85:

**Response Schema:**

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        token | No | string |  |  | 
        ttl | No | integer |  |  | ttl of returned token in seconds


**Example:**

.. code-block:: javascript

    {
        "token": "somestring", 
        "ttl": 1
    }

**404**
^^^^^^^

Invalid Username or Password




  
SECRET
~~~~~~


Operations related to Secrets





DELETE ``/domain/{domainName}/secret/{secretName}``
---------------------------------------------------


Summary
+++++++

Deletes a Secret


Parameters
++++++++++

.. csv-table::
    :delim: |
    :header: "Name", "Located in", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 15, 10, 10, 10, 20, 30

        secretName | path | Yes | string |  |  | Name of Secret to Delete
        domainName | path | Yes | string |  |  | Path to the SecretDomain which contains the Secret


Request
+++++++


Responses
+++++++++

**204**
^^^^^^^

Successful Deletion


**404**
^^^^^^^

Invalid Path or Path not found






GET ``/domain/{domainName}/secret``
-----------------------------------


Summary
+++++++

List secret Names in this domain

Description
+++++++++++

.. raw:: html

    Gets all secret names in this domain

Parameters
++++++++++

.. csv-table::
    :delim: |
    :header: "Name", "Located in", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 15, 10, 10, 10, 20, 30

        domainName | path | Yes | string |  |  | Name of the domain in which to look at


Request
+++++++


Responses
+++++++++

**200**
^^^^^^^

Successful operation


.. _i_1dcddfd6f11cba3fb2516d3a61cd1b77:

**Response Schema:**

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        secretnames | No | array of string |  |  | Array of strings referencing the secret names


**Example:**

.. code-block:: javascript

    {
        "secretnames": [
            "secretname1", 
            "secretname2", 
            "secretname3"
        ]
    }

**404**
^^^^^^^

Invalid Path or Path not found






GET ``/domain/{domainName}/secret/{secretName}``
------------------------------------------------


Summary
+++++++

Find Secret by Name

Description
+++++++++++

.. raw:: html

    Returns a single secret

Parameters
++++++++++

.. csv-table::
    :delim: |
    :header: "Name", "Located in", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 15, 10, 10, 10, 20, 30

        domainName | path | Yes | string |  |  | Name of the domain in which to look at
        secretName | path | Yes | string |  |  | Name of the secret which is needed


Request
+++++++


Responses
+++++++++

**200**
^^^^^^^

successful operation


Type: :ref:`Secret <d_5e5fddd9ede6eb091e8496a9c55b84c3>`

**Example:**

.. code-block:: javascript

    {
        "name": "somestring", 
        "values": {
            "Age": 40, 
            "admin": true, 
            "name": "john"
        }
    }

**404**
^^^^^^^

Invalid Path or Path not found






POST ``/domain/{domainName}/secret``
------------------------------------


Summary
+++++++

Add a new secret


Parameters
++++++++++

.. csv-table::
    :delim: |
    :header: "Name", "Located in", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 15, 10, 10, 10, 20, 30

        domainName | path | Yes | string |  |  | Name of the domain


Request
+++++++



.. _d_5e5fddd9ede6eb091e8496a9c55b84c3:

Body
^^^^

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        name | No | string |  |  | Name of the secret
        values | No | :ref:`values <i_a9213c9639162b77082e257e19cca0d0>` |  |  | Map of key value pairs that constitute the secret

.. _i_a9213c9639162b77082e257e19cca0d0:

**Values schema:**


Map of key value pairs that constitute the secret

Map of {"key":":ref:`values-mapped <m_4d863967ef9a9d9efdadd1b250c76bd6>`"}

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25



.. code-block:: javascript

    {
        "name": "somestring", 
        "values": {
            "Age": 40, 
            "admin": true, 
            "name": "john"
        }
    }

Responses
+++++++++

**201**
^^^^^^^

Successful Creation


**404**
^^^^^^^

Invalid Path or Path not found




  
SYSTEM
~~~~~~


Operations related to quorum client which are not useful to clients





GET ``/status``
---------------


Summary
+++++++

Get backend status

Description
+++++++++++

.. raw:: html

    Gets current backend status. This API is used only by quorum clients


Request
+++++++


Responses
+++++++++

**200**
^^^^^^^

Successful operation


.. _i_ac1bc8e82eadbd8c03f852e15be4d03b:

**Response Schema:**

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        sealstatus | No | string |  |  | seal status of backend


**Example:**

.. code-block:: javascript

    {
        "sealstatus": "somestring"
    }

**404**
^^^^^^^

Invalid Path or Path not found






POST ``/unseal``
----------------


Summary
+++++++

Unseal backend

Description
+++++++++++

.. raw:: html

    Sends unseal shard to unseal if backend is sealed


Request
+++++++



.. _i_9d32e021ba68855cbb6e633520b7cd2d:

Body
^^^^

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        unsealshard | No | string |  |  | Unseal shard that will be used along with other shards to unseal backend

.. code-block:: javascript

    {
        "unsealshard": "somestring"
    }

Responses
+++++++++

**201**
^^^^^^^

Submitted unseal key


**404**
^^^^^^^

Invalid Path or Path not found




  
Data Structures
~~~~~~~~~~~~~~~

.. _d_8e36d758bad367e4538a291a5dd5355f:

Credential Model Structure
--------------------------

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        password | No | string |  |  | 
        username | No | string |  |  | 

.. _d_c7bdcff9aff0692da98e588abdbc895b:

Domain Model Structure
----------------------

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        name | No | string |  |  | Name of the secret domain under which all secrets will be stored
        uuid | No | string |  |  | Optional value provided by user. If user does not provide, server will auto generate

.. _d_5e5fddd9ede6eb091e8496a9c55b84c3:

Secret Model Structure
----------------------

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25

        name | No | string |  |  | Name of the secret
        values | No | :ref:`values <i_a9213c9639162b77082e257e19cca0d0>` |  |  | Map of key value pairs that constitute the secret

.. _i_a9213c9639162b77082e257e19cca0d0:

**Values schema:**


Map of key value pairs that constitute the secret

Map of {"key":":ref:`values-mapped <m_4d863967ef9a9d9efdadd1b250c76bd6>`"}

.. csv-table::
    :delim: |
    :header: "Name", "Required", "Type", "Format", "Properties", "Description"
    :widths: 20, 10, 15, 15, 30, 25



