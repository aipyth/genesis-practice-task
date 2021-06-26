package storage

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
)

const perm = os.FileMode(0660)

type CSVStorage struct {
	filename string
	data map[string][]string
}

func NewCSVStorage(
	filename string) *CSVStorage {
	c := new(CSVStorage)
	c.filename = filename
	c.data = make(map[string][]string)
	return c
}

func (c *CSVStorage) read() {
	file, err := os.OpenFile(c.filename, os.O_RDONLY, perm)
	if err != nil {
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	var record []string
	for err != io.EOF {
		record, err = r.Read()
		if err != io.EOF {
			c.data[record[0]] = record[1:]
		}
	}
}

func (c *CSVStorage) Connect() error {
	c.read()
	return nil
}

func (c *CSVStorage) Disconnect() error {
	return nil
}

func (c *CSVStorage) Set(key string, value interface{}) error {
	val1, ok1 := value.(string)
	if ok1 {
		c.data[key] = []string{val1}
	}
	val2, ok2 := value.([]string)
	if ok2 {
		c.data[key] = val2
	}

	if !(ok1 || ok2) {
		return errors.New("cannot write not string or []string data")
	}
	return nil
}

func (c *CSVStorage) Get(key string) (val interface{}, err error) {
	data := c.data[key]
	if data == nil {
		return nil, nil
	}

	return interface{}(data), nil
}

func (c *CSVStorage) Delete(key string) error {
	c.data[key] = nil
	return nil
}

func (c *CSVStorage) Save() error {
	fmt.Println(c.data)
	return c.flush()
}

func (c *CSVStorage) flush() error {
	file, err := os.OpenFile(c.filename, os.O_WRONLY | os.O_CREATE | os.O_TRUNC,
		perm)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	allRecords := make([][]string, 0)

	for key, val := range c.data {
		if val != nil {
			record := []string{key}
			record = append(record, val...)
			allRecords = append(allRecords, record)
		}
	}

	fmt.Println(allRecords)

	err = writer.WriteAll(allRecords)
	if err != nil {
		return err
	}
	return nil
}
