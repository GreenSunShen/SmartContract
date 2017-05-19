package main

import (
	"fmt"
	"strconv"
	"encoding/json"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"time"
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

//award parties (award party id, award id, role type, account id)
type AwardParties struct {
	GrantorId string `json:"grantorid"`
	GranteeId string `json:"graneeid"`
	SubgranteeId string `json:"subgranteeid"`
	SupplierId string `json:"supplierId"`
}


//award (award id, award name, award status, amount_requested, parent award id(-1))
type Award struct {
	AwardId        string `json:"awardid"`
	AwardName      string `json:"awardname"`
	AwardStatus    string `json:"awardstatus"`
	AwardRequested string `json:"awardrequested"`
	Party AwardParties `json:"awardparties"`
	Expenses []Expenditure `json:"expenditures"`
	Reimburses []Reimbursement `json:"reimburses"`
}

//award amount (award_amount_id, award id, award amount, grantor id)
type AwardAmount struct {
	AwardAmountId string `json:"awardamountid"`
	AwardId       string `json:"awardid"`
	AwardAmount   string `json:"awardamount"`
	GrantorId     string `json:"grantorid"`
}



//reimbursement (reimbursement id, status, award id, amount)
type Reimbursement struct {
	ReimbursementId string `json:"reimbursementid"`
	Amount          string `json:"amount"`
	FromActor        string `json:"fromuser"`
	ToActor          string `json:"touser"`
	Date            string `json:"date"`
	ExpenditureId string `json:"expenditureid"`
}

//expenditure (expenditure id, amount, project id, date, type, reimbursement id)
type Expenditure struct {
	ExpenditureId   string `json:"expenditureid"`
	Amount          string `json:"amount"`
	Date            string `json:"date"`
	Type            string `json:"type"`
	Status          string `json:"status"`
	FromActor        string `json:"fromuser"`
	ToActor          string `json:"touser"`
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
// Call init_account, CreateAward, <not needed>--RequestAward
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) SetUp(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
//------------------------------create roles----------------------------------------------
	// grantor info
	var act1  = make([]string, 8,8)
	act1[0] = "ACT-101"                //ActorId
	act1[1] = "PPM Foundation"     //ActorName
	act1[2] = "125000"             //Comiitted
	act1[3] = "55000"              //Reimbursed
	act1[4] = "-1"                   //Awarded
	act1[5] = "-1"                   //Spent
	act1[6] = "-1"                   //Received
	act1[7] = "-1"                   //Delegated

	// grantee info
	var act2  = make([]string, 8,8)
	act2[0] = "ACT-102"                     //ActorId
	act2[1] = "Stanford University"      //ActorName
	act2[2] = "-1"                        //Comiitted
	act2[3] = "-1"                         //Reimbursed
	act2[4] = "125000"                 //Awarded
	act2[5] = "23000"                  //Spent
	act2[6] = "55000"                 //Received
	act2[7] = "45000"                 //Delegated

	// sub-grantee info
	var act3  = make([]string, 8,8)
	act3[0] = "ACT-103"                //ActorId
	act3[1] = "John Hopkins University"     //ActorName
	act3[2] = "-1"               //Comiitted
	act3[3] = "-1"               //Reimbursed
	act3[4] = "45000"          //Awarded
	act3[5] = "12000"          //Spent
	act3[6] = "25000"          //Received
	act3[7] = "-1"               //Delegated

	// Supplier info -- shows spending form all the grantees and sub-grantees
	var act4  = make([]string, 8,8)
	act4[0] = "ACT-104"                //ActorId
	act4[1] = "Dixon consulting"     //ActorName
	act4[2] = "-1"               //Comiitted
	act4[3] = "-1"               //Reimbursed
	act4[4] = "-1"          //Awarded
	act4[5] = "-1"          //Spent
	act4[6] = "35000"          //Received
	act4[7] = "-1"               //Delegated


	t.init_actor(stub, act1)
	t.init_actor(stub, act2)
	t.init_actor(stub, act3)
	t.init_actor(stub, act4)

	//----------------create expenses----------------------------------------------------------
	// Expense
	var exp1  = make([]string, 7,7)
	exp1[0] = "EXP-201"                //ExpenditureId
	exp1[1] = "3000"     //Amount
	exp1[2] = "05-02-2017"             //Date
	exp1[3] = "Travel"              //Type
	exp1[4] = "Approved"                   //Status
	exp1[5] = "ACT-102"                   //FromActor --Grantee spending
	exp1[6] = "ACT-104"                   //ToActor   --Supplier receiving the spending


	var exp2  = make([]string, 7,7)
	exp2[0] = "EXP-202"                //ExpenditureId
	exp2[1] = "8000"     //Amount
	exp2[2] = "05-03-2017"             //Date
	exp2[3] = "Equipment"              //Type
	exp2[4] = "Pending"                   //Status
	exp2[5] = "ACT-102"                   //FromActor --Grantee spending
	exp2[6] = "ACT-104"                   //ToActor   --Supplier receiving the spending

	var exp3 = make([]string, 7,7)
	exp3[0] = "EXP-203"                //ExpenditureId
	exp3[1] = "4000"     //Amount
	exp3[2] = "05-04-2017"             //Date
	exp3[3] = "Training"              //Type
	exp3[4] = "Approved"                   //Status
	exp3[5] = "ACT-102"                   //FromActor --Grantee spending
	exp3[6] = "ACT-104"

	var exp4 = make([]string, 7,7)
	exp4[0] = "EXP-204"                //ExpenditureId
	exp4[1] = "3000"     //Amount
	exp4[2] = "05-09-2017"             //Date
	exp4[3] = "Software License"              //Type
	exp4[4] = "Approved"                   //Status
	exp4[5] = "ACT-102"                   //FromActor --Grantee spending
	exp4[6] = "ACT-104"

	var exp5 = make([]string, 7,7)
	exp5[0] = "EXP-205"                //ExpenditureId
	exp5[1] = "5000"     //Amount
	exp5[2] = "05-11-2017"             //Date
	exp5[3] = "Specimens"              //Type
	exp5[4] = "Approved"                   //Status
	exp5[5] = "ACT-102"                   //FromActor --Grantee spending
	exp5[6] = "ACT-104"

	var exp6 = make([]string, 7,7)
	exp6[0] = "EXP-206"                //ExpenditureId
	exp6[1] = "2000"     //Amount
	exp6[2] = "05-12-2017"             //Date
	exp6[3] = "Consultancy"              //Type
	exp6[4] = "Approved"                   //Status
	exp6[5] = "ACT-103"                   //FromActor --Grantee spending
	exp6[6] = "ACT-104"

	var exp7 = make([]string, 7,7)
	exp7[0] = "EXP-207"                //ExpenditureId
	exp7[1] = "7500"     //Amount
	exp7[2] = "05-15-2017"             //Date
	exp7[3] = "Equipment"              //Type
	exp7[4] = "Pending"                   //Status
	exp7[5] = "ACT-103"                   //FromActor --Grantee spending
	exp7[6] = "ACT-104"

	var exp8 = make([]string, 7,7)
	exp8[0] = "EXP-208"                //ExpenditureId
	exp8[1] = "1000"     //Amount
	exp8[2] = "05-16-2017"             //Date
	exp8[3] = "Travel"              //Type
	exp8[4] = "Approved"                   //Status
	exp8[5] = "ACT-103"                   //FromActor --Grantee spending
	exp8[6] = "ACT-104"

	var exp9 = make([]string, 7,7)
	exp9[0] = "EXP-209"                //ExpenditureId
	exp9[1] = "1500"     //Amount
	exp9[2] = "05-19-2017"             //Date
	exp9[3] = "Training"              //Type
	exp9[4] = "Approved"                   //Status
	exp9[5] = "ACT-103"                   //FromActor --Grantee spending
	exp9[6] = "ACT-104"

	t.init_expenditure(stub, exp1)
	t.init_expenditure(stub, exp2)
	t.init_expenditure(stub, exp3)
	t.init_expenditure(stub, exp4)
	t.init_expenditure(stub, exp5)
	t.init_expenditure(stub, exp6)
	t.init_expenditure(stub, exp7)
	t.init_expenditure(stub, exp8)
	t.init_expenditure(stub, exp9)
////----------------------------create award -------------------------------------------------
//	var award1 []string{}
//
//	t.CreateAward(stub, award1)


	return nil, nil
}


func (t *SimpleChaincode) init_expenditure(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//       0        1      2..
	// "accountid", "name",  ...

	if len(args) != 7 {
		return nil, errors.New("Incorrect number of arguments. Expecting 7")
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

	/*
	     exp9[0] = "EXP-209"                //ExpenditureId
		exp9[1] = "1500"     //Amount
		exp9[2] = "05-19-2017"             //Date
		exp9[3] = "Training"              //Type
		exp9[4] = "Approved"                   //Status
	    exp9[5] = "ACT-103"                   //FromActor --Grantee spending
	    exp9[6] = "ACT-104"
	*/

	expId := args[0]

	expAmount, err := strconv.ParseFloat(args[1], 64)   //strings.ToLower(args[1])

	if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}

	expDate := args[2]
	expType := args[3]
	expStatus := args[4]
	fromActor := args[5]
	toActor := args[6]



	//check if account already exists
	accountAsBytes, err := stub.GetState(expId)
	if err != nil {
		return nil, errors.New("Failed to get expenditure id")
	}

	exp := Expenditure{}
	json.Unmarshal(accountAsBytes, &exp)
	if exp.ExpenditureId == expId {
		return nil, errors.New("This expenditure arleady exists")
	}


      expAmountStr := strconv.FormatFloat(expAmount, 'f', -1, 64)

	//newActor := Actor{}
	//newActor.ActorId = actorId
	//newActor.ActorName = actorName
	//newActor.Balance = balance

	//	   expDate := arg[2]
	///  expType := arg[3]
	/// expStatus := arg[4]
	//fromActor := arg[5]
	//toActor := arg[6]


	//build the expenditure json string
	str := `{"expenditureid": "` + expId + `", "amount": "` + expAmountStr + `", "date": "` + expDate + `", "type": "` + expType + `", "status": "` + expStatus + `", "fromactor": "` + fromActor + `", "toactor": "` + toActor + `"}`
	//jsonAsBytesActor, _ := json.Marshal(newActor)
	err = stub.PutState(expId, []byte(str))
	if err != nil {
		return nil, err
	}

	return nil, nil
}





// ============================================================================================================================
// ReleaseFund Function - Called when the grantor approves the reimbursement
// Function: update Expenditure struct (status), update Reimbursement struct (status), update Actor struct (transfer balance),
// update Award struct (status changed)
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) ReleaseFund(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//arg[0] actor id
	//arg[1] award id
	//arg[2] ... exp id
	//var expenseIds []string
	//
	//for i := 2; i < len(args); i++{
	//	expenseIds = append(expenseIds, args[i])
	//}
	//
	//expenseNumber := len(expenseIds)






	return nil, nil
}



