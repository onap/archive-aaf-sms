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


from setuptools import setup, find_packages

setup(

    name='onapsmsclient',
    keywords=("secret", "consul", "onap"),
    description="ONAP python SMS client library",
    long_description="python-package onapsmsclient client library for using"
                     " Secret Management Service (SMS) inside ONAP. Refer "
                     "https://wiki.onap.org/display/DW/Secret+Management+Service for more details.",
    version="0.0.3",
    url="https://gerrit.onap.org/r/gitweb?p=aaf%2Fsms.git;a=summary",
    license="Apache 2",
    author="Kiran Kamineni",
    author_email='kiran.k.kamineni@intel.com',
    packages=find_packages(),
    platforms=["all"],
    install_requires=[
        'requests>=2.7.0',
        'urllib3>=1.22'
    ],
    classifiers=[
        "Intended Audience :: Developers",
        "Programming Language :: Python :: 2.7",
        "Programming Language :: Python :: 3.6"
    ]
)
