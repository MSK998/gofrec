package gofrec

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type IParser interface {
	Parse() []interface{}
}

type Parser struct {
	RecordTypes   []interface{}
	Records       []interface{}
	IdentifierMap map[string]reflect.Type
	Lines         []string
	IdStart       int
	IdEnd         int
}

func (p *Parser) BytesToLines(fileContents []byte) (int, error) {
	var sb strings.Builder

	for idx, byte := range fileContents {
		_, err := fmt.Fprintf(&sb, "%c", byte)
		if err != nil {
			return idx, err
		}
	}

	p.Lines = strings.Split(sb.String(), "\n")
	return len(p.Lines), nil
}

func (p *Parser) MapIdentifiers() (int, error) {
	if len(p.RecordTypes) == 0 {
		return 0, errors.New("no record structs to map")
	}

	identifiersMap := make(map[string]reflect.Type)

	for _, t := range p.RecordTypes {
		recordType := reflect.TypeOf(t)
		for i := 0; i < recordType.NumField(); i++ {
			v, ok := recordType.Field(i).Tag.Lookup("Identifier")
			if ok {
				identifiersMap[v] = recordType
			}
		}
	}

	p.IdentifierMap = identifiersMap
	return len(identifiersMap), nil
}

func (p *Parser) MapLine(line string) (interface{}, error) {
	lineId := line[p.IdStart:p.IdEnd]
	recordType := p.IdentifierMap[lineId]
	recordValue := reflect.New(recordType)

	pos := 0
	for i := 0; i < recordType.NumField(); i++ {
		if _, ok := recordType.Field(i).Tag.Lookup("Ignore"); ok {
			continue
		}

		propertyLength, ok := recordType.Field(i).Tag.Lookup("Length")
		if ok {
			propertyLengthInt, _ := strconv.Atoi(propertyLength)
			data := line[(pos):(pos + propertyLengthInt)]
			DynamicType(recordType, i, &recordValue, data)
			//recordValue.Elem().Field(i).SetString(line[(pos):(pos + propertyLengthInt)])
			pos += (propertyLengthInt)
		}
	}
	return interface{}(recordValue.Elem().Interface()), nil
}

func (p *Parser)Parse()(int, error){
	if len(p.RecordTypes) > 0 && len(p.IdentifierMap) == 0 {
		p.MapIdentifiers()
	} else {
		return 0, errors.New("no record types to parse, please supply an array/slice of structs to parse")
	}

	totalRecords := 0
	for i, v := range p.Lines {
		rec, err := p.MapLine(v)
		if err != nil {
			return i, err
		}
		p.Records = append(p.Records, rec)
		totalRecords++
	}

	return totalRecords, nil
}
