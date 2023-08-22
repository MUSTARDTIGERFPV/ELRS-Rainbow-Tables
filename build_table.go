// Builds a table of ExpressLRS binding phrases for a dictionary list
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"crypto/md5"
	"encoding/binary"
	"encoding/csv"
	"math/rand"
)

const PREFIX = "-DMY_BINDING_PHRASE=\""
const SUFFIX = "\""

func computeMD5Hash(input string) []byte {
	runes := []rune(input)
	data := []byte(string(runes))
	hash := md5.Sum(data)
	return hash[:]
}

func getFullUIDFromText(text string) []byte {
	// UIDs are hashed with the prefix and suffix included
	fullString := PREFIX + text + SUFFIX
	fullHash := computeMD5Hash(fullString)
	return fullHash
}

// Returns both the truncated byte array and a uint64 representation of the byte array
func getUIDFromText(text string) ([]byte, uint64) {
	truncatedHash := getFullUIDFromText(text)[:6]
	// Arbitrarily chose LE for the hash key
	return truncatedHash, binary.LittleEndian.Uint64(append(truncatedHash, 0, 0))
}

func generateRandomText(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func dumpMapToCSV(filename string, data map[uint64]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for key, value := range data {
		record := []string{strconv.FormatUint(key, 10), value}
		err := writer.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

func restoreMapFromCSV(filename string) (map[uint64]string, error) {
	data := make(map[uint64]string)

	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		key, err := strconv.ParseUint(record[0], 10, 64)
		if err != nil {
			return nil, err
		}
		value := record[1]
		data[key] = value
	}

	return data, nil
}

func readWordlist(filename string) (*bufio.Scanner, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	return scanner, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	found, err := restoreMapFromCSV("found.csv")
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}
	importedCount := len(found)

	fmt.Printf("Restored %d entries from CSV\n", importedCount)

	scanner, err := readWordlist("words.txt")
	if err != nil {
		panic(err)
	}
	for scanner.Scan() {
		bindingPhrase := scanner.Text()
		_, hashKey := getUIDFromText(bindingPhrase)
		if _, ok := found[hashKey]; !ok {
			found[hashKey] = bindingPhrase
		}
	}

	fmt.Printf("Discovered %d new entries\n", len(found) - importedCount)
	err = dumpMapToCSV("found.csv", found)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %d entries to CSV\n", len(found))
}
