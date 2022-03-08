package draft

type ThumbReturn struct {
	MediaId string        `json:"media_id"`
	Url     string        `json:"url"`
	Item    []interface{} `json:"item"`
}
type Articles map[string][]Article

type Article struct {
	Title              string `json:"title"`
	Author             string `json:"author"`
	Digest             string `json:"digest"`
	Content            string `json:"content"`
	ContentSourceUrl   string `json:"content_source_url"`
	ThumbMediaId       string `json:"thumb_media_id"`
	NeedOpenComment    int    `json:"need_open_comment"`
	OnlyFansCanComment int    `json:"only_fans_can_comment"`
}

type ConfigFile struct {
	MarkdownFilePath string `json:"markdown_file_path"`
	CssFilePath      string `json:"css_file_path"`
	AssetsPath       string `json:"assets_path"` //文中图片路径
	ImagePath        string `json:"image_path"`  //标题图片路径
	Title            string `json:"title"`
	Author           string `json:"author"`
	Digest           string `json:"digest"`
}
