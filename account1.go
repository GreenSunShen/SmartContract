package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
)

//==============================================================================================================================
//	Structure Definitions
//==============================================================================================================================
//	Chaincode - A blank struct for use with Shim (A HyperLedger included go file used for get/put state
//				and other HyperLedger functions)
//==============================================================================================================================
type SimpleChaincode struct {
}

//==============================================================================================================================
//	Account - Defines the structure for an account object. JSON on right tells it what JSON fields to map to
//			  that element when reading a JSON object into the struct e.g. JSON currency -> Struct Currency
//==============================================================================================================================

type Actor struct {
	ActorId    string `json:"actorid"`
	ActorName  string `json:"actorname"`
	Committed  string  `json:"committed"`
	Reimbursed string `json:"reimbursed"`
	Awarded    string `json:"awarded"`
	Spent      string `json:"spent"`
	Received   string `json:"received"`
	Delegated  string `json:"delegated"`
}

//award (award id, award name, award status, amount_requested, parent award id(-1))
type Award struct {
	AwardId        string `json:"awardid"`
	AwardName      string `json:"awardname"`
	AwardStatus    string `json:"awardstatus"`
	AwardRequested int `json:"awardrequested"`
	ParentAwardId  string `json:"parentawardid"`
}

//award amount (award_amount_id, award id, award amount, grantor id)
type AwardAmount struct {
	AwardAmountId string `json:"awardamountid"`
	AwardId       string `json:"awardid"`
	AwardAmount   int `json:"awardamount"`
	GrantorId     string `json:"grantorid"`
}

//award parties (award party id, award id, role type, account id)
type AwardParties struct {
	AwardPartyId string `json:"awardpartyid"`
	AwardId      string `json:"awardid"`
	RoleType     string `json:"roletype"`
	AccountId    string `json:"accountid"`
}

//reimbursement (reimbursement id, status, award id, amount)
type Reimbursement struct {
	ReimbursementId string `json:"reimbursementid"`
	Status          string `json:"status"`
	AwardId         string `json:"awardid"`
	Amount          int `json:"amount"`
	FromUser        string `json:"fromuser"`
	ToUser          string `json:"touser"`
}

