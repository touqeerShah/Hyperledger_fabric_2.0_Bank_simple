
# Bring the test network down
pushd ./artifacts/channel
rm -r crypto-config/*
rm genesis.block
rm mychannel.tx
rm Org1MSPanchors.tx
rm Org2MSPanchors.tx
./create-artifacts.sh
popd
