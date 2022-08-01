package parser_test

import (
	"fmt"
	"gofrec"
	"reflect"
	"strconv"
	"testing"
)

type T1 struct {
	Prop1      string `Identifier:"001" Length:"3"`
	Prop2      string `Length:"8"`
	Prop3      string `Length:"10"`
}

type T2 struct {
	Prop1 string `Identifier:"002" Length:"3"`
}

func TestBytesToLines(t *testing.T) {

	gofrec := gofrec.Parser{}
	fileContent := []byte("001BLAHBLAH1234567890\n001FOO BAR 1234567890")

	fileLines, err := gofrec.BytesToLines(fileContent)

	if err != nil {
		t.Errorf("Expected %s but got %v", "nil", err)
	}
	if fileLines != 2 {
		t.Errorf("Expected: %d\nActual: %d", 2, len(gofrec.Lines))
	}
}

func TestMapIdentifiers(t *testing.T) {
	recordTypes := []interface{}{
		T1{},
		T2{},
	}

	gofrecParser := gofrec.Parser{RecordTypes: recordTypes}

	numTypes, err := gofrecParser.MapIdentifiers()

	if err != nil {
		t.Errorf("Expected nil error but got %v", err)
	}

	if numTypes != 2 {
		t.Errorf("Expected: %d\nActual: %d", 2, numTypes)
	}

}

func TestLineToType(t *testing.T) {
	recordTypes := []interface{}{
		T1{},
		T2{},
	}

	gofrecParser := gofrec.Parser{
		RecordTypes: recordTypes,
		IdStart:     0,
		IdEnd:       3,
	}

	gofrecParser.MapIdentifiers()

	gofrecParser.BytesToLines([]byte("001BLAHBLAH1234567890\n001FOO BAR 1234567890"))

	//gofrecParser.MapLine("001BLAHBLAH1234567890")

}

func TestMapLine(t *testing.T) {
	recordTypes := []interface{}{
		T1{},
		T2{},
	}

	par := gofrec.Parser{
		RecordTypes: recordTypes,
		IdStart:     0,
		IdEnd:       3,
	}

	par.MapIdentifiers()
	par.BytesToLines([]byte("001BLAHBLAH1234567890\n001FOO BAR 1234567890"))
	rec, err := par.MapLine(par.Lines[0])
	
	if err != nil{
		t.Error(err, rec)
	}

	s := reflect.ValueOf(rec)
	for i := 0; i < s.NumField(); i++ {
		t.Log(s.Field(i).String())
	}
}

func TestMapper(t *testing.T) {
	identifiersMap := make(map[string]reflect.Type)

	recordTypes := []interface{}{
		T1{},
		T2{},
	}

	for _, t := range recordTypes {
		recordType := reflect.TypeOf(t)
		for i := 0; i < recordType.NumField(); i++ {
			v, ok := recordType.Field(i).Tag.Lookup("Identifier")
			if ok {
				identifiersMap[v] = recordType
			}
		}
	}

	line := "001BLAHBLAH1234567890"
	lineId := line[0:3]
	recordType := identifiersMap[lineId]
	typeVal := reflect.New(recordType)
	// recordValue := reflect.ValueOf(reflect.New(reflect.TypeOf(&T1{}))).Elem()
	// rT := recordValue.Type()

	pos := 0
	for i := 0; i < recordType.NumField(); i++ {

		if _, ok := recordType.Field(i).Tag.Lookup("Ignore"); ok == true {
			continue
		}

		propertyLength, ok := recordType.Field(i).Tag.Lookup("Length")
		if ok {
			propertyLengthInt, _ := strconv.Atoi(propertyLength)
			//val := reflect.ValueOf(line[(pos):(pos + propertyLengthInt)])
			typeVal.Elem().Field(i).SetString(line[(pos):(pos + propertyLengthInt)])

			//typeVal.Elem().Field(i)
			//println(val.String())
			pos += (propertyLengthInt)
		} else {
			fmt.Printf("Length tag doesn't exist on type %q:%q",
				recordType.String(),
				recordType.Field(i).Name)
		}
	}

	//finalRecord := typeVal.Elem().Interface()

}
