package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

const (
	configPath    = "/home/sagar/Me/Work/hlf-sdk-go/sdk/a.yml"
	channelName   = "mychannel"
	chaincodeName = "basic"
	n             = 2
)

type sdk struct {
	fabricsdk      *fabsdk.FabricSDK
	requestOptions []channel.RequestOption
	client         *channel.Client
	counter        int32
}

func newSdk(user, org string) *sdk {

	sdkk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		log.Fatalln("fail to create new sdk", err)
	}

	channelContext := sdkk.ChannelContext(channelName, fabsdk.WithUser(user), fabsdk.WithOrg(org))

	if channelContext == nil {
		log.Fatalf("Failed to create channel context")
	}

	channelClient, err := channel.New(channelContext)

	if err != nil {
		log.Fatalln("failed to create new channel client, err:", err)
	}

	return &sdk{
		fabricsdk: sdkk,
		client:    channelClient,
		counter:   0,
	}
}

func (s *sdk) setRequestOptions() {
	peer0Org1 := getPeerFromConfig(s.fabricsdk, "peer0.org1.example.com")
	peer0Org2 := getPeerFromConfig(s.fabricsdk, "peer0.org2.example.com")

	s.requestOptions = append(s.requestOptions, channel.WithTargets(peer0Org1))
	s.requestOptions = append(s.requestOptions, channel.WithTargets(peer0Org2))
}

func (s *sdk) commit(request *channel.Request) (string, error) {

	response, err := s.client.Execute(*request)
	if err != nil {
		return "", err
	}

	return string(response.TransactionID), nil
}

// func query(client *channel.Client,
// 	request *channel.Request,
// 	targetPeers *channel.RequestOption) (string, error) {

// 	response, err := client.Query(*request, *targetPeers)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(response.Payload), nil
// }

func main() {

	var wg sync.WaitGroup

	s := newSdk("User1", "Org1")
	defer s.fabricsdk.Close()
	s.setRequestOptions()

	wg.Add(n)

	start := time.Now()

	for i := 0; i < n; i++ {
		go func() {

			if s.counter == 0 {
				atomic.StoreInt32(&s.counter, 1)
			} else {
				atomic.StoreInt32(&s.counter, 0)
			}

			fmt.Println("counter", s.counter)

			defer wg.Done()
			assetId := fmt.Sprintf("asset%s", uuid.New())

			request := channel.Request{
				ChaincodeID: chaincodeName,
				Fcn:         "AssetContract:Create",
				Args:        [][]byte{[]byte("Electronic"), []byte(assetId), []byte("FAN")},
			}

			res, err := s.commit(&request)
			if err != nil {
				log.Fatal("error while commiting, err:", err)
			}

			fmt.Println("Transaction successful, result: ", res)
		}()
	}

	wg.Wait()

	fmt.Println("time took:", time.Since(start))

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
