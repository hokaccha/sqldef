package postgres

import (
	"os"
	"testing"

	"github.com/sqldef/sqldef/database"
	"github.com/sqldef/sqldef/parser"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParse(t *testing.T) {
	tests, err := readTests("tests.yml")
	if err != nil {
		t.Fatal(err)
	}

	genericParser := database.NewParser(parser.ParserModePostgres)
	postgresParser := NewParser()
	postgresParser.testing = true
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			psqlResult, err := postgresParser.Parse(test.SQL)
			if err != nil {
				t.Fatal(err)
			}

			if !test.CompareWithGenericParser {
				return
			}

			genericResult, err := genericParser.Parse(test.SQL)
			if err != nil {
				t.Fatal(err)
			}

			// pp.Println(genericResult)
			// pp.Println(psqlResult)
			assert.Equal(t, genericResult, psqlResult)
		})
	}
}

// func TestOfTest(t *testing.T) {
// 	tests, err := readTests("tests.yml")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	sql1 := tests["CreateViewWithCaseWhen"].SQL
// 	sql2 := tests["CreateViewWithCaseWhen2"].SQL

// 	postgresParser := NewParser()
// 	postgresParser.testing = true

// 	// genericParser := database.NewParser(parser.ParserModePostgres)
// 	// g1, err := genericParser.Parse(sql1)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// g2, err := genericParser.Parse(sql2)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// pp.Println("generic")
// 	// assert.Equal(t, g1, g2)

// 	p1, err := postgresParser.Parse(sql1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	p2, err := postgresParser.Parse(sql2)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	pp.Println("postgres")
// 	assert.Equal(t, p1[0].Statement, p2[0].Statement)
// }

type TestCase struct {
	SQL                      string
	CompareWithGenericParser bool `yaml:"compare_with_generic_parser"`
}

func readTests(file string) (map[string]TestCase, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var tests map[string]TestCase
	err = yaml.UnmarshalStrict(buf, &tests)
	if err != nil {
		return nil, err
	}

	return tests, nil
}
