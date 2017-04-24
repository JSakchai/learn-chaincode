/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"encoding/json"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}
type Marble struct {
	Name string `json:"name"`
	Color string `json:"color"`
	Size int `json:"size"`
	User string `json:"user"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}
	err := stub.PutState("hello world" , []byte(args[0]))
	if err != nil {
		return  nil,err
	}
	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {													//initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else  if function == "write" {
		return  t.write(stub,args)
	}
	fmt.Println("invoke did not find func: " + function)					//error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" {											//read a variable
		fmt.Println("hi there " + function)						//error
		return nil, nil;
	}
	fmt.Println("query did not find func: " + function)						//error

	return nil, errors.New("Received unknown function query: " + function)
}
// write is put data to block
func (t *SimpleChaincode) write(stup shim.ChaincodeStubInterface, args []string) ([]byte ,error){
	var key, value string
	var err error
	fmt.Println("running write()")
	if len(args) != 2 {
		return  nil, errors.New("inconnect number of argument. Expecting 2. name of the key and value to set")
	}
	key = args[0]
	value =args[1]
	err = stup.PutState(key, []byte(value))
	if err != nil {
		return  nil, err
	}
	return nil ,nil
}
func (t *SimpleChaincode) read(stup shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}
	key = args[0]
	valAsbytes, err := stup.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil,errors.New(jsonResp)
	}
	return  valAsbytes , nil
}
func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte,error){
	var err error
	// 0       1
	// "name"  "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of argument. Expecting")
	}
	fmt.Println("-start set user")
	fmt.Println(args[0] + " - " + args[1])
	marbleAsByte, err := stub.GetState(args[0])
	if err != nil {
		return  nil, errors.New("Fail to get thing")
	}
	res := Marble{}
	json.Unmarshal(marbleAsByte, &res)  //un stringify it aka JSON.parse()
	res.User = args[1]
	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes)  //rewrite the marble with id as key
	if err != nil{
		return nil,err
	}
	fmt.Println("- end set user")
	return nil,nil
}