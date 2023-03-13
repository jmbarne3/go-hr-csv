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
	args := os.Args[1:]

	if len(args) == 0 {
		log.Fatalln("You must include a file path as an argument")
	}

	xlfpath := args[0]

	xlf, err := excelize.OpenFile(xlfpath)
	if err != nil {
		log.Fatal(err)
	}

	rows, err := xlf.Rows("Sheet1")
	if err != nil {
		log.Fatal(err)
	}

	var records []record

	process_list(rows, &records)
	no_ops := filter_results(&records)

	var wg sync.WaitGroup
	wg.Add(2)

	go write_csv_file(&no_ops, "./without.csv", &wg)
	go write_csv_file(&records, "./with.csv", &wg)

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
		if len(row) < 14 {
			// Not enough columns, sp we're missing data
			i++
			continue
		}

		first_name := row[3]
		last_name := row[2]
		email := row[13]
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

func filter_results(records *[]record) []record {
	jfg_map := map[string]bool{
		"Faculty":                       true,
		"OPS":                           false,
		"Administrative & Professional": true,
		"Contingent Workers":            false,
		"Executive Service":             true,
		"UCF Athletic Association":      true,
		"USPS":                          true,
	}

	var retval []record

	for _, r := range *records {
		if v, exists := jfg_map[r.job_family]; exists && v {
			retval = append(retval, r)
		}
	}

	return retval
}

func write_csv_file(records *[]record, filepath string, wg *sync.WaitGroup) {

	f, err := os.Create(filepath)

	if err != nil {
		log.Fatalln("Failed to open file", err)
	}

	defer f.Close()

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
