package builder

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path"
	"time"
)

type sitemapURL struct {
	Loc     string `xml:"loc"`
	Lastmod string `xml:"lastmod"`
}

type sitemap struct {
	XMLName xml.Name
	Urlset  []sitemapURL `xml:"url"`
}

func loadSitemap() (sm *sitemap, err error) {

	smPath := path.Join(cfg.ContentFolder, "sitemap.xml")

	sm = &sitemap{}
	sm.XMLName.Local = "urlset"
	sm.XMLName.Space = "http://www.sitemaps.org/schemas/sitemap/0.9"

	smBytes, err := ioutil.ReadFile(smPath)
	if os.IsNotExist(err) {
		err = nil
		return
	}
	if err != nil {
		return
	}

	err = xml.Unmarshal(smBytes, &sm)
	if err != nil {
		return
	}

	sm.XMLName.Local = "urlset"
	sm.XMLName.Space = "http://www.sitemaps.org/schemas/sitemap/0.9"

	return
}

func (s *sitemap) save() (err error) {

	sb := bytes.Buffer{}

	sb.Write([]byte(xml.Header))

	smBytes, err := xml.Marshal(s)
	if err != nil {
		return
	}

	sb.Write(smBytes)

	smBytes = sb.Bytes()

	smPath := path.Join(cfg.ContentFolder, "sitemap.xml")
	err = ioutil.WriteFile(smPath, smBytes, 0777)
	if err != nil {
		return
	}

	smPath = path.Join(cfg.PublicFolder, "sitemap.xml")
	err = ioutil.WriteFile(smPath, smBytes, 0777)
	if err != nil {
		return
	}

	return
}

func (s *sitemap) Ensure(url string) {
	for _, u := range s.Urlset {
		if u.Loc == url {
			return
		}
	}
	s.Urlset = append(s.Urlset,
		sitemapURL{
			Loc:     url,
			Lastmod: time.Now().UTC().Format("2006-01-02"),
		},
	)
}

func (s *sitemap) Update(url string) {
	for i, u := range s.Urlset {
		if u.Loc == url {
			s.Urlset[i].Lastmod = time.Now().UTC().Format("2006-01-02")
			return
		}
	}
	s.Urlset = append(s.Urlset,
		sitemapURL{
			Loc:     url,
			Lastmod: time.Now().UTC().Format("2006-01-02"),
		},
	)
}

func (s *sitemap) Delete(url string) {
	for i, u := range s.Urlset {
		if u.Loc == url {
			s.Urlset = append(s.Urlset[:i], s.Urlset[i+1:]...)
			return
		}
	}
}
