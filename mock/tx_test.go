package mock

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/binance-chain/go-sdk/client"
	"github.com/binance-chain/go-sdk/client/rpc"
	"github.com/binance-chain/go-sdk/client/transaction"
	"github.com/binance-chain/go-sdk/common/types"
	"github.com/binance-chain/go-sdk/keys"
	"github.com/binance-chain/go-sdk/types/msg"
	"github.com/binance-chain/go-sdk/types/tx"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	tmtypes "github.com/tendermint/tendermint/types"
	"math"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	ksFilePath = "./bd@Admino0o0oo0-ks.json"
	ksAuth     = "bd@Admino0o0oo0"

	networkType = types.TestNetwork
	dexUrl      = "testnet-dex.binance.org"
	nodeUrl     = "tcp://seed-pre-s3.binance.org:80"

	//networkType = types.ProdNetwork
	//dexUrl      = "dex.binance.org"
	//nodeUrl     = "tcp://dataseed1.ninicoin.io:80"

	dexClient  client.DexClient
	rpcClient  *rpc.HTTP
	onceClient = sync.Once{}
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

		//rpcClient = rpc.NewRPCClient(nodeUrl, networkType)
		//if _, err := rpcClient.Status(); err != nil {
		//	fmt.Printf("init rpc client fail, err is %s\n", err.Error())
		//	panic(err)
		//}
	}
	m.Run()
}

func defaultClient() *rpc.HTTP {
	onceClient.Do(func() {
		c := rpc.NewRPCClient(nodeUrl, networkType)
		if _, err := c.Status(); err != nil {
			panic(err)
		} else {
			rpcClient = c
		}
	})
	return rpcClient
}

