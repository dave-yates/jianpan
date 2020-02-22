package chinese

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
)

func TestNewItem(t *testing.T) {
	testItem := Item{Pinyin: "shi", Simplified: []rune("是")[0], Traditional: []rune("是")[0]}

	resultItem := NewItem("shi", []rune("是")[0], []rune("是")[0])

	if testItem != resultItem {
		t.Errorf("NewItem was incorrect. Expected %v, got: %v", testItem, resultItem)
	}
}

func TestNewResultItem(t *testing.T) {
	testItem := Item{ID: 123, Pinyin: "shi", Simplified: []rune("是")[0], Traditional: []rune("是")[0]}

	resultItem := newResultItem(123, "shi", []rune("是")[0], []rune("是")[0])

	if testItem != resultItem {
		t.Errorf("newResultItem was incorrect. Expected %v, got: %v", testItem, resultItem)
	}
}

// func TestTranslate(t *testing.T) {
// 	tables := []struct {
// 		search string
// 		result Result
// 	}{
// 		{"wo3", Result{Search: "wo3", Chars: []string{"我"}}},
// 		//{"wo", Result{Search: "wo", Chars: []string{"我", "握", "臥", "窩", "沃"}}},
// 		//{"what", Result{Search: "what", Chars: []string{}}},
// 	}

// 	for _, table := range tables {
// 		result, err := translate(table.search, mockdbGetTranslations)

// 		if err != nil {
// 			t.Errorf("translate was incorrect. For search %v, expected: %v, got error: %v", table.result.Search, result.Chars, err)
// 			continue
// 		}

// 		if len(table.result.Chars) != len(result.Chars) {
// 			t.Errorf("translate was incorrect. Expected %v results, but got %v results for search: %v", len(table.result.Chars), len(result.Chars), table.result.Search)
// 			continue
// 		}

// 		for i, char := range table.result.Chars {
// 			if char != result.Chars[i] {
// 				t.Errorf("translate was incorrect. Expected: %v, got: %v. Full results %v", char, result.Chars[i], result)
// 			}
// 		}
// 	}

// }

func mockdbGetTranslations(ctx context.Context, search string) ([]bson.D, error) {

	var results []bson.D
	if search == "wo3" {

		// chars := []string{"我"}

		// for _, char := range chars {

		// 	value := bson.D{
		// 		{Key: "_id", Value: int(11)},
		// 		{Key: "pinyin", Value: "wo3"},
		// 		{Key: "simplified", Value: char},
		// 		{Key: "traditional", Value: char},
		// 	}
		// 	results = append(results, value)
		// }

		return results, nil
	}
	if search == "wo" {
		return results, nil
	}
	if search == "what" {
		return results, nil
	}

	return nil, fmt.Errorf("Network error")
}

func TestValidateSearch(t *testing.T) {
	tables := []struct {
		search string
		result bool
	}{
		{"w", true},
		{"wo", true},
		{"wo3", true},
		{"", false},
	}

	for _, table := range tables {
		ok, _ := validateSearch(table.search)
		if ok != table.result {
			t.Errorf("validateSearch is incorrect. Expected %v, got: %v", table.result, ok)
		}
	}
}

func TestGetResults_Traditional(t *testing.T) {
	tables := []struct {
		trad   rune
		simp   rune
		pinyin string
		itemID int
	}{
		{[]rune("這")[0], []rune("这")[0], "zhe4", 4},
		{[]rune("是")[0], []rune("是")[0], "shi4", 3},
		{[]rune("我")[0], []rune("我")[0], "wo3", 1},
		{[]rune("的")[0], []rune("的")[0], "de5", 6},
		{[]rune("你")[0], []rune("你")[0], "ni3", 0},
		{[]rune("很")[0], []rune("很")[0], "hen3", 7},
		{[]rune("國")[0], []rune("国")[0], "guo2", 5},
		{[]rune("日")[0], []rune("日")[0], "ri4", 2},
	}

	var translations Translations
	var testResults Result

	for _, table := range tables {
		translations.Items = append(translations.Items, Item{ID: table.itemID, Pinyin: table.pinyin, Simplified: table.simp, Traditional: table.trad})
		testResults.Chars = append(testResults.Chars, string(table.trad))
	}

	results := getResults(translations, false)

	for i, char := range testResults.Chars {
		if char != results.Chars[i] {
			t.Errorf("getResults is incorrect. Expected %v, got: %v", testResults, results)
			break
		}
	}

}

