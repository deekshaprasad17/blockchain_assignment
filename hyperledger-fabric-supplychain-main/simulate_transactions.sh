#!/bin/bash
# assumes you’ve already set CORE_PEER_ENV vars for each org
# and the chaincode name is 'scm'

echo "=== Manufacturer creates product ==="
peer chaincode invoke -C channel-supply -n scm -c '{"Args":["CreateProduct","P001","Laptop","Intel i5"]}' --waitForEvent

echo "=== Manufacturer ships to Distributor ==="
peer chaincode invoke -C channel-md -n scm -c '{"Args":["CreateShipment","S001","P001","DistributorMSP"]}' --waitForEvent

echo "=== Distributor receives ==="
peer chaincode invoke -C channel-md -n scm -c '{"Args":["ReceiveShipment","S001"]}' --waitForEvent

echo "=== Distributor ships to Retailer ==="
peer chaincode invoke -C channel-dr -n scm -c '{"Args":["CreateShipment","S002","P001","RetailerMSP"]}' --waitForEvent

echo "=== Retailer receives ==="
peer chaincode invoke -C channel-dr -n scm -c '{"Args":["ReceiveShipment","S002"]}' --waitForEvent

echo "=== Query Product History ==="
peer chaincode query -C channel-supply -n scm -c '{"Args":["GetHistory","P001"]}'
