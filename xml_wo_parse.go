package main

import (
	"encoding/xml"
	"fmt"
	"time"
)

// 要素を格納する変数名は必ず「大文字」で始める
// スコープと関係があるのか？
type Item struct {
	Title   string    `xml:"title"`
	PubDate time.Time `xml:"pubDate"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type RSS struct {
	Channel Channel `xml:"channel"`
}

// 時刻情報は直接time.Timeに格納してくれないので、
// 専用の解析処理を書く必要がある。
func (item *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// 一時的に名前データを格納する場所。
	raw := struct {
		Title   string `xml:"title"`
		PubDate string `xml:"pubDate"`
	}{}
	// 一旦素の解析処理を呼び出して、日付は文字列で取得
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}
	// 文字列で撮ってきた時刻情報を個別にtime.Time化
	t, err := time.Parse(time.RFC1123Z, raw.PubDate)
	if err != nil {
		return err
	}
	// 解析結果を格納。
	*item = Item{
		raw.Title,
		t,
	}
	return nil
}

func (p *RSS) String() string {
	s := ``
	s += `channel > title > ` + p.Channel.Title + "\n"
	for _, v := range p.Channel.Items {
		s += `channel > items > title > ` + v.Title + "\n"
		s += `channel > items > time > ` + v.PubDate.Format(time.RFC822) + "\n"
	}
	return s
}

func main() {
	s := `
<?xml version="1.0" ?>
<rss>
  <channel>
    <title>たいとる</title>
    <item>
      <title>記事のタイトル</title>
      <pubDate>Thu, 14 Feb 2019 07:00:20 +0000</pubDate>
    </item>
  </channel>
</rss>
    `
	var rss RSS
	err := xml.Unmarshal([]byte(s), &rss)
	if err != nil {
		panic(err)
	}
	fmt.Println(rss.String())

}
