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

	"github.com/openblockchain/obc-peer/openchain/chaincode/shim"
)

const   PRODUCER = "FISHERMAN"
const   SHIPPING = "SHIPPINGCo"
const   RETAILER = "RETAILER"
const 	CONSUMER = "CONSUMER"
const 	CERTIFIER = "CERTIFIER"

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Batch struct {
	Id   		string  `json:"id"`
	BType  		string  `json:"bType"`
	Owner  		string  `json:"owner"`
	Quantity	int 	`json:"quantity"`
	Quality	 	string 	`json:"quality"`
	Signature	string 	`json:"signature"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	VDate 		string   `json:"vDate"`
	Location  	string   `json:"location"`
	TType 		string   `json:"ttype"`
}

type AllBatches struct{
	Batches []string `json:"batches"`
}

// ============================================================================================================================
// Init 
// ============================================================================================================================
func (t *SimpleChaincode) init(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	
	var batches AllBatches
	jsonAsBytes, _ := json.Marshal(batches)
	err = stub.PutState("allBatches", jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}



// ============================================================================================================================
// Run - Our entry point
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Run is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state
		return t.init(stub, args)
	} else if function == "createBatch" {											//create a batch
		return t.createBatch(stub, args)
	} else if function == "claimBatch" {											//claim ownership of a batch
		return t.claimBatch(stub, args)
	} else if function == "transferBatch" {											//transfer ownership of a batch
		return t.transferBatch(stub, args)
	} else if function == "sellBatchItem" {											//sell a batch item
		return t.sellBatchItem(stub, args)
	} else if function == "updateBatchQuality" {									//update batch quality
		return t.updateBatchQuality(stub, args)
	}
	fmt.Println("run did not find func: " + function)						//error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - read a variable from chaincode state - (aka read)
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {

	if len(args) != 2 { return nil, errors.New("Incorrect number of arguments passed") }

	if args[0] != "getBatch" && args[0] != "getAllBatches"{
		return nil, errors.New("Invalid query function name.")
	}

	if args[0] == "getBatch" { return t.getBatch(stub, args[1]) }
	if args[0] == "getAllBatches" { return t.getAllBatches(stub, args[1]) }

	return nil, nil										
}


// ============================================================================================================================
// Get Batch Details
// ============================================================================================================================
func (t *SimpleChaincode) getBatch(stub *shim.ChaincodeStub, batchId string)([]byte, error){
	
	fmt.Println("Start find Batch")
	fmt.Println("Looking for Batch #" + batchId);

	//get the batch index
	bAsBytes, err := stub.GetState(batchId)
	if err != nil {
		return nil, errors.New("Failed to get Batch #" + batchId)
	}

	return bAsBytes, nil
	
}


// ============================================================================================================================
// Get All Batches - only if user is CERTIFIER
// ============================================================================================================================
func (t *SimpleChaincode) getAllBatches(stub *shim.ChaincodeStub, user string)([]byte, error){
	
	fmt.Println("Start find getAllBatches ")
	fmt.Println("Looking for All Batches " + user);

	if user != CERTIFIER { return nil, errors.New("You are not allowed to retrieve all batches") }

	//get the AllBatches index
	allBAsBytes, err := stub.GetState("allBatches")
	if err != nil {
		return nil, errors.New("Failed to get all Batches")
	}

	return allBAsBytes, nil
	
}



// ============================================================================================================================
// Create new Batch of Items
// type Batch struct {
// 	Id   		string  `json:"id"`
// 	BType  		string  `json:"bType"`
// 	Owner  		string  `json:"owner"`
// 	Quantity	int 	`json:"quantity"`
// 	Quality	 	string 	`json:"quality"`
// 	Signature	string 	`json:"signature"`
// 	Transactions []Transaction `json:"transactions"`
// }

// type Transaction struct {
// 	VDate 		string   `json:"vDate"`
// 	Location  	string   `json:"location"`
// 	TType 		string   `json:"ttype"`
// }
// ============================================================================================================================
func (t *SimpleChaincode) createBatch(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running createBatch")

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments. Expecting 6")
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	if args[2] != PRODUCER { return nil, errors.New("You are not allowed to create a new batch") }

	var bt Batch
	bt.Id 			= args[0]
	bt.BType		= args[1]
	bt.Owner		= args[2]
	quantityValue, err := strconv.Atoi(args[3])
	if err != nil {
		return nil, errors.New("Invalid Quantity")
	}else{
		bt.Quantity = quantityValue
	}
	bt.Quality 		= "OK"
	bt.Signature 	= ""

	var tx Transaction
	tx.VDate		= args[4]
	tx.Location 	= args[5]
	tx.TType 		= "CREATE"

	bt.Transactions = append(bt.Transactions, tx)

	//Commit batch to ledger
	fmt.Println("createBatch Commit Batch To Ledger");
	btAsBytes, _ := json.Marshal(bt)
	err = stub.PutState(bt.Id, btAsBytes)	
	if err != nil {
		return nil, err
	}

	//Update All Batches Array
	allBAsBytes, err := stub.GetState("allBatches")
	if err != nil {
		return nil, errors.New("Failed to get all Batches")
	}
	var allb AllBatches
	err = json.Unmarshal(allBAsBytes, &allb)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Batches")
	}
	allb.Batches = append(allb.Batches,bt.Id)

	allBuAsBytes, _ := json.Marshal(allb)
	err = stub.PutState("allBatches", allBuAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}


// ============================================================================================================================
// Claim a batch - can only be done by shipping company
// ============================================================================================================================
func (t *SimpleChaincode) claimBatch(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running claimBatch")

	if len(args) != 4 {
		fmt.Println("Incorrect number of arguments. Expecting 4 (BatchId, user, date, location)")
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	if args[1] != SHIPPING { return nil, errors.New("You are not allowed to claim a batch") }

	//Update Batch owner
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Batch #" + args[0])
	}
	var bch Batch
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Batch #" + args[0])
	}
	bch.Owner = args[1]

	var tx Transaction
	tx.VDate		= args[2]
	tx.Location 	= args[3]
	tx.TType 		= "CLAIM"

	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	fmt.Println("claimBatch Commit Updates To Ledger");
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.Id, btAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Transfer a batch - can only be done by shipping company
// ============================================================================================================================
func (t *SimpleChaincode) transferBatch(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running transferBatch")

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments. Expecting 6 (BatchId, user, date, location, newOwner, signature)")
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	if args[1] != SHIPPING { return nil, errors.New("You are not allowed to transfer a batch") }

	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Batch #" + args[0])
	}
	var bch Batch
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Batch #" + args[0])
	}
	bch.Owner = args[4]
	bch.Signature = args[5]

	var tx Transaction
	tx.VDate		= args[2]
	tx.Location 	= args[3]
	tx.TType 		= "TRANSFER"

	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	fmt.Println("transferBatch Commit Updates To Ledger");
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.Id, btAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}


// ============================================================================================================================
// Sell batch item - can only be done by retailer (it changes the quantity in the batch object)
// ============================================================================================================================
func (t *SimpleChaincode) sellBatchItem(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running sellBatchItem")

	if len(args) != 6 {
		fmt.Println("Incorrect number of arguments. Expecting 6 (BatchId, user, date, location, quantity, newOwner)")
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	if args[1] != RETAILER { return nil, errors.New("You are not allowed to sell a batch item") }

	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Batch #" + args[0])
	}
	var bch Batch
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Batch #" + args[0])
	}
	bch.Owner = args[5]

	if(bch.Quantity-1 < 0) { return nil, errors.New("You can't sell anymore item from this batch") }
	
	bch.Quantity = bch.Quantity-1

	var tx Transaction
	tx.VDate		= args[2]
	tx.Location 	= args[3]
	tx.TType 		= "SELL"

	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	fmt.Println("sellBatchItem Commit Updates To Ledger");
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.Id, btAsBytes)	
	if err != nil {
		return nil, err
	}

	return nil, nil
}


// ============================================================================================================================
// Update batch quality - this is done through the shipping company via IOT (if the temperature of food truck is above a certain value it will call this method)
// ============================================================================================================================
func (t *SimpleChaincode) updateBatchQuality(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running updateBatchQuality")

	if len(args) != 5 {
		fmt.Println("Incorrect number of arguments. Expecting 5 (BatchId, user, date, location, message)")
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}

	if args[1] != SHIPPING { return nil, errors.New("You are not allowed to update batch quality") }

	//Update Batch data
	bAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get Batch #" + args[0])
	}
	var bch Batch
	err = json.Unmarshal(bAsBytes, &bch)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Batch #" + args[0])
	}
	bch.Quality = args[4]

	var tx Transaction
	tx.VDate		= args[2]
	tx.Location 	= args[3]
	tx.TType 		= "UPDATE QUALITY"

	bch.Transactions = append(bch.Transactions, tx)

	//Commit updates batch to ledger
	fmt.Println("updateBatchQuality Commit Updates To Ledger");
	btAsBytes, _ := json.Marshal(bch)
	err = stub.PutState(bch.Id, btAsBytes)	
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