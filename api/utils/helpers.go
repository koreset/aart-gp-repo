package utils

import (
	"api/models"
	"bytes"
	"compress/gzip"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/gobuffalo/envy"
	ztable "github.com/gregscott94/z-table-golang"
	"github.com/rs/zerolog/log"
	"gonum.org/v1/gonum/mathext"
)

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i

}

func StringToFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func ParseDateString(dateString string) (time.Time, error) {
	layouts := []string{"2006/01/02",
		"2006/01/2",
		"2006/1/2",
		"2006/1/02",
		"2006-01-02",
		"2006-1-2",
		"2006-01-2",
		"2006-1-02"}

	var parsedTime time.Time
	var err error
	for _, layout := range layouts {
		parsedTime, err = time.Parse(layout, dateString)
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{}, errors.New("failed to parse date string")
}

func PrintProjection(projection models.Projection) {
	s := reflect.ValueOf(projection)
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("%s %s = %v ; \n", typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
	fmt.Println()
}

func FloatPrecision(x float64, precision float64) float64 {
	return math.Round(x*math.Pow(10, precision)) / math.Pow(10, precision)
}

func Max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// ReadCSV to read the content of CSV File
func ReadCSV(file *multipart.FileHeader) ([]byte, string) {
	mps, err := file.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("The file is not found || wrong root")
	}

	defer mps.Close()

	copyMps, err := file.Open()
	if err != nil {
		log.Fatal().Err(err).Msg("The file is not found || wrong root")
	}

	defer copyMps.Close()

	return ConvertToStruct(mps, copyMps)
}

func ConvertToStruct(mps multipart.File, copyMps multipart.File) ([]byte, string) {

	delimiter, err := GetDelimiter(copyMps)

	reader := csv.NewReader(mps)
	reader.TrimLeadingSpace = true
	reader.Comma = delimiter
	content, _ := reader.ReadAll()

	if len(content) < 1 {
		log.Fatal().Msg("Something wrong, the file maybe empty or length of the lines are not the same")
	}

	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		headE = strings.TrimSpace(headE)
		headE = strings.Replace(headE, " ", "_", -1)
		headE = strings.Replace(headE, "@", "_at_", -1)
		headE = strings.Replace(headE, "/", "_", -1)
		headE = strings.ToLower(headE)
		headE = strings.TrimSpace(headE)

		headersArr = append(headersArr, headE)
	}

	//Remove the header row
	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, d := range content {
		buffer.WriteString("{")
		for j, y := range d {
			buffer.WriteString(`"` + headersArr[j] + `":`)
			_, iErr := strconv.ParseInt(y, 0, 64)
			_, fErr := strconv.ParseFloat(y, 64)
			_, bErr := strconv.ParseBool(y)

			if iErr == nil {
				buffer.WriteString(y)
			} else if bErr == nil {
				//Special case for F = Female
				if y == "F" {
					buffer.WriteString((`"` + y + `"`))
				} else {
					buffer.WriteString(strings.ToLower(y))
				}
			} else if fErr == nil {
				buffer.WriteString(strings.ToLower(y))
			} else {
				buffer.WriteString((`"` + y + `"`))
			}
			//end of property
			if j < len(d)-1 {
				buffer.WriteString(",")
			}

		}
		//end of object of the array
		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(`]`)
	rawMessage := json.RawMessage(buffer.String())
	x, err := json.MarshalIndent(rawMessage, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	cwdPath, _ := os.Getwd()
	path := cwdPath + "/exported"
	newFileName := filepath.Base(path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(path)
	return x, filepath.Join(r, newFileName)
}

func GetDelimiter(copyMps multipart.File) (rune, error) {
	// Reset the file pointer to the beginning
	if _, err := copyMps.Seek(0, io.SeekStart); err != nil {
		return 0, fmt.Errorf("failed to seek to beginning of file: %v", err)
	}

	var reader io.Reader = copyMps

	// Check if the file is gzipped
	magic := make([]byte, 2)
	if n, err := copyMps.Read(magic); err == nil && n == 2 && magic[0] == 0x1f && magic[1] == 0x8b {
		if _, err := copyMps.Seek(0, io.SeekStart); err != nil {
			return 0, fmt.Errorf("failed to seek to beginning of file after magic check: %v", err)
		}
		gzReader, err := gzip.NewReader(copyMps)
		if err != nil {
			return 0, fmt.Errorf("failed to create gzip reader: %v", err)
		}
		defer gzReader.Close()
		reader = gzReader
	} else {
		// Reset if it wasn't gzipped or there was an error reading magic bytes
		if _, err := copyMps.Seek(0, io.SeekStart); err != nil {
			return 0, fmt.Errorf("failed to seek to beginning of file: %v", err)
		}
	}

	someBytes := make([]byte, 200)

	_, err := reader.Read(someBytes)
	if err != nil && err != io.EOF {
		return 0, fmt.Errorf("failed to read from file: %v", err)
	}
	var delimiter rune
	if strings.Contains(string(someBytes), ";") {
		delimiter = ';'
	}

	if strings.Contains(string(someBytes), ",") {
		delimiter = ','
	}

	if delimiter == 0 {
		// File does not contain the required delimeters... return error
		// Only commas and semicolons required.
		delimiter = ','
	}

	// Reset the file pointer to the beginning for subsequent reads
	if _, err := copyMps.Seek(0, io.SeekStart); err != nil {
		return 0, fmt.Errorf("failed to seek to beginning of file at end: %v", err)
	}

	return delimiter, nil
}

func ValidateCSVHeaders(headers []string, target interface{}) error {
	v := reflect.TypeOf(target)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("target must be a struct or a pointer to a struct")
	}

	// Create a map for quick lookup of headers
	headerMap := make(map[string]bool)
	for _, h := range headers {
		headerMap[strings.TrimSpace(h)] = true
	}

	var missingHeaders []string

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		csvTag := field.Tag.Get("csv")

		// If the field has a csv tag and it's not "-", check if it exists in the headers
		if csvTag != "" && csvTag != "-" {
			// Some tags might have options like "fieldname,omitempty", we only need the first part
			expectedHeader := strings.Split(csvTag, ",")[0]

			// Fields like created_by, creation_date or fields that have to do with time should be excluded
			// We skip these as they are usually system-populated and not required in CSV uploads
			lowerHeader := strings.ToLower(expectedHeader)
			if lowerHeader == "created_by" || lowerHeader == "creation_date" ||
				strings.Contains(lowerHeader, "time") || strings.Contains(lowerHeader, "date") {
				continue
			}

			if !headerMap[expectedHeader] {
				missingHeaders = append(missingHeaders, expectedHeader)
			}
		}
	}

	if len(missingHeaders) > 0 {
		return fmt.Errorf("missing required columns: %s", strings.Join(missingHeaders, ", "))
	}

	return nil
}

// SaveFile Will Save the file, magic right?
func SaveFile(myFile []byte, path string) {
	if err := os.WriteFile(path, myFile, os.FileMode(0644)); err != nil {
		fmt.Println(err)
	}
}

//func SaveAsCsv()

func ReadCsv(path string) ([]byte, string) {
	csvFile, err := os.Open(path)
	if err != nil {
		log.Fatal().Msg("The file is not found || wrong root")
	}
	defer csvFile.Close()

	return ConvertToStruct(csvFile, csvFile)

	////var mps []models.ModelPoint
	//modelpoints, err := csv.NewReader(csvFile).ReadAll()
	//for i, modelpointRow := range modelpoints {
	//	//var mp models.ModelPoint
	//	if i == 0 {
	//		fmt.Println("Header: Will not print", modelpointRow)
	//		continue
	//	}
	//	fmt.Println(modelpointRow)
	//}
}

func SendEmailNotification(product models.Product, job *models.ProjectionJob, email string) {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, "https://api.postmarkapp.com/email", bytes.NewBuffer([]byte(fmt.Sprintf("{From: 'aart@adsolutions.co.za', To: '%s', Subject: 'Result Runs from ADS', HtmlBody: '<strong>Hello</strong> The processing Job run number <b>%d</b> has completed in %d minutes. You can download the results using this link, %s/projections/job/%d'}", email, job.ID, int(job.RunTime), envy.Get("CLIENT_HOST", "http://localhost:8080"), job.ID))))
	req.Header.Add("X-Postmark-Server-Token", "b09eb5ee-fe99-42e2-ab79-c3a95ed55d4d")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error().Err(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error().Err(err).Msg("PostMark Sending error")
	}
	fmt.Println(string(body))
}

