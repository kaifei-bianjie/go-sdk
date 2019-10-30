package mock

import (
	"fmt"
	"github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/transaction"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"math"
	"testing"
	"time"
)

var (
	ksFilePath  = "./BD-testnet-ks_2_Admin@123456.txt"
	ksAuth      = "Admin@123456"
	networkType = types.TestNetwork
	dexUrl      = "testnet-dex.binance.org"
	dexClient   client.DexClient
)

func TestMain(m *testing.M) {
	if km, err := keys.NewKeyStoreKeyManager(ksFilePath, ksAuth); err != nil {
		panic(err)
	} else {
		if c, err := client.NewDexClient(dexUrl, networkType, km); err != nil {
			panic(err)
		} else {
			fmt.Printf("km address is %s\n", km.GetAddr().String())
			dexClient = c
		}
	}
	m.Run()
}

func TestRecoverFromKeyStore(t *testing.T) {
	file := "./BD-testnet-ks_2_Admin@123456.txt"
	auth := "Admin@123456"
	if km, err := keys.NewKeyStoreKeyManager(file, auth); err != nil {
		t.Fatal(err)
	} else {
		content := []byte("Testing")
		if signedBytes, err := km.GetPrivKey().Sign(content); err != nil {
			t.Fatal(err)
		} else {
			if km.GetPrivKey().PubKey().VerifyBytes(content, signedBytes) {
				t.Log("verify signed bytes success")
			}
		}
	}

}

func TestSendToken(t *testing.T) {
	var (
		msgs []msg.Transfer
	)
	denom := "IRISB-9FD"

	receivers := map[string]float64{
		"tbnb1rlkq9zx6cze6umc8r57r5grww7tunxsl7d24ws": 121,
		"tbnb1udy08aymj03zwuhxgumlzg9ddmn67xpqs9yhaz": 3.12,
	}

	for k, v := range receivers {
		if toAddr, err := types.AccAddressFromBech32(k); err != nil {
			t.Fatalf("invalid addr: %s\n", toAddr)
		} else {
			var coins []types.Coin
			coin := types.Coin{
				Denom:  denom,
				Amount: int64(v * math.Pow10(8)),
			}
			coins = append(coins, coin)

			msgSendToken := msg.Transfer{
				ToAddr: toAddr,
				Coins:  coins,
			}
			msgs = append(msgs, msgSendToken)
		}
	}

	for _, v := range msgs {
		msgSendToken := []msg.Transfer{v}
		option := transaction.WithMemo("congratulation, you got airdrop from IRISnet")

		if sendResult, err := dexClient.SendToken(msgSendToken, true, option); err != nil {
			fmt.Printf("send token occur error, toAddr is %s, err is %s\n", v.ToAddr.String(), err.Error())
		} else {
			if sendResult.Ok {
				fmt.Printf("send token success, toAddr is %s, txHash is %s\n", v.ToAddr.String(), sendResult.Hash)
			} else {
				fmt.Printf("send token fail, toAddr is %s, txHash is %s, log is %s\n", v.ToAddr.String(),
					sendResult.Hash, sendResult.Log)
			}
		}
		fmt.Println("now sleep 5 seconds")
		time.Sleep(time.Duration(5) * time.Second)
	}
}

func TestGetTxDetail(t *testing.T) {

}
