package main

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
)

var bibleToken string

const bibleVersion = `eng-KJV`
const bibleURL = `https://bibles.org/v2`

type BibleBooks struct {
	Response struct {
		Books []struct {
			VersionID   string `json:"version_id"`
			Name        string `json:"name"`
			Abbr        string `json:"abbr"`
			Ord         string `json:"ord"`
			BookGroupID string `json:"book_group_id"`
			Testament   string `json:"testament"`
			ID          string `json:"id"`
			OsisEnd     string `json:"osis_end"`
			Parent      struct {
				Version struct {
					Path string `json:"path"`
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"version"`
			} `json:"parent"`
			Next struct {
				Book struct {
					Path string `json:"path"`
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"book"`
			} `json:"next,omitempty"`
			Copyright string `json:"copyright"`
			Previous  struct {
				Book struct {
					Path string `json:"path"`
					Name string `json:"name"`
					ID   string `json:"id"`
				} `json:"book"`
			} `json:"previous,omitempty"`
		} `json:"books"`
		Meta struct {
			Fums          string `json:"fums"`
			FumsTid       string `json:"fums_tid"`
			FumsJsInclude string `json:"fums_js_include"`
			FumsJs        string `json:"fums_js"`
			FumsNoscript  string `json:"fums_noscript"`
		} `json:"meta"`
	} `json:"response"`
}

func bibleBooks(s *discordgo.Session, m *discordgo.MessageCreate) {
	u, err := url.Parse(bibleURL)
	if err != nil {
		panic(err)
	}
	u.Path += "/versions/" + bibleVersion + "/books.js"
	req := &http.Request{Method: "GET",
		URL: u,
		Header: http.Header{
			"User-Agent": {userAgent},
		},
	}
	//Use only the token
	req.SetBasicAuth(bibleToken, "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", http.StatusText(resp.StatusCode))
		s.ChannelMessageSend(m.ChannelID, "Sorry, we could not fetch the bible books at this moment.")
		return
	}
	var books BibleBooks
	if err := json.NewDecoder(resp.Body).Decode(&books); err != nil {
		//Invalid JSON
		panic(err)
	}
	if len(books.Response.Books) <= 0 {
		log.Println("We failed to get any books")
		s.ChannelMessageSend(m.ChannelID, "Sorry, we could not fetch the bible books at this moment.")
		return
	}

	var msg []string
	var bookList []string
	for _, book := range books.Response.Books {
		bookList = append(bookList, book.Name)
	}
	msg = append(msg, strings.Join(bookList, ", "))
	//Copyright message assumed to be the same for all books
	//copyright := books.Response.Books[0].Copyright
	//copyright = strings.Replace(copyright, `<p>`, "*", -1)
	//copyright = strings.Replace(copyright, `</p>`, "*", -1)
	//msg = append(msg, copyright)

	s.ChannelMessageSend(m.ChannelID, strings.Join(msg, "\n"))
}

type BibleSearch struct {
	Response struct {
		Search struct {
			Result struct {
				Type    string `json:"type"`
				Summary struct {
					Query      string   `json:"query"`
					Start      int      `json:"start"`
					Total      int      `json:"total"`
					Rpp        string   `json:"rpp"`
					Sort       string   `json:"sort"`
					Versions   []string `json:"versions"`
					Testaments []string `json:"testaments"`
				} `json:"summary"`
				Passages []struct {
					Display             string `json:"display"`
					Version             string `json:"version"`
					VersionAbbreviation string `json:"version_abbreviation"`
					Path                string `json:"path"`
					StartVerseID        string `json:"start_verse_id"`
					EndVerseID          string `json:"end_verse_id"`
					Text                string `json:"text"`
					Copyright           string `json:"copyright"`
				} `json:"passages"`
				Verses []struct {
					Auditid    string `json:"auditid"`
					Verse      string `json:"verse"`
					Lastverse  string `json:"lastverse"`
					ID         string `json:"id"`
					OsisEnd    string `json:"osis_end"`
					Label      string `json:"label"`
					Reference  string `json:"reference"`
					PrevOsisID string `json:"prev_osis_id"`
					NextOsisID string `json:"next_osis_id"`
					Text       string `json:"text"`
					Parent     struct {
						Chapter struct {
							Path string `json:"path"`
							Name string `json:"name"`
							ID   string `json:"id"`
						} `json:"chapter"`
					} `json:"parent"`
					Next struct {
						Verse struct {
							Path string `json:"path"`
							Name string `json:"name"`
							ID   string `json:"id"`
						} `json:"verse"`
					} `json:"next"`
					Previous struct {
						Verse struct {
							Path string `json:"path"`
							Name string `json:"name"`
							ID   string `json:"id"`
						} `json:"verse"`
					} `json:"previous"`
					Copyright string `json:"copyright"`
				} `json:"verses"`
			} `json:"result"`
		} `json:"search"`
		Meta struct {
			Fums          string `json:"fums"`
			FumsTid       string `json:"fums_tid"`
			FumsJsInclude string `json:"fums_js_include"`
			FumsJs        string `json:"fums_js"`
			FumsNoscript  string `json:"fums_noscript"`
		} `json:"meta"`
	} `json:"response"`
}