//expenditure (expenditure id, amount, project id, date, type, reimbursement id)
type Expenditure struct {
	ExpenditureId   string `json:"expenditureid"`
	Amount          int `json:"amount"`
	ProjectId       string `json:"projectid"`
	Date            string `json:"date"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	ReimbursementId string `json:"reimbursementid"`
	FromUser        string `json:"fromuser"`
	ToUser          string `json:"touser"`
}

var accountIndexStr = "_accountindex" // Define an index variable to track all the accounts stored in the world state

// ============================================================================================================================
//  Main - main - Starts up the chaincode
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// SetUp Function - Called after the user deploys the chain code, before demo
// Function: create 4 actors, update AwardParty struct, update Award struct
// Call init_account, CreateAward, RequestAward
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) SetUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// CommitFund Function - Called when the grantor decided to commit the award
// Function: update AwardAmount struct, update Actor struct (transfer balance), update Award struct (status changed),
// update AwardParty struct (add grantor role)
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) CommitFund(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// ReleaseFund Function - Called when the grantor approves the reimbursement
// Function: update Expenditure struct (status), update Reimbursement struct (status), update Actor struct (transfer balance),
// update Award struct (status changed)
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) ReleaseFund(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// RequestReimbursement Function - Called when the grantee approves the reimbursement
// Function: update Expenditure struct (status), update Reimbursement struct (status)
// Needs from user and to user
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) RequestReimbursement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// CreateAward Function - Called when the grantee creates an award
// Function: update Award struct, update AwardParty struct
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) CreateAward(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// RequestAward Function - Called when the grantee request funding of an award
// Function: update Award struct (status)
// Necessary?
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) RequestAward(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// AllocateSubAward Function - Called when the grantee commits sub-award to sub-grantee
// Function: update Award struct (create a new sub-award, parent award id needs to be passed), update Account struct (balance transfer),
// update AwardParty struct (add sub-grantee), ?? AwardAmount() ??
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) AllocateSubAward(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// Spend Function - Called when the grantee or sub-grantee has an expenditure
// Function: update Expenditure struct (create a new one), update Account struct (balance transfer),
// Need from user and to user
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) Spend(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// Init Function - Called when the user deploys the chaincode
// stub -- name/alias of ChaincodeStubInterface
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting a single integer")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for testing the blockchain network")
	}

	// Write the state to the ledger, test the network
	err = stub.PutState("test_key", []byte(strconv.Itoa(Aval)))
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the account index
	err = stub.PutState(accountIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Invoke - Called on chaincode invoke. Takes a function name passed and calls that function. Converts some
//		    initial arguments passed to other things for use in the called function.
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "delete" {
		return t.Delete(stub, args)
	} else if function == "write" {
		return t.Write(stub, args)
	} else if function == "init_actor" {
		return t.init_actor(stub, args)
	} else if function == "CommitFund" {
		return t.CommitFund(stub, args)
	}
	//else if function == "transfer_balance" {
	//	return t.transfer_balance(stub, args)
	//}

	return nil, errors.New("Received unknown function invocation: " + function)
}

// ============================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "read" {
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query " + function)
}

// ============================================================================================================================
// Read - read a variable from chaincode world state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

// ============================================================================================================================
// Delete - remove a key/value pair from the world state
// ============================================================================================================================
func (t *SimpleChaincode) Delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	name := args[0]
	err := stub.DelState(name) //remove the key from chaincode state
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	//get the account index
	accountsAsBytes, err := stub.GetState(accountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get account index")
	}
	var accountIndex []string
	json.Unmarshal(accountsAsBytes, &accountIndex)

	//remove account from index
	for i, val := range accountIndex {
		if val == name { //find the correct account
			accountIndex = append(accountIndex[:i], accountIndex[i+1:]...) //remove it
			break
		}
	}
	jsonAsBytes, _ := json.Marshal(accountIndex) //save the new index
	err = stub.PutState(accountIndexStr, jsonAsBytes)
	return nil, nil
}

// ============================================================================================================================
// Write - directly write a variable into chaincode world state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string
	var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]
	value = args[1]
	err = stub.PutState(name, []byte(value))
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init account - create a new account, store into chaincode world state, and then append the account index
// ============================================================================================================================
func (t *SimpleChaincode) init_actor(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//       0        1      2..
	// "accountid", "name",  ...

	if len(args) != 8 {
		return nil, errors.New("Incorrect number of arguments. Expecting 8")
	}

	//input sanitation
	fmt.Println("- start init acount")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return nil, errors.New("6th argument must be a non-empty string")
	}
	if len(args[6]) <= 0 {
		return nil, errors.New("7th argument must be a non-empty string")
	}
	if len(args[7]) <= 0 {
		return nil, errors.New("8th argument must be a non-empty string")
	}

	actorId := args[0]

	actorName := strings.ToLower(args[1])

	committed, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	reimbursed, err := strconv.ParseFloat(args[3], 64)
	if err != nil {
		return nil, errors.New("4th argument must be a numeric string")
	}
	awarded, err := strconv.ParseFloat(args[4], 64)
	if err != nil {
		return nil, errors.New("5th argument must be a numeric string")
	}
	spent, err := strconv.ParseFloat(args[5], 64)
	if err != nil {
		return nil, errors.New("6th argument must be a numeric string")
	}
	received, err := strconv.ParseFloat(args[6], 64)
	if err != nil {
		return nil, errors.New("7th argument must be a numeric string")
	}
	delegated, err := strconv.ParseFloat(args[7], 64)
	if err != nil {
		return nil, errors.New("8th argument must be a numeric string")
	}

	//check if account already exists
	accountAsBytes, err := stub.GetState(actorId)
	if err != nil {
		return nil, errors.New("Failed to get account number")
	}

	res := Actor{}
	json.Unmarshal(accountAsBytes, &res)
	if res.ActorId == actorId {
		return nil, errors.New("This account arleady exists")
	}
	committedStr := strconv.FormatFloat(committed, 'f', -1, 64)
	reimbursedStr := strconv.FormatFloat(reimbursed, 'f', -1, 64)
	awardedStr := strconv.FormatFloat(awarded, 'f', -1, 64)
	spentStr := strconv.FormatFloat(spent, 'f', -1, 64)
	receivedStr := strconv.FormatFloat(received, 'f', -1, 64)
	delegatedStr := strconv.FormatFloat(delegated, 'f', -1, 64)

	//newActor := Actor{}
	//newActor.ActorId = actorId
	//newActor.ActorName = actorName
	//newActor.Balance = balance

	//build the account json string 
	str := `{"actorid": "` + actorId + `", "actorName": "` + actorName + `", "comitted": "` + committedStr + `", "reimbursed": "` + reimbursedStr + `", "awarded": "` + awardedStr + `", "spent": "` + spentStr + `", "received": "` + receivedStr + `", "delegated": "` + delegatedStr + `"}`
	//jsonAsBytesActor, _ := json.Marshal(newActor)
	err = stub.PutState(actorId, []byte(str))
	if err != nil {
		return nil, err
	}

	//get the account index
	accountsAsBytes, err := stub.GetState(accountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get account index")
	}
	var accountIndex []string
	json.Unmarshal(accountsAsBytes, &accountIndex)

	//append the index 
	accountIndex = append(accountIndex, actorId)
	jsonAsBytes, _ := json.Marshal(accountIndex)
	err = stub.PutState(accountIndexStr, jsonAsBytes)

	return nil, nil
}

// ============================================================================================================================
// Transfer Balance - Create a transaction between two accounts, transfer a certain amount of balance
// ============================================================================================================================
/*func (t *SimpleChaincode) transfer_balance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	//       0           1         2
	// "accountA", "accountB", "100.20"

	var err error
	var newAmountA, newAmountB float64

	if len(args) < 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	amount,err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	accountAAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get the first account")
	}
	resA := Account{}
	json.Unmarshal(accountAAsBytes, &resA)								
	
	accountBAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get the second account")
	}
	resB := Account{}
	json.Unmarshal(accountBAsBytes, &resB)											
	
	BalanceA,err := strconv.ParseFloat(resA.Balance, 64)
	if err != nil {
		return nil, err
	}
	BalanceB,err := strconv.ParseFloat(resB.Balance, 64)
	if err != nil {
		return nil, err
	}

	//Check if accountA has enough balance to transact or not
	if (BalanceA - amount) < 0 {
		return nil, errors.New(args[0] + " doesn't have enough balance to complete transaction")
	}

	newAmountA = BalanceA - amount
	newAmountB =  BalanceB + amount
	newAmountStrA := strconv.FormatFloat(newAmountA, 'E', -1, 64)
	newAmountStrB := strconv.FormatFloat(newAmountB, 'E', -1, 64)

	resA.Balance = newAmountStrA
	resB.Balance = newAmountStrB

	jsonAAsBytes, _ := json.Marshal(resA)
	err = stub.PutState(args[0], jsonAAsBytes)								
	if err != nil {
		return nil, err
	}

	jsonBAsBytes, _ := json.Marshal(resB)
	err = stub.PutState(args[1], jsonBAsBytes)								
	if err != nil {
		return nil, err
	}
	
	return nil, nil
}
*/
