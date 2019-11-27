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
	"github.com/binance-chain/go-sdk/client/rpc"
	"encoding/hex"
	"encoding/json"
)

var (
	ksFilePath  = "./BD-ks-airdrop-Admin@123456.txt"
	ksAuth      = "Admin@123456"
	networkType = types.ProdNetwork
	dexUrl      = "dex.binance.org"
	nodeUrl     = "tcp://seed1.4leapbnb.com:80"
	dexClient   client.DexClient
	rpcClient   rpc.Client
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

		rpcClient = rpc.NewRPCClient(nodeUrl, networkType)
		if _, err := rpcClient.Status(); err != nil {
			fmt.Printf("init rpc client fail, err is %s\n", err.Error())
			panic(err)
		}
	}
	m.Run()
}

func TestRecoverFromKeyStore(t *testing.T) {
	file := "./BD-ks-airdrop-Admin@123456.txt"
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

			if phrase, err := km.ExportAsPrivateKey(); err != nil {
				t.Fatal(err)
			} else {
				t.Logf("phrase is: %s\n", phrase)
			}
		}
	}

}

func TestSendToken(t *testing.T) {
	var (
		msgs []msg.Transfer
	)
	denom := "IRIS-D88"

	//receivers := map[string]float64{
	//	"bnb1zffz72cfe56zhru8qsx0dxy49vd64uurwxmxdg": 1500,
	//	"bnb1vq0felttf4xywck80mjhsnz55gpv674c3dtt6p": 1420.6,
	//	"bnb1cyejwpv203jyyapdltzlmsseen6mvk42ec9a39": 1500,
	//	"bnb1h408kq2euyzyftcw89safl2pt9hq8l0n5xd3j4": 1500,
	//	"bnb1a0shu67rvsdyedr6rw5sfay3xf5p4jh82r0rt3": 66.8,
	//	"bnb16g9lgmfkycsgzwssmfwkvfp5n8nutx2p9cpmt4": 1500,
	//	"bnb12dsxhqkdlnnld7jsfnd22c30g27hncwmywre8k": 1500,
	//	"bnb163pak9zglyncfu8na84dl43hezc6egf90cw9d3": 1500,
	//	"bnb1x3v6wlfdw8u78r9w455wetk2uxdxazekswyq4p": 1500,
	//	"bnb1tf5ngk9ra9s9lwz0rj5qx8lu2w4hsjjf5e4hsa": 1500,
	//	"bnb1es5app887rpfvhchjdcdmyyll330p92qyz9mdw": 1500,
	//	"bnb13uh6wfjfxn3kp3zj6g55tp6ak4eyvea2gq6zuz": 1500,
	//	"bnb18wajhf0mxnffkjutyr80zfmrandwxt53rfned5": 79.96,
	//	"bnb1h6hx9dyy37yhhftqft8h8xslnnwmgjss0ztx9u": 1500,
	//	"bnb1ps4jhdf53g4cm0jm2rvq02y0ucymavxzhscn38": 1500,
	//	"bnb173ks5el7fze729j0h7ftl0zgq8w2nvrh9ku46t": 1500,
	//	"bnb1x3590yxs84cpvxp25jmrtryd3t4jdp3xwj25qg": 1500,
	//	"bnb18qfqfa8x6z2ds9ufjtzsfwq3rvwdcx7zq3d4jl": 1500,
	//	"bnb1sd2cfg0yeec0tyvly4pt72juhthxl9hctymayp": 1500,
	//	"bnb12hj4tru4mfp4yqk6tzckfruvshp7g6vg5hm30c": 1500,
	//	"bnb16zf5ww2camfaygq0xtg3srtnd6r2r4mjm520a6": 1500,
	//	"bnb1gkqed0vud2atcxlju2krat37zafnxmxqxul5p4": 1400.3,
	//	"bnb15e0ymru8nawpkpjr0xv26yk4m9ysfpkqwpncfq": 205,
	//	"bnb17mwlw7653anr60wckahxsq4e02fzmggc9wltk2": 580.04,
	//	"bnb17q3cs933tevtuzhsw2ty86j6lazpzrnjashlqc": 280,
	//	"bnb1apuhczfpf5vr3a9qth9fvmrgcfse63m47pkzmd": 1500,
	//	"bnb1usgzttfpk7t0jd9v95thumw27vjr6cls5tn23r": 1077.2,
	//	"bnb1j3ar0tj3uhcyvqqumqm4seyv3zjh36mh6tpqlg": 1500,
	//	"bnb1rl3j5f6spldm9384ykcv6dyhe02gl8yct7j9js": 1500,
	//	"bnb1lqkzujtzvkm8nyz866wuq7ul9ugh85z9pfgxek": 110.4,
	//	"bnb13kmf39jp9pllmnnzp6a9jlrz6qx4fe08f8c6ug": 866.3,
	//	"bnb1d8rdcd7t8c4sdtssan9w6ks0cvyk5cqg6ssydy": 1500,
	//	"bnb1lys498aqr6f5k64x6l02nrdwqhz8mza0e0u8gy": 1500,
	//	"bnb1sz7znvzfuazwl38s26upgupmszsvupd9c9gr9e": 1500,
	//	"bnb1akran4uh3w7f4xtjg475m9p6pcnxrsml705n3f": 1500,
	//	"bnb1dy2vz9ppdnel6uy2qtla7xq8kyvz8wftyun5rz": 46,
	//	"bnb1yu87ql8l9phwefy43jkq29ldvjlsdkpjd9rjjy": 602.4,
	//	"bnb1r9kexl2n5ttt4fqsm34e5pkk9hmdqy46r99yaf": 216,
	//	"bnb1yx9vvrjxn43ycwv9stkkklhwxqt0dapfhrx77v": 784.66,
	//	"bnb19c0c8crgmms68caztjel4m8waeg8xxkk7gel9n": 1500,
	//	"bnb1zcls06dh9jm205m2fgtkrt4p8w9h8nd5dhqvp8": 90,
	//	"bnb1qdazrdrgue5chl9zeukhqs9tr9yaav9unve0lh": 1500,
	//}
	receivers := map[string]float64{
		"bnb15e0ymru8nawpkpjr0xv26yk4m9ysfpkqwpncfq": 205,
		"bnb1rl3j5f6spldm9384ykcv6dyhe02gl8yct7j9js": 1500,
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
	txHash := "193B0E7C89D5EC3C16FF81949F1EB471C2E0E4245EAA21DD6018F9ADBE36E9F7"
	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		t.Fatal(err)
	}
	resultTx, err := rpcClient.Tx(txHashBytes, false)
	if err != nil {
		t.Fatal(err)
	}

	if resBytes, err := json.Marshal(resultTx); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(resBytes))
	}
}
