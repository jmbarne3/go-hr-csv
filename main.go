package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/xuri/excelize/v2"
)

type record struct {
	first_name     string
	last_name      string
	job_family     string
	email          string
	preferred_name string
}

func main() {
	xlf, err := excelize.OpenFile("./example-list.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := xlf.Rows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	var records []record

	process_list(rows, &records)

	var wg sync.WaitGroup
	wg.Add(2)

	go create_without_ops_list(&records, &wg)
	go create_with_ops_list(&records, &wg)

	wg.Wait()

	fmt.Println("All done!")
}

func process_list(rows *excelize.Rows, records *[]record) {
	var i int = 0

	for rows.Next() {
		if i >= 0 && i <= 2 {
			i++
			continue
		}

		row, _ := rows.Columns()
		if len(row) < 13 {
			// Not enough columns, sp we're missing data
			i++
			continue
		}

		first_name := row[3]
		last_name := row[2]
		email := row[12]
		preferred_name := row[5]
		job_family := row[8]

		nRecord := record{
			first_name:     first_name,
			last_name:      last_name,
			job_family:     job_family,
			email:          email,
			preferred_name: preferred_name,
		}

		*records = append(*records, nRecord)
	}
}

func create_without_ops_list(records *[]record, wg *sync.WaitGroup) {
	f, err := os.Create("./without.csv")
	defer f.Close()

	if err != nil {
		log.Fatalln("Failed to open file", err)
	}

	jfg_map := map[string]bool{
		"Faculty":                       true,
		"OPS":                           false,
		"Administrative & Professional": true,
		"Contingent Workers":            false,
		"Executive Service":             true,
		"UCF Athletic Association":      true,
		"USPS":                          true,
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{
		"first_name",
		"last_name",
		"email",
		"preferred_name",
	}

	w.Write(headers)

	for _, r := range *records {
		if v, exists := jfg_map[r.job_family]; exists && v {
			row := []string{
				r.first_name,
				r.last_name,
				r.email,
				r.preferred_name,
			}

			if err := w.Write(row); err != nil {
				log.Fatalln("error writing record to file", err)
			}
		}
	}

	wg.Done()
}

func create_with_ops_list(records *[]record, wg *sync.WaitGroup) {
	f, err := os.Create("./with.csv")
	defer f.Close()

	if err != nil {
		log.Fatalln("Failed to open file", err)
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{
		"first_name",
		"last_name",
		"email",
		"preferred_name",
	}

	w.Write(headers)

	for _, r := range *records {
		row := []string{
			r.first_name,
			r.last_name,
			r.email,
			r.preferred_name,
		}

		if err := w.Write(row); err != nil {
			log.Fatalln("error writing record to file", err)
		}
	}

	wg.Done()
}
