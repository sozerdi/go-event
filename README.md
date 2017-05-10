## Installation

The Project has been implemented in Go Language 1.8.1. It is expected to be installed in order to run the project.
There is a need of MySQL database to store events.

## Prerequisites

After installations have been done, there are some other issues which are needed to complete prior
to use the API or histogram. At first "go get" command should be invoked to download libraries.
Then database configuration should be set via properties.ini file. If the events table does not exist in database, it should be created by running table.go.
Then the main.go should be run in order to wake the system up.

## Usage

There are two features inside the project.

1. There is a RESTful API which consumes user generated events. The API may be invoked on <ip:port>/event 
   and it consumes application/json. There are three required fields:
   API_KEY (Possible values: 1,2,3), USER_ID, TIMESTAMP. More fields may be included as optional fields in request body.
   Request bodies will be stored in database with the API_KEY as a key.

2. The other feature of the project is Histogram. Response times of the API can be monitorized on <ip:port>/histogram

## Technology

The Project has been implemented in Go Language 1.8.1 which is the recent version.
MySQL has been preferred as database because of reliability and easy usage.

