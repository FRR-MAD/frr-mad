# Containerlab Deployment

## Requirements
- Docker
- Containerlab
- amd64 architecture (or docker-in-docker)
- (make)

## Setup
A detail installation process can be found here: https://containerlab.dev/install/


## Usage

If the requirements from above are met and make is installed, you can use the makefile to deploy the environment locally.
- creates docker images for frr and build host
- it deploys the topology found in [devenvironment](../.devenvironment/)
- restarts or cleans up the topology


**make hmr/docker**
This command build the docker images requried for the topology. The version of FRR can be adjusted, but currently only 8.5.4 and 10.3.0 are working.

**make hmr/run**
To start the topology simply run the command above.

**make hmr/stop**
To stop the topology simply run the command above.

**make hmr/restart**
To restart the topology simply run the command above. Important to note, this does not clean up the topology. If there are caching issues, use the clean command below.

**make hmr/clean**
Destroys and cleans up the topology. Cleaning up any outstanding artifacts too.
