package html_check

import (
	"bufio"
	"strings"

	//"fmt"
	"io"
	"os"
	"regexp"
)

type Info struct { // Структура, которая хранит данные о некорректных url
	NumOfLine   int
	NumOfColumn int
	Description string
}

func hasUrl(line string) (bool, error) { // Функция осуществляет проверку содержится ли в строке тег href или src
	matchedUrls, err := regexp.MatchString(`<[^>]*href=[^<]*>|<[^>]*src=[^<]*>`, line)
	return matchedUrls, err

}

func isValidProtocol(url string) (bool, error) { // Функция осуществляет проверку валидности протокола url
	matchedValidProtocol, err := regexp.MatchString(`^http:\/\/|^https:\/\/`, url)
	return matchedValidProtocol, err
}

func tooManyDoubleSlashes(url string) bool { // Функция осуществляет проверку на колличество двойных слешей
	doubleSlashesTmp := regexp.MustCompile(`//`)
	if len(doubleSlashesTmp.FindAllStringIndex(url, -1)) > 1 {
		return true
	}
	return false
}

func hasInvalidSymbols(url string) (bool, error) { // Функция проверяет url на содержание невалидных символов
	matchedInvSymbols, err := regexp.MatchString(`[^A-Z|a-z|/|:|\.]+`, url) //крайне костыльный парсинг, переделать!
	if matchedInvSymbols {
		return true, err
	}
	return false, err
}

func parseLine(line string, idxLine int, errs *[]Info) { // В этой функции происходит парсинг строки и создание ошибок
	tagTemplate := regexp.MustCompile(`<[^>]*href=[^<]*>|<[^>]*src=[^<]*>`)
	urlTemplate := regexp.MustCompile(`href=['"][^'"]*|src=['"][^'"]*`)
	idxUrls := tagTemplate.FindAllStringIndex(line, -1)
	urls := tagTemplate.FindAllString(line, -1)
	for idx, url := range urls {
		index := urlTemplate.FindStringIndex(url)
		url = urlTemplate.FindString(url)
		idxUrls[idx][0] += index[0]
		url = strings.Replace(url, "href=", "", 1)
		url = strings.Replace(url, "src=", "", 1)
		url = strings.TrimLeft(url, `'"`)
		urls[idx] = url
	}
	for idx, url := range urls {
		if ans, err := isValidProtocol(url); !ans && err == nil {
			*errs = append(*errs, Info{idxLine + 1, idxUrls[idx][0],
				"Invalid protocol!"})
		} else if tooManyDoubleSlashes(url) {

			*errs = append(*errs, Info{idxLine + 1, idxUrls[idx][0],
				"Too many double slashes!"})

		} else if ans, err := hasInvalidSymbols(url); ans && err == nil {
			*errs = append(*errs, Info{idxLine + 1, idxUrls[idx][0],
				"Invalid symbols!"})
		}

	}
}

func HtmlCheck(filename string) (error, []Info) { // Стартовая функция, открывает файл и читает из него построчно,
	// запуская остальные обработки
	file, err := os.OpenFile(filename, 'r', 0600)
	if err != nil {
		return err, nil
	}
	reader := bufio.NewReader(file)
	var errs []Info
	idx := 0
	for line, err := reader.ReadString('\n'); err != io.EOF; line, err = reader.ReadString('\n') {
		if ans, err := hasUrl(line); ans && err == nil {
			parseLine(line, idx, &errs)
		} else if err != nil {
			return err, nil
		}
		idx++
	}

	return nil, errs

}

//
//func main()  {
//	ans, err := HtmlCheck("test1.txt")
//	fmt.Print(ans, err)
//}
