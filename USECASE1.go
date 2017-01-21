/*
Copyright Capgemini India. 2016 All Rights Reserved.
*/

package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
)

// Region Chaincode implementation
type RegionChaincode struct {
}

var regionIndexTxStr = "_regionIndexTxStr"

type RegionData struct{
	REGION_NAME string `json:"REGION_NAME"`
	INSURED_ID string `json:"INSURED_ID"`
	INSURED string `json:"INSURED"`
	BUSINESS_AREA string `json:"BUSINESS_AREA"`
	LINE_OF_BUSINESS_ID string `json:"LINE_OF_BUSINESS_ID"`
	LINE_OF_BUSINESS string `json:"LINE_OF_BUSINESS"`
	POLICY string `json:"POLICY"`
	DEAL_ID string `json:"DEAL_ID"`
	DEAL_NUM string `json:"DEAL_NUM"`
	BROKER_ID string `json:"BROKER_ID"`
	BROKER string `json:"BROKER"`
	INCEPTION_DATE string `json:"INCEPTION_DATE"`
	EXPIRATION_DATE string `json:"EXPIRATION_DATE"`
	CARRIER_CD string `json:"CARRIER_CD"`
	CARRIER string `json:"CARRIER"`
}


func (t *RegionChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	
	var err error
	// Initialize the chaincode
	
	fmt.Printf("Deployment of Loyalty is completed\n")
	
	var emptyPolicyTxs []RegionData
	jsonAsBytes, _ := json.Marshal(emptyPolicyTxs)
	err = stub.PutState(regionIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	
	
	return nil, nil
}

// Add region data for the policy
func (t *RegionChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == regionIndexTxStr {		
		return t.RegisterPolicy(stub, args)
	}
	return nil, nil
}
 
func (t *RegionChaincode)  RegisterPolicy(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	
	var RegionDataObj RegionData
	var RegionDataList []RegionData
	var err error

	if len(args) != 15 {
		return nil, errors.New("Incorrect number of arguments. Need 14 arguments")
	}

	// Initialize the chaincode
	RegionDataObj.REGION_NAME = args[0]
	RegionDataObj.INSURED_ID = args[1]
	RegionDataObj.INSURED = args[2]
	RegionDataObj.BUSINESS_AREA = args[3]
	RegionDataObj.LINE_OF_BUSINESS_ID = args[4]
	RegionDataObj.LINE_OF_BUSINESS = args[5]
	RegionDataObj.POLICY = args[6]
	RegionDataObj.DEAL_ID = args[7]
	RegionDataObj.DEAL_NUM = args[8]
	RegionDataObj.BROKER_ID = args[9]
	RegionDataObj.BROKER = args[10]
	RegionDataObj.INCEPTION_DATE = args[11]
	RegionDataObj.EXPIRATION_DATE = args[12]
	RegionDataObj.CARRIER_CD = args[13]
	RegionDataObj.CARRIER = args[14]
	
	fmt.Printf("Input from user:%s\n", RegionDataObj)
	
	regionTxsAsBytes, err := stub.GetState(regionIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get consumer Transactions")
	}
	json.Unmarshal(regionTxsAsBytes, &RegionDataList)
	
	RegionDataList = append(RegionDataList, RegionDataObj)
	jsonAsBytes, _ := json.Marshal(RegionDataList)
	
	err = stub.PutState(regionIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *RegionChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {
	
	var PolicyId string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the person to query")
	}

	PolicyId = args[0]
	
	resAsBytes, err = t.GetPolicyDetails(stub, PolicyId)
	
	fmt.Printf("Query Response:%s\n", resAsBytes)
	
	if err != nil {
		return nil, err
	}
	
	return resAsBytes, nil
}

func (t *RegionChaincode)  GetPolicyDetails(stub shim.ChaincodeStubInterface, PolicyId string) ([]byte, error) {
	
	//var requiredObj RegionData
	var objFound bool
	PolicyTxsAsBytes, err := stub.GetState(regionIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Merchant Transactions")
	}
	var PolicyTxObjects []RegionData
	var PolicyTxObjects1 []RegionData
	json.Unmarshal(PolicyTxsAsBytes, &PolicyTxObjects)
	length := len(PolicyTxObjects)
	fmt.Printf("Output from chaincode: %s\n", PolicyTxsAsBytes)
	
	if PolicyId == "" {
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
		if PolicyId == obj.POLICY {
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
	err := shim.Start(new(RegionChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
