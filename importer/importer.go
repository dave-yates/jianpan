package importer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/dave-yates/jianpan/chinese"
	"github.com/dave-yates/jianpan/db"
	"go.mongodb.org/mongo-driver/bson"
)

var counter int

//Import reads text files in the resources folder and inserts into mongodb
func Import(ctx context.Context) error {

	counter = 0

	//import
	for i := 1; i <= 6; i++ {
		err := importFromFile(ctx, fmt.Sprintf("resources/HSK%v.txt", i))
		if err != nil {
			return err
		}
	}

	return db.EnsureIndex(ctx)
}

func importFromFile(ctx context.Context, filename string) error {

	fmt.Printf("reading from file: %v\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "\t")
		translations := processEntry(s)
		addTranslations(ctx, translations)
	}

	return nil
}

func addTranslations(ctx context.Context, translations chinese.Translations) {
	for _, item := range translations.Items {

		value := bson.D{
			{Key: "_id", Value: counter},
			{Key: "pinyin", Value: item.Pinyin},
			{Key: "simplified", Value: item.Simplified},
			{Key: "traditional", Value: item.Traditional},
		}

		counter++
		db.Insert(ctx, value)
	}
}

type input struct {
	entries []entry
}

type entry struct {
	pinyin   string
	simpChar string
	tradChar string
}

func processEntry(s []string) chinese.Translations {

	var translations chinese.Translations
	translations.Items = make([]chinese.Item, 0)

	var results input
	results.entries = make([]entry, 0)

	//split on multiple pronuciations (i.e comma in the pinyin)
	results = processPronunciations(s)

	//split on phrases (i.e entry is multiple chars and not a single character)
	translations = processPhrases(results)

	//validate final result
	//number of runes
	//regex see processpronunciation commented out bit
	//see google doc

	//return results
	return translations
}

func processPronunciations(s []string) input {

	var results input

	//split on comma
	re := regexp.MustCompile(",")
	pinyin := re.Split(s[2], -1)

	for _, v := range pinyin {
		results.entries = append(results.entries, entry{v, s[0], s[1]})
	}

	return results
}

func processPhrases(in input) chinese.Translations {

	var results chinese.Translations

	for _, e := range in.entries {

		simpChars := []rune(e.simpChar)
		tradChars := []rune(e.tradChar)

		length := len(simpChars)

		//split pinyin into single words
		pinyins := processPinyin(e.pinyin, length)

		if length != len(tradChars) || length != len(pinyins) {
			panic(fmt.Sprintf("invalid input. Number of characters and number of pinyin translations do not match"))
		}

		for i := 0; i < length; i++ {
			results.Items = append(results.Items, chinese.NewItem(pinyins[i], simpChars[i], tradChars[i]))
		}
	}

	return results
}

func processPinyin(romanisation string, length int) []string {

	romanisation = strings.ToLower(strings.TrimSpace(romanisation))

	if length == 1 {
		return []string{romanisation}
	}

	//split on number to separate pinyin for characters
	re := regexp.MustCompile("[0-9]+")

	//this split leaves an empty trailing slice entry
	pinyin := re.Split(romanisation, length+1)
	tones := re.FindAllString(romanisation, length)

	var results []string

	for i := 0; i < length; i++ {

		//panic if the input contains any strange characters
		// matched, _ := regexp.MatchString(`^[a-z0-9]*$`, pinyin[i])
		// if !matched {
		// 	panic(fmt.Sprintf("invalid character in import file %v", pinyin))
		// }

		results = append(results, pinyin[i]+tones[i])
	}
	return results
}
