package chinese

import "testing"

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
		char   rune
		pinyin string
		itemID int
	}{
		{[]rune("這")[0], "zhe4", 4},
		{[]rune("是")[0], "shi4", 3},
		{[]rune("我")[0], "wo3", 1},
		{[]rune("的")[0], "de5", 6},
		{[]rune("你")[0], "ni3", 0},
		{[]rune("很")[0], "hen3", 7},
		{[]rune("國")[0], "guo2", 5},
		{[]rune("日")[0], "ri4", 2},
	}

	var translations Translations

	for _, table := range tables {
		translations.Items = append(translations.Items, Item{ID: table.itemID, Pinyin: table.pinyin, Simplified: table.char, Traditional: table.char})
	}

	sortByFrequency(translations)

	for i, item := range translations.Items {
		if i != item.ID {
			t.Errorf("sortbyFrequency was incorrect. for character %v, expected position %v, got %v", item.Traditional, item.ID, i)
		}
	}

}