func bibleSearch(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	u, err := url.Parse(bibleURL)
	if err != nil {
		panic(err)
	}
	u.Path += "/search.js"
	q := u.Query()
	q.Set("query", query)
	q.Set("version", bibleVersion)
	u.RawQuery = q.Encode()

	req := &http.Request{Method: "GET",
		URL: u,
		Header: http.Header{
			"User-Agent": {userAgent},
		},
	}
	//Use only the token
	req.SetBasicAuth(bibleToken, "")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println("Error:", http.StatusText(resp.StatusCode))
		s.ChannelMessageSend(m.ChannelID, "Sorry, we could not fetch the bible books at this moment.")
		return
	}
	if resp.StatusCode == http.StatusNotFound {
		s.ChannelMessageSend(m.ChannelID, query+":"+http.StatusText(http.StatusNotFound))
		log.Println(req.URL)
		s.ChannelMessageSend(m.ChannelID, "Sorry, we could not fetch the bible books at this moment.")
		return
	}
	var result BibleSearch
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		//Invalid JSON
		panic(err)
		return
	}
	//s.ChannelMessageSend(m.ChannelID,result.Response.Search.Result.Summary.Query)
	if result.Response.Search.Result.Type == "passages" {
		log.Println(result.Response.Search.Result)
		if result.Response.Search.Result.Passages == nil || len(result.Response.Search.Result.Passages) <= 0 {
			s.ChannelMessageSend(m.ChannelID, query+":"+http.StatusText(http.StatusNotFound))
			return
		}
		verse := result.Response.Search.Result.Passages[rand.Intn(len(result.Response.Search.Result.Passages))].Text
		verse = html2markup(verse)
		//DAS RITE
		if strings.Contains(strings.ToLower(query), `lord`) {
			if strings.ToLower(query) == `lord` {
				verse = strings.Replace(verse, `LORD`, "__**LORD**__", -1)
			} else {
				verse = strings.Replace(verse, `LORD`, "* __***LORD***__ *", -1)
			}
		} else {
			verse = strings.Replace(verse, `LORD`, "__***LORD***__", -1)
		}

		lines := strings.Split(verse, "\n")
		for _, line := range lines {
			_, err = s.ChannelMessageSend(m.ChannelID, line)
			if err != nil {
				log.Println(err)
			}
		}
		log.Println(req.URL)
	}
	if result.Response.Search.Result.Type == "verses" {
		log.Println(result.Response.Search.Result)
		if result.Response.Search.Result.Verses == nil || len(result.Response.Search.Result.Verses) <= 0 {
			s.ChannelMessageSend(m.ChannelID, query+":"+http.StatusText(http.StatusNotFound))
			return
		}
		verse := result.Response.Search.Result.Verses[rand.Intn(len(result.Response.Search.Result.Verses))].Text
		verse = html2markup(verse)
		//DAS RITE
		if strings.Contains(strings.ToLower(query), `lord`) {
			if strings.ToLower(query) == `lord` {
				verse = strings.Replace(verse, `LORD`, "__**LORD**__", -1)
			} else {
				verse = strings.Replace(verse, `LORD`, "* __***LORD***__ *", -1)
			}
		} else {
			verse = strings.Replace(verse, `LORD`, "__***LORD***__", -1)
		}

		_, err = s.ChannelMessageSend(m.ChannelID, verse)
		if err != nil {
			log.Println(err)
		}
		log.Println(req.URL)
	}
}

func html2markup(text string) string {
	text = strings.Replace(text, `<em>`, "*", -1)
	text = strings.Replace(text, `</em>`, "*", -1)
	text = strings.Replace(text, `<sup>`, "**", -1)
	text = strings.Replace(text, `</sup> `, "** ", -1)
	text = strings.Replace(text, `</sup>`, "** ", -1)
	return text
}
