package main

import (
	"encoding/xml"
	"fmt"
	"time"
)

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

func (item *Item) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	raw := struct {
		Title   string `xml:"title"`
		PubDate string `xml:"pubDate"`
	}{}
	if err := d.DecodeElement(&raw, &start); err != nil {
		return err
	}
	t, err := time.Parse(time.RFC1123Z, raw.PubDate)
	if err != nil {
		return err
	}
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
