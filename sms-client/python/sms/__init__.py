# Copyright 2018 Intel Corporation, Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#       http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import json
import requests
import requests.exceptions
import urlparse

class InvalidRequestException(Exception):
    pass

class InternalServerError(Exception):
    pass

class UnexpectedError(Exception):
    pass

class Client(object):
    """Python Client for Secret Management Service"""

    def __init__(self, url='http://localhost:10443', timeout=30, cacert=None):
        """Creates a new SMS client instance

        Args:
            url (str): Base URL with port pointing to the SMS service
            timeout (int): Number of seconds before aborting the connection
            cacert (str): Path to the cacert that will be used to verify
                If this is None, verify will be False and the server cert
                is not verified by the client.
        Returns:
            A Client object which can be used to interact with SMS service
        """

        self.base_url = url
        self.timeout = timeout
        self.cacert = cacert
        self.session = requests.Session()

        self._base_api_url = '/v1/sms'

    def _urlJoin(self, *urls):
        """Joins given urls into a single url

        Args:
            urls (str): url fragments to be combined.

        Returns:
            str: Joined URL
        """

        return '/'.join(urls)

    def _raiseException(self, statuscode, errors=None):
        """ Handles Exception Raising based on statusCode

        Args:
            statuscode (int): status code returned by the server
            errors (str): list of strings containing error messages

        Returns:
            exception: An exception is raised based on error message
        """

        if statuscode == 400:
            raise InvalidRequestException(errors)
        if statuscode == 500:
            raise InternalServerError(errors)

        raise UnexpectedError(errors)

    def _request(self, method, url, headers=None, **kwargs):
        """Handles request for all the client methods

        Args:
            method (str): type of HTTP method (get, post or delete).
            url (str): api URL.
            headers (dict): custom headers if any.
            **kwargs: various args supported by requests library

        Returns:
            requests.Response: An object containing status_code and returned
                               json data is returned here.
        """

        if headers is None:
            headers = {
                'content-type': "application/json",
                'Accept': "application/json"
            }

        #Verify the server or not based on the cacert argument
        if self.cacert is None:
            verify = False
        else:
            verify = self.cacert

        url = urlparse.urljoin(self.base_url, url)
        response = self.session.request(method, url, headers=headers,
                                        allow_redirects=False, verify=verify,
                                        timeout = self.timeout, **kwargs)

        errors = None
        if response.status_code >= 400 and response.status_code < 600:
            #Request Failed. Raise Exception.
            errors = response.text
            self._raiseException(response.status_code, errors)

        return response

    def getStatus(self):
        """Returns Status of SMS Service

        Returns:
            bool: True or False depending on if SMS Service is ready.
        """

        url = self._urlJoin(self._base_api_url, 'quorum', 'status')

        response = self._request('get', url)
        return response.json()['sealstatus']

    def createDomain(self, domainName):
        """Creates a Secret Domain

        Args:
            domainName (str): Name of the secret domain to create

        Returns:
            string: UUID of the created domain name
        """


        domainName = domainName.strip()
        data = {"name": domainName}
        url = self._urlJoin(self._base_api_url, 'domain')

        response = self._request('post', url, json = data)
        return response.json()['uuid']

    def deleteDomain(self, domainName):
        """Deletes a Secret Domain

        Args:
            domainName (str): Name of the secret domain to delete

        Returns:
            bool: True. An exception will be raised if delete failed.
        """

        domainName = domainName.strip()
        url = self._urlJoin(self._base_api_url, 'domain', domainName)

        self._request('delete', url)
        return True

    def getSecretNames(self, domainName):
        """Get all Secret Names in Domain

        Args:
            domainName (str): Name of the secret domain

        Returns:
            string[]: List of strings each corresponding to a
                      Secret's Name in this Domain.
        """

        domainName = domainName.strip()
        url = self._urlJoin(self._base_api_url, 'domain', domainName,
                            'secret')

        response = self._request('get', url)
        return response.json()['secretnames']


    def storeSecret(self, domainName, secretName, values):
        """Store a Secret in given Domain

        Args:
            domainName (str): Name of the secret domain
            secretName (str): Name for the Secret
            values (dict): A dict containing name-value pairs which
                           form the secret

        Returns:
            bool: True. An exception will be raised if store failed.
        """

        domainName = domainName.strip()
        secretName = secretName.strip()
        url = self._urlJoin(self._base_api_url, 'domain', domainName,
                            'secret')

        if not isinstance(values, dict):
            raise TypeError('Input values is not a dictionary')

        data = {"name": secretName, "values": values}
        self._request('post', url, json = data)
        return True

    def getSecret(self, domainName, secretName):
        """Get a particular Secret from Domain.

        Args:
            domainName (str): Name of the secret domain
            secretName (str): Name of the secret

        Returns:
            dict: dictionary containing the name-value pairs
                  which form the secret
        """

        domainName = domainName.strip()
        secretName = secretName.strip()
        url = self._urlJoin(self._base_api_url, 'domain', domainName,
                            'secret', secretName)

        response = self._request('get', url)
        return response.json()['values']

    def deleteSecret(self, domainName, secretName):
        """Delete a particular Secret from Domain.

        Args:
            domainName (str): Name of the secret domain
            secretName (str): Name of the secret

        Returns:
            bool: True. An exception will be raised if delete failed.
        """

        domainName = domainName.strip()
        secretName = secretName.strip()
        url = self._urlJoin(self._base_api_url, 'domain', domainName,
                            'secret', secretName)

        self._request('delete', url)
        return True