.. This work is licensed under a Creative Commons Attribution 4.0 International License.
.. http://creativecommons.org/licenses/by/4.0
.. Copyright 2018 Intel Corporation, Inc

Installation
------------

**Kubernetes**

The Secret Management Service project is a sub-project of AAF and will be deployed via Helm on Kubernetes
under the OOM Project umbrella. It will be automatically installed when the AAF chart is installed.

**Standalone Install on Bare-Metal or VM**

A script for doing a standalone install is provided in the repository
Run it as below:

.. code-block:: console

    cd sms-service/bin/deploy
    sms.sh start

.. end