// ============================================================================================================================
// CreateAward Function - Called when the grantee creates an award
// Function: update Award struct, update AwardParty, update AwardAmount struct
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) CreateAward(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	return nil, nil
}

// ============================================================================================================================
// Query Function - Called when query all expenditure
// Function: query all the expenditures of this award
// Query
// ============================================================================================================================
func (t *SimpleChaincode) QueryAllExpenses(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// 0 award id
	if len(args) != 1{
		return nil, errors.New("Number of arguments is not correct. Expected 1")
	}
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}

	awardAsBytes, err := stub.GetState(args[0])
	if err != nil{
		return nil, errors.New("Failed to get award.")
	}

	//awardPlain := Award{}
	//json.Unmarshal(awardAsBytes, &awardPlain)
	//
	//expenses := awardPlain.Expenses
	//awardAsBytes, _ = json.Marshal(expenses)


	awardAsBytes = []byte("string")

	return awardAsBytes, nil
}


//2. pending expenditure
//3. all expenditure , reimbursement
//4. all actor balance




// ============================================================================================================================
// Spend Function - Called when the grantee or sub-grantee has an expenditure
// Function: update Expenditure struct (create a new one), update Account struct (balance transfer),
// Need from user and to user
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) Spend(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//     0           1        2        3         4
	// "from id"   "to id"   "amount"  "type"  "award id"

	//get from actor
	accountAAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get the first account")
	}
	resA := Actor{}
	json.Unmarshal(accountAAsBytes, &resA)

	accountBAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get the second account")
	}

	//get to actor
	resB := Actor{}
	json.Unmarshal(accountBAsBytes, &resB)


	//get amount
	amount, err := strconv.ParseFloat(args[2], 64)
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	//get award
	awardAsbytes, err := stub.GetState(args[4])
	awardPlain := Award{}
	json.Unmarshal(awardAsbytes, &awardPlain)

	//get date
	current_time := time.Now().Local()

	// populate expenditure
	exp := Expenditure{}
	exp.Amount = strconv.FormatFloat(amount, 'f', -1, 64)
	exp.Date = current_time.String()
	exp.Type = args[3]
	exp.ExpenditureId = "ex1"
	exp.FromActor = resA.ActorId
	exp.ToActor = resB.ActorId

	// compare with threshold to determine status
	if amount > 6000{
		exp.Status = "Pending"
	}else{
		exp.Status = "Auto"
	}

	// assign the expenditure to the award
	expenses := awardPlain.Expenses
	expenses = append(expenses, exp)
	awardPlain.Expenses = expenses

	//save to world state
	awardAsbytes, _ = json.Marshal(awardPlain) //save the new index
	err = stub.PutState(args[4], awardAsbytes)
	if err != nil {
		return nil, err
	}


	t.transfer_balance(stub, []string{args[0], args[1], args[2], "spend"})

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
	} else if function == "initactor" {
		return t.init_actor(stub, args)
	} else if function == "setup"{
		return t.SetUp(stub, args)
	} else if function == "spend"{
		return t.Spend(stub, args)
	} else if function == "releasefund"{
		return t.ReleaseFund(stub, args)
	}
	//else if function == "transferbalance" {
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
func (t *SimpleChaincode) transfer_balance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	//     0         1         2         3
	// "actorA", "actorB", "100.20"  "function"
	var err error
	var newAmountA, newAmountB float64

	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	amount,err := strconv.ParseFloat(args[2], 64)

	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}

	accountAAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get the first account")
	}
	resA := Actor{}
	json.Unmarshal(accountAAsBytes, &resA)								
	
	accountBAsBytes, err := stub.GetState(args[1])
	if err != nil {
		return nil, errors.New("Failed to get the second account")
	}
	resB := Actor{}
	json.Unmarshal(accountBAsBytes, &resB)

	switch args[3] {
		case "spend" :
			fmt.Println("INSIDE CASE SPEND==========================")
			AwardA, err := strconv.ParseFloat(resA.Awarded, 64)
			if err != nil {
				return nil, err
			}
			BalanceA,err := strconv.ParseFloat(resA.Spent, 64)
			if err != nil {
				return nil, err
			}
			BalanceB,err := strconv.ParseFloat(resB.Received, 64)
			if err != nil {
				return nil, err
			}
			//Check if accountA has enough balance to transact or not
			if ( AwardA - amount) < 0 {
				return nil, errors.New(args[0] + " doesn't have enough balance to complete transaction")
			}

			newAmountA = BalanceA + amount
			newAmountB =  BalanceB + amount
			newAmountStrA := strconv.FormatFloat(newAmountA, 'f', -1, 64)
			newAmountStrB := strconv.FormatFloat(newAmountB, 'f', -1, 64)

			resA.Spent = newAmountStrA
			resB.Received = newAmountStrB
	default:
	}
	/*
	BalanceA,err := strconv.ParseFloat(resA.b, 64)
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
*/
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

