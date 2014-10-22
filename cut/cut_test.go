package main

import "testing"
import "os"
import "bytes"
import "fmt"
import "reflect"
import "strings"

func equal(t *testing.T, expected interface{}, actual interface{}) {

	if !reflect.DeepEqual(expected, actual) {
		t.Error(
			"Expected", fmt.Sprintf("[%v]", expected),
			"\n",
			"Actual", fmt.Sprintf("[%v]", actual))
	}
}

func assert(t *testing.T, assertion interface{}) {
	equal(t, true, assertion)
}

func TestArgumentParsingFailures(t *testing.T) {
	_, msg := parseArguments([]string{"-z"})
	equal(t, "Invalid argument -z", msg)

	_, msg = parseArguments([]string{"idontexist"})
	equal(t, "open idontexist: no such file or directory", msg)
}

func TestArgumentParsingByteMode(t *testing.T) {
	variations := [][]string{
		[]string{"-b1-2"},
		[]string{"-b", "1-2"},
		[]string{"--bytes", "1-2"},
		[]string{"--bytes=1-2"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		assert(t, parameters.mode == byteMode)
		equal(t, []Range{Range{start: 1, end: 2}}, parameters.ranges)
	}
}

func TestArgumentParsingCharacterMode(t *testing.T) {
	variations := [][]string{
		[]string{"-c1-2"},
		[]string{"-c", "1-2"},
		[]string{"--characters", "1-2"},
		[]string{"--characters=1-2"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		assert(t, parameters.mode == characterMode)
		equal(t, []Range{Range{start: 1, end: 2}}, parameters.ranges)
	}
}

func TestArgumentParsingDelimiter(t *testing.T) {
	variations := [][]string{
		[]string{"-d;"},
		[]string{"-d", ";"},
		[]string{"--delimiter", ";"},
		[]string{"--delimiter=;"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		equal(t, ";", parameters.inputDelimiter)
	}
}

func TestArgumentParsingCSVMode(t *testing.T) {
	variations := [][]string{
		[]string{"-e1-2"},
		[]string{"-e", "1-2"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		assert(t, parameters.mode == csvMode)
		equal(t, []Range{Range{start: 1, end: 2}}, parameters.ranges)
	}
}

func TestArgumentParsingFieldMode(t *testing.T) {
	variations := [][]string{
		[]string{"-f1-2"},
		[]string{"-f", "1-2"},
		[]string{"--fields", "1-2"},
		[]string{"--fields=1-2"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		assert(t, parameters.mode == fieldMode)
		equal(t, []Range{Range{start: 1, end: 2}}, parameters.ranges)
	}
}

func TestArgumentParsingIgnored(t *testing.T) {
	_, messages := parseArguments([]string{"-n"})
	equal(t, "", messages)
}

func TestArgumentParsingComplement(t *testing.T) {
	parameters, messages := parseArguments([]string{"--complement"})
	equal(t, "", messages)
	assert(t, parameters.complement == true)
}

func TestArgumentParsingOnlyDelimited(t *testing.T) {
	variations := [][]string{
		[]string{"-s"},
		[]string{"--only-delimited"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		assert(t, parameters.delimitedOnly)
	}
}

func TestArgumentParsingOutputDelimiter(t *testing.T) {
	variations := [][]string{
		[]string{"--output-delimiter=|"},
		[]string{"--output-delimiter", "|"},
	}

	for _, variation := range variations {
		parameters, messages := parseArguments(variation)
		equal(t, "", messages)
		equal(t, "|", parameters.outputDelimiter)
	}
}

func TestArgumentParsingHelp(t *testing.T) {
	parameters, messages := parseArguments([]string{"--help"})
	equal(t, "", messages)
	assert(t, parameters.printUsage)
}

func TestArgumentParsingVersion(t *testing.T) {
	parameters, messages := parseArguments([]string{"--version"})
	equal(t, "", messages)
	assert(t, parameters.printVersion)
}

func TestFieldsArgumentParsing(t *testing.T) {
	arguments, _ := parseArguments([]string{"-f1,3,5"})
	equal(t, []Range{NewRange(1, 1), NewRange(3, 3), NewRange(5, 5)}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "1,3,5"})
	equal(t, []Range{NewRange(1, 1), NewRange(3, 3), NewRange(5, 5)}, arguments.ranges)

	arguments, _ = parseArguments([]string{})
	equal(t, []Range{}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "1-3"})
	equal(t, []Range{NewRange(1, 3)}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "1-"})
	equal(t, []Range{NewRange(1, 0)}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "1"})
	equal(t, []Range{NewRange(1, 1)}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "-23"})
	equal(t, []Range{NewRange(0, 23)}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "1-3,5"})
	equal(t, []Range{NewRange(1, 3), NewRange(5, 5)}, arguments.ranges)

	arguments, _ = parseArguments([]string{"-f", "1-3,-5,23,42-"})
	equal(t, []Range{NewRange(1, 3), NewRange(0, 5), NewRange(23, 23), NewRange(42, 0)}, arguments.ranges)
}