// Split splits the camelcase word and returns a list of words. It also
// supports digits. Both lower camel case and upper camel case are supported.
// For more info please check: http://en.wikipedia.org/wiki/CamelCase
//
// Examples
//
//	"" =>                     [""]
//	"lowercase" =>            ["lowercase"]
//	"Class" =>                ["Class"]
//	"MyClass" =>              ["My", "Class"]
//	"MyC" =>                  ["My", "C"]
//	"HTML" =>                 ["HTML"]
//	"PDFLoader" =>            ["PDF", "Loader"]
//	"AString" =>              ["A", "String"]
//	"SimpleXMLParser" =>      ["Simple", "XML", "Parser"]
//	"vimRPCPlugin" =>         ["vim", "RPC", "Plugin"]
//	"GL11Version" =>          ["GL", "11", "RiskRateCode"]
//	"99Bottles" =>            ["99", "Bottles"]
//	"May5" =>                 ["May", "5"]
//	"BFG9000" =>              ["BFG", "9000"]
//	"BöseÜberraschung" =>     ["Böse", "Überraschung"]
//	"Two  spaces" =>          ["Two", "  ", "spaces"]
//	"BadUTF8\xe2\xe2\xa1" =>  ["BadUTF8\xe2\xe2\xa1"]
//
// Splitting rules
//
//  1. If string is not valid UTF-8, return it without splitting as
//     single item array.
//  2. Assign all unicode characters into one of 4 sets: lower case
//     letters, upper case letters, numbers, and all other characters.
//  3. Iterate through characters of string, introducing splits
//     between adjacent characters that belong to different sets.
//  4. Iterate through array of split strings, and if a given string
//     is upper case:
//     if subsequent string is lower case:
//     move last character of upper case string to beginning of
//     lower case string
func Split(src string) (output string) {
	// don't split invalid utf8
	if !utf8.ValidString(src) {
		//return []string{src}, ""
		return src
	}
	var entries []string
	var runes [][]rune
	lastClass := 0
	class := 0
	// split into fields based on class of unicode character
	for _, r := range src {
		switch true {
		case unicode.IsLower(r):
			class = 1
		case unicode.IsUpper(r):
			class = 2
		case unicode.IsDigit(r):
			class = 3
		default:
			class = 4
		}
		if class == lastClass {
			runes[len(runes)-1] = append(runes[len(runes)-1], r)
		} else {
			runes = append(runes, []rune{r})
		}
		lastClass = class
	}
	// handle upper case -> lower case sequences, e.g.
	// "PDFL", "oader" -> "PDF", "Loader"
	for i := 0; i < len(runes)-1; i++ {
		if unicode.IsUpper(runes[i][0]) && unicode.IsLower(runes[i+1][0]) {
			runes[i+1] = append([]rune{runes[i][len(runes[i])-1]}, runes[i+1]...)
			runes[i] = runes[i][:len(runes[i])-1]
		}
	}
	// construct []string from results
	for _, s := range runes {
		if len(s) > 0 {
			entries = append(entries, string(s))
		}
	}
	output = join(entries...)
	return
}

