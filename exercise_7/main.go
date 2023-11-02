package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type IncomingJson struct {
	ID        string          `json:"id"`
	Latitude  float64         `json:"latitude"`
	Longitude float64         `json:"longitude"`
	Tags      json.RawMessage `json:"tags"`
}

type Config struct {
	path   string
	output string
}

func main() {
	var config Config
	flag.StringVar(&config.path, "config", ".json", "path to file.json")
	flag.StringVar(&config.output, "output", ".csv", "name of file to output")
	flag.Parse()

	//checking path for default value then read from executable directory or from stdin
	if config.path == ".json" {
		fileName, err := isExecutableDirectory()
		if err != nil {
			log.Printf("Error: %v", err)

			jsonByte, err := readFromStdIN()
			if err != nil {
				log.Printf("%v", err)
			}

			jsonData, err := parseJson(jsonByte)
			if err != nil {
				log.Printf("%v", err)
			}

			tagsMap := countTags(jsonData)
			outPut(config, tagsMap)

			return
		}

		//if .json file present in executable directory
		config.path = fileName
		goto readJson
	} else if config.path != ".json" {
		err := isFileExist(config.path)
		if err != nil {
			fileName, err := isExecutableDirectory()
			if err != nil {
				log.Printf("Error: %v", err)

				jsonByte, err := readFromStdIN()
				if err != nil {
					log.Printf("%v", err)
				}

				jsonData, err := parseJson(jsonByte)
				if err != nil {
					log.Printf("%v", err)
				}

				tagsMap := countTags(jsonData)
				outPut(config, tagsMap)

				return
			}

			//if .json file present in executable directory
			config.path = fileName
		}
	}

readJson:
	byteValue, err := readJsonFile(config.path)
	if err != nil {
		log.Printf("%v", err)
	}

	var users IncomingJson

	err = json.Unmarshal(byteValue, &users)
	if err != nil {
		fmt.Println(err)
	}

	tagsMap := countTags(users)
	outPut(config, tagsMap)
}

// isExecutableDirectory
func isExecutableDirectory() (string, error) {
	mainFileDir := filepath.Dir(".")

	files, err := ioutil.ReadDir(mainFileDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.ToLower(filepath.Ext(file.Name())) == ".json" {
			return file.Name(), nil
		}
	}

	return "", errors.New("there are no matching files for read")
}

// readFromStdIN
func readFromStdIN() ([]byte, error) {
	fmt.Println("Insert JSON:")

	jsonByte, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}

	return jsonByte, nil
}

// parseJson
func parseJson(jsonByte []byte) (IncomingJson, error) {
	var jsonData IncomingJson

	err := json.Unmarshal(jsonByte, &jsonData)
	if err != nil {
		return jsonData, err
	}

	return jsonData, nil
}

// countTags
func countTags(jsonData IncomingJson) map[string][]interface{} {
	tagsMap := make(map[string][]interface{})

	result := gjson.Parse(string(jsonData.Tags))

	result.ForEach(func(key, value gjson.Result) bool {
		tagName := key.String()
		tagValue := value.Value()

		tagsMap[tagName] = append(tagsMap[tagName], tagValue)

		return true
	})

	return tagsMap
}

// isFileExist
func isFileExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	extension := strings.ToLower(filepath.Ext(path))
	if extension == ".json" {
		return nil
	}

	return err
}

// readJsonFile
func readJsonFile(path string) ([]byte, error) {
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	log.Printf("Successfully Opened %v", jsonFile.Name())

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return byteValue, nil
}

// outPut
func outPut(config Config, tagsMap map[string][]interface{}) {
	if config.output == ".csv" {
		mainFileDir := filepath.Dir(".")

		files, err := ioutil.ReadDir(mainFileDir)
		if err != nil {
			log.Printf("%v", err)
		}

		for _, file := range files {
			if strings.ToLower(filepath.Ext(file.Name())) == ".csv" {
				config.output = file.Name()
				break
			}
		}
	}

	_, err := os.Stat(config.output)
	if os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, err)

		fmt.Println("tag,value,count")
		for tag, value := range tagsMap {
			if len(value) > 1 {
				localMap := make(map[interface{}]int)

				for _, key := range value {
					localMap[key]++
				}

				for i, j := range localMap {
					total := strconv.Itoa(j)

					switch i.(type) {
					case string:
						fmt.Println(tag, i.(string), total)
					case int:
						fmt.Println(tag, i.(int), total)
					case float64:
						fmt.Println(tag, strconv.FormatFloat(i.(float64), 'g', -1, 64), total)
					case float32:
						fmt.Println(tag, strconv.FormatFloat(float64(i.(float32)), 'g', -1, 32), total)
					case bool:
						fmt.Println(tag, strconv.FormatBool(i.(bool)), total)
					}
				}
				continue
			}

			switch value[0].(type) {
			case string:
				fmt.Println(tag, value[0].(string), "1")
			case int:
				fmt.Println(tag, strconv.Itoa(value[0].(int)), "1")
			case float64:
				fmt.Println(tag, strconv.FormatFloat(value[0].(float64), 'g', -1, 64), "1")
			case float32:
				fmt.Println(tag, strconv.FormatFloat(float64(value[0].(float32)), 'g', -1, 32), "1")
			case bool:
				fmt.Println(tag, strconv.FormatBool(value[0].(bool)), "1")
			}
		}
		return
	}

	outputFile, err := os.Create(config.output)
	if err != nil {
		fmt.Fprintln(os.Stderr, "file is not exist\n", err)
	}
	defer outputFile.Close()

	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()

	csvWriter.Write([]string{"tag", "value", "count"})
	for tag, value := range tagsMap {
		if len(value) > 1 {
			localMap := make(map[interface{}]int)

			for _, key := range value {
				localMap[key]++
			}

			for i, j := range localMap {
				total := strconv.Itoa(j)

				switch i.(type) {
				case string:
					err := csvWriter.Write([]string{tag, i.(string), total})
					if err != nil {
						fmt.Println(err)
					}
				case int:
					err := csvWriter.Write([]string{tag, strconv.Itoa(i.(int)), total})
					if err != nil {
						fmt.Println(err)
					}
				case float64:
					err := csvWriter.Write([]string{tag, strconv.FormatFloat(i.(float64), 'g', -1, 64), total})
					if err != nil {
						fmt.Println(err)
					}
				case float32:
					err := csvWriter.Write([]string{tag, strconv.FormatFloat(float64(i.(float32)), 'g', -1, 32), total})
					if err != nil {
						fmt.Println(err)
					}
				case bool:
					err := csvWriter.Write([]string{tag, strconv.FormatBool(i.(bool)), total})
					if err != nil {
						fmt.Println(err)
					}
				}
			}
			continue
		}

		switch value[0].(type) {
		case string:
			err := csvWriter.Write([]string{tag, value[0].(string), "1"})
			if err != nil {
				fmt.Println(err)
			}
		case int:
			err := csvWriter.Write([]string{tag, strconv.Itoa(value[0].(int)), "1"})
			if err != nil {
				fmt.Println(err)
			}
		case float64:
			err := csvWriter.Write([]string{tag, strconv.FormatFloat(value[0].(float64), 'g', -1, 64), "1"})
			if err != nil {
				fmt.Println(err)
			}
		case float32:
			err := csvWriter.Write([]string{tag, strconv.FormatFloat(float64(value[0].(float32)), 'g', -1, 32), "1"})
			if err != nil {
				fmt.Println(err)
			}
		case bool:
			err := csvWriter.Write([]string{tag, strconv.FormatBool(value[0].(bool)), "1"})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
