package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"encoding/csv"
	"crypto/md5"
	"encoding/binary"
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

func getHashKey(text string) uint64 {
	parts := strings.Split(text, ",")

	var numValues []uint8
	for _, numStr := range parts {
		num, err := strconv.ParseUint(strings.TrimSpace(numStr), 10, 8)
		if err != nil {
			panic(err)
		}
		numValues = append(numValues, uint8(num))
	}
	numValues = append(numValues, 0)
	numValues = append(numValues, 0)

	return binary.LittleEndian.Uint64([]byte(numValues))
}

func main() {

	fmt.Printf("[...] Loading table from CSV, please wait\n")
	found, err := restoreMapFromCSV("found.csv")
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	}
	importedCount := len(found)

	fmt.Printf("[...] Loaded %d entries from CSV\n", importedCount)

	uidBytes := ""
	if len(os.Args[1:]) > 0 {
		uidBytes = os.Args[1]
		fmt.Printf("[...] Looking up binding phrase for %s\n", uidBytes)
	}

	// Continuous mode
	if uidBytes == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Printf("Enter binding phrase: ")
		for scanner.Scan() {
			uidBytes = scanner.Text()
			hashKey := getHashKey(uidBytes)
			if val, ok := found[hashKey]; ok {
				fmt.Printf("[+++] Found binding phrase for %s: %s\n", uidBytes, val)
			} else {
				fmt.Printf("[---] Unable to find binding phrase for %s\n", uidBytes)
			}
			fmt.Printf("Enter binding phrase: ")
		}
	// One-shot
	} else {
		hashKey := getHashKey(uidBytes)
		if val, ok := found[hashKey]; ok {
			fmt.Printf("[+++] Found binding phrase for %s: %s\n", uidBytes, val)
			os.Exit(0)
		}
		fmt.Printf("[---] Unable to find binding phrase for %s\n", uidBytes)
	}
}
