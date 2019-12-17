package chinese

import (
	"context"
	"encoding/json"

	"github.com/dave-yates/jianpan/db"
)

//Translations is a slice of items. Each item is the pinyin and a slice of possible simplified and traditional characters
type Translations struct {
	Items []Item `bson:"items" json:"items"`
}

//var Phrases []Phrase

//Item contains a pinyin string and a slice of characters
type Item struct {
	Pinyin string `bson:"pinyin" json:"pinyin"`
	//Frequency   int    `bson:"_id" json:"frequency"`
	Simplified  rune `bson:"simplified" json:"simplified"`
	Traditional rune `bson:"traditional" json:"traditional"`
}

//NewItem create a new Item from the input
func NewItem(pinyin string, simp rune, trad rune) Item {
	item := Item{Pinyin: pinyin, Simplified: simp, Traditional: trad}
	return item
}

//TRANSLATION

//Result is the format for the json returned to the browser. The search string and an array of characters
type Result struct {
	Search string   `json:"search"`
	Chars  []string `json:"chars"`
}

//Translate takes a search string plus a context and returns a slice of chinese characters
func Translate(ctx context.Context, search string) ([]byte, error) {

	//validate input
	if search == "" {
		return nil, nil
	}

	translations := getTranslation(ctx, search)

	//setup character type properly
	simplified := false
	results := sortTranslations(translations, simplified)

	return jsonConvert(results)
}

func getTranslation(ctx context.Context, search string) Translations {

	var results Translations

	resultbson := db.GetTranslations(ctx, search)

	for _, item := range resultbson {
		results.Items = append(results.Items, NewItem(item[1].Value.(string), item[2].Value.(rune), item[3].Value.(rune)))
	}
	return results
}

func sortTranslations(translations Translations, simplified bool) Result {
	var results Result

	//sortByFrequency(translations)

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

// //TRANSLATIONS HELPERS
func newCharacter(character string, chars *[]string) bool {
	for _, char := range *chars {
		if char == character {
			return false
		}
	}
	return true
}

// //SORTING
// type byPinyin []Item

// func (a byPinyin) Len() int           { return len(a) }
// func (a byPinyin) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a byPinyin) Less(i, j int) bool { return a[i].Pinyin < a[j].Pinyin }

// //SortByPinyin sorts the dictionary items from a-z by pinyin
// func (dict *Translations) SortByPinyin() {
// 	sort.Sort(byPinyin(dict.Items))
// }

// type byFreq []Char

// func (a byFreq) Len() int           { return len(a) }
// func (a byFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
// func (a byFreq) Less(i, j int) bool { return a[i].ID < a[j].ID }

// func sortByFrequency(translations []Char) {
// 	sort.Sort(byFreq(translations))
// }