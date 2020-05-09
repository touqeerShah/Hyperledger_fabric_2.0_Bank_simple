set -ex

. ./createChannel.sh


# Bring the test network down
pushd ./artifacts
docker-compose up -d   ca-org1  ca-org2
sleep 5
docker-compose up -d    orderer.example.com   orderer2.example.com   orderer3.example.com
sleep 5
docker-compose up -d couchdb0 couchdb1 couchdb2 couchdb3
sleep 3
docker-compose up -d  peer0.org1.example.com peer0.org2.example.com
sleep 2
docker-compose up -d peer1.org1.example.com peer1.org2.example.com

popd

#
setup_channel
