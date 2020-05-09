# Hyperledger_fabric_2.0_Bank_simple
This Project is example of hyperledger fabric 2.0 with updated network configration ,install of chaincode and new writen chaincode in Golang 

**This project contain four folder**

* Artifacts
* Channel-artifacts
* Connection
* Nodejs

**Network structure **

* Two organzation 
* Three Orderer Srevices 
* Two CA
* Two Peer in each org
* Four Coucdb

Before Start insatll updated version of golang (go version go1.14.2 linux/amd64) otherwise it will careate problem and setup bin folder path in you .bashrc file.

**To Generater Certifacte and Crypto Matrial**

Run ./generate.sh this fill create all necessary file and folder inside the artifact/channel/

**To Start Network**

Run ./start_network.sh file this file run your hyperledger fabric network.

**To install Chaincode**

Run ./installchaincode.sh file this file install chaincode by calling deployChaincode.sh and pass argumnet of chaincode name
This project contain Two different  chaincode
* User
* Teansaction


**To Remove and stop network**

Run ./teardown.sh file this file remove all docker container.


**Environment Variables**

envVar.sh help as to set environment variables while creating channel and install chaincode


Once You will install and run above command **cd nodejs** follow this steps
* Remove older keys from wallet
* Run node enrollAdmin.js
* Run node server.js
* Open Browser localhost:3000/ you can see the web interface of application

Once the applicaion  is start system have no user so before start click on create user button create user

This system have tow type of user 
* Employee who will do following thing
  * Create Transaction which is validated by his respected boss
  * View all Transaction created by him/her
  * View all Transaction in company under one boss 
  * To see details of Transaction he have to send request to the boss he/she will accepte or reject
* Boss who can Validate Transaction in system
  * View all new Transaction
  * View all transaction Validated by him/her
  * View Request for View Transaction Details

when we create Boss the role will be boss and who is boss will none in case of epmployee we have to tell who is the boss

