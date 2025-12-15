package process

import (
	"encoding/csv"
	"fmt"
	"os"
	"tayrosagr/utility"
)

func (k *Krinica) ReadFile(f string) error {
	if err := k.workTxtFile(f); err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}

type Record struct {
	Cis   *utility.CisInfo
	Korob string
}

func NewRecord(row []string) (*Record, error) {
	if len(row) != 2 {
		return nil, fmt.Errorf("wrong record in csv file")
	}
	cis, err := utility.ParseCisInfo(row[1])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return &Record{
		Cis:   cis,
		Korob: row[0],
	}, nil
}

func (k *Krinica) workTxtFile(file string) error {
	ar, err := readCsvFile(file)
	if err != nil {
		return fmt.Errorf("%w", err)
	}
	for _, record := range ar {
		rec, err := NewRecord(record)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
		if k.gtin == "" {
			k.gtin = rec.Cis.Gtin
		}
		if k.gtin != rec.Cis.Gtin {
			return fmt.Errorf("коробка %s марка %s GTIN %s отличается от других %s", rec.Korob, rec.Cis.Cis, rec.Cis.Gtin, k.gtin)
		}

		if _, exist := k.CisAll[rec.Cis.Cis]; exist {
			return fmt.Errorf("коробка %s марка %s дубль", rec.Korob, rec.Cis.Cis)
		}
		if _, exist := k.Pallet[rec.Korob]; !exist {
			k.Pallet[rec.Korob] = make([]*utility.CisInfo, 0)
		}
		k.Pallet[rec.Korob] = append(k.Pallet[rec.Korob], rec.Cis)
		k.Cis = append(k.Cis, rec.Cis)
	}
	for key, korob := range k.Pallet {
		if len(korob) != 24 {
			return fmt.Errorf("коробка %s кол-во единиц %d не равно 24", key, len(korob))
		}
	}
	return nil
}

func readCsvFile(filePath string) ([][]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.Comma = '\t'
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	return records, nil
}