func TestRecoverFromKeyStore(t *testing.T) {
	file := "./bd@Admino0o0oo0-ks.json"
	auth := "bd@Admino0o0oo0"
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
	receivers := make(map[string]float64)

	jsonStr := `{"bnb108hndds8fvpe5am5jxpcan8ejwvmxxe6epfwef":344.74,"bnb108wur49ndty8rmx066hagrl2vg2q94cudveul0":344.74,"bnb10ffgjcexxxkggx83wxa8zkvargr2m0asskd8c6":344.74,"bnb10gmnn5k9pgsmlmvf7sl3r32kll0lu2ududcwtu":1292.77,"bnb10jvlys8cmnuls9kux6gwp39c4lhklf733vwqjy":344.74,"bnb10p2r5v7fhdcsxf6qhl8exvkw3x5a9lucmml6wl":344.74,"bnb10w7cwxe09kwe2ssausscahgqn22t88u4mrqdng":344.74,"bnb10xg327kqhz62rwjwt4cykmvc5um8yhusy5vskh":344.74,"bnb124llarx0nhl0n2cvfquaeh9tq06gu0zfvszznt":344.74,"bnb12hhz9jdgga4w76z7dfkgfcykzxrmdkpq9s3s0e":2585.54,"bnb12mahfshmuz0hhutlmg4wl4wdh5p0av6jps8dch":344.74,"bnb12qs8xk8peguw7cuuq5ujvqxds9fvenzpdvskgj":1292.77,"bnb12u4a79x7lhes5s5kme32p9yradm5rzamfwdp63":344.74,"bnb12u5y400zkyqxuhksz7y5j98rjvsasqek8q2y35":344.74,"bnb134gzsvzrq3avds83w92arat78nv8qpwut0vhwc":1292.77,"bnb13ez7wzhvwfa73m2kwxa7nfhvr3957pmt3uf9mu":344.74,"bnb14ckm9uamehk22rhvqwc87427nk546xj9qr8wjw":344.74,"bnb14rqmpv08wllupu3gsee9fk03hk06epzpryj64d":344.74,"bnb14xft5tgzlc470ld33pyumsvp7z494kl6wpcszz":2671.72,"bnb15azhkyq3dlzx5306730uadkdanlqfznu4dkmxz":344.74,"bnb15kkhdtj34cjujmkx4mzx0fhcn4arpqvx5atr9n":2671.72,"bnb165362zqn4tv8ezmhccjap5jvqzmh8vz556d0mc":344.74,"bnb16g32uka0szqp9uzkfhlp53hgkuul7a5ds3cn80":2671.72,"bnb16p0qz09qr739rxee8nt4txxayg4ujf9evqt5ag":1206.5900000000001,"bnb16p2xpr8f35nwe45zc0z2set0jk72zw5zu7rhhw":861.85,"bnb16pwxffhxadrzjf2pugfqzm0ljyet5x29dwcp9x":2585.54,"bnb170c5jj7yzhdsx0rrnwukdnqg4tun24677q5jxj":2671.72,"bnb1773rvcayu4rfedel5p0u27cwh00gjy5yg3t2th":344.74,"bnb183rzm3jfufukgw6ysf0clsdexqey50hj8rpcgp":344.74,"bnb18cjynn8tu8amuqju70pk8smtwmgnnt8x6m5hs6":344.74,"bnb18lductqygfk5vwag44ws9uyd7x8ec4p6eurqsl":344.74,"bnb18vfmh0fckwffp2gqtr48seeu7xzcenhkutf2rh":1292.77,"bnb190yu8eu8qttnxhkv24uqqheg3qreerdsgr47re":344.74,"bnb1a0etk9846as3668wgg5gvgffjrsu2nukck6gae":2671.72,"bnb1a0wu625s8x35ry9tu0khuxnkar73kxwh6q0sqx":2671.72,"bnb1adzex94jcg6skyqlnsm4uretj6e5f8qhrqmh8h":2671.72,"bnb1arxsqsd7pk3kzfjsqttvvcuydg3a080p9mgql0":344.74,"bnb1as8azxj9dl8xguvhn6j9zv5th7cg52tgxr4e4f":2671.72,"bnb1d5qwvjw7aa9n8wwnfdhvwc2js5vku42tseazqj":344.74,"bnb1d6v5rhtlfu0pvmxa5mrn2rqgkgcxxhmw2vag8p":2671.72,"bnb1daj7wg0wpp4v0klve9nttnmtq8g4pqfyrztvlr":2585.54,"bnb1dcpt3ezadrnz3rc757uem0ggwk59g7ya4kp9gt":344.74,"bnb1egfuqvegtrwe2875cplsph7q2ncqtaket95ysd":344.74,"bnb1ejquen0drv8uzckxysxq8fm27q6spgccqy2fvl":861.85,"bnb1fwkhgwtafwkk8ueta5zfzadxkvsh8yfzp9g7as":2585.54,"bnb1grys2mxpm7nxhsh9w6qexhvt9hcvdwvq5sqdfm":344.74,"bnb1gvmp6ww322eyde2lgqvuu2vc5hvx256w79yake":1292.77,"bnb1h93c5u9usavyce7kgryy27uwku4ddpz6d7k4u2":2585.54,"bnb1jap72s4zrdccdr0h0d5rqkp07ch5zg2x7pn6eg":344.74,"bnb1jkdxzjrhjvfjhae34hm5hjtfzsc3cjq2r8jft3":861.85,"bnb1jr56zc8g85p9fdxsamkffxf2gwa7arzrq2k5az":344.74,"bnb1jt8e76rx8cvlealya48ct20g9ehej4g8snh265":2671.72,"bnb1jtr4jugme5u6a2mwr9nvcdsf683gk6zelz4lzz":344.74,"bnb1juf57k99wn2yyf2y24kj4ta0pr02qsusuh85v2":344.74,"bnb1jz5a6k2tjmaqpu43zq2fye4p3jetnn6x2a8kpq":344.74,"bnb1k8yzpmzd6jrtrule8p77ecuh60j2luxgh3lzwe":344.74,"bnb1knnhy0qedkggueya5qc5xfdqkzgadnfryueeez":2671.72,"bnb1kqtg4hk64vznyjv8ujct5l45rjedrtc6juyngl":1206.5900000000001,"bnb1kxvguwlncf7mpj72vgm298geyct28lk2hwck3a":344.74,"bnb1l6s3uncuanydnwr8a7fc8s7gz2vmjtsd9nx0ml":2671.72,"bnb1l7a58wt675wlfky0druyqwma62dm59d3jplg4p":2671.72,"bnb1l9ayjvk99e9vlv0gu2kc9wz0jec0agxryyz5e5":344.74,"bnb1lclgynfvt6jjvlcw9nf8ahz40r2xad2cuu90s7":344.74,"bnb1ls6jqm3zeh5exspeql59jzyt2r38kltzvp5u59":861.85,"bnb1lzp75c36sksphqjrtlzx4xwkej4dhf4za5929u":344.74,"bnb1m0l27q6f7n9zpclkdtfszgjffhl8fh9hwydd6w":861.85,"bnb1m73w2vl8k7kpnmu89qpg4gl6uhscdlnz7evpjv":2585.54,"bnb1m7m72dzylpsclr5jtz2qhj8t42staqnar9c20m":344.74,"bnb1mlejyeeyhzm550dz0vch7qf87nuqfj3nnmt9x5":861.85,"bnb1mqqsmm8es90z5l6y0gjejchwjltqjjmmhzmgp5":344.74,"bnb1mu0g7fzuhj0cl89n4eulh97q4fgp3u43hg6zq2":344.74,"bnb1mxpzrgrfd2zrk8c08vrlkuqmjz46h0787y9ld0":2671.72,"bnb1nn4dcvqfx8625syfry0un2hwrvwzu0n47wzr0k":1292.77,"bnb1nreregu3kvexwqfyaz3t655vmtmwy7lkvefxzx":689.48,"bnb1nts4r6wfxgvq2r0n577azvuk2yv6xfedcazrdm":1292.77,"bnb1p0an9ttvd33xysssc4nt3glt4sm6a3gmdwvtas":1292.77,"bnb1p53w8pusm8sdzd4rpedyr5lv2tfhk7flw73tjj":344.74,"bnb1pd35qkpzmkgpsv9gxjucx20fpmuk99hw8g4hje":344.74,"bnb1pkv4utc9ewxxggr06zuj6em7gaghef57g5x5xn":344.74,"bnb1q4vw2l5hg6dykqs8ngjmexyzdzcfs65rp9f48k":344.74,"bnb1q9uzmf0ufc0umyf2hjne4ajrg66ndugd5s65sn":689.48,"bnb1qd78yr7mdjaw2z7je6rgt2jzhuts0r24kh36l4":2585.54,"bnb1qpm3qs902y2mmvcfrnfgwvhqhx9trp8wp4w9fm":2671.72,"bnb1r36z4qnu2eglyacpa7h54tc8w2dasgxh989lft":344.74,"bnb1r5rkkeewazxves43acwzyjfyzmt6xms6xwfwq6":2671.72,"bnb1rxmcjf2slgjdjd0ldqad95ydhzr8grk3j5u0mm":344.74,"bnb1skgywtay55j9tt4ejexftqs8utw6hjcrfrgxes":1292.77,"bnb1suyhudla39rvaldwqqs9n5346q8g25y6yw5c6w":2671.72,"bnb1t2z9p0qah2tsatu76ghm30vh4cvshnpqlm9v6j":344.74,"bnb1tv3xc805w89m3h5gy0r3au72w3wez6krc7363r":861.85,"bnb1u4lf027sha92kf0aw6rear5hd4trt72lzn7kyw":1292.77,"bnb1uf0h5mqc8kgykhjzmwh5zxcvlaec0d3847hz8q":1292.77,"bnb1ujzn2lwjwk82rq7nanuqzf0c4lvc843hfmg9fc":2930.2799999999997,"bnb1uwevr8eeyqxq9u5p36jaj3npsn00ffrka86qya":344.74,"bnb1v4xnwl2n9m2cpyzjcplrfajgesh4s6he0fjsyc":344.74,"bnb1v7np68j0mh4xg53rhyrze7shfkgj4tt8psyggk":2671.72,"bnb1v9kve6evmxy5xgy83qj65nw6d5ln6ctceeszp3":1292.77,"bnb1vqq77mmr3d3kljklqc4fsr2q49xewmtdxptm2w":2671.72,"bnb1vv4wlfnptaxurfuc7lplzn88ggudfe92srkpew":344.74,"bnb1vzs70slkdhnxalg0ttdrwd932xysjpextvs8nj":2585.54,"bnb1w8dezrf7p5a6sggscec9ufqnj70qzs8y9w2h0c":344.74,"bnb1wgecw07u76mwczgkfg4qw63z5w2jaezhnyekcf":2671.72,"bnb1x0m098cktr07f3c3ppgtxfvy959f4hvduuj3cr":344.74,"bnb1x3v2e3l028sflpfj7am5fmy6rxt28mvywp3gl8":344.74,"bnb1xa0g586un43rms7g8pznxj6ehhr0txfg4xetal":344.74,"bnb1y20hcr7nr2hglqtnjf2xd3atvtfszl4e29km2g":344.74,"bnb1y8pj62wgcp9v8vzqrt2ukla38tcsrhh255fee3":344.74,"bnb1y9gyamnx7jgel5sjvfnlpdh450s5hmfejfmtsy":2930.2799999999997,"bnb1yh44ye7zzcwwuupnmwnl4eqrnz9ucegn2wax05":344.74,"bnb1ykl62mn8cnc65yg689pfylnh059922wu2stfxw":344.74,"bnb1zvd8dlq944ddncsdn0dmkjwda6cew3xx6gkwdh":344.74}`
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &receivers); err != nil {
			t.Fatalf("unmarshal json fail, err is %s\n", err.Error())
		}
	}

	for k, v := range receivers {
		if toAddr, err := types.AccAddressFromBech32(k); err != nil {
			t.Fatalf("invalid addr: %s\n", k)
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
		option := transaction.WithMemo("airdrop from irisnet ama")

		if sendResult, err := dexClient.SendToken(msgSendToken, true, option); err != nil {
			fmt.Printf("send token occur error, toAddr is %s, err is %s\n", v.ToAddr.String(), err.Error())
		} else {
			amt := float64(v.Coins.AmountOf(denom)) / math.Pow10(8)
			if sendResult.Ok {
				fmt.Printf("success, %s:%f, txHash: %s\n", v.ToAddr.String(), amt, sendResult.Hash)
			} else {
				fmt.Printf("failed, %s:%f, txHash: %s, log: %s\n", v.ToAddr.String(), amt, sendResult.Hash, sendResult.Log)
			}
		}
		fmt.Println("now sleep 1 seconds")
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func TestGetTxDetail(t *testing.T) {
	rpcClient = defaultClient()
	txHash := "F9AE66EC24E9D90394CB87AE4D19D98A67F3062139F63F1EF45FED140CB631B1"
	txHashBytes, err := hex.DecodeString(txHash)
	if err != nil {
		t.Fatal(err)
	}
	resultTx, err := rpcClient.Tx(txHashBytes, false)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("tx result is:\n %s\n", marshalJsonIgnoreError(resultTx))
	}

	txDataStr := getTxDataStr(resultTx.Tx.String())
	txData, err := hex.DecodeString(txDataStr)
	if err != nil {
		t.Fatal(err)
	}

	if stdTx, err := parseTxToStdTx(txData); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("txdata, %s\n", marshalJsonIgnoreError(stdTx))
		msgs := stdTx.Msgs
		if len(msgs) > 0 {
			txMsg := msgs[0]
			switch txMsg.(type) {
			case msg.SendMsg:
				txMsg := txMsg.(msg.SendMsg)
				t.Logf("tx msg is:\n %s\n", marshalJsonIgnoreError(txMsg))
				break
			default:
				t.Log("unknown tx msg")
			}
		}
	}
}

