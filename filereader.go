package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
)

type Path struct {
	id_ride   int
	lat, lng  float64
	timestamp int64
}

func GetPathFromRecord(record []string) Path {
	id_ride, err := strconv.ParseInt(record[0], 10, 0)
	checkError(fmt.Sprintf("Invalid data in input CSV when parsing 'id_ride': '%v'", record), err)
	lat, err := strconv.ParseFloat(record[1], 64)
	checkError(fmt.Sprintf("Invalid data in input CSV when parsing 'lat': '%v'", record), err)
	lng, err := strconv.ParseFloat(record[2], 64)
	checkError(fmt.Sprintf("Invalid data in input CSV when parsing 'lng': '%v'", record), err)
	timestamp, err := strconv.ParseInt(record[3], 10, 64)
	checkError(fmt.Sprintf("Invalid data in input CSV when parsing 'timestamp': '%v'", record), err)
	return Path{int(id_ride), lat, lng, timestamp}
}

var results [][2]string

func ReadCsv(r *csv.Reader) {
	paths := make(map[int]Path)
	count := 0
	last_path := Path{id_ride: -1}
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		path := GetPathFromRecord(record)

		// data guaranteed to not be multiplexed, so different id means new block
		if last_path.id_ride != -1 && path.id_ride != last_path.id_ride {
			go CalculateResult(last_path, paths)
			count = 0
			paths = make(map[int]Path)
		}
		paths[count] = path
		last_path = path
		count += 1
	}
}

func CalculateResult(last_path Path, paths map[int]Path) {
	res := [2]string{strconv.Itoa(last_path.id_ride), fmt.Sprintf("%.2f", float32(EstimateFare(FilterPaths(&paths), last_path.id_ride))/100.0)}
	results = append(results, res)
}

func ProcessFile(path string) error {
	inFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer inFile.Close()

	r := csv.NewReader(inFile)
	ReadCsv(r)

	WriteResults()

	return nil
}

func WriteResults() {
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	sort.Slice(results, func(i, j int) bool { return results[i][0] < results[j][0] })

	for i := 0; i < len(results); i++ {
		err := writer.Write(results[i][:])
		checkError("Cannot write to file", err)
	}
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
