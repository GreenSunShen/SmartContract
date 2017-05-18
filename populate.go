package main
func 

/*
	0 ActorId    string `json:"actorid"`
	1 ActorName  string `json:"actorname"`
	2 Committed  string  `json:"committed"`
	3 Reimbursed string `json:"reimbursed"`
	4 Awarded    string `json:"awarded"`
	5 Spent      string `json:"spent"`
	6 Received   string `json:"received"`
	7 Delegated  string `json:"delegated"`
*/

	// grantor info
	var act1  = make([]string{}, 8,8)
	act1[0] = "101"                //ActorId
	act1[1] = "PPM Foundation"     //ActorName
	act1[2] = "125000"             //Comiitted
	act1[3] = "55000"              //Reimbursed
	act1[4] = ""                   //Awarded
	act1[5] = ""                   //Spent
    act1[6] = ""                   //Received 
    act1[7] = ""                   //Delegated

	// grantee info
	var act2  = make([]string{}, 8,8)
 	act2[0] = "102"                     //ActorId
	act2[1] = "Stanford University"      //ActorName
	act2[2] = ""                        //Comiitted
	act2[3] = ""                         //Reimbursed
	act2[4] = "125000"                 //Awarded
	act2[5] = "23000"                  //Spent
    act2[6] = "55000"                 //Received 
    act2[7] = "45000"                 //Delegated

	// sub-grantee info
	var act3  = make([]string{}, 8,8)
 	act3[0] = "103"                //ActorId
	act3[1] = "John Hopkins University"     //ActorName
	act3[2] = ""               //Comiitted
	act3[3] = ""               //Reimbursed
	act3[4] = "45000"          //Awarded
	act3[5] = "12000"          //Spent
    act3[6] = "25000"          //Received 
    act3[7] = ""               //Delegated

    // Supplier info -- shows spending form all the grantees and sub-grantees
	var act4  = make([]string{}, 8,8)
 	act4[0] = "104"                //ActorId
	act4[1] = "Dixon consulting"     //ActorName
	act4[2] = ""               //Comiitted
	act4[3] = ""               //Reimbursed
	act4[4] = ""          //Awarded
	act4[5] = ""          //Spent
    act4[6] = "35000"          //Received 
    act4[7] = ""               //Delegated


	t.init_actor(stub, act1)
	t.init_actor(stub, act2)
	t.init_actor(stub, act3)
	t.init_actor(stub, act4)