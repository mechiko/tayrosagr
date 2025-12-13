package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"tayrosagr/utility"
)

const inDir = `IN`
const outDir = `OUT`

func main() {
	var err error
	var files []string

	err = os.RemoveAll(outDir)
	if err != nil {
		fmt.Printf("Error removing directory: %v\n", err)
		return
	}

	exitMsg := func(title, msg string) {
		utility.MessageBox(title, msg)
		os.Exit(1)
	}

	// 2. Recreate the empty directory
	err = os.MkdirAll(outDir, 0755) // 0755 provides read/write/execute permissions for owner, read/execute for group and others
	if err != nil {
		fmt.Printf("Error recreating directory: %v\n", err)
		return
	}
	ok := utility.MessageBox32("Ошибка", "test32")
	utility.MessageBox("Ошибка", "test")
	fmt.Println(ok)
	if _, err := utility.DialogOpenFile([]utility.FileType{utility.Csv, utility.Txt}, "", "."); err != nil {
		exitMsg("Ошибка", err.Error())
	}

	_, err = startDialog()
	if err != nil {
		exitMsg("Ошибка", err.Error())
	}
	re, err := regexp.Compile(`.*\.csv$`)
	if err != nil {
		exitMsg("Ошибка", err.Error())
	}
	if files, err = utility.FilteredSearchOfDirectoryTree(re, inDir); err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	for _, file := range files {
		if err := workTxtFile(file); err != nil {
			log.Fatal(err)
		}
	}

}

var Koroba map[string][]string
var CIS map[string]string

type Record struct {
	Cis   string
	Korob string
}

func NewRecord(row []string) (*Record, error) {
	if len(row) != 2 {
		return nil, fmt.Errorf("wrong record in csv file")
	}
	cis := row[1]
	// indexGS := strings.Index(cis, "\x1D")
	// if indexGS > 0 {
	// 	cis = cis[:indexGS]
	// }
	if len(cis) < 26 {
		return nil, fmt.Errorf("wrong cis lenght %d %s", len(cis), cis)
	}
	// cis = cis[:25]
	cis = cis[:len(cis)-6]
	return &Record{
		Cis:   cis,
		Korob: row[0],
	}, nil
}

func trimCis(s string) string {
	indexGS := strings.Index(s, "\x1D")
	cis := s
	if indexGS > 0 {
		cis = s[:indexGS]
	}
	return cis
}

func workTxtFile(file string) error {
	name := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	arrKorob := map[string][]string{}
	ar := readCsvFile(file)
	arrCsv := make([][]string, 0)
	for i, record := range ar {
		if i == 0 {
			continue
		}
		rec, err := NewRecord(record)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if _, exist := CIS[rec.Cis]; exist {
			fmt.Printf("file %s korob %s KM %s doube\n", name, rec.Korob, rec.Cis)
		}
		if _, exist := arrKorob[rec.Korob]; !exist {
			arrKorob[rec.Korob] = make([]string, 0)
		}
		arrKorob[rec.Korob] = append(arrKorob[rec.Korob], rec.Cis)
		arrCsv = append(arrCsv, []string{rec.Korob, rec.Cis})
	}
	for key, korob := range arrKorob {
		if len(korob) != 24 {
			fmt.Printf("file %s korob %s len %d\n", name, key, len(korob))
		}
	}
	saveCsvCustom("koroba_"+name+".csv", arrCsv)
	return nil
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func readStringArray(filePath string) []string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	arr := make([]string, 0)
	scanner := bufio.NewScanner(f)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		txt := scanner.Text()
		arr = append(arr, txt)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return arr
}
