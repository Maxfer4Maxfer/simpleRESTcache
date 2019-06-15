# Deploy service

In this deployment [Docker Machine](https://docs.docker.com/machine) and [Google Cloud Platform](https://cloud.google.com) will be used as a platform for Transaction App deployment.

## Prerequsites

Set up environment variables:
```bash
GOOGLE_PROJECT=<GOOGLE_PROJECT>
GOOGLE_ZONE=<GOOGLE_ZONE>
```
Create a virtual machine in GCP
```bash
docker-machine create --driver google \
--google-project $GOOGLE_PROJECT \
--google-machine-image https://www.googleapis.com/compute/v1/projects/ubuntu-os-cloud/global/images/family/ubuntu-1604-lts \
--google-machine-type n1-standard-1 \
--google-zone $GOOGLE_ZONE \
--google-open-port 80/tcp \
--google-open-port 8081/tcp \
vm1

eval $(docker-machine env vm1)
```

## Clone the git repository 
```bash
git clone https://github.com/Maxfer4Maxfer/simpleRESTcache.git
cd ./simpleRESTcache
```

## Copy configs and source files to the docker machine
```bash
docker-machine ssh vm1 "sudo rm -fR ~/*"
docker-machine scp -r $(pwd) vm1:~
docker-machine ssh vm1 "sudo mv ~/simpleR*/* ~/"
```
### Run docker-compose
```bash
docker-compose up -d
```

## Clean up installation
Stop application and delete a docker virtual machine
```bash
docker-compose down
docker-machine rm vm1
```