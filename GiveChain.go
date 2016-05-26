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
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"

)

// const   PRODUCER = "FISHCO"


// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type Donation struct {
	Id   		string  `json:"id"`
	Owner  		string  `json:"owner"`
	Amount		int64 	`json:"amount"`
	ProjectID	string 	`json:"projectID"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Id   		string  `json:"id"`
	tDate		string	`json:"tdate"`
	Amount		int64 	`json:"amount"`
	ProjectID	string 	`json:"projectID"`
	TType 		string   `json:"ttype"`
	Details 	string	`json:"details"`
	SupplierID	string	`json:"SupplierID"`
}

type AllDonations struct{
	Donations []string `json:"donations"`
}

type AllDonationsDetails struct{
	Donations []Donation `json:"donations"`
}

type ProjectAmount struct {
	Name	string	`json:"name"`
	Amount	int64	`json:"amount"`
}

type DonationAmount struct {
	Name	string	`json:"name"`
	Amount	int64	`json:"amount"`
}

type Supplier struct {
	ID	string	`json:"id"`
	Transactions []Transaction `json:"transactions"`
}

type AllSuppliers struct{
	Suppliers []string `json:"suppliers"`
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

	var suppliers AllSuppliers
	supAsBytes, _ := json.Marshal(suppliers)
	err = stub.PutState("allSuppliers", supAsBytes)
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
	} else if function == "transfer" {											//create a batch
		return t.transfer(stub, args)
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
	if function == "getDonation" { return t.getDonation(stub, args[0]) }
	if function == "getTransactions" { return t.getTransactions(stub, args[0]) }
	if function == "getAmount" { return t.getAmount(stub, args[0]) }
	if function == "getSupplier" {return t.getSupplier(stub, args[0])}
	if function == "getAllSupplier" {return t.getAllSuppliers(stub)}
	if function == "getAllDonations" {return t.getAllDonations(stub)}
	if function == "getAmountD" { return t.getAmountD(stub, args[0]) }

	return nil, nil
}

// Get all donations by user id

func (t *SimpleChaincode) getAllDonationsByUserId(stub *shim.ChaincodeStub, userID string)([]byte, error){

	fmt.Println("Start find all Donations ")
	fmt.Println("Looking for all Donations " + userID);

	//get the AllBatches index
	allDonsBytes, err := stub.GetState("allDonations")
	if err != nil {
		return nil, errors.New("Failed to get all Donations")
	}

	var res AllDonations
	err = json.Unmarshal(allDonsBytes, &res)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Donations")
	}

	var rad AllDonations

	for i := range res.Donations{

		sdASBytes, err := stub.GetState(res.Donations[i])
		if err != nil {
			return nil, errors.New("Failed to get Donation")
		}
		var sd Donation
		json.Unmarshal(sdASBytes, &sd)

		if(sd.Owner == userID) {
			rad.Donations = append(rad.Donations, sd.Id);
		}

	}

	radAsBytes, _ := json.Marshal(rad)

	return radAsBytes, nil

}

// get all donations
func (t *SimpleChaincode) getAllDonations(stub *shim.ChaincodeStub)([]byte, error){

	fmt.Println("Start find all Donations ")
	fmt.Println("Looking for all Donations ");

	//get the AllBatches index
	allDonsBytes, err := stub.GetState("allDonations")
	if err != nil {
		return nil, errors.New("Failed to get all Donations")
	}

	return allDonsBytes, nil

}

// get all donations
func (t *SimpleChaincode) getAllSuppliers(stub *shim.ChaincodeStub)([]byte, error){

	fmt.Println("Start find all Suppliers ")
	fmt.Println("Looking for all Suppliers ");


	allDonsBytes, err := stub.GetState("allSuppliers")
	if err != nil {
		return nil, errors.New("Failed to get all suppliers")
	}

	return allDonsBytes, nil

}
// Create a donation and store in the block chain

func (t *SimpleChaincode) createDonation(stub *shim.ChaincodeStub, args []string) ([]byte, error) {

	var err error
	fmt.Println("Running create Donation")

	if len(args) != 5 {
		fmt.Println("Incorrect number of arguments. Expecting 5")
		return nil, errors.New("Incorrect number of arguments. Expecting 5")
	}


	var do Donation
	do.Id 			= args[3]
	//do.Amount		= args[0]
	do.Owner		= args[2]
	do.ProjectID		= args[1]

	var tx Transaction
	tx.ProjectID		= args[1]
	tx.TType 		= "CREATE"
	intAmount, err := strconv.Atoi(args[0])
	tx.Amount		= int64(intAmount)
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
	alld.Donations = append(alld.Donations,do.Id)

	allDoAsBytes, _ := json.Marshal(alld)
	err = stub.PutState("allDonations", allDoAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// get donations based on ID
func (t *SimpleChaincode) getDonation(stub *shim.ChaincodeStub, donationID string)([]byte, error){

	fmt.Println("Start find Donation")
	fmt.Println("Looking for Donation #" + donationID);

	//get the batch index
	dAsBytes, err := stub.GetState(donationID)
	if err != nil {
		return nil, errors.New("Failed to get Batch #" + donationID)
	}

	return dAsBytes, nil

}

//get all transactions based on donation ID
func (t *SimpleChaincode) getTransactions(stub *shim.ChaincodeStub, donationID string)([]byte, error){

	fmt.Println("Start find Donation")
	fmt.Println("Looking for Donation #" + donationID);


	dAsBytes, err := stub.GetState(donationID)
	if err != nil {
		return nil, errors.New("Failed to get Donation #" + donationID)
	}

	var do Donation
	err = json.Unmarshal(dAsBytes, &do)

	tAsBytes, err := json.Marshal(do.Transactions)

	return tAsBytes, nil

}

//get all transactions based on donation ID
func (t *SimpleChaincode) getSupplierTranssactions(stub *shim.ChaincodeStub, supplierID string)([]byte, error){

	fmt.Println("Start find Supplier")
	fmt.Println("Looking for Supplier #" + supplierID);


	dAsBytes, err := stub.GetState(supplierID)
	if err != nil {
		return nil, errors.New("Failed to get supplier #" + supplierID)
	}

	var do Supplier
	err = json.Unmarshal(dAsBytes, &do)

	tAsBytes, err := json.Marshal(do.Transactions)

	return tAsBytes, nil

}

func (t *SimpleChaincode) getAmount(stub *shim.ChaincodeStub, projectID string)([]byte, error){

	fmt.Println("Start find all Donations ")
	fmt.Println("Looking for all Donations with project id " + projectID);

	//get the AllBatches index
	allDonsBytes, err := stub.GetState("allDonations")
	if err != nil {
		return nil, errors.New("Failed to get all Donations")
	}

	var res AllDonations
	err = json.Unmarshal(allDonsBytes, &res)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal all Donations")
	}

	var pAmount ProjectAmount
	pAmount.Name = projectID


	for i := range res.Donations{

		sdASBytes, err := stub.GetState(res.Donations[i])
		if err != nil {
			return nil, errors.New("Failed to get Donation")
		}
		var sd Donation
		json.Unmarshal(sdASBytes, &sd)

		if(sd.ProjectID == projectID) {
			for x := range sd.Transactions{
				var tx Transaction
				tx = sd.Transactions[x]
				pAmount.Amount = pAmount.Amount + tx.Amount
			}

		}

	}

	paAsBytes, _ := json.Marshal(pAmount)

	return paAsBytes, nil

}



func (t *SimpleChaincode) transfer(stub *shim.ChaincodeStub, args []string) ([]byte, error) {


	var supplierID string = args[0]
	var projectID string = args[1]
	intAmount, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("Failed to get Amount")
	}
	var amount int64 = int64(intAmount)
	var originalAmount int64 = int64(intAmount)
	fmt.Println("calling getAmount function")
	pAmountAsByte , err := t.getAmount(stub,projectID)
	if err != nil {
		return nil, errors.New("Failed to gtun getAmount Function")
	}

	var pAmount ProjectAmount
	fmt.Println("unmarshall spamount")
	json.Unmarshal(pAmountAsByte, &pAmount)

	if pAmount.Amount > amount {

		//CHECK IF SUPPLIER EXISTS
		fmt.Println("getting suppler")
		supAsByte , err := stub.GetState(supplierID)

			//CREATE SUPPLIER AND ADD TO LEDGER AND ALLSUPPLIERS STATE

			var supplier Supplier
			json.Unmarshal(supAsByte,&supplier)
			supplier.ID = supplierID



			fmt.Println("getting allSuppliers")
			allSUsBytes, err := stub.GetState("allSuppliers")
			if err != nil {
				return nil, errors.New("Failed to get all Suppliers")
			}
			var allsup AllSuppliers
			err = json.Unmarshal(allSUsBytes, &allsup)
			if err != nil {
				return nil, errors.New("Failed to Unmarshal all Suppliers")
			}
			allsup.Suppliers = append(allsup.Suppliers,supplier.ID)

			allSupAsBytes, _ := json.Marshal(allsup)
			err = stub.PutState("allSuppliers", allSupAsBytes)
			if err != nil {
				return nil, errors.New("Failed to save all Suppliers to ledger")
				var ttt Supplier
				json.Unmarshal(supAsByte, &ttt)
			}


		allDonsBytes, err := stub.GetState("allDonations")
		if err != nil {
			return nil, errors.New("Failed to get all Donations")
		}

		var res AllDonations
		err = json.Unmarshal(allDonsBytes, &res)
		if err != nil {
			return nil, errors.New("Failed to Unmarshal all Donations")
		}

		for i := range res.Donations{

			sdASBytes, err := stub.GetState(res.Donations[i])
			if err != nil {
				return nil, errors.New("Failed to get Donation")
			}
			var sd Donation
			json.Unmarshal(sdASBytes, &sd)

			if(sd.ProjectID == projectID) {
				var dAmount int64 = 0
				var tCount int
				for x := range sd.Transactions{
					var tx Transaction
					tx = sd.Transactions[x]
					dAmount = dAmount + tx.Amount
					tCount = x
				}

				tCount = tCount + 1
				fmt.Println(dAmount)
				fmt.Println(amount)
				if dAmount < amount {
					fmt.Println("not enough funds")
					var tempamount int64
					tempamount = amount - dAmount
					amount = tempamount
					var tx Transaction
					tx.ProjectID		= projectID
					tx.TType 		= "TRANSFER"
					tx.Amount		= -dAmount
					tTime 			:= time.Now()
					tx.tDate		= tTime.Format("2006-01-02 15:04:05")
					tx.Id			= sd.Id + "T" + strconv.Itoa(tCount)
					tx.SupplierID		= supplier.ID
					if len(args) == 4 {
						tx.Details = args[3]
					}

					sd.Transactions = append(sd.Transactions, tx)

					//Commit batch to ledger
					fmt.Println("Create tx To Ledger");
					doAsBytes, _ := json.Marshal(sd)
					err = stub.PutState(sd.Id, doAsBytes)
					if err != nil {
						return nil, errors.New("Failed to put tx")
					}

					amount = amount - tempamount
					fmt.Println("tempampunt")
					fmt.Println(tempamount)
				} else {
					fmt.Println("enough funds")
					var tx Transaction
					tx.ProjectID		= projectID
					tx.TType 		= "TRANSFER"
					tx.Amount		= -amount
					tTime 			:= time.Now()
					tx.tDate		= tTime.Format("2006-01-02 15:04:05")
					tx.Id			= sd.Id + "T" + strconv.Itoa(tCount)
					tx.SupplierID		= supplier.ID

					if len(args) == 4 {
						tx.Details = args[3]
					}

					sd.Transactions = append(sd.Transactions, tx)


					fmt.Println("tx To Ledger");
					doAsBytes, _ := json.Marshal(sd)
					err = stub.PutState(sd.Id, doAsBytes)
					if err != nil {
						return nil, errors.New("Failed to put tx")
					}
					amount = amount - amount
				}

				if amount == 0 {
					fmt.Println("amount is 0")
					var sTx Transaction
					sTx.ProjectID		= projectID
					sTx.TType 		= "TRANSFER"
					sTx.Amount		= originalAmount
					tTime 			:= time.Now()
					sTx.tDate		= tTime.Format("2006-01-02 15:04:05")
					sTx.Id			= sd.Id + "T" + strconv.Itoa(tCount) + "S"
					sTx.SupplierID		= supplier.ID

					if len(args) == 4 {
						sTx.Details = args[3]
					}

					supplier.Transactions = append(supplier.Transactions, sTx)

					fmt.Println("createDonation Commit Donation To Ledger");
					doAsBytes, _ := json.Marshal(supplier)
					err = stub.PutState(supplier.ID, doAsBytes)
					if err != nil {
						return nil, errors.New("Failed to put tx")
					}
					return nil, nil
				}

			}

		}


	}

return nil, nil

}


func (t *SimpleChaincode) getAmountD(stub *shim.ChaincodeStub, donationID string)([]byte, error){


	fmt.Println("Looking for Donation with  id " + donationID);


	sDonAsBytes, err := stub.GetState(donationID)
	if err != nil {
		return nil, errors.New("Failed to get Donation")
	}

	var don Donation
	err = json.Unmarshal(sDonAsBytes, &don)
	if err != nil {
		return nil, errors.New("Failed to Unmarshal Donation")
	}

	var dAmount DonationAmount
	dAmount.Name = donationID

			for x := range don.Transactions{
				var tx Transaction
				tx = don.Transactions[x]
				dAmount.Amount = dAmount.Amount + tx.Amount
			}

	daAsBytes, _ := json.Marshal(dAmount)

	return daAsBytes, nil

}

// get donations based on ID
func (t *SimpleChaincode) getSupplier(stub *shim.ChaincodeStub, supplierID string)([]byte, error){

	fmt.Println("Start find Supplier")
	fmt.Println("Looking for supplier #" + supplierID);

	//get the batch index
	dAsBytes, err := stub.GetState(supplierID)
	if err != nil {
		return nil, errors.New("Failed to get supplier #" + supplierID)
	}

	return dAsBytes, nil

}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
