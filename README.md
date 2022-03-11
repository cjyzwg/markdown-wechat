# markdown-wechat
markdown文件转成微信公众号格式，根据自己的css添加

### Demo[链接](https://github.com/cjyzwg/markdown-wechat-demo)

#### **微信公众号配置项**
- cp .env.example .env
- 将.env配置修改成微信公众号对应的值
```
APPID=xxxxxxx
SECRET=xxxxxxxxxxxxx
TOKEN=xxxxx
EncodingAESKey=xxxxxxxxxxxxxxxxxxxx

```

#### **引用步骤**

```
import (
	"github.com/cjyzwg/markdown-wechat/draft"
	"github.com/fastwego/offiaccount"
	"github.com/spf13/viper"
)
```
```
res, err := draft.DraftRun(config_file, App)
```

#### 参考链接
1.chroma代码块高亮:https://www.zupzup.org/go-markdown-syntax-highlight-chroma/  
2.chroma代码块高亮:https://github.com/zupzup/markdown-code-highlight-chroma  
3.chroma高亮的css:https://xyproto.github.io/splash/docs/index.html  
4.goquery方法:https://blog.csdn.net/skh2015java/article/details/72998418  
5.goquery详细方法:https://cloud.tencent.com/developer/article/1196783



