package main

import (
	"fmt"
	"time"
	"encoding/json"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"io/ioutil"
	"strings"
)

type Resp struct {
	Stock []struct{
		Number int `json:"StoreNumber"`
		MintStockCSS string `json:"MintStockStatusCss"`
		MintStock string `json:"MintStockStatus"`
		MintStockDesc string `json:"MintStockStatusDescription"`
		PreownedCSS string `json:"PreownedStockStatusCss"`
		Preowned string `json:"PreownedStockStatus"`
		PreownedDesc string `json:"PreownedStockStatusDescription"`
		Name string `json:"StoreName"`
		Distance string `json:"StoreMilesText"`
		Selected bool `json:"IsSelected"`
		Link string `json:"StoreFinderLink"`
		CanReserve bool `json:"ClickAndReserveStoreExclusion"`
		MintDifference bool `json:"MintClickAndReserveStorePriceDifferent"`
		PreownedDifference bool `json:"PreownedClickAndReserveStorePriceDifferent"`
		MintMessage string `json:"MintStockMessage"`
		PreownedMessage string `json:"PreowedStockMessage"`
	} `json:"StoreStockDetails"`
	Store struct {
		Number int `json:"StoreNumber"`
		Latitude float32 `json:"Latitude"`
		Longitude float32 `json:"Longitude"`
		Name string `json:"StoreName"`
		Address string `json:"Address"`
		Times []struct{
			Open string `json:"OpeningTime"`
			Close string `json:"ClosingTime"`
			Day string `json:"Day"`
		} `json:"StoreTimings"`
	} `json:"SelectedStoreDetails"`
}

func main() {
	os.Exit(_main(os.Args[1:]))
}

func _main(args []string) int {
	if len(args) >= 1 {
		if args[0] == "localDev" {
			err := godotenv.Load()
			if err != nil {
				fmt.Println(fmt.Errorf(".env: %w", err))
				return 1
			}
		}
	}

	err := runStock()
	if err != nil {
		fmt.Println(fmt.Errorf("stock: %w", err))
		return 1
	}	

	return 0
}

func getPostCodes(startNum int, limit int) []string {
	postCodes := []string{
		"AB10 1AN",
		"BB4 5DD",
		"KA6 5HR",
		"OX15 6EF",
		"LL57 1AS",
		"EX31 1BD",
		"LA13 0AQ",
		"CF62 5AE",
		"RG21 3BA",
		"B1 1DA",
		"PO20 1LN",
	}

	ret := []string{}

	for i := startNum; i < limit; i++ {
		ret = append(ret, postCodes[i])
	}

	return ret
}

func runStock() error {
// 	https://powerup.game.co.uk/StockCheckerUI/ExternalStockChecker/PostcodeStockDetailsJson?mintSku=265901&preownedSku=&postcode=CF10%202AR
	
	for j := 0; j < 500; j++ {
		codes := getPostCodes(j, 5)
	
		for i := 0; i < len(codes); i++ {
			sender := fmt.Sprintf("https://powerup.game.co.uk/StockCheckerUI/ExternalStockChecker/PostcodeStockDetailsJson?mintSku=265901&preownedSku=&postcode=%s", strings.Replace(codes[i], " ", "%20", -1))
			req , err := http.NewRequest("GET", sender, nil)
			if err != nil {
				return fmt.Errorf("runStock: %w", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("runStock: %w", err)
			}

			defer resp.Body.Close()
		
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("runStock: %w", err)
			}

			if resp.StatusCode == 200 {
				lr := &Resp{}
				jErr := json.Unmarshal(body, &lr)
				if jErr != nil {
					return fmt.Errorf("runStock: %w", jErr)
				}
	
				stockCheck(*lr)
			}

			time.Sleep(1 * time.Second)

		}

		j += 5
		time.Sleep(20 * time.Second)
	}

	return nil
}

func stockCheck(r Resp) {
	for i := 0; i < len(r.Stock); i++ {
		rp := r.Stock[i]
		if rp.MintStock != "NoStock" {
			fmt.Println(fmt.Sprintf("%v", rp))
		}
	}
}