func TestGetResults_Simplified(t *testing.T) {
	tables := []struct {
		trad   rune
		simp   rune
		pinyin string
		itemID int
	}{
		{[]rune("這")[0], []rune("这")[0], "zhe4", 4},
		{[]rune("是")[0], []rune("是")[0], "shi4", 3},
		{[]rune("我")[0], []rune("我")[0], "wo3", 1},
		{[]rune("的")[0], []rune("的")[0], "de5", 6},
		{[]rune("你")[0], []rune("你")[0], "ni3", 0},
		{[]rune("很")[0], []rune("很")[0], "hen3", 7},
		{[]rune("國")[0], []rune("国")[0], "guo2", 5},
		{[]rune("日")[0], []rune("日")[0], "ri4", 2},
	}

	var translations Translations
	var testResults Result

	for _, table := range tables {
		translations.Items = append(translations.Items, Item{ID: table.itemID, Pinyin: table.pinyin, Simplified: table.simp, Traditional: table.trad})
		testResults.Chars = append(testResults.Chars, string(table.simp))
	}

	results := getResults(translations, true)

	for i, char := range testResults.Chars {
		if char != results.Chars[i] {
			t.Errorf("getResults is incorrect. Expected %v, got: %v", testResults, results)
			break
		}
	}

}

func TestJsonConvert(t *testing.T) {

	var results Result
	results.Search = "sh"
	results.Chars = []string{"是", "時", "事", "市"}

	expJSON, _ := json.Marshal(results)

	json, err := jsonConvert(results)

	for i, b := range json {

		if b != expJSON[i] {
			t.Errorf("Json result incorrect, expected: %v, got: %v", expJSON, json)
			break
		}
	}

	if err != nil {
		t.Errorf("Error returned. Expected no error. Error message: %v", err)
	}
}

//inject function to test failure??
// func TestJsonConvert_Error(t *testing.T) {

// 	var results Result
// 	results.Search = "sh{{"
// 	results.Chars = []string{"是", "時", "事", "市"}

// 	expJSON, _ := json.Marshal(results)

// 	json, err := jsonConvert(results)

// 	for i, b := range json {

// 		if b != expJSON[i] {
// 			t.Errorf("Json result incorrect, expected: %v, got: %v", expJSON, json)
// 			break
// 		}
// 	}

// 	if err != nil {
// 		t.Errorf("Error returned. Expected no error. Error message: %v", err)
// 	}
// }

func TestNewCharacter(t *testing.T) {
	tables := []struct {
		char     string
		expected bool
	}{
		{"這", false},
		{"是", false},
		{"我", false},
		{"的", false},
		{"你", true},
		{"很", true},
		{"國", true},
		{"日", true},
	}

	chars := []string{"這", "是", "我", "全", "部", "的", "子"}

	for _, table := range tables {
		result := newCharacter(table.char, &chars)

		if result != table.expected {
			t.Errorf("newCharacter was incorrect, got: %v, want: %v.", result, table.expected)
		}
	}

}

func TestSortByFrequency(t *testing.T) {
	tables := []struct {
		trad   rune
		simp   rune
		pinyin string
		itemID int
	}{
		{[]rune("這")[0], []rune("这")[0], "zhe4", 4},
		{[]rune("是")[0], []rune("是")[0], "shi4", 3},
		{[]rune("我")[0], []rune("我")[0], "wo3", 1},
		{[]rune("的")[0], []rune("的")[0], "de5", 6},
		{[]rune("你")[0], []rune("你")[0], "ni3", 0},
		{[]rune("很")[0], []rune("很")[0], "hen3", 7},
		{[]rune("國")[0], []rune("国")[0], "guo2", 5},
		{[]rune("日")[0], []rune("日")[0], "ri4", 2},
	}

	var translations Translations

	for _, table := range tables {
		translations.Items = append(translations.Items, Item{ID: table.itemID, Pinyin: table.pinyin, Simplified: table.trad, Traditional: table.simp})
	}

	sortByFrequency(translations)

	for i, item := range translations.Items {
		if i != item.ID {
			t.Errorf("sortbyFrequency was incorrect. for character %v, expected position %v, got %v", item.Traditional, item.ID, i)
		}
	}

}
