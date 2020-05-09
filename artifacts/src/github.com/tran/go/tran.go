/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"


	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

// Car describes basic details of what makes up a car
type Tran struct {
	TransactionId       string `json:"transactionid"`
	To                  string `json:"to"`
	From                string `json:"from"`
	Date                string `json:"date"`
	Amount              string `json:"amount"`
	TransactionCreateBy string `json:"transactioncreateby"` // this will tell who create the transaction
	ValidatedBy         string `json:"validatedby"`         // this will tell who is the Validator
	Validate            string `json:"validate"`            // this will tell transaction is valid or not
	ViewBy              string `json:"viewby"`              // this will tell who have permissions to View the transaction Details
}

type TransactionViewRequest struct { // this struct is user to store request data into SmartContract
	TransactionId  string `json:"transactionid"`
	RequestId      string `json:"requestid"`
	RequestTo      string `json:"requestto"`
	RequestBy      string `json:"requestby"`
	RequestProcess string `json:"requestprocess"`
}
//
// type QueryResult struct {
// 	Key    string
// }

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	_ = ctx.GetStub().PutState("transationCount", []byte(strconv.Itoa(1)))
	_ = ctx.GetStub().PutState("requestCount", []byte(strconv.Itoa(1)))

	return nil
}

// CreateCar adds a new car to the world state with given details
func (s *SmartContract) CreateTransaton(ctx contractapi.TransactionContextInterface, transactionId string, to string, from string, date string, amount string, transactioncreateby string, validatedby string) error {
	transationCountBytes, _ := ctx.GetStub().GetState("transationCount")
	transationCount, _ := strconv.Atoi(string(transationCountBytes))
	transationCount = transationCount + 1

	var transaction = Tran{
		TransactionId:       transactionId,
		To:                  to,
		From:                from,
		Date:                amount,
		Amount:              date,
		TransactionCreateBy: transactioncreateby,
		ValidatedBy:         validatedby,
		Validate:            `false`,
		ViewBy:              ``}
	transationAsBytes, _ := json.Marshal(transaction)
	id :=  transactionId + `-` + string(transationCountBytes)

	_ = ctx.GetStub().PutState("transationCount", []byte(strconv.Itoa(transationCount)))
	return ctx.GetStub().PutState(id, transationAsBytes)
}

/*
this function required two parameter one is transaction id and other will be userid to check the permission
*/

func (s *SmartContract) QueryTransation(ctx contractapi.TransactionContextInterface, transactionId string, viewBy string) *Tran {
	transationAsBytes, err := ctx.GetStub().GetState(transactionId)

	if err != nil {
		return nil
	}
	if transationAsBytes == nil {
		return nil
	}

	transaction := new(Tran)
	_ = json.Unmarshal(transationAsBytes, transaction)


	if transaction.ValidatedBy == viewBy {
		return transaction
	} else {
		index := strings.Index(transaction.ViewBy, viewBy)
		if index > -1 {
			return transaction
		} else {
			return nil
		}
	}
	return nil
}

/*
this function will get singel argument UserId which must be Boss id and return all the transaction id which are waiting for validation
*/
func (s *SmartContract) QueryAllUnvalidatedTransationId(ctx contractapi.TransactionContextInterface, validatedBy string) ([]string, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []string{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		if queryResponse.Key == `transationCount` || queryResponse.Key == `requestCount` {
			continue
		}

		transaction := new(Tran)
		_ = json.Unmarshal(queryResponse.Value, transaction)
		if transaction.Validate == `false` && transaction.ValidatedBy == validatedBy {
			results = append(results, queryResponse.Key)
		}
	}

	return results, nil
}

/*
this function will get singel argument UserId which must be Boss id and return all the transaction id which are Validated by Me
*/
func (s *SmartContract) QueryAllValidatedTransationIdByMe(ctx contractapi.TransactionContextInterface, validatedby string) ([]string, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []string{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		if queryResponse.Key == `transationCount` || queryResponse.Key == `requestCount` {
			continue
		}

		transaction := new(Tran)
		_ = json.Unmarshal(queryResponse.Value, transaction)
		if transaction.ValidatedBy == validatedby && transaction.Validate == `true` {
			results = append(results, queryResponse.Key)
		}
	}

	return results, nil
}

/*
this function will take 3 argument in parameter and change the transaction status to validated
first will user id which is logIn to check the validator is come who have permissions to make transaction validated
secand will be transactionid
third will be tell either transaction is valid or not
*/
func (s *SmartContract) ValidateTransation(ctx contractapi.TransactionContextInterface, transactionid string, userid string, validate string) error {
	transaction:= s.QueryTransation(ctx, transactionid, userid)


	if transaction.ValidatedBy == userid {
		transaction.Validate = validate
	}

	transactionAsBytes, _ := json.Marshal(transaction)

	return ctx.GetStub().PutState(transactionid, transactionAsBytes)
}

