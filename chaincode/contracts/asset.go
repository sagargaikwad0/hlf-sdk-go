package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-samples/chaincode/fabcar/go/models"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type AssetContract struct {
	contractapi.Contract
}

// -------------------Invoke Methods---------------

func (ac *AssetContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	assets := []models.Asset{
		{DocType: "Asset", Id: "1", Name: "Car"},
		{DocType: "Asset", Id: "2", Name: "Mobile"},
		{DocType: "Asset", Id: "3", Name: "Laptop"},
	}

	for _, asset := range assets {
		assetBytes, err := json.Marshal(asset)
		if err != nil {
			return fmt.Errorf("error while marshalling, err: %v", err)
		}

		err = ctx.GetStub().PutState(asset.Id, assetBytes)
		if err != nil {
			return fmt.Errorf("error while put state, err: %v", err)
		}
	}
	return nil
}

func (ac *AssetContract) Create(ctx contractapi.TransactionContextInterface, assetType, assetId, name string) error {

	clientId, err := ctx.GetStub().GetCreator()
	if err != nil {
		return err
	}

	fmt.Println("cleintId", string(clientId))

	asset := models.Asset{
		DocType: assetType,
		Id:      assetId,
		Name:    name,
	}

	assetBytes, err := json.Marshal(asset)
	if err != nil {
		return fmt.Errorf("error while marshalling, err: %v", err)
	}

	err = ctx.GetStub().PutState(assetId, assetBytes)
	if err != nil {
		return fmt.Errorf("error while put state, err: %v", err)
	}

	return nil
}

// ---------------Query Methods----------------

func (ac *AssetContract) QueryById(ctx contractapi.TransactionContextInterface, id string) (*models.Asset, error) {
	asset := new(models.Asset)

	assetBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("error while getting state for assetId:%s, err: %v", id, err)
	}

	err = json.Unmarshal(assetBytes, asset)
	if err != nil {
		return nil, fmt.Errorf("error while unmarshalling state for assetId:%s, err: %v", id, err)
	}

	return asset, nil
}

func (ac *AssetContract) QueryAll(ctx contractapi.TransactionContextInterface) ([]models.Asset, error) {
	// assets := new([]models.Asset)
	var assets []models.Asset

	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("error while getting assets, err: %v", err)
	}

	defer iterator.Close()

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("error while getting assets, err: %v", err)
		}
		asset := new(models.Asset)
		err = json.Unmarshal(kv.Value, asset)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling state for assetId:%s, err: %v", kv.Key, err)
		}
		assets = append(assets, *asset)
	}

	return assets, nil
}

func (ac *AssetContract) QueryByType(ctx contractapi.TransactionContextInterface) ([]models.Asset, error) {
	var assets []models.Asset
	queryString := `{"selector":{"docType":"Asset"}}`

	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("error while getting assets, err: %v", err)
	}
	defer iterator.Close()

	for iterator.HasNext() {
		kv, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("error while getting assets, err: %v", err)
		}
		asset := new(models.Asset)
		err = json.Unmarshal(kv.Value, asset)
		if err != nil {
			return nil, fmt.Errorf("error while unmarshalling state for assetId:%s, err: %v", kv.Key, err)
		}
		assets = append(assets, *asset)
	}

	return assets, nil
}