func TestDelimiterArgumentParsing(t *testing.T) {
	arguments, _ := parseArguments([]string{"-d", ","})
	equal(t, ",", arguments.inputDelimiter)

	arguments, _ = parseArguments([]string{"-d,"})
	equal(t, ",", arguments.inputDelimiter)

	arguments, _ = parseArguments([]string{})
	equal(t, ",", arguments.inputDelimiter)
}

func TestFileNameArgumentParsing(t *testing.T) {
	arguments, _ := parseArguments([]string{"sample.csv"})
	equal(t, "sample.csv", arguments.input[0].Name())

	arguments, _ = parseArguments([]string{})
	equal(t, []*os.File{os.Stdin}, arguments.input)

	arguments, _ = parseArguments([]string{"-"})
	equal(t, []*os.File{os.Stdin}, arguments.input)
}

var fullFile = `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch`

var cutTests = []struct {
	parameters []string
	input      string
	expected   string
}{
	{ // full file when no delimiter XXXX
		parameters: []string{"-dx", "-f1"},
		input:      fullFile,
		expected: `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch
`,
	},
	{ // cutting first column
		parameters: []string{"-d,", "-f1"},
		input:      fullFile,
		expected: `first name
hans
peter
`,
	},
	{ // cutting second column
		parameters: []string{"-d,", "-f2"},
		input:      fullFile,
		expected: `last name
hansen
petersen
`,
	},
	{ // cutting third column
		parameters: []string{"-d,", "-f3"},
		input:      fullFile,
		expected: `favorite pet
moose
monarch
`,
	},
	{ // inversing range
		parameters: []string{"-d,", "-f-2", "--complement"},
		input:      fullFile,
		expected: `favorite pet
moose
monarch
`,
	},
	{ // cutting first and third column
		parameters: []string{"-d,", "-f1,3"},
		input:      fullFile,
		expected: `first name,favorite pet
hans,moose
peter,monarch
`,
	},
	{ // cutting first and second column via range
		parameters: []string{"-d,", "-f1-2"},
		input:      fullFile,
		expected: `first name,last name
hans,hansen
peter,petersen
`,
	},
	{ // cutting all via a range
		parameters: []string{"-d,", "-f1-"},
		input:      fullFile,
		expected: `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch
`,
	},
	{ // cutting all via a range
		parameters: []string{"-d,", "-f-3"},
		input:      fullFile,
		expected: `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch
`,
	},
	{ // cutting all via a range
		parameters: []string{"-d,", "-f1-3,3"},
		input:      fullFile,
		expected: `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch
`,
	},
	{ // cutting fields with multi-byte delimiter
		parameters: []string{"-d€", "-f2"},
		input: "first name€last name€favorite pet\x0a" +
			"hans€hansen€moose\x0a" +
			"peter€petersen€monarch\x0a",
		expected: `last name
hansen
petersen
`,
	},
	{ // cutting fields separater by spaces
		parameters: []string{"-d ", "-f2"},
		input: "first second third\x0a" +
			"a b c\x0a" +
			"d e f\x0a",
		expected: "second\x0ab\x0ae\x0a",
	},
	{ // cutting fields separater by quotes
		parameters: []string{"-d'", "-f2"},
		input: "first'second'third\x0a" +
			"a'b'c\x0a" +
			"d'e'f\x0a",
		expected: "second\x0ab\x0ae\x0a",
	},
	{ // cutting fields separater by double quotes
		parameters: []string{"-d\"", "-f2"},
		input: "first\"second\"third\x0a" +
			"a\"b\"c\x0a" +
			"d\"e\"f\x0a",
		expected: "second\x0ab\x0ae\x0a",
	},
	{ // cutting csv values with LF rather than CRLF line ending
		parameters: []string{"-e2-", "--line-end=LF"},
		input: "first a,last b,favorite pet\x0a" +
			"hans,hansen,moose\x0a" +
			"peter,petersen,monarch\x0a",
		expected: "last b,favorite pet\x0a" +
			"hansen,moose\x0a" +
			"petersen,monarch\x0a",
	},
	{ // cutting csv values with CRLF explicitly
		parameters: []string{"-e2-", "--line-end=CRLF"},
		input: "first a,last b,favorite pet\x0d\x0a" +
			"hans,hansen,moose\x0d\x0a" +
			"peter,petersen,monarch\x0d\x0a",
		expected: "last b,favorite pet\x0d\x0a" +
			"hansen,moose\x0d\x0a" +
			"petersen,monarch\x0d\x0a",
	},
	{ // cutting csv values
		parameters: []string{"-e2-"},
		input: "first a,last a,favorite pet\x0d\x0a" +
			"hans,hansen,moose\x0d\x0a" +
			"peter,petersen,monarch\x0d\x0a",
		expected: "last a,favorite pet\x0d\x0a" +
			"hansen,moose\x0d\x0a" +
			"petersen,monarch\x0d\x0a",
	},
	{ // cutting csv values with custom input delimiters
		parameters: []string{"-e2-", "-d;"},
		input: "first a;last a;favorite pet\x0d\x0a" +
			"hans;hansen;moose\x0d\x0a" +
			"peter;petersen;monarch\x0d\x0a",
		expected: "last a;favorite pet\x0d\x0a" +
			"hansen;moose\x0d\x0a" +
			"petersen;monarch\x0d\x0a",
	},
	{ // cutting csv values with custom multi-byte input delimiters
		parameters: []string{"-e2-", "-d€", "--output-delimiter=;"},
		input: "first a€last a€favorite pet\x0d\x0a" +
			"hans€hansen€moose\x0d\x0a" +
			"peter€petersen€monarch\x0d\x0a",
		expected: "last a;favorite pet\x0d\x0a" +
			"hansen;moose\x0d\x0a" +
			"petersen;monarch\x0d\x0a",
	},
	{ // cutting csv values with custom input and output delimiters
		parameters: []string{"-e2-", "-d;", "--output-delimiter=|"},
		input: "first a;last a;favorite pet\x0d\x0a" +
			"hans;hansen;moose\x0d\x0a" +
			"peter;petersen;monarch\x0d\x0a",
		expected: "last a|favorite pet\x0d\x0a" +
			"hansen|moose\x0d\x0a" +
			"petersen|monarch\x0d\x0a",
	},
	{ // cutting csv values that are escaped
		parameters: []string{"-e2-3"},
		input: "first name,last name,\"favorite pet\"\x0d\x0a" +
			"\"hans\",hansen,\"moose,goose\"\x0d\x0a" +
			"peter,\"petersen,muellersen\",monarch\x0d\x0a",
		expected: "last name,\"favorite pet\"\x0d\x0a" +
			"hansen,\"moose,goose\"\x0d\x0a" +
			"\"petersen,muellersen\",monarch\x0d\x0a",
	},
	{ // cutting csv values that are escaped and contain new lines
		parameters: []string{"-e2-3"},
		input: "first name,last name,\"\x0d\x0afavorite pet\"\x0d\x0a" +
			"\"hans\",hansen,\"moose,goose\"\x0d\x0a" +
			"peter,\"petersen,muellersen\x0d\x0a\",monarch\x0d\x0a",
		expected: "last name,\"\x0d\x0afavorite pet\"\x0d\x0a" +
			"hansen,\"moose,goose\"\x0d\x0a" +
			"\"petersen,muellersen\x0d\x0a\",monarch\x0d\x0a",
	},
	{ // cutting csv values that are doubly escaped
		parameters: []string{"-e2-3"},
		input: "first name,last name,\"favorite\"\" pet\"\x0d\x0a" +
			"\"hans\",hansen,\"moose,goose\"\x0d\x0a" +
			"peter,\"petersen,\"\"\"\"\"\"\"\"muellersen\",monarch\x0d\x0a",
		expected: "last name,\"favorite\"\" pet\"\x0d\x0a" +
			"hansen,\"moose,goose\"\x0d\x0a" +
			"\"petersen,\"\"\"\"\"\"\"\"muellersen\",monarch\x0d\x0a",
	},
	{ // select bytes
		parameters: []string{"-b-2"},
		input:      `€foo`,
		expected:   "\xe2\x82\x0a",
	},
	{ // select characters / runes
		parameters: []string{"-c-2"},
		input:      `€foo`,
		expected: `€f
`,
	},
	{ // select characters / runes with custom separator
		parameters: []string{"-c-2", "--output-delimiter", "x"},
		input:      `€foo`,
		expected: `€xf
`,
	},
	{ // select characters / runes with custom separator with different argument formatting
		parameters: []string{"-c-2", "--output-delimiter=x"},
		input:      `€foo`,
		expected: `€xf
`,
	},
	{ // include lines that don't contain delimiter by default
		parameters: []string{"-d,", "-f2"},
		input: `first name,last name
no delimiter here
same name,and another`,
		expected: `last name
no delimiter here
and another
`,
	},
	{ // include exclude lines without delimiter
		parameters: []string{"-d,", "-f2", "-s"},
		input: `first name,last name
no delimiter here
same name,and another`,
		expected: `last name
and another
`,
	},
	{ // include exclude lines without delimiter
		parameters: []string{"-d,", "--only-delimited", "-f2"},
		input: `first name,last name
no delimiter here
same name,and another`,
		expected: `last name
and another
`,
	},
	{ // include exclude lines without delimiter
		parameters: []string{"--output-delimiter", "x", "-d,", "--only-delimited", "-f1,2"},
		input: `first name,last name
no delimiter here
same name,and another`,
		expected: `first namexlast name
same namexand another
`,
	},
	{ // ignore -n
		parameters: []string{"--output-delimiter", "x", "-n", "-d,", "--only-delimited", "-f1,2"},
		input: `first name,last name
no delimiter here
same name,and another`,
		expected: `first namexlast name
same namexand another
`,
	},
}

