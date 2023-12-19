package modules

import (
	"fmt"
	gsrpc "github.com/centrifuge/go-substrate-rpc-client/v4"
	"github.com/centrifuge/go-substrate-rpc-client/v4/signature"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"time"
)

func AvailModule(daData string) (bool, string) {
	rpcUrl := "wss://kate.avail.tools:443/ws"
	seedValue := "gravity ocean disease tent fitness stereo canal angry pill energy clutch rubber"

	//* rpc api ws check function
	api, err := gsrpc.NewSubstrateAPI(rpcUrl)
	if err != nil {
		panic(fmt.Sprintf("cannot create api:%v", err))
	}

	fmt.Print("Connected to ", api)

	// metadata function
	meta, err := api.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(fmt.Sprintf("cannot get metadata:%v", err))
	}

	//* data submission call
	newCall, err := types.NewCall(meta, "DataAvailability.submit_data", types.NewBytes([]byte(daData)))
	if err != nil {
		panic(fmt.Sprintf("cannot create new call:%v", err))
	}

	//* extrinsic function
	ext := types.NewExtrinsic(newCall)

	genesisHash, err := api.RPC.Chain.GetBlockHash(0)
	if err != nil {
		panic(fmt.Sprintf("cannot get block hash:%v", err))
	}

	rv, err := api.RPC.State.GetRuntimeVersionLatest()
	if err != nil {
		panic(fmt.Sprintf("cannot get latest runtime version:%v", err))
	}

	keyringPair, err := signature.KeyringPairFromSecret(seedValue, 42)
	if err != nil {
		panic(fmt.Sprintf("cannot create KeyPair:%v", err))
	}

	key, err := types.CreateStorageKey(meta, "System", "Account", keyringPair.PublicKey)
	if err != nil {
		panic(fmt.Sprintf("cannot create storage key:%v", err))
	}

	var accountInfo types.AccountInfo
	ok, err := api.RPC.State.GetStorageLatest(key, &accountInfo)
	if err != nil || !ok {
		panic(fmt.Sprintf("cannot get latest storage:%v", err))
	}

	nonce := uint32(accountInfo.Nonce)
	options := types.SignatureOptions{
		BlockHash:          genesisHash,
		Era:                types.ExtrinsicEra{IsMortalEra: false},
		GenesisHash:        genesisHash,
		Nonce:              types.NewUCompactFromUInt(uint64(nonce)),
		SpecVersion:        rv.SpecVersion,
		Tip:                types.NewUCompactFromUInt(100),
		TransactionVersion: rv.TransactionVersion,
		AppID:              types.NewUCompactFromUInt(0),
	}

	err = ext.Sign(keyringPair, options)
	if err != nil {
		panic(fmt.Sprintf("cannot sign:%v", err))
	}

	//* check extrinsic status
	sub, err := api.RPC.Author.SubmitAndWatchExtrinsic(ext)
	if err != nil {
		panic(fmt.Sprintf("cannot submit extrinsic:%v", err))
	}

	defer sub.Unsubscribe()

	timeout := time.After(15 * time.Second)
	for {
		select {
		case status := <-sub.Chan():
			if status.IsInBlock {
				return true, status.AsInBlock.Hex()
			}
		case <-timeout:
			return false, "timeout of 15 seconds reached without getting finalized status for extrinsic"
		}
	}

}
