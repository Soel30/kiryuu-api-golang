package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

type Posts struct {
	Title   string    `json:"title"`
	Alt     string    `json:"alt"`
	Thumb   string    `json:"thumb"`
	Status  string    `json:"status"`
	Desc    string    `json:"desc"`
	Type    string    `json:"type"`
	Release string    `json:"release"`
	Authors string    `json:"authors"`
	Artist  string    `json:"artists"`
	Score   string    `json:"score"`
	Genres  []Genres  `json:"genres"`
	Chapter []Chapter `json:"chapters"`
}

type Genres struct {
	Name string `json:"name"`
}

type Chapter struct {
	Title string `json:"title"`
	Img   string `json:"img"`
}

type Link struct {
	Data []Url `json:"data"`
}

type Url struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func crawll(c *gin.Context) {
	D := colly.NewCollector()
	D.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	chpmage := D.Clone()

	D.OnHTML("#content > div > div.terebody > div.postbody.seriestu.seriestere > article", func(h *colly.HTMLElement) {
		post := Posts{}
		post.Title = strings.ReplaceAll(h.ChildText("div.seriestucon > div.seriestuheader > h1"), " Bahasa Indonesia", "")
		post.Alt = h.ChildText("div.seriestucon > div.seriestuheader > div")
		post.Thumb = h.ChildAttr("div.seriestucon > div.seriestucontent > div.seriestucontl > div.thumb > img", "data-lazy-src")
		post.Desc = h.ChildText("div.entry-content.entry-content-single > p")
		post.Status = h.ChildText("div.seriestucontentr > div.seriestucont > div > table > tbody > tr:nth-child(1) > td:nth-child(2)")
		post.Release = h.ChildText("div.seriestucontentr > div.seriestucont > div > table > tbody > tr:nth-child(3) > td:nth-child(2)")
		post.Type = h.ChildText("div.seriestucontentr > div.seriestucont > div > table > tbody > tr:nth-child(2) > td:nth-child(2)")
		post.Authors = h.ChildText("div.seriestucontentr > div.seriestucont > div > table > tbody > tr:nth-child(4) > td:nth-child(2)")
		post.Artist = h.ChildText("div.seriestucontent > div.seriestucontentr > div.seriestucont > div > table > tbody > tr:nth-child(4) > td:nth-child(2)")
		post.Score = h.ChildText("div.seriestucontent > div.seriestucontl > div.rating.bixbox > div > div.num")

		h.ForEach("div.seriestugenre a", func(_ int, x *colly.HTMLElement) {
			genre := Genres{}
			genre.Name = x.Text
			post.Genres = append(post.Genres, genre)
		})

		h.ForEach("#chapterlist > ul > li", func(_ int, s *colly.HTMLElement) {
			chapter := Chapter{}
			chapter.Title = s.ChildText("div > div.eph-num > a > span.chapternum")
			url := s.ChildAttr("div > div.eph-num > a", "href")
			url = s.Request.AbsoluteURL(url)

			chpmage.OnHTML("#readerarea", func(v *colly.HTMLElement) {
				chapter.Img = v.Text
			})

			chpmage.OnRequest(func(r *colly.Request) {
				fmt.Println("Visiting ", url)
			})

			chpmage.Visit(url)
			post.Chapter = append(post.Chapter, chapter)
		})

		c.JSON(200, post)
	})

	D.Visit("https://kiryuu.id/manga/" + c.Param("slug"))
}

func sjw(c *gin.Context) {
	D := colly.NewCollector()

	D.OnHTML("#content > div.wrapper > div.postbody > div.bixbox.seriesearch > div.mrgn ", func(h *colly.HTMLElement) {
		var link Link
		h.ForEach("div.soralist > div", func(_ int, p *colly.HTMLElement) {
			p.ForEach("ul > li", func(_ int, z *colly.HTMLElement) {
				var url Url
				url.Url = strings.ReplaceAll(z.ChildAttr("a", "href"), "https://kiryuu.id", "")
				url.Title = z.ChildText("a")
				link.Data = append(link.Data, url)
			})
		})
		c.JSON(200, link)

	})

	D.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})
	D.Visit("https://kiryuu.id/manga/list-mode")
}

func main() {
	port := os.Getenv("PORT")
	r := gin.Default()
	r.GET("/manga/:slug", crawll)
	r.GET("/sjw", sjw)
	r.Run(port)
}