func TestCutFile(t *testing.T) {
	for _, data := range cutTests {
		parameters, _ := parseArguments(data.parameters)
		input := strings.NewReader(data.input)
		output := bytes.NewBuffer(nil)

		cutFile(input, output, parameters)

		equal(t, data.expected, output.String())
	}
}

func TestCut(t *testing.T) {
	fileName := "sample.csv"
	contents := `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch
`

	input, _ := os.Open(fileName)
	defer input.Close()
	output := bytes.NewBuffer(nil)

	cut([]string{fileName, "--line-end=LF"}, output)

	equal(t, string(contents), output.String())
}

func TestCuttingMultipleFiles(t *testing.T) {
	fileName := "sample.csv"
	contents := `first name,last name,favorite pet
hans,hansen,moose
peter,petersen,monarch
`

	input, _ := os.Open(fileName)
	defer input.Close()
	output := bytes.NewBuffer(nil)

	cut([]string{fileName, fileName, "--line-end=LF"}, output)

	equal(t, fmt.Sprint(string(contents), string(contents)), output.String())
}

func TestPrintingUsageInformation(t *testing.T) {
	output := bytes.NewBuffer(nil)
	cut([]string{"--help"}, output)

	equal(t, true, strings.HasPrefix(output.String(), "Usage: "))
}
