/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

)

// const   PRODUCER = "FISHCO"


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Domation struct {
	Id   		string  `json:"id"`
	Owner  		string  `json:"owner"`
	Amount		long 	`json:"amount"`
	ProjectID	long 	`json:"projectID"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id   		string  `json:"id"`
	tDate		string	`json:"tdate"`
	Destination  	string  `json:"destination"`
	Amount		int 	`json:"amount"`
	ProjectID	long 	`json:"projectID"`
	TType 		string   `json:"ttype"`
}

type AllDonations struct{
	Donations []string `json:"donations"`
}

type AllDonationsDetails struct{
	Donations []Donation `json:"donations"`
}

// ============================================================================================================================
// Init 
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	var err error
	
	var donations AllDonations
	jsonAsBytes, _ := json.Marshal(donations)
	err = stub.PutState("allDonations", jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}



// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Run is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state
		return t.Init(stub, "init", args)
	} else if function == "createDonation" {											//create a batch
		return t.createDonation(stub, args)
	}

	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) != 1 { return nil, errors.New("Incorrect number of arguments passed") }

	//if function != "getBatch" && function != "getAllBatches" && function != "getAllBatchesDetails" && function != "getNbItems"{
		//return nil, errors.New("Invalid query function name.")
	//}

	if function == "getAllDonationsByUserId" { return t.getAllDonationsByUserId(stub, args[0]) }


	return nil, nil										
}

// Get all donations by user id

func (t *SimpleChaincode) getAllDonationsByUserId(stub *shim.ChaincodeStub, userID string)([]byte, error){

	fmt.Println("Start find all Donations ")
	fmt.Println("Looking for all Donations " + userID);

	//get the AllBatches index
	allDonsBytes, err := stub.GetState("allDonations")
	if err != nil {
		return nil, errors.New("Failed to get all Batches")
	}

	var res AllDonations
	err = json.Unmarshal(allDonsBytes, &res)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Donations")
	}

	var rad AllDonations

	for i := range res.Donations{

		sdASBytes, err := stub.GetState(res.Batches[i])
		if err != nil {
			return nil, errors.New("Failed to get Donation")
		}
		var sd Donation
		json.Unmarshal(sdASBytes, &sd)

		if(sb.Owner == userID) {
			rad.Donations = append(rad.Donations,sd.Id);
		}

	}

	radAsBytes, _ := json.Marshal(rad)

	return radAsBytes, nil

}


// Create a donation and store in the block chain

func (t *SimpleChaincode) createDonation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running create Donation")

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments. Expecting 6")
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}


	var do Domation
	do.Id 			= args[3]
	//do.Amount		= args[0]
	do.Owner		= args[2]
	do.ProjectID		= args[1]

	var tx Transaction
	tx.ProjectID		= args[1]
	tx.TType 		= "CREATE"
	tx.Amount		= args[0]
	tx.tDate		= args[4]
	tx.Id			= args[3] + "T1"

	do.Transactions = append(do.Transactions, tx)

	//Commit batch to ledger
	fmt.Println("createDonation Commit Donation To Ledger");
	doAsBytes, _ := json.Marshal(do)
	err = stub.PutState(do.Id, doAsBytes)
	if err != nil {
		return nil, err
	}

	//Update All Batches Array
	allDOsBytes, err := stub.GetState("allDonations")
	if err != nil {
		return nil, errors.New("Failed to get all Donations")
	}
	var alld AllDonations
	err = json.Unmarshal(allDOsBytes, &alld)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Donations")
	}
	alld.Donations = append(alld.Batches,do.Id)

	allDoAsBytes, _ := json.Marshal(alld)
	err = stub.PutState("allDonations", allDoAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}






func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
