// Copyright 2016 HeadwindFly. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in the LICENSE file.

package ini

import (
	"bufio"
	"io"
	"os"
	"regexp"
	"strings"
	"strconv"
	"fmt"
)

var (
	regSection *regexp.Regexp
	regParam *regexp.Regexp
)

func init() {
	var err error
	regSection, err = regexp.Compile(`^\[([a-z0-9-]+)\]$`)
	if err != nil {
		panic(err)
	}

	pattern := `[0-9a-zA-Z-_,:;\|\./]`

	regParam, err = regexp.Compile(`^(` + pattern + `+)(\s)*=(\s)*(` + pattern + `*)$`)
	if err != nil {
		panic(err)
	}
}

// Create an ini config instance.
func NewConfig(filename string) *Config {
	config := &Config{
		Filename:filename,
		Sections:make(map[string]*Section, 0),
		currentKey:"",
	}

	section := NewSection()
	config.Sections[config.currentKey] = section

	config.readContent()
	return config
}

type Config struct {
	Filename   string
	Sections   map[string]*Section
	currentKey string
}

// Read ini config file's content.
func (this *Config) readContent() {
	f, err := os.Open(this.Filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		line, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		this.readLine(line)
	}
}

// Get section by key.
func (this *Config) GetSection(args ... string) (*Section, error) {
	var key = ""
	if (len(args) > 0) {
		key = args[0]
	}
	section, ok := this.Sections[key]
	if ok {
		return section, nil
	}
	return nil, fmt.Errorf("The section named %s does not exist.", key)
}

// Read content line by line.
func (this *Config) readLine(line []byte) {
	lineStr := string(line)
	lineStr = strings.TrimSpace(lineStr)
	if len(lineStr) == 0 {
		return
	}

	if regSection.MatchString(lineStr) {
		submatch := regSection.FindStringSubmatch(lineStr)
		sectionKey := submatch[len(submatch) - 1]
		section := NewSection()
		this.Sections[sectionKey] = section
		this.currentKey = sectionKey
	} else if regParam.MatchString(lineStr) {
		submatch := regParam.FindStringSubmatch(lineStr)
		key := submatch[1]
		value := submatch[len(submatch) - 1]
		this.Sections[this.currentKey].Params[key] = value
	}
}

// Create a Section pointer instance.
func NewSection() *Section {
	return &Section{
		Params: make(Params, 0),
	}
}

type Section struct {
	Params Params
}

type Params map[string]string

// Return string.
// The param is set, it will be returned.
// If args[0] is set, it will be returned default.
// The default returns empty string and error.
func (this *Section) GetString(key string, args ... string) (string, error) {
	value, ok := this.Params[key]
	if ok && (len(value) > 0) {
		return value, nil
	}
	if len(args) > 0 {
		return args[0], nil
	}
	return "", fmt.Errorf("The param named %s does not exist.", key)
}

// Return int.
// The param will be converted to integer type before returning.
// If args[0] is set, it will be returned default.
// The default returns zero and error.
func (this *Section) GetInt(key string, args ... int) (int, error) {
	value, ok := this.Params[key]
	if ok {
		retVal, err := strconv.Atoi(value)
		if err != nil {
			return 0, err
		}
		return retVal, nil
	}
	if len(args) > 0 {
		return args[0], nil
	}
	return 0, fmt.Errorf("The param named %s does not exist.", key)
}

// Return boolean.
// If the param is set to 'on' or 'true', true will be returned.
// If the param is set to the other value, false will be returned.
// If args[0] is set, it will be returned default.
// The default returns false and error.
func (this *Section) GetBool(key string, args ... bool) (bool, error) {
	value, ok := this.Params[key]
	if ok {
		if strings.EqualFold("on", value) || strings.EqualFold("true", value) {
			return true, nil
		}
		return false, nil
	}
	if len(args) > 0 {
		return args[0], nil
	}
	return false, fmt.Errorf("The param named %s does not exist.", key)
}

// Return slice.
// The param will be converted to slice type by method named strings.Split() before returning.
// If args[0] is set, it will be returned default.
// The default returns an empty empty slice and error.
func (this *Section) GetSlice(key, sep string, args ... []string) ([]string, error) {
	value, ok := this.Params[key]
	if ok {
		return strings.Split(value, sep), nil
	}
	if len(args) > 0 {
		return args[0], nil
	}
	return []string{}, fmt.Errorf("The param named %s does not exist.", key)
}