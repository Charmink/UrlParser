package html_check

import (
	"bufio"
	//"fmt"
	"io"
	"os"
	"regexp"
)

type Info struct {  // Структура, которая хранит данные о некорректных url
	NumOfLine int
	NumOfColumn int
	Description string
}
func hasUrl(line string) (bool, error){ // Функция осуществляет проверку содержится ли в строке тег href или src
	matchedHref, err := regexp.MatchString(`href="[^"]*"`, string(line))
	matchedSrc, err := regexp.MatchString(`src="[^"]*"`, string(line))
	return matchedHref || matchedSrc, err;

}

func isValidProtocol(url string) (bool, error){ // Функция осуществляет проверку валидности протокола url
	matchedHttp, err := regexp.MatchString(`"http://`, url)
	matchedHttps, err := regexp.MatchString(`"https://`, url)
	return matchedHttp || matchedHttps, err;
}

func tooManyDoubleSlashes(url string) bool{ // Функция осуществляет проверку на колличество двойных слешей
	doubleSlashesTmp := regexp.MustCompile(`//`)
	if len(doubleSlashesTmp.FindAllStringIndex(url, -1)) > 1{
		return true
	}
	return false
}

func hasInvalidSymbols(url string) (bool, error){ // Функция проверяет url на содержание невалидных символов
	matchedInvSymbols, err := regexp.MatchString(`[^A-Z|a-z|/|"|=|:|\.]`, url)
	if matchedInvSymbols{
		return true, err
	}
	return false, err
}

func parseLine(line string, idx_line int, errs *[]Info){ // В этой функции происходит парсинг строки и создание ошибок
	hrefUrlTemplate := regexp.MustCompile(`href="[^"]*"`)
	srcUrlTemplate := regexp.MustCompile(`src="[^"]*"`)
	idxColoumnHref := hrefUrlTemplate.FindAllStringIndex(line, -1)
	idxColoumnSrc := srcUrlTemplate.FindAllStringIndex(line, -1)
	hrefUrls := hrefUrlTemplate.FindAllString(line, -1)
	srcUrls := srcUrlTemplate.FindAllString(line, -1)
	urls := append(hrefUrls, srcUrls...)
	idx_urls := append(idxColoumnHref, idxColoumnSrc...)
	for idx, url := range urls{
		if ans, err := isValidProtocol(url); !ans && err == nil{

			*errs = append(*errs, Info{idx_line + 1, idx_urls[idx][0],
				"Invalid protocol!"})
		}else if tooManyDoubleSlashes(url){

			*errs = append(*errs, Info{idx_line + 1, idx_urls[idx][0],
				"Too many double slashes!"})

		}else if ans, err := hasInvalidSymbols(url); ans && err == nil {
			*errs = append(*errs, Info{idx_line + 1, idx_urls[idx][0],
				"Invalid symbols!"})
		}

	}
}

func HtmlCheck(filename string) (error, []Info){ // Стартовая функция, открывает файл и читает из него построчно,
	// запуская остальные обработки
	file, err := os.OpenFile(filename, 'r', 0600)
	if err != nil{
		return err, nil
	}
	reader := bufio.NewReader(file)
	var errs []Info
	idx := 0
	for line, err := reader.ReadString('\n'); err != io.EOF; line, err = reader.ReadString('\n'){
		if ans, err := hasUrl(line); ans && err == nil{
			parseLine(line, idx, &errs)
		}else if err != nil{
			return err, nil
		}
		idx ++
	}
	return nil, errs

}
//
//func main()  {
//	ans, err := HtmlCheck("test1.txt")
//	fmt.Print(ans, err)
//}