func join(strs ...string) string {
	if len(strs) == 1 {
		return strs[0]
	}
	var sb strings.Builder
	for i, str := range strs {
		sb.WriteString(str)
		if i < len(strs)-1 {
			sb.WriteString(" ")
		}
	}
	return sb.String()
}

func Snakify(value string) string {
	value = strings.TrimSpace(value)
	value = strings.Replace(value, " ", "_", -1)
	value = strings.ToLower(value)
	return value
}

func GenerateColumns() []string {
	var alphabets = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	var template []string

	for _, alpha := range alphabets {
		template = append(template, alpha)
	}

	for i, _ := range alphabets {
		for j, _ := range alphabets {
			fmt.Println(i + j)
			template = append(template, alphabets[i]+alphabets[j])
		}
	}
	return template
}

func RoundUp(x float64) float64 {
	if math.IsNaN(x) {
		return x
	}
	if x == 0.0 {
		return x
	}
	roundFn := math.Ceil
	if math.Signbit(x) {
		roundFn = math.Floor
	}
	xOrig := x
	x -= math.Copysign(0.5, x)
	if x == 0 || math.Signbit(x) != math.Signbit(xOrig) {
		return math.Copysign(0.0, xOrig)
	}
	if x == xOrig-math.Copysign(1.0, x) {
		return xOrig
	}
	r := roundFn(x)
	if r != x {
		return r
	}
	return roundFn(x*0.5) * 2.0
}

