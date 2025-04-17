package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// LetterOfCredit defines the structure for the Letter of Credit
type LetterOfCredit struct {
	LOCID            string   `json:"locId"`
	Buyer            string   `json:"buyer"`
	Seller           string   `json:"seller"`
	IssuingBank      string   `json:"issuingBank"`
	AdvisingBank     string   `json:"advisingBank"`
	Amount           string   `json:"amount"`
	Currency         string   `json:"currency"`
	ExpiryDate       string   `json:"expiryDate"`
	GoodsDescription string   `json:"goodsDescription"`
	Status           string   `json:"status"`
	DocumentHashes   []string `json:"documentHashes"`
	History          []string `json:"history"`
}

// SmartContract provides functions for managing the Letter of Credit
type SmartContract struct {
	contractapi.Contract
}

// InitLedger initializes the chaincode (optional)
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// TODO: Initialization code if needed
	return nil
}

// RequestLOC creates a new LoC request
func (s *SmartContract) RequestLOC(ctx contractapi.TransactionContextInterface, locID string, buyer string, seller string, issuingBank string, advisingBank string, amount string, currency string, expiryDate string, goodsDescription string) error {
mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get MSP ID: %v", err)
	}
	if mspID != "TataMotorsMSP" {
		return fmt.Errorf("only TataMotorsMSP is authorized to request a Letter of Credit")
	}

	existing, err := ctx.GetStub().GetState(locID)
	if err != nil {
		return fmt.Errorf("failed to read from ledger: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("LoC with ID %s already exists", locID)
	}


func (s *SmartContract) IssueLOC(ctx contractapi.TransactionContextInterface, locID string) error {
	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	if mspID != "SBIBankMSP" {
		return fmt.Errorf("only SBIBankMSP (Issuing Bank) can issue the LoC")
	}

	locBytes, err := ctx.GetStub().GetState(locID)
	if err != nil || locBytes == nil {
		return fmt.Errorf("LoC not found")
	}

	var loc LetterOfCredit
	_ = json.Unmarshal(locBytes, &loc)

	if loc.Status != "Requested" {
		return fmt.Errorf("LoC must be in 'Requested' state to be issued")
	}

	loc.Status = "Issued"
	loc.History = append(loc.History, "LoC issued by Issuing Bank")

	newBytes, _ := json.Marshal(loc)
	return ctx.GetStub().PutState(locID, newBytes)
}

// AcceptLOC - Seller accepts the issued LoC
func (s *SmartContract) AcceptLOC(ctx contractapi.TransactionContextInterface, locID string) error {
	mspID, _ := ctx.GetClientIdentity().GetMSPID()
	if mspID != "SellerOrgMSP" {
		return fmt.Errorf("only SellerOrgMSP can accept the LoC")
	}

	locBytes, err := ctx.GetStub().GetState(locID)
	if err != nil || locBytes == nil {
		return fmt.Errorf("LoC not found")
	}

	var loc LetterOfCredit
	_ = json.Unmarshal(locBytes, &loc)

	if loc.Status != "Issued" {
		return fmt.Errorf("LoC must be 'Issued' to be accepted")
	}

	loc.Status = "Accepted"
	loc.History = append(loc.History, "LoC accepted by Seller")

	newBytes, _ := json.Marshal(loc)
	return ctx.GetStub().PutState(locID, newBytes)
}

func main() {
	chaincode, err := contractapi.NewChaincode(&SmartContract{})
	if err != nil {
		fmt.Printf("Error creating loc chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting loc chaincode: %s", err.Error())
	}
}