//this function will get singel argument UserId which must be Employee id and return all the transaction id which are Created by that Employee

func (s *SmartContract) QueryAllTransationCreatedByMe(ctx contractapi.TransactionContextInterface, TransactionCreateBy string) ([]string, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []string{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		if queryResponse.Key == `transationCount` || queryResponse.Key == `requestCount` {
			continue
		}
		transaction := new(Tran)
		_ = json.Unmarshal(queryResponse.Value, transaction)
		index := strings.Index(queryResponse.Key, `request`)
		if index == -1 {
			if transaction.TransactionCreateBy == TransactionCreateBy {
				results = append(results, queryResponse.Key)
			}
		}

	}

	return results, nil
}

/*
this function Return all id's of transaction which is created under some Boss
It required two parameter one boss id and other will be Employee id
*/
func (s *SmartContract) QueryAllTransationInCompany(ctx contractapi.TransactionContextInterface, ValidatedBy string, userid string) ([]string, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []string{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		if queryResponse.Key == `transationCount` || queryResponse.Key == `requestCount` {
			continue
		}
		transaction := new(Tran)
		_ = json.Unmarshal(queryResponse.Value, transaction)
		index := strings.Index(queryResponse.Key, `request`)
		if index == -1 {
			if transaction.ValidatedBy == ValidatedBy {
				data := queryResponse.Key
				index := strings.Index(transaction.ViewBy, userid+`,`) // here we check dose viewby string contains my user id if yes then i have permission to view the transaction details otherwise no i have to send requested for view

				if index > -1 {
					data += ":Yes"
				} else {
					data += ":No"
				}
				results = append(results, data)
			}
		}

	}

	return results, nil
}

func (s *SmartContract) CreateViwRequest(ctx contractapi.TransactionContextInterface, requestid string, transactionId string, requestto string, requestby string) error {
	requestCountBytes, _ := ctx.GetStub().GetState("requestCount")
	requestCount, _ := strconv.Atoi(string(requestCountBytes))
	requestCount = requestCount + 1
	var transactionViewRequest = TransactionViewRequest{
		TransactionId:  transactionId,
		RequestId:      requestid,
		RequestTo:      requestto,
		RequestBy:      requestby,
		RequestProcess: ``,
	}

	requestAsBytes, _ := json.Marshal(transactionViewRequest)
	id := requestid + `-` + string(requestCountBytes)

	_ = ctx.GetStub().PutState("requestCount", []byte(strconv.Itoa(requestCount)))
	return ctx.GetStub().PutState(id, requestAsBytes)
}

/*
this function will get singel argument UserId which must be Boss id and return all the Request To view Tran

*/

func (s *SmartContract) QueryAllRequestToViewTransion(ctx contractapi.TransactionContextInterface, RequestTo string) ([]string, error) {
	startKey := ""
	endKey := ""

	resultsIterator, err := ctx.GetStub().GetStateByRange(startKey, endKey)

	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	results := []string{}

	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()

		if err != nil {
			return nil, err
		}
		if queryResponse.Key == `transationCount` || queryResponse.Key == `requestCount` {
			continue
		}
		request := new(TransactionViewRequest)
		_ = json.Unmarshal(queryResponse.Value, request)

		if request.RequestTo == RequestTo && request.RequestProcess == `` {
			data := queryResponse.Key + ":" + request.TransactionId
			results = append(results, data)
		}

	}

	return results, nil
}

/*
This function two parameter request id and change the status of request to approved or deny
and secand will be allow or deny
*/
func (s *SmartContract) RequestProcess(ctx contractapi.TransactionContextInterface, requestid string, permission string) error {
	requestAsBytes, _ := ctx.GetStub().GetState(requestid)
	request := new(TransactionViewRequest)
	_ = json.Unmarshal(requestAsBytes, request)


	transactionAsBytes, _ := ctx.GetStub().GetState(request.TransactionId)
	transaction := new(Tran)
	_ = json.Unmarshal(transactionAsBytes, transaction)


	request.RequestProcess = `true` // change the RequestProcess true that mean request is process

	if permission == `allow` { // if Boss allow
		transaction.ViewBy += request.RequestBy + `,`     // then we add the requestby append to the viewby add allow user to view the transaction
	}
	transactionAsBytes, _ = json.Marshal(transaction) // update the Tran Record
	_ = ctx.GetStub().PutState(request.TransactionId, transactionAsBytes)

	requestAsBytes, _ = json.Marshal(request) // update The Request Record
	_ = ctx.GetStub().PutState(requestid, requestAsBytes)

	return nil
}

func main() {

	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create fabcar chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting fabcar chaincode: %s", err.Error())
	}
}
