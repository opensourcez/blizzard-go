package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/FuzzyStatic/blizzard/v1"
	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var CLIENT *blizzard.Client
var DB *sql.DB

func createTable() {
	createStudentTableSQL := `CREATE TABLE auctions (
		"id" integer NOT NULL PRIMARY KEY,		
		"price" integer,
		"quantity" integer,
		"item" integer,
		"created_at" integer, 		
		"updated_at" integer
	  );` // SQL Statement for Create Table

	statement, err := DB.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
}
func SaveAuction(id int, price int, quantity int, itemID int) {

	insertStudentSQL := `INSERT INTO auctions(id, price, quantity, item, created_at, updated_at) VALUES (?, ?, ?,?,?,?)`
	statement, err := DB.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		// log.Println("INSERT ERROR....")
		// log.Println(err)

		return
	}
	_, err = statement.Exec(id, price, quantity, itemID, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		insertStudentSQL2 := `UPDATE auctions SET updated_at = ? WHERE id = ? `
		statement2, err := DB.Prepare(insertStudentSQL2) // Prepare statement.
		_, err = statement2.Exec(time.Now().Unix(), id)
		if err != nil {
			// log.Println("POST STATE ERROR")
			log.Println("level one: ", err)
			return
		}
		return
		// log.Println("POST STATE ERROR")
	}
}

func getAuctions() {
	row, err := DB.Query("SELECT * FROM auctions ORDER BY created_at DESC")
	if err != nil {
		log.Println("SELECT ERROR....")
		log.Println(err)
		return
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var price int
		var quantity int
		var item int
		row.Scan(&id, &price, &quantity, &item)
		log.Println("Auction: ", id, " p:", price, " q:", quantity, " id:", item)
	}
}

// power 42319
// stamina 42321

// agility 42327
// stamina pot 42331
// mana 42326
// healing 42323
// int 42330
// strength 42332
// fire 42356
// fixation 42358
// anima 42358
func main() {
	CLIENT = blizzard.NewClient("caa67cccbc9a4babaee0e2f6e105fd81", "wO4iICGbHAmS1qB0iVSDd6PkSdGK6bfC", blizzard.EU, blizzard.EnUS)

	// log.Println("Creating sqlite-database.db...")
	// file, err := os.Create("sqlite-database.db") // Create SQLite file
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// file.Close()
	// var err error
	// DB, err = sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	// if err != nil {
	// 	panic(err)
	// }
	// defer DB.Close() // Defer Closing the database
	// // createTable()
	recipeMap := make(map[string]int)
	recipeMap["POWER FLASK"] = 42319
	recipeMap["STAMINA FLASK"] = 42321
	recipeMap["AGILITY POT"] = 42327
	recipeMap["STAMINA POT"] = 42331
	recipeMap["MANA POT"] = 42326
	recipeMap["HEALING POT"] = 42323
	recipeMap["INT POT"] = 42330
	recipeMap["STRENGTH POT"] = 42332
	recipeMap["PHANTOM FIRE"] = 42356
	recipeMap["FIXATION"] = 42358

	matsMap := make(map[int]int)

	xxx, _, _ := WoWConnectedRealm()
	for i, v := range recipeMap {
		res, _, _ := recipie(strconv.Itoa(v))
		log.Println()
		log.Println("===============================================================")
		log.Println("ITEM:", i)
		log.Println("..................................")
		var totalPrice float64 = 0
		for _, re := range res.Reagents {
			if re.Reagent.ID == 180732 {
				// ignore rune etched vial
				continue
			}
			matsMap[re.Reagent.ID] = re.Quantity
			price := GetItemPrice(xxx, re.Reagent.ID, 3)
			rePrice := price * float64(re.Quantity)
			log.Println("TOTAL:", rePrice, " // PER: ", price, " // QTY: ", re.Quantity, " // NAME: ", re.Reagent.Name)
			totalPrice += rePrice
		}

		ahPrice := GetItemPrice(xxx, res.CraftedItem.ID, 0)
		log.Println("..................................")
		log.Println("COSTS / SELLS:", totalPrice, ahPrice)
		ahCut := (ahPrice / 100) * 5
		ahprofit := (ahPrice - totalPrice) - ahCut
		log.Println("CUT / PROFIT:", ahCut, ahprofit)
		if ahprofit > 50 {
			log.Println("OMFG YOU CAN MAKE 50G PER ITEM HERE !!!!!!!!!!!!!!")
			log.Println("OMFG YOU CAN MAKE 50G PER ITEM HERE !!!!!!!!!!!!!!")
			log.Println("OMFG YOU CAN MAKE 50G PER ITEM HERE !!!!!!!!!!!!!!")
		}

		// FindPrice(matsMap, xxx, res.CraftedItem.ID)
		log.Println("===============================================================")
	}
	os.Exit(1)

	// for _, v := range xxx.Auctions {
	// 	time.Sleep(70 * time.Millisecond)
	// 	if v.UnitPrice < 1 {
	// 		SaveAuction(v.ID, v.Buyout, v.Quantity, v.Item.ID)
	// 	} else {
	// 		SaveAuction(v.ID, v.UnitPrice, v.Quantity, v.Item.ID)
	// 	}
	// }

	// matsMap := make(map[int]int)
	// matsMap[171315] = 3 // nightshade
	// matsMap[168586] = 4 // rising
	// matsMap[168589] = 4 // marrow
	// matsMap[168583] = 4 // widow
	// matsMap[170554] = 4 // vigil

	// log.Println("....................POWER......................")
	// FindPrice(matsMap, xxx, 171276)
	// log.Println("..........................................")

	// log.Println("....................STAMINA FLASK......................")
	// matsMap2 := make(map[int]int)
	// matsMap2[171315] = 1 // nightshade
	// matsMap2[168586] = 3 // rising
	// matsMap2[168589] = 3 // marrow
	// FindPrice(matsMap2, xxx, 171278)
	// log.Println("..........................................")

	// log.Println("....................... STRENGTH POT...................")
	// matsMap3 := make(map[int]int)
	// matsMap3[168586] = 5 // rising
	// FindPrice(matsMap3, xxx, 171275)
	// log.Println("..........................................")

	// log.Println("....................... AGI POT...................")
	// matsMap4 := make(map[int]int)
	// matsMap4[168583] = 5 // widow
	// FindPrice(matsMap4, xxx, 171270)
	// log.Println("..........................................")

	// log.Println("....................... EMBALM OIL...................")
	// matsMap5 := make(map[int]int)
	// matsMap5[169701] = 2 // widow
	// FindPrice(matsMap5, xxx, 171286)
	// log.Println("..........................................")

	// log.Println("....................... SHADOWCORE OIL...................")
	// matsMap6 := make(map[int]int)
	// matsMap6[169701] = 2 // widow
	// FindPrice(matsMap6, xxx, 171285)
	// log.Println("..........................................")
}

