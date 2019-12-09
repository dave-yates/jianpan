package dictionary

import (
	"encoding/json"
	"sort"
	"strings"
)

//Dictionary is a slice of items. Each item is the pinyin and a slice of possible simplified or traditional characters
type Dictionary struct {
	Items      []Item
	Simplified bool
	charCount  int
}

//Item contains a pinyin string and a slice of characters
type Item struct {
	Pinyin string
	Chars  []Char
}

//Char contains a rune for the simplified and traditional characters
type Char struct {
	Frequency   int
	Simplified  rune
	Traditional rune
}

//change data structure to this?
//or use external DB?
/*type Dictionary struct {
	Items      []Item
	Simplified bool
	charCount  int
}

type Item struct {
	Pinyin string
	Frequency   int
	Simplified  rune
	Traditional rune
}*/

//TRANSLATION

//Translate takes a letter and returns a slice of chinese characters
func (dict *Dictionary) Translate(input string) ([]byte, error) {

	//validate input
	if input == "" {
		return nil, nil
	}

	translations := dict.getTranslation(input)

	sorted := sortTranslations(translations, dict.Simplified)

	return jsonConvert(input, sorted)
}

func (dict *Dictionary) getTranslation(input string) []Char {

	var results []Char

	for _, item := range dict.Items {
		if strings.HasPrefix(item.Pinyin, input) {
			for _, char := range item.Chars {
				//continue to next char if char is already in results
				if resultExists(char.Traditional, &results) {
					continue
				}
				results = append(results, char)
			}
		}
	}

	return results
}

//TRANSLATIONS HELPERS
func resultExists(tradChar rune, results *[]Char) bool {
	for _, char := range *results {
		if char.Traditional == tradChar {
			return true
		}
	}
	return false
}

func sortTranslations(translations []Char, simplified bool) []string {
	sorted := []string{}

	sortByFrequency(translations)

	if simplified {
		for _, char := range translations {
			sorted = append(sorted, string(char.Simplified))
		}
	} else {
		for _, char := range translations {
			sorted = append(sorted, string(char.Traditional))
		}
	}

	return sorted
}

type translation struct {
	Pinyin string   `json:"pinyin"`
	Chars  []string `json:"chars"`
}

func jsonConvert(pinyin string, chars []string) ([]byte, error) {

	out := translation{pinyin, chars}
	json, err := json.Marshal(out)

	if err != nil {
		return nil, err
	}

	return json, nil
}

//ADD TO DICTIONARY

//AddtoDictionary takes a string and two runes and creates a dictionary entry
func (dict *Dictionary) AddtoDictionary(pinyin string, simpChar rune, tradChar rune) {

	if i, found := dict.pinyinExists(pinyin); found {
		if !dict.charExists(i, tradChar) {
			dict.addCharacter(i, simpChar, tradChar)
		}
		return
	}

	newItem := Item{Pinyin: pinyin}
	newItem.Chars = make([]Char, 0)
	dict.Items = append(dict.Items, newItem)
	//possibly strange
	dict.AddtoDictionary(pinyin, simpChar, tradChar)
}

func (dict *Dictionary) addCharacter(i int, simpChar rune, tradChar rune) {

	dict.charCount++

	newChar := Char{dict.charCount, simpChar, tradChar}
	dict.Items[i].Chars = append(dict.Items[i].Chars, newChar)
}

//ADD TO DICTIONARY HELPERS
func (dict *Dictionary) pinyinExists(pinyin string) (int, bool) {
	for i := range dict.Items {
		if dict.Items[i].Pinyin == pinyin {
			return i, true
		}
	}
	return 0, false
}

func (dict *Dictionary) charExists(i int, tradChar rune) bool {
	for _, char := range dict.Items[i].Chars {
		if char.Traditional == tradChar {
			return true
		}
	}
	return false
}

//SORTING
type byPinyin []Item

func (a byPinyin) Len() int           { return len(a) }
func (a byPinyin) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPinyin) Less(i, j int) bool { return a[i].Pinyin < a[j].Pinyin }

//SortByPinyin sorts the dictionary items from a-z by pinyin
func (dict *Dictionary) SortByPinyin() {
	sort.Sort(byPinyin(dict.Items))
}

type byFreq []Char

func (a byFreq) Len() int           { return len(a) }
func (a byFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byFreq) Less(i, j int) bool { return a[i].Frequency < a[j].Frequency }

func sortByFrequency(translations []Char) {
	sort.Sort(byFreq(translations))
}
