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

//reimbursement (reimbursement id, status, award id, amount)
type Reimbursement struct {
	ReimbursementId string `json:"reimbursementid"`
	Amount          string `json:"amount"`
	FromActor       string `json:"fromactor"`
	ToActor         string `json:"toactor"`
	Date            string `json:"date"`
	ExpenditureId   string `json:"expenditureid"`
}

//expenditure (expenditure id, amount, project id, date, type, reimbursement id)
type Expenditure struct {
	ExpenditureId string `json:"expenditureid"`
	Amount        string `json:"amount"`
	Date          string `json:"date"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	FromActor     string `json:"fromactor"`
	ToActor       string `json:"toactor"`
}

var accountIndexStr = "_accountindex" // Define an index variable to track all the actors stored in the world state
var expIndexStr = "_expindex"         // Define an index variable to track all the expenditures stored in the world state
var reimbIndexStr = "_reimbindex"     // Define an index variable to track all the reimbursements stored in the world state
var expNumber int = 0
var reimbNumber int = 0

var act1 = make([]string, 8, 8)
var act2 = make([]string, 8, 8)
var act3 = make([]string, 8, 8)
var act4 = make([]string, 8, 8)
var exp1 = make([]string, 7, 7)
var exp2 = make([]string, 7, 7)
var exp3 = make([]string, 7, 7)
var exp4 = make([]string, 7, 7)
var exp5 = make([]string, 7, 7)
var exp6 = make([]string, 7, 7)
var exp7 = make([]string, 7, 7)
var exp8 = make([]string, 7, 7)
var exp9 = make([]string, 7, 7)
var rem1 = make([]string, 6, 6)
var rem2 = make([]string, 6, 6)
var rem3 = make([]string, 6, 6)
var rem4 = make([]string, 6, 6)
var rem5 = make([]string, 6, 6)
var rem6 = make([]string, 6, 6)
var rem7 = make([]string, 6, 6)

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
	act1[0] = "ACT-101"        //ActorId
	act1[1] = "PPM Foundation" //ActorName
	act1[2] = "125000"         //Comiitted
	act1[3] = "19500"          //Reimbursed
	act1[4] = "-1"             //Awarded
	act1[5] = "-1"             //Spent
	act1[6] = "-1"             //Received
	act1[7] = "-1"             //Delegated

	// grantee info
	act2[0] = "ACT-102"             //ActorId
	act2[1] = "Stanford University" //ActorName
	act2[2] = "-1"                  //Comiitted
	act2[3] = "-1"                  //Reimbursed
	act2[4] = "125000"              //Awarded
	act2[5] = "23000"               //Spent
	act2[6] = "10000"               //Received
	act2[7] = "45000"               //Delegated

	// sub-grantee info
	act3[0] = "ACT-103"                 //ActorId
	act3[1] = "John Hopkins University" //ActorName
	act3[2] = "-1"                      //Comiitted
	act3[3] = "-1"                      //Reimbursed
	act3[4] = "45000"                   //Awarded
	act3[5] = "12000"                   //Spent
	act3[6] = "9500"                    //Received
	act3[7] = "-1"                      //Delegated

	// Supplier info -- shows spending form all the grantees and sub-grantees
	act4[0] = "ACT-104"          //ActorId
	act4[1] = "Dixon consulting" //ActorName
	act4[2] = "-1"               //Comiitted
	act4[3] = "-1"               //Reimbursed
	act4[4] = "-1"               //Awarded
	act4[5] = "-1"               //Spent
	act4[6] = "35000"            //Received
	act4[7] = "-1"               //Delegated

	t.Init_actor(stub, act1)
	t.Init_actor(stub, act2)
	t.Init_actor(stub, act3)
	t.Init_actor(stub, act4)

	//----------------create expenses----------------------------------------------------------
	// Expense

	exp1[0] = "EXP-201"    //ExpenditureId
	exp1[1] = "3000"       //Amount
	exp1[2] = "05-02-2017" //Date
	exp1[3] = "Travel"     //Type
	exp1[4] = "Approved"   //Status
	exp1[5] = "ACT-102"    //FromActor --Grantee spending
	exp1[6] = "ACT-104"    //ToActor   --Supplier receiving the spending

	exp2[0] = "EXP-202"    //ExpenditureId
	exp2[1] = "8000"       //Amount
	exp2[2] = "05-03-2017" //Date
	exp2[3] = "Equipment"  //Type
	exp2[4] = "Pending"    //Status
	exp2[5] = "ACT-102"    //FromActor --Grantee spending
	exp2[6] = "ACT-104"    //ToActor   --Supplier receiving the spending

	exp3[0] = "EXP-203"    //ExpenditureId
	exp3[1] = "4000"       //Amount
	exp3[2] = "05-04-2017" //Date
	exp3[3] = "Training"   //Type
	exp3[4] = "Approved"   //Status
	exp3[5] = "ACT-102"    //FromActor --Grantee spending
	exp3[6] = "ACT-104"

	exp4[0] = "EXP-204"          //ExpenditureId
	exp4[1] = "3000"             //Amount
	exp4[2] = "05-09-2017"       //Date
	exp4[3] = "Software License" //Type
	exp4[4] = "Approved"         //Status
	exp4[5] = "ACT-102"          //FromActor --Grantee spending
	exp4[6] = "ACT-104"

	exp5[0] = "EXP-205"    //ExpenditureId
	exp5[1] = "5000"       //Amount
	exp5[2] = "05-11-2017" //Date
	exp5[3] = "Specimens"  //Type
	exp5[4] = "Approved"   //Status
	exp5[5] = "ACT-102"    //FromActor --Grantee spending
	exp5[6] = "ACT-104"

	exp6[0] = "EXP-206"     //ExpenditureId
	exp6[1] = "2000"        //Amount
	exp6[2] = "05-12-2017"  //Date
	exp6[3] = "Consultancy" //Type
	exp6[4] = "Approved"    //Status
	exp6[5] = "ACT-103"     //FromActor --Grantee spending
	exp6[6] = "ACT-104"

	exp7[0] = "EXP-207"    //ExpenditureId
	exp7[1] = "7500"       //Amount
	exp7[2] = "05-15-2017" //Date
	exp7[3] = "Equipment"  //Type
	exp7[4] = "Pending"    //Status
	exp7[5] = "ACT-103"    //FromActor --Grantee spending
	exp7[6] = "ACT-104"

	exp8[0] = "EXP-208"    //ExpenditureId
	exp8[1] = "1000"       //Amount
	exp8[2] = "05-16-2017" //Date
	exp8[3] = "Travel"     //Type
	exp8[4] = "Approved"   //Status
	exp8[5] = "ACT-103"    //FromActor --Grantee spending
	exp8[6] = "ACT-104"

	exp9[0] = "EXP-209"    //ExpenditureId
	exp9[1] = "1500"       //Amount
	exp9[2] = "05-19-2017" //Date
	exp9[3] = "Training"   //Type
	exp9[4] = "Approved"   //Status
	exp9[5] = "ACT-103"    //FromActor --Grantee spending
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

	//---------------------create reimbursement------------------------------------------------
	// Reimbursement

	rem1[0] = "REM-301"    //ReimbursementId
	rem1[1] = "3000"       //Amount
	rem1[2] = "ACT-101"    //FromActor
	rem1[3] = "ACT-102"    //ToActor
	rem1[4] = "05-02-2017" //Date
	rem1[5] = "EXP-201"    //ExpenditureId

	rem2[0] = "REM-302"    //ReimbursementId
	rem2[1] = "4000"       //Amount
	rem2[2] = "ACT-101"    //FromActor
	rem2[3] = "ACT-102"    //ToActor
	rem2[4] = "05-04-2017" //Date
	rem2[5] = "EXP-203"    //ExpenditureId

	rem3[0] = "REM-303"    //ReimbursementId
	rem3[1] = "3000"       //Amount
	rem3[2] = "ACT-101"    //FromActor
	rem3[3] = "ACT-102"    //ToActor
	rem3[4] = "05-09-2017" //Date
	rem3[5] = "EXP-204"    //ExpenditureId

	rem4[0] = "REM-304"    //ReimbursementId
	rem4[1] = "5000"       //Amount
	rem4[2] = "ACT-101"    //FromActor
	rem4[3] = "ACT-102"    //ToActor
	rem4[4] = "05-11-2017" //Date
	rem4[5] = "EXP-205"    //ExpenditureId

	rem5[0] = "REM-305"    //ReimbursementId
	rem5[1] = "2000"       //Amount
	rem5[2] = "ACT-102"    //FromActor
	rem5[3] = "ACT-103"    //ToActor
	rem5[4] = "05-12-2017" //Date
	rem5[5] = "EXP-206"    //ExpenditureId

	rem6[0] = "REM-306"    //ReimbursementId
	rem6[1] = "1000"       //Amount
	rem6[2] = "ACT-102"    //FromActor
	rem6[3] = "ACT-103"    //ToActor
	rem6[4] = "05-16-2017" //Date
	rem6[5] = "EXP-208"    //ExpenditureId

	rem7[0] = "REM-307"    //ReimbursementId
	rem7[1] = "1500"       //Amount
	rem7[2] = "ACT-102"    //FromActor
	rem7[3] = "ACT-103"    //ToActor
	rem7[4] = "05-19-2017" //Date
	rem7[5] = "EXP-209"    //ExpenditureId

	t.init_reimbursement(stub, rem1)
	t.init_reimbursement(stub, rem2)
	t.init_reimbursement(stub, rem3)
	t.init_reimbursement(stub, rem4)
	t.init_reimbursement(stub, rem5)
	t.init_reimbursement(stub, rem6)
	t.init_reimbursement(stub, rem7)

	return nil, nil
}

func (t *SimpleChaincode) init_reimbursement(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
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

	remId := args[0]

	remAmount, err := strconv.ParseFloat(args[1], 64) //strings.ToLower(args[1])

	if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}

	remFromActor := args[2]
	remToActor := args[3]
	remDate := args[4]
	remExpId := args[5]

	//check if account already exists
	accountAsBytes, err := stub.GetState(remId)
	if err != nil {
		return nil, errors.New("Failed to get expenditure id")
	}

	rem := Reimbursement{}
	json.Unmarshal(accountAsBytes, &rem)
	if rem.ReimbursementId == remId {
		return nil, errors.New("This reimbursement arleady exists")
	}

	remAmountStr := strconv.FormatFloat(remAmount, 'f', -1, 64)

	//build the expenditure json string
	str := `{"reimbursementid": "` + remId + `", "amount": "` + remAmountStr + `", "fromactor": "` + remFromActor + `", "toActor": "` + remToActor + `", "date": "` + remDate + `", "expenditureid": "` + remExpId + `"}`
	//jsonAsBytesActor, _ := json.Marshal(newActor)
	err = stub.PutState(remId, []byte(str))
	if err != nil {
		return nil, err
	}

	//get the reimb index
	reimbAsBytes, err := stub.GetState(reimbIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get reimbursement index")
	}
	var reimbIndex []string
	json.Unmarshal(reimbAsBytes, &reimbIndex)

	//append the index
	reimbIndex = append(reimbIndex, remId)
	jsonAsBytes, _ := json.Marshal(reimbIndex)
	err = stub.PutState(reimbIndexStr, jsonAsBytes)

	reimbNumber++

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

	expAmount, err := strconv.ParseFloat(args[1], 64) //strings.ToLower(args[1])

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

	//get the exp index
	expsAsBytes, err := stub.GetState(expIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get expenditure index")
	}
	var expIndex []string
	json.Unmarshal(expsAsBytes, &expIndex)

	//append the index
	expIndex = append(expIndex, expId)
	jsonAsBytes, _ := json.Marshal(expIndex)
	err = stub.PutState(expIndexStr, jsonAsBytes)

	expNumber++

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
	//arg[1] ... exp id

	// get all expense ids
	var expenseIds []string
	for i := 1; i < len(args); i++ {
		expenseIds = append(expenseIds, args[i])
	}

	// get Actor A -- from actor
	fromActorAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get from actor")
	}
	fromActor := Actor{}
	json.Unmarshal(fromActorAsBytes, &fromActor)

	//get the exp index array
	expsIndexAsBytes, err := stub.GetState(expIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get expenditure index")
	}
	var expIndex []string
	json.Unmarshal(expsIndexAsBytes, &expIndex)

	//loop through to get exp.fromActor as toActor for reimbursement
	for i := 0; i < len(expenseIds); i++ {
		for j := 0; j < len(expIndex); j++ {
			if expIndex[j] == expenseIds[i] {
				// get to actor ID
				oneExpAsByte, err1 := stub.GetState(expIndex[j])
				if err1 != nil {
					return nil, errors.New("Failed to get expenditure")
				}
				oneExp := Expenditure{}
				json.Unmarshal(oneExpAsByte, &oneExp)

				//TODO CHECK IF THE EXP STATUS IS "Pending"

				// transfer balance
				t.Transfer_balance(stub, []string{args[0], oneExp.FromActor, oneExp.Amount, "fund"})

				// change exp status
				oneExp.Status = "Approved"
				oneExpAsByte, _ = json.Marshal(oneExp)
				err = stub.PutState(oneExp.ExpenditureId, oneExpAsByte)

				// create a new reimbursement
				var remid string = "REM-"
				ii := strconv.Itoa(reimbNumber + 301)
				remid += ii

				current_time := time.Now().Local()

				t.init_reimbursement(stub, []string{remid, oneExp.Amount, args[0], oneExp.FromActor, current_time.String(), expIndex[j]})

			}
		}
	}

	return nil, nil
}

// ============================================================================================================================
// Query Function - Called when query all expenditure
// Function: query all the expenditures of this award
// Query
// ============================================================================================================================
func (t *SimpleChaincode) QueryAllExpenses(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//args[0] = ""

	//get the exp index
	expsIndexAsBytes, err := stub.GetState(expIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get expenditure index")
	}
	var expIndex []string
	json.Unmarshal(expsIndexAsBytes, &expIndex)

	var expenses []Expenditure
	for i := 0; i < len(expIndex); i++ {
		expAsBytes, err := stub.GetState(expIndex[i])
		if err != nil {
			return nil, errors.New("Failed to get expenditure")
		}
		oneExpense := Expenditure{}
		json.Unmarshal(expAsBytes, &oneExpense)
		expenses = append(expenses, oneExpense)
	}

	expsAsBytes, _ := json.Marshal(expenses)

	//awardAsBytes = []byte("string")

	return expsAsBytes, nil
}

// ============================================================================================================================
// Query Function - Called when query pending expenditure
// Function: query all the pending expenditures of this award
// Query
// ============================================================================================================================
func (t *SimpleChaincode) QueryPendingExpenses(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	expsIndexAsBytes, err := stub.GetState(expIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get expenditure index")
	}
	var expIndex []string
	json.Unmarshal(expsIndexAsBytes, &expIndex)

	var expenses []Expenditure
	for i := 0; i < len(expIndex); i++ {
		expAsBytes, err := stub.GetState(expIndex[i])
		if err != nil {
			return nil, errors.New("Failed to get expenditure")
		}
		oneExpense := Expenditure{}
		json.Unmarshal(expAsBytes, &oneExpense)

		if oneExpense.Status == "Pending" {
			expenses = append(expenses, oneExpense)
		}
	}

	expsAsBytes, _ := json.Marshal(expenses)

	//awardAsBytes = []byte("string")

	return expsAsBytes, nil

}


// ============================================================================================================================
// Query Function - Called when query block chain diagram
// Function: query all the transactions of this award before certain date
// Query
// ============================================================================================================================
func (t *SimpleChaincode) QueryBlockChain(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//threshold date
	thresholdDate := time.Date(
		2017, 05, 14, 20, 34, 58, 651387237, time.UTC)

	// get all expenditure index
	expsIndexAsBytes, err := stub.GetState(expIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get expenditure index")
	}
	var expIndex []string
	json.Unmarshal(expsIndexAsBytes, &expIndex)
        var res []byte
	//var expenses []Expenditure
	//for i := 0; i < len(expIndex); i++ {
		expAsBytes, err := stub.GetState(expIndex[0])
		if err != nil {
			return nil, errors.New("Failed to get expenditure")
		}
		oneExpense := Expenditure{}
		json.Unmarshal(expAsBytes, &oneExpense)
		dateStr := oneExpense.Date
		date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, errors.New(dateStr)
	}
		diff := date.Sub(thresholdDate)
		if  diff > 0{
			//expenses = append(expenses, oneExpense)
			res = []byte("true")
		} else {
			res = []byte("false")
		}


	//}
	//expsAsBytes, _ := json.Marshal(expenses)
/*
	// get all reimbursement index
	reimbsIndexAsBytes, err := stub.GetState(reimbIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get reimbursement index")
	}
	var reimbIndex []string
	json.Unmarshal(reimbsIndexAsBytes, &reimbIndex)

	var reimbursements []Reimbursement
	for i := 0; i < len(reimbIndex); i++ {
		reimbAsBytes, err := stub.GetState(reimbIndex[i])
		if err != nil {
			return nil, errors.New("Failed to get reimbursement")
		}
		oneReimburse := Reimbursement{}
		json.Unmarshal(reimbAsBytes, &oneReimburse)
		dateStr := oneReimburse.Date
		date, err := time.Parse("2006-01-02", dateStr)
		diff := date.Sub(thresholdDate)
		if diff > 0{
			reimbursements = append(reimbursements, oneReimburse)
		}


	}

	reimbsAsBytes, _ := json.Marshal(reimbursements)
	expsAsBytes = append(expsAsBytes, reimbsAsBytes...)
*/

	//return expsAsBytes, nil
	return res, err
}


// ============================================================================================================================
// Query Function - Called when query actors' wallets
// Function: query the balance of all actors
// Query
// ============================================================================================================================
func (t *SimpleChaincode) QueryWallet(stub shim.ChaincodeStubInterface, args []string) ([]byte, error){
	actorIndexAsBytes, err := stub.GetState(accountIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get actor index")
	}
	var actorIndex []string
	json.Unmarshal(actorIndexAsBytes, &actorIndex)

	var resultAsBytes []byte

	for i := 0; i < len(actorIndex); i++{
		actorAsBytes, err := stub.GetState(actorIndex[i])
		if err != nil{
			return nil, errors.New("Failed to get actor.")
		}
		resultAsBytes = append(resultAsBytes, actorAsBytes...)
	}

	return resultAsBytes, nil
}



// ============================================================================================================================
// Spend Function - Called when the grantee or sub-grantee has an expenditure
// Function: update Expenditure struct (create a new one), update Account struct (balance transfer),
// Need from user and to user
// Invoke
// ============================================================================================================================
func (t *SimpleChaincode) Spend(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//     0           1        2        3
	// "from id"   "to id"   "amount"  "type"

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

	//get date
	current_time := time.Now().Local()

	// populate
	var expid string = "EXP-"
	ii := strconv.Itoa(expNumber + 201)
	expid += ii

	var expstatus string

	// compare with threshold to determine status
	if amount > 6000 {
		expstatus = "Pending"
	} else {
		expstatus = "Approved"
	}

	t.init_expenditure(stub, []string{expid, strconv.FormatFloat(amount, 'f', -1, 64), current_time.String(), args[3], expstatus, resA.ActorId, resB.ActorId})
	//exp := Expenditure{}
	//exp.Amount = strconv.FormatFloat(amount, 'f', -1, 64)
	//exp.Date = current_time.String()
	//exp.Type = args[3]
	//exp.ExpenditureId = "ex1"
	//exp.FromActor = resA.ActorId
	//exp.ToActor = resB.ActorId

	t.Transfer_balance(stub, []string{args[0], args[1], args[2], "spend"})

	/*If the status of this exp is "Approved", then a reimbursement gonna auto generated*/
	//TODO CALL INIT_REIMBURSEMENT TO GENERATE A NEW REIMBURSEMENT

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

	err = stub.PutState(expIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	err = stub.PutState(reimbIndexStr, jsonAsBytes)
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
		return t.Init_actor(stub, args)
	} else if function == "setup" {
		return t.SetUp(stub, args)
	} else if function == "spend" {
		return t.Spend(stub, args)
	} else if function == "releasefund" {
		return t.ReleaseFund(stub, args)
	} else if function == "transferbalance" {
		return t.Transfer_balance(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

// ============================================================================================================================
//	Query - Called on chaincode query. Takes a function name passed and calls that function. Passes the
//  		initial arguments passed are passed on to the called function.
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "read" {
		return t.read(stub, args)
	} else if function == "queryallexpenses" {
		return t.QueryAllExpenses(stub, args)
	} else if function == "querypendingexpenses" {
		return t.QueryPendingExpenses(stub, args)
	} else if function == "queryblockchain"{
		return t.QueryBlockChain(stub, args)
	} else if function == "querywallet"{
		return t.QueryWallet(stub, args)
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
func (t *SimpleChaincode) Init_actor(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
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
	str := `{"actorid": "` + actorId + `", "actorName": "` + actorName + `", "committed": "` + committedStr + `", "reimbursed": "` + reimbursedStr + `", "awarded": "` + awardedStr + `", "spent": "` + spentStr + `", "received": "` + receivedStr + `", "delegated": "` + delegatedStr + `"}`
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
func (t *SimpleChaincode) Transfer_balance(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//     0         1         2         3
	// "actorA", "actorB", "100.20"  "function"
	var err error
	var newAmountA, newAmountB float64

	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	amount, err := strconv.ParseFloat(args[2], 64)

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

	case "spend":
		fmt.Println("INSIDE CASE SPEND==========================")
		AwardA, err := strconv.ParseFloat(resA.Awarded, 64)
		if err != nil {
			return nil, err
		}
		BalanceA, err := strconv.ParseFloat(resA.Spent, 64)
		if err != nil {
			return nil, err
		}
		BalanceB, err := strconv.ParseFloat(resB.Received, 64)
		if err != nil {
			return nil, err
		}
		//Check if accountA has enough balance to transact or not
		if ( AwardA - amount) < 0 {
			return nil, errors.New(args[0] + " doesn't have enough balance to complete transaction")
		}

		newAmountA = BalanceA + amount
		newAmountB = BalanceB + amount
		newAmountStrA := strconv.FormatFloat(newAmountA, 'f', -1, 64)
		newAmountStrB := strconv.FormatFloat(newAmountB, 'f', -1, 64)

		resA.Spent = newAmountStrA
		resB.Received = newAmountStrB

	case "fund":
		AwardA, err := strconv.ParseFloat(resA.Committed, 64)
		if err != nil {
			return []byte("error in resA.Committed"), err
		}
		BalanceA, err:= strconv.ParseFloat(resA.Reimbursed, 64)
		if err != nil {
			return nil, err
		}
		BalanceB, err := strconv.ParseFloat(resB.Received, 64)
		if err != nil {
			return nil, err
		}
		//Check if accountA has enough balance to transact or not
		if  AwardA - amount < 0 {
			return nil, errors.New(args[0] + " doesn't have enough balance to complete transaction")
		}


		newAmountA = BalanceA + amount
		newAmountB = BalanceB + amount
		newAmountStrA := strconv.FormatFloat(newAmountA, 'f', -1, 64)
		newAmountStrB := strconv.FormatFloat(newAmountB, 'f', -1, 64)

		resA.Reimbursed = newAmountStrA
		resB.Received = newAmountStrB

	}

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