func ConvertNumberWordToString(number string) string {
	number = strings.TrimSpace(number)
	number = strings.Replace(number, " ", "", -1)
	number = strings.ToLower(number)
	number = strings.Replace(number, "rd0", "0", -1)
	number = strings.Replace(number, "rd1", "1", -1)
	number = strings.Replace(number, "rd2", "2", -1)
	number = strings.Replace(number, "rd3", "3", -1)
	number = strings.Replace(number, "rd4", "4", -1)
	number = strings.Replace(number, "rd5", "5", -1)
	number = strings.Replace(number, "rd6", "6", -1)
	number = strings.Replace(number, "rd7", "7", -1)
	number = strings.Replace(number, "rd8", "8", -1)
	number = strings.Replace(number, "rd9", "9", -1)
	number = strings.Replace(number, "rd10", "10", -1)
	number = strings.Replace(number, "rd11", "11", -1)
	number = strings.Replace(number, "rd12", "12", -1)
	number = strings.Replace(number, "rd13", "13", -1)
	number = strings.Replace(number, "rd14", "14", -1)
	number = strings.Replace(number, "rd15", "15", -1)
	number = strings.Replace(number, "rd16", "16", -1)
	number = strings.Replace(number, "rd17", "17", -1)
	number = strings.Replace(number, "rd18", "18", -1)
	number = strings.Replace(number, "rd19", "19", -1)
	number = strings.Replace(number, "rd20", "20", -1)
	number = strings.Replace(number, "rd21", "21", -1)
	number = strings.Replace(number, "rd22", "22", -1)
	number = strings.Replace(number, "rd23", "23", -1)
	number = strings.Replace(number, "rd24", "24", -1)
	number = strings.Replace(number, "rd25", "25", -1)
	number = strings.Replace(number, "rd26", "26", -1)
	number = strings.Replace(number, "rd27", "27", -1)
	number = strings.Replace(number, "rd28", "28", -1)
	number = strings.Replace(number, "rd29", "29", -1)
	number = strings.Replace(number, "rd30", "30", -1)
	number = strings.Replace(number, "rd31", "31", -1)
	number = strings.Replace(number, "rd32", "32", -1)
	number = strings.Replace(number, "rd33", "33", -1)
	number = strings.Replace(number, "rd34", "34", -1)
	number = strings.Replace(number, "rd35", "35", -1)
	number = strings.Replace(number, "rd36", "36", -1)
	number = strings.Replace(number, "rd37", "37", -1)
	number = strings.Replace(number, "rd38", "38", -1)
	number = strings.Replace(number, "rd39", "39", -1)
	number = strings.Replace(number, "rd40", "40", -1)
	number = strings.Replace(number, "rd41", "41", -1)
	number = strings.Replace(number, "rd42", "42", -1)
	number = strings.Replace(number, "rd43", "43", -1)
	number = strings.Replace(number, "rd44", "44", -1)
	number = strings.Replace(number, "rd45", "45", -1)
	number = strings.Replace(number, "rd46", "46", -1)
	number = strings.Replace(number, "rd47", "47", -1)
	number = strings.Replace(number, "rd48", "48", -1)
	number = strings.Replace(number, "rd49", "49", -1)
	number = strings.Replace(number, "rd50", "50", -1)
	number = strings.Replace(number, "rd51", "51", -1)
	number = strings.Replace(number, "rd52", "52", -1)
	number = strings.Replace(number, "rd53", "53", -1)
	number = strings.Replace(number, "rd54", "54", -1)
	number = strings.Replace(number, "rd55", "55", -1)
	number = strings.Replace(number, "rd56", "56", -1)
	number = strings.Replace(number, "rd57", "57", -1)
	number = strings.Replace(number, "rd58", "58", -1)
	number = strings.Replace(number, "rd59", "59", -1)
	number = strings.Replace(number, "rd60", "60", -1)
	number = strings.Replace(number, "rd61", "61", -1)
	number = strings.Replace(number, "rd62", "62", -1)
	number = strings.Replace(number, "rd63", "63", -1)
	number = strings.Replace(number, "rd64", "64", -1)
	number = strings.Replace(number, "rd65", "65", -1)
	number = strings.Replace(number, "rd66", "66", -1)
	number = strings.Replace(number, "rd67", "67", -1)
	number = strings.Replace(number, "rd68", "68", -1)
	number = strings.Replace(number, "rd69", "69", -1)
	number = strings.Replace(number, "rd70", "70", -1)
	number = strings.Replace(number, "rd71", "71", -1)
	number = strings.Replace(number, "rd72", "72", -1)
	number = strings.Replace(number, "rd73", "73", -1)
	number = strings.Replace(number, "rd74", "74", -1)
	number = strings.Replace(number, "rd75", "75", -1)
	number = strings.Replace(number, "rd76", "76", -1)
	number = strings.Replace(number, "rd77", "77", -1)
	number = strings.Replace(number, "rd78", "78", -1)
	number = strings.Replace(number, "rd79", "79", -1)
	number = strings.Replace(number, "rd80", "80", -1)
	number = strings.Replace(number, "rd81", "81", -1)
	number = strings.Replace(number, "rd82", "82", -1)
	number = strings.Replace(number, "rd83", "83", -1)
	number = strings.Replace(number, "rd84", "84", -1)
	number = strings.Replace(number, "rd85", "85", -1)
	number = strings.Replace(number, "rd86", "86", -1)
	number = strings.Replace(number, "rd87", "87", -1)
	number = strings.Replace(number, "rd88", "88", -1)
	number = strings.Replace(number, "rd89", "89", -1)
	number = strings.Replace(number, "rd90", "90", -1)
	number = strings.Replace(number, "rd91", "91", -1)
	number = strings.Replace(number, "rd92", "92", -1)
	number = strings.Replace(number, "rd93", "93", -1)
	number = strings.Replace(number, "rd94", "94", -1)
	number = strings.Replace(number, "rd95", "95", -1)
	number = strings.Replace(number, "rd96", "96", -1)
	number = strings.Replace(number, "rd97", "97", -1)
	number = strings.Replace(number, "rd98", "98", -1)
	number = strings.Replace(number, "rd99", "99", -1)
	number = strings.Replace(number, "rd100", "100", -1)
	number = strings.Replace(number, "rd101", "101", -1)
	number = strings.Replace(number, "rd102", "102", -1)
	number = strings.Replace(number, "rd103", "103", -1)
	number = strings.Replace(number, "rd104", "104", -1)
	number = strings.Replace(number, "rd105", "105", -1)
	number = strings.Replace(number, "rd106", "106", -1)
	number = strings.Replace(number, "rd107", "107", -1)
	number = strings.Replace(number, "rd108", "108", -1)
	number = strings.Replace(number, "rd109", "109", -1)
	number = strings.Replace(number, "rd110", "110", -1)
	number = strings.Replace(number, "rd111", "111", -1)
	number = strings.Replace(number, "rd112", "112", -1)
	number = strings.Replace(number, "rd113", "113", -1)
	number = strings.Replace(number, "rd114", "114", -1)
	number = strings.Replace(number, "rd115", "115", -1)
	number = strings.Replace(number, "rd116", "116", -1)
	number = strings.Replace(number, "rd117", "117", -1)
	number = strings.Replace(number, "rd118", "118", -1)
	number = strings.Replace(number, "rd119", "119", -1)
	number = strings.Replace(number, "rd120", "120", -1)

	return number
}

