package importer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/dave-yates/jianpan/dictionary"
)

//ImportData reads text files in the resources folder and builds the dictionary
func ImportData() dictionary.Dictionary {
	var dict dictionary.Dictionary

	//change to user controlled later
	dict.Simplified = false

	for i := 1; i <= 6; i++ {
		importFromFile(fmt.Sprintf("resources/HSK%v.txt", i), &dict)
	}

	//dict.SortByPinyin()

	//print for testing
	// for _, item := range dict.Items {
	// 	fmt.Printf("%v\n", item.Pinyin)
	// 	for _, char := range item.Chars {
	// 		fmt.Printf("\t%c\t%v\n", char.Traditional, char.Frequency)
	// 	}
	// }

	return dict
}

func importFromFile(filename string, dict *dictionary.Dictionary) {

	fmt.Printf("reading from file: %v\n", filename)

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "\t")
		processEntry(s, dict)

	}
}

func processEntry(s []string, dict *dictionary.Dictionary) {

	for i := range s {
		s[i] = strings.ToLower(strings.TrimSpace(s[i]))
	}

	//add characters with multiple pronunciations separately
	if strings.Contains(s[2], ",") {
		re := regexp.MustCompile(",")
		pinyin := re.Split(s[2], -1)

		for i := range pinyin {
			newStr := []string{s[0], s[1], pinyin[i]}
			addEntry(newStr, dict)
		}
		return
	}

	addEntry(s, dict)
}

func addEntry(s []string, dict *dictionary.Dictionary) {

	//clean up validation
	//how many runes?
	//check pinyin for lenth ==1 route
	//what about multiple pronunciations and multi chars?

	simpChars := []rune(s[0])
	tradChars := []rune(s[1])

	length := len(simpChars)
	//if it's only one character then add
	if length == 1 {
		dict.AddtoDictionary(s[2], simpChars[0], tradChars[0])
		return
	}

	//otherwise separate into single characters before adding
	//pinyin first
	romanisation := processPinyin(s, length)

	if len(simpChars) != len(tradChars) || len(simpChars) != len(romanisation) {
		panic(fmt.Sprintf("invalid input. Characters and translation lengths don't match"))
	}

	for i := 0; i < len(tradChars); i++ {
		dict.AddtoDictionary(romanisation[i], simpChars[i], tradChars[i])
	}

}

func processPinyin(s []string, length int) []string {

	//split on number to separate pinyin for characters
	re := regexp.MustCompile("[0-9]+")

	//this split leaves an empty trailing slice entry
	pinyin := re.Split(s[2], length+1)
	tones := re.FindAllString(s[2], length)

	var romanisation []string

	for i := 0; i < length; i++ {

		//panic if the input contains any strange characters
		// matched, _ := regexp.MatchString(`^[a-z0-9]*$`, pinyin[i])
		// if !matched {
		// 	panic(fmt.Sprintf("invalid character in import file %v", pinyin))
		// }

		romanisation = append(romanisation, pinyin[i]+tones[i])
	}
	return romanisation
}
