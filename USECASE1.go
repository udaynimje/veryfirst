/*
Copyright Capgemini India. 2016 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
    //"github.com/hyperledger/fabric/chaincode/"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
)

// User Chaincode implementation
type UserChaincode struct {
}
var userIndexTxStr="_userInitState"
var userTxnDetails = "_USERTXN"
var userRegister = "_USERREG"
var userTransfer = "_USERTRANSFER"
var getTxnDetails="_USERTXNDETAILS"

type UserData struct{
	User_NAME string `json:"USER_NAME"`
	User_ID string `json:"USER_ID"`
	Password string `json:"PASSWORD"`
	MID string `json:"MID"`
	CONTACT_NO string `json:"CONTACT_NO"`
	TRANS_AMT string `json:"TRANS_AMOUNT"`
	TRANS_POINT string `json:"TRANS_POINT"`
	REG_DATE string `json:"REGISTRATION_DATE"`
	TRANS_DATE string `json:"TRANSACTION_DATE"`
}


func (t *UserChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var err error
	// Initialize the chaincode
	
	fmt.Printf("Loyalty Rewards is implemented\n")
	
	var user []UserData
	jsonAsBytes, _ := json.Marshal(user)
	err = stub.PutState(userRegister, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	
	return jsonAsBytes, nil
}

// Add user data
func (t *UserChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {


	if function == userRegister {		
		return t.RegisterPolicy(stub, args)
	}  else if  function == userTransfer {		
		return t.TransferPoints(stub, args)
	} else if function == getTxnDetails {
        return t.Query(stub,function,args)
    }
   jasonAsBytes, _:=json.Marshal(args)
	return jasonAsBytes, nil
}

func (t *UserChaincode)  TransferPoints(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
    fmt.Printf("INSIDE TransferPoints")
return nil,nil
}

func (t *UserChaincode)  RegisterPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var UserDataObj UserData
	var UserDataList []UserData
	var err error

	if len(args) != 9 {
		return nil, errors.New("Incorrect number of arguments. Need 8 arguments")
	}

	// Initialize the chaincode
	UserDataObj.User_NAME = args[0]
	UserDataObj.User_ID = args[1]
	UserDataObj.Password = args[2]
	UserDataObj.MID = args[3]
	UserDataObj.CONTACT_NO = args[4]
	UserDataObj.TRANS_AMT = args[5]
	UserDataObj.TRANS_POINT = args[6]
	UserDataObj.REG_DATE = args[7]
	UserDataObj.TRANS_DATE = args[8]
	
	
	fmt.Printf("Input from user:%s\n", UserDataObj)
	
	userTxsAsBytes, err := stub.GetState(userIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get user Details")
	}
	json.Unmarshal(userTxsAsBytes, &UserDataList)
	
	UserDataList = append(UserDataList, UserDataObj)
	jsonAsBytes, _ := json.Marshal(UserDataList)
	
	err = stub.PutState(userIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
    return jsonAsBytes, nil
}

// Query callback representing the query of a chaincode
func (t *UserChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	
	var PolicyId string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	PolicyId = args[0]
	
	//resAsBytes, err = t.GetPolicyDetails(stub, PolicyId)
	readAsBytes,_=json.Marshal(args)
	fmt.Printf("Query Response:%s\n", resAsBytes)
	
	if err != nil {
		return nil, err
	}
	
	return resAsBytes, nil
}

func (t *UserChaincode)  GetPolicyDetails(stub shim.ChaincodeStubInterface, User_ID string) ([]byte, error) {
	
	//var requiredObj UserData
	var objFound bool
	PolicyTxsAsBytes, err := stub.GetState(userIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	var PolicyTxObjects []UserData
	var PolicyTxObjects1 []UserData
	json.Unmarshal(PolicyTxsAsBytes, &PolicyTxObjects)
	length := len(PolicyTxObjects)
	fmt.Printf("Output from chaincode: %s\n", PolicyTxsAsBytes)
	
	if User_ID == "" {
		res, err := json.Marshal(PolicyTxObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
	
	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := PolicyTxObjects[i]
		if User_ID == obj.User_ID {
			PolicyTxObjects1 = append(PolicyTxObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}
	
	if objFound {
		res, err := json.Marshal(PolicyTxObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}
}

func main() {
	err := shim.Start(new(UserChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
