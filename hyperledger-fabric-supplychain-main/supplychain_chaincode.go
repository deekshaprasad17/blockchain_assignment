package main

import (
    "encoding/json"
    "fmt"
    "time"

    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Product struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Creator  string `json:"creator"`
    Owner    string `json:"owner"`
    Metadata string `json:"metadata"`
}

type Shipment struct {
    ShipmentID string `json:"shipmentId"`
    ProductID  string `json:"productId"`
    From       string `json:"from"`
    To         string `json:"to"`
    Timestamp  string `json:"timestamp"`
    Status     string `json:"status"`
}

type SmartContract struct {
    contractapi.Contract
}

func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, id, name, metadata string) error {
    clientMSP, _ := ctx.GetClientIdentity().GetMSPID()
    exists, _ := ctx.GetStub().GetState(id)
    if exists != nil {
        return fmt.Errorf("product already exists")
    }
    p := Product{ID: id, Name: name, Creator: clientMSP, Owner: clientMSP, Metadata: metadata}
    data, _ := json.Marshal(p)
    return ctx.GetStub().PutState(id, data)
}

func (s *SmartContract) CreateShipment(ctx contractapi.TransactionContextInterface, shipID, productID, to string) error {
    clientMSP, _ := ctx.GetClientIdentity().GetMSPID()
    prodBytes, _ := ctx.GetStub().GetState(productID)
    if prodBytes == nil {
        return fmt.Errorf("product not found")
    }

    var p Product
    _ = json.Unmarshal(prodBytes, &p)
    if p.Owner != clientMSP {
        return fmt.Errorf("only current owner can ship product")
    }

    sh := Shipment{
        ShipmentID: shipID, ProductID: productID, From: clientMSP, To: to,
        Timestamp: time.Now().UTC().Format(time.RFC3339), Status: "created",
    }
    data, _ := json.Marshal(sh)
    return ctx.GetStub().PutState("SHIPMENT_"+shipID, data)
}

func (s *SmartContract) ReceiveShipment(ctx contractapi.TransactionContextInterface, shipID string) error {
    clientMSP, _ := ctx.GetClientIdentity().GetMSPID()
    data, _ := ctx.GetStub().GetState("SHIPMENT_" + shipID)
    if data == nil {
        return fmt.Errorf("shipment not found")
    }

    var sh Shipment
    _ = json.Unmarshal(data, &sh)
    if sh.To != clientMSP {
        return fmt.Errorf("only recipient can confirm receipt")
    }

    sh.Status = "received"
    sh.Timestamp = time.Now().UTC().Format(time.RFC3339)
    newData, _ := json.Marshal(sh)
    _ = ctx.GetStub().PutState("SHIPMENT_"+shipID, newData)

    prodBytes, _ := ctx.GetStub().GetState(sh.ProductID)
    var p Product
    _ = json.Unmarshal(prodBytes, &p)
    p.Owner = clientMSP
    pb, _ := json.Marshal(p)
    return ctx.GetStub().PutState(p.ID, pb)
}

func main() {
    chaincode, _ := contractapi.NewChaincode(new(SmartContract))
    if err := chaincode.Start(); err != nil {
        fmt.Printf("Error starting chaincode: %s", err)
    }
}
