package main

import (
	"os"
	"bufio"
	"bytes"
	"io"
	"fmt"
	"strings"
	"strconv"
	"sort"
)

type AverageData struct {
	Sum float64
	Counter   float64
	UserMap map[string]int
	Amount []float64

}

type AverageType map[string]*AverageData

func readLines(path string) (lines []string, err error) {
	var (
		file *os.File
		part []byte
		prefix bool
	)
	if file, err = os.Open(path); err != nil {
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	buffer := bytes.NewBuffer(make([]byte, 0))
	for {
		if part, prefix, err = reader.ReadLine(); err != nil {
			break
		}
		buffer.Write(part)
		if !prefix {
			lines = append(lines, buffer.String())
			buffer.Reset()
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

func writeLines(lines [4]string, path string) (err error) {
	var (
		file *os.File
	)

	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	for _,item := range lines {
		fmt.Println(item)
		_, err := file.WriteString(strings.TrimSpace(item) + "\n")
		if err != nil {
			fmt.Println(err)
			break
		}
	}
	return
}



func proccessLog(lines []string) [4]string {
	typeMap := make(AverageType)
	typeMap["pago"] = &AverageData{0, 0, make(map[string]int), 	[]float64{}}
	typeMap["cobro"] = &AverageData{0, 0, make(map[string]int), []float64{}}
	typeMap["descuento"] = &AverageData{0, 0, make(map[string]int), []float64{}}
	typeMap["inversión"] = &AverageData{0, 0, make(map[string]int),[]float64{}}

	for _, line := range lines {
		if strings.Contains(line, "ammount:") && strings.Contains(line, "user:") &&
			strings.Contains(line, "type"){
			var splittedLine []string = strings.Split(line,"]")
			var users []string = strings.Split(splittedLine[0],":")
			var types []string = strings.Split(splittedLine[1],":")
			var amounts []string = strings.Split(splittedLine[2],":")
			if amount, err := strconv.ParseFloat(amounts[1], 64); err == nil {
				typeMap[types[1]].Counter++
				typeMap[types[1]].Sum += amount
				typeMap[types[1]].UserMap[users[1]]++
				typeMap[types[1]].Amount = append(typeMap[types[1]].Amount, amount)
			}
		}
	}
	var results [4]string = proccessResults(typeMap)
	return results
}

func proccessResults(averageType AverageType) [4]string {
	i := 0
	var results [4]string
	var max int
	var user string


	for key, value := range averageType {
		max = 0
		user = ""
		sort.Float64s(value.Amount)
		indexPercentil := len(value.Amount) * 95 / 100
		//fmt.Printf("percentil: %f", ))

		for keyUser, valueUser := range value.UserMap {
			if max < valueUser{
				max = valueUser
				user = keyUser
			}
		}
		results[i] = fmt.Sprintf("Average for type %s is = %f\nUser with max number of transactions = " +
			"User: %s. Transactions: %d.\nPercentil 95 = %f\n\n", key, float64(value.Sum) / float64(value.Counter),
			user, max, float64(value.Amount[indexPercentil]))
		i++
	}

	return results
}

func main() {

	lines, err := readLines("movements.log")

	results := proccessLog(lines)

	err = writeLines(results, "movements_result.log")
	if err != nil {
		fmt.Println(err)
	}

}