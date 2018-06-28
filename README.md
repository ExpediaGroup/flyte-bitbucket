## Overview

![Build Status](https://travis-ci.org/HotelsDotCom/flyte-bitbucket.svg?branch=master)
[![Docker Stars](https://img.shields.io/docker/stars/hotelsdotcom/flyte-bitbucket.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-bitbucket)
[![Docker Pulls](https://img.shields.io/docker/pulls/hotelsdotcom/flyte-bitbucket.svg)](https://hub.docker.com/r/hotelsdotcom/flyte-bitbucket)

The bitbucket pack provides the ability to create a new repo.

## Build & Run
### Command Line
To build and run from the command line:
* Clone this repo
* Run `dep ensure` (must have [dep](https://github.com/golang/dep) installed )
* Run `go build`
* Run `FLYTE_API_URL=http://someflyteurl.com/ BITBUCKET_HOST=http://somebitbucketurl.com BITBUCKET_USER=user1 BITBUCKET_PASS=password ./flyte-bitbucket`
* Fill in this command with the relevant API url, bitbucket host, bitbucket user and bitbucket password environment variables

### Docker
To build and run from docker
* Run `docker build -t flyte-bitbucket .`
* Run `docker run -e FLYTE_API_URL=http://someurl.com/ -e BITBUCKET_HOST=http://bitbuckethost.com -e BITBUCKET_USER=user1 -e BITBUCKET_PASS=password flyte-bitbucket`
* All of these environment variables need to be set

## Commands
### createRepo command
This command creates a bitbucket repo.
#### Input
This commands input is the project and repo name:
```
"input": {
    "project": "FLYTE",
    "name": "newRepo"
    }
```
#### Output
This command returns either a `createRepo` event or a `createRepoFailure` event. 
##### createRepo event
This contains the url of the new repo.
```
"payload": {
    "project": "FLYTE",
    "name": "newRepo",
    "url": "http://somebitbuckethost.com/projects/FLYTE/repos/newRepo/browse"
}
```
##### createRepoFailure event
This contains the error message of why the repo could not be created
```
"payload": {
    "project": "FLYTE",
    "name": "newRepo",
    "error": "fail: status code 400"
}
```
### Rules
An example rule can be seen below:
```
{
  "id": "hipchatCreateRepo",
  "name": "createRepo",
  "description": "",
  "labels": {},
  "event": {
    "pack-id": "Hipchat",
    "name": "ReceivedMessage"
  },
  "command": {
    "pack-id": "BitBucket",
    "name": "createRepo",
    "input": {
        "project": "FLYTE",
        "name": "newRepo"
    }
  },
  "links": ...
}
```