//https://eu.api.blizzard.com/data/wow/connected-realm/1306/auctions?namespace=dynamic-eu&locale=en_US&access_token=US7LSne6c3iEOwoOfvSAjqfvn6smcOZWRY
type xx struct {
	Auctions []Auction `json:"auctions"`
}
type Auction struct {
	ID        int    `json:"id"`
	Buyout    int    `json:"buyout"`
	UnitPrice int    `json:"unit_price"`
	Quantity  int    `json:"quantity"`
	TimeLeft  string `json:"time_left"`
	Item      struct {
		ID         int   `json:"id"`
		Context    int   `json:"context"`
		BonusLists []int `json:"bonus_lists"`
	} `json:"item"`
}

type Recipe struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Media struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		ID int `json:"id"`
	} `json:"media"`
	CraftedItem struct {
		Key struct {
			Href string `json:"href"`
		} `json:"key"`
		Name string `json:"name"`
		ID   int    `json:"id"`
	} `json:"crafted_item"`
	Reagents []struct {
		Reagent struct {
			Key struct {
				Href string `json:"href"`
			} `json:"key"`
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"reagent"`
		Quantity int `json:"quantity"`
	} `json:"reagents"`
	CraftedQuantity struct {
		Value int `json:"value"`
	} `json:"crafted_quantity"`
}

func WoWConnectedRealm() (*xx, []byte, error) {
	var (
		dat xx
		b   []byte
		err error
	)

	b, err = getURLBody("https://eu.api.blizzard.com"+fmt.Sprintf("/data/wow/connected-realm/%d/auctions?locale=%s", 1084, "en_US"), "dynamic-eu")
	if err != nil {
		return &dat, b, err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return &dat, b, err
	}

	return &dat, b, nil
}
func recipie(id string) (*Recipe, []byte, error) {
	var (
		dat Recipe
		b   []byte
		err error
	)

	b, err = getURLBody("https://eu.api.blizzard.com/data/wow/recipe/"+id+"?locale=en_US", "static-eu")
	if err != nil {
		return &dat, b, err
	}

	err = json.Unmarshal(b, &dat)
	if err != nil {
		return &dat, b, err
	}

	return &dat, b, nil
}

func getURLBody(url, namespace string) ([]byte, error) {
	var (
		req  *http.Request
		res  *http.Response
		body []byte
		err  error
	)

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	if namespace != "" {
		req.Header.Set("Battlenet-Namespace", namespace)
	}

	res, err = CLIENT.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	return body, nil
}

func FindPrice(matsMap map[int]int, xxx *xx, itemID int) {
	// matsCost := 0
	auctionMap := make(map[int][]Auction)
	// f, _ := os.Create("ITEMS.json")
	// f2, _ := os.Create("DATA.json")
	// defer f.Close()
	// _, _ = f2.WriteString(string(xx2))
	itemSellPrice := 999999999
	cost := 0
	for _, v := range xxx.Auctions {
		if v.Item.ID == itemID {
			if itemSellPrice > v.UnitPrice {
				itemSellPrice = v.UnitPrice
			}
			continue
		}
		_, ok := matsMap[v.Item.ID]
		if !ok {
			continue
		}
		auctionMap[v.Item.ID] = append(auctionMap[v.Item.ID], v)

	}
	for _, v := range auctionMap {
		var id int
		var price int = 999999999
		for _, iv := range v {
			id = iv.Item.ID
			if iv.Quantity < 3 {
				continue
			}
			if iv.UnitPrice < price {
				price = iv.UnitPrice
			}
		}
		cost += matsMap[id] * price
		log.Println("MATS ID:", id, " // COUNT NEEDED:", matsMap[id], " // COSTS PER ITEM:", price/10000)
	}
	log.Println("Cost to make:", cost/10000, " gold")
	log.Println("Sells for:", itemSellPrice/10000, " gold")
	log.Println("AH CUT:", ((itemSellPrice/10000)/100)*5, " gold")
	log.Println("PROFIT:", ((itemSellPrice/10000)-cost/10000)-((itemSellPrice/10000)/100)*5, " gold")
}
func GetItemPrice(xxx *xx, itemID int, minQuantity int) float64 {
	itemSellPrice := 999999999
	for _, v := range xxx.Auctions {
		if v.Item.ID == itemID {
			if v.Quantity < minQuantity {
				continue
			}
			if itemSellPrice > v.UnitPrice {
				itemSellPrice = v.UnitPrice
			}
		}

	}
	return (float64(itemSellPrice) / 10000)
}
