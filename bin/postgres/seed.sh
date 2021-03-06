#!/bin/bash

# Copyright 2016 - 2018 Crunchy Data Solutions, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# seed the database after provisioning
# we will execute this script on each host since
# we dont' want to open up the superuser 'postgres' access 
# to outside of the container

source /opt/cpm/bin/common_lib.sh
enable_debugging

psql -U postgres postgres < /opt/cpm/bin/setup.sql
