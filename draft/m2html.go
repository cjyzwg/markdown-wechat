package draft

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/fastwego/offiaccount"
	"github.com/russross/blackfriday"
	"github.com/tdewolff/minify/v2"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/vanng822/go-premailer/premailer"
)

func WriteCodeCss(theme *chroma.Style) string {
	// write css
	hlbuf := bytes.Buffer{}
	hlw := bufio.NewWriter(&hlbuf)
	formatter := html.New(html.WithClasses(true))
	if err := formatter.WriteCSS(hlw, theme); err != nil {
		panic(err)
	}
	hlw.Flush()
	return hlbuf.String()
}
func ReplaceCodeParts(doc *goquery.Document) (string, error) {
	// find code-parts via selector and replace them with highlighted versions
	var hlErr error
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		if hlErr != nil {
			return
		}
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")
		oldCode := s.Text()
		lexer := lexers.Get(lang)
		formatter := html.New(html.WithClasses(true))
		iterator, err := lexer.Tokenise(nil, string(oldCode))
		if err != nil {
			hlErr = err
			return
		}
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		if err := formatter.Format(buf, styles.GitHub, iterator); err != nil {
			hlErr = err
			return
		}
		if err := buf.Flush(); err != nil {
			hlErr = err
			return
		}
		s.SetHtml(b.String())
	})
	if hlErr != nil {
		return "", hlErr
	}
	new, err := doc.Html()
	if err != nil {
		return "", err
	}
	// replace unnecessarily added html tags
	return new, nil
}

func MarkdownParse(path string) string {
	mdFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	htmlSrc := blackfriday.MarkdownCommon(mdFile)
	return string(htmlSrc)
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
func ChangeLine(content string) string {
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	dom.Find("pre>code").Each(func(i int, selection *goquery.Selection) {
		textstring := selection.Text()
		textstring = strings.Replace(textstring, "\n", "g^g+;", -1)
		selection.SetText(textstring)
	})
	str, _ := dom.Html()
	//移除所有换行
	str = strings.Replace(str, "\n", "", -1)
	return str
}
func RemoveLine(dom *goquery.Document) *goquery.Document {
	dom.Find("pre>code").Each(func(i int, selection *goquery.Selection) {
		textstring := selection.Text()
		textstring = strings.Replace(textstring, "g^g+;", "\n", -1)
		selection.SetText(textstring)
	})
	return dom
}
func MarkdownRun(md_file string, css_file string, App *offiaccount.OffiAccount) string {

	md := MarkdownParse(md_file)

	content := AddHtmlTag(md)
	//将\n换行,更换成
	// content = ChangeLine(content)

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		panic(err)
	}
	// dom = RemoveLine(dom)
	//将代码块部分修改
	ReplaceCodeParts(dom)

	//代码高亮css
	code_css := WriteCodeCss(styles.MonokaiLight)

	//正常css
	css := ReadFile(css_file)
	dom.Find("title").Each(func(i int, selection *goquery.Selection) {
		selection.AfterHtml("<style>" + code_css + css + "</style>")
	})

	// dom.Find("p:contains(a)").Each(func(i int, selection *goquery.Selection) {
	// 	selection.SetText(strings.Replace(selection.Text(), " ", "@s-;", -1))

	// })

	dom.Find("img").Each(func(i int, selection *goquery.Selection) {
		img_url, existed := selection.Attr("src")
		if existed {
			fmt.Println("当前图文信息中的url地址为：" + img_url)
			resp, err := MediaUploadImgUrl(App, img_url)
			if err != nil {
				panic(err)
			}
			var img_return ThumbReturn
			if err = json.Unmarshal([]byte(resp), &img_return); err == nil {
				selection.SetAttr("src", img_return.Url)
			} else {
				fmt.Println(err)
			}
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
		text := selection.Text()
		ishttp := IsHttp(text)
		if !ishttp {
			//非http链接才可以添加
			refindex += 1
			link, exists := selection.Attr("href")
			if exists {
				ref_indexstr := strconv.Itoa(refindex)
				new := "<span class=\"footnote-word\" style=\"color: #1e6bb8; font-weight: bold;\"></span><sup class=\"footnote-ref\" style=\"line-height: 0; color: #1e6bb8; font-weight: bold;\">[" + ref_indexstr + "]</sup>"
				selection.AppendHtml(new)
				ref += "<span id=\"fn" + ref_indexstr + "\" class=\"footnote-item\" style=\"display: flex;\"><span class=\"footnote-num\" style=\"display: inline; background: none; font-size: 80%; opacity: 0.6; line-height: 26px; font-family: ptima-Regular, Optima, PingFangSC-light, PingFangTC-light, PingFang SC, Cambria, Cochin, Georgia, Times, Times New Roman, serif;\">[" + ref_indexstr + "]</span><p style=\"padding-top: 8px; padding-bottom: 8px; display: inline; font-size: 14px; width: 90%; padding: 0px; margin: 0; line-height: 26px; color: black; word-break: break-all; width: calc(100%-50);\"> " + text + "&nbsp: <em style=\"font-style: italic; color: black;\">" + link + "</em></p></span>\n"
			}
		}

	})
	if refindex > 0 {
		ref_part := refheader + ref + "</section>"
		parsedom.Find("body").AppendHtml(ref_part)
	}

	str, _ := parsedom.Html()
	return str

}

// 压缩html文件
func HtmlMinifyFile(filename string) (string, error) {

	m := minify.New()
	m.Add("text/html", &mhtml.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
	})

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}
	mb, err := m.String("text/html", string(b))
	if err != nil {
		return "", err
	}

	return mb, err

}
func HtmlMinify(htmlstring string) (string, error) {

	m := minify.New()
	m.Add("text/html", &mhtml.Minifier{
		KeepDefaultAttrVals: true,
		KeepDocumentTags:    true,
		KeepEndTags:         true,
	})

	mb, err := m.String("text/html", htmlstring)
	if err != nil {
		return "", err
	}

	return mb, err

}
func IsHttp(text string) bool {
	myRegex, _ := regexp.Compile("^(http|https)://")
	found := myRegex.FindStringIndex(text)
	if found == nil {
		return false
	} else {
		return true
	}

}
