package draft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fastwego/offiaccount"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/vanng822/go-premailer/premailer"
)

func MarkdownToHTML(md string) string {
	myHTMLFlags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	myExtensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_HARD_LINE_BREAK

	renderer := blackfriday.HtmlRenderer(myHTMLFlags, "", "")
	bytes := blackfriday.MarkdownOptions([]byte(md), renderer, blackfriday.Options{
		Extensions: myExtensions,
	})
	theHTML := string(bytes)
	return bluemonday.UGCPolicy().Sanitize(theHTML)
}

func MarkdownParse(path string) string {
	res := ReadFile(path)
	result := MarkdownToHTML(res)
	return result
}
func ReadFile(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fi.Close()
	fd, _ := ioutil.ReadAll(fi)
	return string(fd)
}
func AddHtmlTag(input string) string {
	//拼凑HTML页面，需要先导入Strings包
	s1 := "<!DOCTYPE html><html><head><meta charset=\"UTF-8\"><title></title></head><body>"
	s2 := "</body></html>"
	var build strings.Builder
	build.WriteString(s1)
	build.WriteString(input)
	build.WriteString(s2)
	s3 := build.String()
	return s3
}

func ParseInlineCss(content string) string {
	prem, err := premailer.NewPremailerFromString(content, premailer.NewOptions())
	if err != nil {
		panic(err)
	}

	html, err := prem.Transform()
	if err != nil {
		panic(err)
	}
	return html
}
func RepImage(htmls string) string {
	var imgRE = regexp.MustCompile(`<img[^>]+\bsrc=["']([^"']+)["']`)
	imgs := imgRE.FindAllStringSubmatch(htmls, -1)
	return imgs[0][1]
}
func MarkdownRun(md_file string, css_file string, App *offiaccount.OffiAccount) string {

	md := MarkdownParse(md_file)

	content := AddHtmlTag(md)

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	css := ReadFile(css_file)
	dom.Find("title").Each(func(i int, selection *goquery.Selection) {
		selection.AfterHtml("<style>" + css + "</style>")
	})

	dom.Find("p").Each(func(i int, selection *goquery.Selection) {
		html, _ := selection.Html()
		pos := strings.Index(html, "<img")
		if pos > -1 {
			//存在img标签，即图片
			// firstpos := strings.Index(html, content_img_path)
			// lastpos := strings.Index(html, "alt")
			// img_path := html[firstpos:(lastpos - 2)]
			// resp, err := material.MediaUploadImg(App, img_path)
			img_url := RepImage(html)
			fmt.Println("当前图文信息中的url地址为：" + img_url)
			resp, err := MediaUploadImgUrl(App, img_url)
			if err != nil {
				panic(err)
			}
			var img_return ThumbReturn
			if err = json.Unmarshal([]byte(resp), &img_return); err == nil {
				img_url := img_return.Url
				html_tag := "<img src=\"" + img_url + "\" alt=\"avatar\"/>"
				selection.SetHtml(html_tag)
			} else {
				fmt.Println(err)
			}

		} else {
			selection.SetText(strings.Replace(selection.Text(), " ", "@s-;", -1))
		}
	})

	dom.Find("br:not(code)").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})

	dom_content, _ := dom.Html()
	parse_inline_html := ParseInlineCss(dom_content)

	parsedom, err := goquery.NewDocumentFromReader(strings.NewReader(parse_inline_html))
	if err != nil {
		panic(err)
	}
	parsedom.Find("style").Each(func(i int, selection *goquery.Selection) {
		selection.Remove()
	})

	refheader := `<h3 class="footnotes-sep" style="margin-top: 30px; margin-bottom: 15px; padding: 0px; font-weight: bold; color: black; font-size: 20px;">
					<span style="display: block;">参考:</span>
					</h3>`
	ref := `<section class="footnotes">`
	refindex := 0
	parsedom.Find("a").Each(func(i int, selection *goquery.Selection) {
		refindex += 1
		text := selection.Text()
		link, exists := selection.Attr("href")
		if exists {
			ref_indexstr := strconv.Itoa(refindex)
			new := "<span class=\"footnote-word\" style=\"color: #1e6bb8; font-weight: bold;\"></span><sup class=\"footnote-ref\" style=\"line-height: 0; color: #1e6bb8; font-weight: bold;\">[" + ref_indexstr + "]</sup>"
			selection.AppendHtml(new)
			ref += "<span id=\"fn" + ref_indexstr + "\" class=\"footnote-item\" style=\"display: flex;\"><span class=\"footnote-num\" style=\"display: inline; background: none; font-size: 80%; opacity: 0.6; line-height: 26px; font-family: ptima-Regular, Optima, PingFangSC-light, PingFangTC-light, PingFang SC, Cambria, Cochin, Georgia, Times, Times New Roman, serif;\">[" + ref_indexstr + "]</span><p style=\"padding-top: 8px; padding-bottom: 8px; display: inline; font-size: 14px; width: 90%; padding: 0px; margin: 0; line-height: 26px; color: black; word-break: break-all; width: calc(100%-50);\"> " + text + "&nbsp: <em style=\"font-style: italic; color: black;\">" + link + "</em></p></span>\n"
		}

	})
	if refindex > 0 {
		ref_part := refheader + ref + "</section>"
		parsedom.Find("body").AppendHtml(ref_part)
	}

	str, _ := parsedom.Html()
	return str

}
