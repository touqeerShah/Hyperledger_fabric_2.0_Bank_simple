

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing a car
type SmartContract struct {
	contractapi.Contract
}

type Users struct {
	Name     string `json:"name"`
	UserId   string `json:"userid"`
	Password string `json:"password"`
	Role     string `json:"role"`
	WhoBoss  string `json:"whoboss"` // this will tell as who is your boss if rule is employe
}

// InitLedger adds a base set of cars to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
  return nil

}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, name string, userid string, password string, role string, whoboss string) error {
	user := Users{Name: name,
		UserId:   userid,
		Password: password,
		Role:     role,
		WhoBoss:  whoboss} // create object
	userAsBytes, _ := json.Marshal(user)

	return ctx.GetStub().PutState(userid, userAsBytes)
}

func (s *SmartContract) QueryUser(ctx contractapi.TransactionContextInterface, userId string) (*Users, error) {
	userAsBytes, err := ctx.GetStub().GetState(userId)

	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if userAsBytes == nil {
		return nil, fmt.Errorf("%s does not exist", userId)
	}

	user := new(Users)
	_ = json.Unmarshal(userAsBytes, user)

	return user, nil
}

func (s *SmartContract) QueryAllBossId(ctx contractapi.TransactionContextInterface) ([]string, error) {
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

		results = append(results,  queryResponse.Key)
	}

	return results, nil
}

func (s *SmartContract) LogIn(ctx contractapi.TransactionContextInterface, userId string, password string) string {
	user, err := s.QueryUser(ctx, userId)

	if err != nil {
		return `error`
	}
  check := ``
  if user.UserId == userId && user.Password == password {	// check the Value
    check = `true`
  } else {
    check = `fasle`
  }
	return check

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