func TestSubscribe(t *testing.T) {
	rpcClient = defaultClient()
	//receiverAddr := "faa1meawt87fugh4m040zwktuma26xr0e9wl2laa5e"
	//q := fmt.Sprintf("tm.event = 'Tx' AND recipient = '%s'", receiverAddr)
	//rpcClient := defaultClient()
	q := fmt.Sprintf("tm.event = 'Tx'")

	fmt.Println(q)
	fmt.Println(query.MustParse(q).String())
	outChainRes, err := rpcClient.Subscribe(query.MustParse(q).String(), 10)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			select {
			case res := <-outChainRes:
				//eventDataTx := res.Data.(tmtypes.EventDataNewBlock)
				//fmt.Printf("subscribe a new block: %d\n", eventDataTx.Block.Height)
				eventDataTx := res.Data.(tmtypes.EventDataTx)
				fmt.Printf("txData is:\n %s\n", marshalJsonIgnoreError(eventDataTx))

				txDataStr := getTxDataStr(eventDataTx.TxResult.Tx.String())
				txData, err := hex.DecodeString(txDataStr)
				if err != nil {
					fmt.Println(err.Error())
					break
				}

				if stdTx, err := parseTxToStdTx(txData); err != nil {
					fmt.Println(err.Error())
					break
				} else {
					msgs := stdTx.Msgs
					if len(msgs) > 0 {
						txMsg := msgs[0]
						switch txMsg.(type) {
						case msg.SendMsg:
							txMsg := txMsg.(msg.SendMsg)
							fmt.Printf("tx msg is:\n %s\n", marshalJsonIgnoreError(txMsg))
							break
						default:
							fmt.Printf("unknown tx msg\n")
						}
					}
				}
			}
		}
	}()
	time.Sleep(10 * time.Minute)
}

func getTxDataStr(txStr string) string {
	prefix := "Tx{"
	suffix := "}"
	return strings.TrimSuffix(strings.TrimPrefix(txStr, prefix), suffix)
}

func parseTxToStdTx(txBytes []byte) (tx.StdTx, error) {
	var txInfo tx.StdTx
	txStructure, err := rpc.ParseTx(tx.Cdc, txBytes)
	if err != nil {
		return txInfo, err
	}

	switch txStructure.(type) {
	case tx.StdTx:
		txInfo = txStructure.(tx.StdTx)
		return txInfo, nil
	default:
		return txInfo, fmt.Errorf("unkonwn txStructure")
	}
}

func marshalJsonIgnoreError(v interface{}) string {
	data, _ := json.Marshal(v)
	return string(data)
}