func Shuffle(slice []float64) {
	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		fmt.Println("j:", j)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func NormalInverse(mu float64, sigma float64, zScoreTable *ztable.ZTable) float64 { // z= (x-mu)/sigma
	prob := rand.Float64()
	var usedProb, result float64
	//if prob < 0.5 {
	//	usedProb = prob
	//}else{
	//	usedProb = prob -0.5
	//}

	usedProb = prob
	zScore, _ := zScoreTable.FindZScore(usedProb)
	//if prob < 0.5 {
	//	result = -zScore * sigma + mu
	//}else{
	//	result = zScore *sigma + mu
	//}

	result = zScore*sigma + mu
	return result //rand.Float64() * sigma + mu
}

func LogNormalInverse(mu float64, std float64, zScoreTable *ztable.ZTable) float64 { //X is lognormally distributed if Y=ln(X) is normally distributed
	prob := rand.Float64()
	var usedProb, result float64
	//if prob <0.5 {
	//	usedProb = prob
	//}else{
	//	usedProb = prob -0.5
	//}
	usedProb = prob
	zScore, _ := zScoreTable.FindZScore(usedProb)
	//if prob < 0.5 {
	//	result = -zScore * std + mu
	//}else{
	//	result = zScore *std + mu
	//}
	result = zScore*std + mu
	return math.Exp(result)
}

func GammaInverse(mu float64, std float64) float64 { //X follows a Gamma distribution
	//return math.Exp(rand.Float64() * sigma + mu)
	//α=E^2[X]/Var(x),
	//β=E[X]/Var(x).
	//sigmaSq=alpha/Beta^2
	if std > 0 {
		alpha := math.Pow(mu/std, 2.0)
		beta := mu / math.Pow(std, 2) //mu/(alpha-1)
		prob := rand.Float64()
		//g:= distuv.InverseGamma{alpha,beta,prob}
		return mathext.GammaIncRegCompInv(alpha, prob) * beta
	} else {
		return 1
	}
}

func ColIndexToExcelColName(index int) string {
	result := ""
	for index >= 0 {
		remainder := index % 26
		result = string(rune('A'+remainder)) + result
		index = (index / 26) - 1
	}
	return result
}

func IsValidRSAID(id string) bool {
	id = strings.TrimSpace(id)
	if len(id) != 13 {
		return false
	}

	// Check if all characters are digits
	for _, r := range id {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	// Validate date part
	monthPart := id[2:4]
	dayPart := id[4:6]

	month, _ := strconv.Atoi(monthPart)
	day, _ := strconv.Atoi(dayPart)

	if month < 1 || month > 12 || day < 1 || day > 31 {
		return false
	}

	// Basic check for days in month
	if month == 2 {
		if day > 29 {
			return false
		}
	} else if month == 4 || month == 6 || month == 9 || month == 11 {
		if day > 30 {
			return false
		}
	}

	// Luhn Algorithm validation
	sum := 0
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(id[12-i]))
		if i%2 == 1 { // 2nd, 4th, 6th... from the right
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}

	return sum%10 == 0
}

func ValidMemberIdColValues(members []models.GPricingMemberData) error {
	existing := make(map[string]struct{})

	for _, member := range members {
		//  TODO: Implement Checks for lengths and formats of different countries
		if member.MemberIdType == "RSA_ID" || member.MemberIdType == "ID" {
			if !IsValidRSAID(member.MemberIdNumber) {
				return fmt.Errorf("ID validation failed for %s", member.MemberIdNumber)
			}
		} else if member.MemberIdType == "PASSPORT" {
			if !regexp.MustCompile(`^[A-Za-z]\d{8}$`).MatchString(member.MemberIdNumber) {
				return fmt.Errorf("Passport number format not correct for %s", member.MemberIdNumber)
			}
		}

		if _, exits := existing[member.MemberIdNumber]; exits {
			return fmt.Errorf("Duplicate ID Number %s", member.MemberIdNumber)
		}
		existing[member.MemberIdNumber] = struct{}{}
	}
	return nil
}
