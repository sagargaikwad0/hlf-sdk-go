package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	configPath    = "/home/sagar/Me/Work/hlf-sdk-go/a.yml"
	channelName   = "mychannel"
	chaincodeName = "basic"
)

func query(client *channel.Client,
	request *channel.Request,
	targetPeers *channel.RequestOption) (string, error) {

	response, err := client.Query(*request, *targetPeers)
	if err != nil {
		return "", err
	}

	return string(response.Payload), nil
}

func commit(client *channel.Client,
	request *channel.Request,
	targetPeers *channel.RequestOption) (string, error) {

	response, err := client.Execute(*request, *targetPeers)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

func main() {

	var wg sync.WaitGroup

	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		log.Fatalln("fail to create new sdk", err)
	}
	defer sdk.Close()

	channelContext := sdk.ChannelContext(channelName, fabsdk.WithUser("User1"), fabsdk.WithOrg("Org1"))

	if channelContext == nil {
		log.Fatalf("Failed to create channel context")
	}

	channelClient, err := channel.New(channelContext)

	if err != nil {
		log.Fatalln("failed to create new channel client, err:", err)
	}

	peer1 := getPeerFromConfig(sdk, "peer0.org1.example.com")
	peer2 := getPeerFromConfig(sdk, "peer0.org2.example.com")
	reqPeers := channel.WithTargets(peer1, peer2)

	wg.Add(200)

	for i := 0; i < 200; i++ {
		go func() {
			defer wg.Done()
			assetId := fmt.Sprintf("asset%s", uuid.New())

			request := channel.Request{
				ChaincodeID: chaincodeName,
				Fcn:         "AssetContract:InitLedger",
				Args:        [][]byte{[]byte("Electronic"), []byte(assetId), []byte("FAN")},
			}

			res, err := commit(channelClient, &request, &reqPeers)
			if err != nil {
				log.Fatal("error while commiting, err:", err)
			}

			fmt.Println("Transaction successful, result: ", res)
		}()
	}

	wg.Wait()

}

func getPeerFromConfig(sdk *fabsdk.FabricSDK, peerName string) fab.Peer {
	ctxProvider := sdk.Context()

	ctx, err := ctxProvider()
	if err != nil {
		log.Fatalln("failed to get sdk context, err:", err)
	}

	peerConfig, ok := ctx.EndpointConfig().PeerConfig(peerName)
	if !ok {
		log.Fatalln("failed to get peer config, err:", err)
	}

	fmt.Println("peerConfig.URL", peerConfig.URL)

	networkPeerConfig := fab.NetworkPeer{
		MSPID:      "Org1MSP",
		PeerConfig: *peerConfig,
	}

	peer, err := ctx.InfraProvider().CreatePeerFromConfig(&networkPeerConfig)
	if err != nil {
		log.Fatalln("failed to get peer, err:", err)
	}

	return peer
}
