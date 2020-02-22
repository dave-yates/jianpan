package chinese

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/dave-yates/jianpan/db"
	"go.mongodb.org/mongo-driver/bson"
)

//Translations is a slice of items. Each item is the pinyin and a slice of possible simplified and traditional characters
type Translations struct {
	Items []Item `bson:"items" json:"items"`
}

//Item contains an ID,  pinyin string and a slice of characters
type Item struct {
	ID          int    `bson:"_id" json:"itemid"`
	Pinyin      string `bson:"pinyin" json:"pinyin"`
	Simplified  rune   `bson:"simplified" json:"simplified"`
	Traditional rune   `bson:"traditional" json:"traditional"`
}

//NewItem creates a new Item from the input
func NewItem(pinyin string, simp rune, trad rune) Item {
	item := Item{Pinyin: pinyin, Simplified: simp, Traditional: trad}
	return item
}

//newResultItem creates a new Item including the ID field. Unexported.
func newResultItem(id int, pinyin string, simp rune, trad rune) Item {
	item := Item{ID: id, Pinyin: pinyin, Simplified: simp, Traditional: trad}
	return item
}

//TRANSLATION

//Result is the format for the json returned to the browser. The search string and an array of characters
type Result struct {
	Search string   `json:"search"`
	Chars  []string `json:"chars"`
}

//Translate takes a search string plus a context and returns a slice of chinese characters
func Translate(search string) ([]byte, error) {

	//validate input
	ok, error := validateSearch(search)

	if !ok {
		return nil, error
	}

	results, error := translate(search, db.GetTranslations)
	if error != nil {
		return nil, error
	}

	return jsonConvert(results)
}

func translate(search string, dbGetFunc func(context.Context, string) ([]bson.D, error)) (Result, error) {

	//context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	//translations := getTranslations(ctx, search)

	var translations Translations
	resultbson, error := dbGetFunc(ctx, search)
	if error != nil {
		return Result{}, error
	}

	for _, item := range resultbson {
		translations.Items = append(translations.Items, newResultItem(int(item[0].Value.(int32)), item[1].Value.(string), item[2].Value.(rune), item[3].Value.(rune)))
	}

	//setup character type properly
	simplified := false
	sortByFrequency(translations)
	results := getResults(translations, simplified)

	results.Search = search

	return results, nil
}

func validateSearch(search string) (bool, error) {

	if search == "" {
		return false, fmt.Errorf("Invalid Search. Search cannot be blank")
	}

	return true, nil
}

func getResults(translations Translations, simplified bool) Result {
	var results Result

	if simplified {
		for _, item := range translations.Items {
			character := string(item.Simplified)
			if newCharacter(character, &results.Chars) {
				results.Chars = append(results.Chars, character)
			}
		}
	} else {
		for _, item := range translations.Items {
			character := string(item.Traditional)
			if newCharacter(character, &results.Chars) {
				results.Chars = append(results.Chars, character)
			}
		}
	}

	return results
}

func jsonConvert(results Result) ([]byte, error) {

	json, err := json.Marshal(results)

	if err != nil {
		return nil, err
	}

	return json, nil
}

//TRANSLATIONS HELPERS
func newCharacter(character string, chars *[]string) bool {
	for _, char := range *chars {
		if char == character {
			return false
		}
	}
	return true
}

//SORTING
type byFreq []Item

func (a byFreq) Len() int           { return len(a) }
func (a byFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFreq) Less(i, j int) bool { return a[i].ID < a[j].ID }

func sortByFrequency(translations Translations) {
	sort.Sort(byFreq(translations.Items))
}
