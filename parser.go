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
		return 0, errors.New("No Record Structs to Map")
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
		if _, ok := recordType.Field(i).Tag.Lookup("Ignore"); ok == true {
			continue
		}
		
		propertyLength, ok := recordType.Field(i).Tag.Lookup("Length")
		if ok {
			propertyLengthInt, _ := strconv.Atoi(propertyLength)
			recordValue.Elem().Field(i).SetString(line[(pos):(pos + propertyLengthInt)])
			pos += (propertyLengthInt)
		} 
		// else {
		// 	return nil, errors.New(
		// 		fmt.Sprintf("Length tag doesn't exist on type %q:%q",
		// 			recordType.String(),
		// 			recordType.Field(i).Name))
		// }
	}
	return interface{}(recordValue.Elem().Interface()), nil
}
