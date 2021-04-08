package html_check

import (
	"bufio"
	"io"
	"os"
	"regexp"
)

type Info struct {  // Struct of error in url
	Num_of_line int
	Num_of_column int
	Description string
}
func has_url(line string) bool{ // This function find urls in string and if urls are detected reterns true else returns false
	matched_href, err := regexp.MatchString(`href="[^"]*"`, string(line))
	if err != nil{
		panic(err)
	}
	matched_src, err := regexp.MatchString(`src="[^"]*"`, string(line))
	if err != nil{
		panic(err)
	}
	return matched_href || matched_src;

}

func is_valid_protocol(url string) bool{ // This function check protocol in url, if url contains "https" or "http"
	// returns true else returns false
	matched_http, err := regexp.MatchString(`"http://`, url)
	if err != nil{
		panic(err)
	}
	matched_https, err := regexp.MatchString(`"https://`, url)
	if err != nil{
		panic(err)
	}
	return matched_http || matched_https;
}

func too_many_double_slashes(url string) bool{ // This function check how much "//" contains the url
	double_slashes_tmp := regexp.MustCompile(`//`)
	if len(double_slashes_tmp.FindAllStringIndex(url, -1)) > 1{
		return true
	}
	return false
}

func has_invalid_symbols(url string) bool{ // This function find invalid symbols in the url
	matched_inv_symbols, err := regexp.MatchString(`[^A-Z|a-z|/|"|=|:|\.]`, url)
	if err != nil{
		panic(err)
	}
	if matched_inv_symbols{
		return true
	}
	return false
}

func parse_line(line string, idx_line int, errs *[]Info){ // Here the string is parsed
	href_url_template := regexp.MustCompile(`href="[^"]*"`)
	src_url_template := regexp.MustCompile(`src="[^"]*"`)
	idx_coloumn_href := href_url_template.FindAllStringIndex(line, -1)
	idx_coloumn_src := src_url_template.FindAllStringIndex(line, -1)
	href_urls := href_url_template.FindAllString(line, -1)
	src_urls := src_url_template.FindAllString(line, -1)
	urls := append(href_urls, src_urls...)
	idx_urls := append(idx_coloumn_href, idx_coloumn_src...)
	for idx, url := range urls{
		if !is_valid_protocol(url){

			*errs = append(*errs, Info{idx_line + 1, idx_urls[idx][0],
				"Invalid protocol!"})
		}else if too_many_double_slashes(url){

			*errs = append(*errs, Info{idx_line + 1, idx_urls[idx][0],
				"Too many double slashes!"})

		}else if has_invalid_symbols(url){
			*errs = append(*errs, Info{idx_line + 1, idx_urls[idx][0],
				"Invalid symbols!"})
		}
	}
}

func Html_check(filename string) (error, []Info){ // Main function that iterates strings in the gotten file
	file, err := os.OpenFile(filename, 'r', 0600)
	if err != nil{
		return err, []Info{}
	}
	reader := bufio.NewReader(file)
	var errs []Info
	idx := 0
	for line, err := reader.ReadString('\n'); err != io.EOF; line, err = reader.ReadString('\n'){
		if has_url(line){

			parse_line(line, idx, &errs)
		}
		idx ++
	}
	return nil, errs

}