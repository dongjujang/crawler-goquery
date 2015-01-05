package main

import (
      "fmt"
      "log"
      "strings"
      "github.com/PuerkitoBio/goquery"
      "gopkg.in/mgo.v2"
)

type Movie struct {
  Name string
  Url string
}

func main() {
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
    panic(err)
  }
  defer session.Close()
  session.SetMode(mgo.Monotonic, true)

  c := session.DB("movie").C("korea")
  
  doc, err := goquery.NewDocument("http://www.torrentbest.net/bbs/board.php?bo_table=torrent_movie_kor")
  if err != nil {
    log.Fatal(err)
  }

  doc.Find("td.subject").Each(func(i int, s *goquery.Selection) {
    subject := s.Find("a").Text()
    val, exists := s.Find("a").Attr("href")

    str := "http://www.torrentbest.net"
    substr := string([]byte(val[2:]))
    url := str + substr

    fmt.Printf(" %d: %s - %s    %t\n", i, subject, url, exists)

    doc2, err := goquery.NewDocument(url)
    if err != nil {
      log.Fatal(err)
    }
    doc2.Find("td.view_file").Each(func(j int, s2 *goquery.Selection) {
      magnet, exist := s2.Find("a").Attr("href")
      if strings.Contains(magnet, "magnet") {
        c.Insert(&Movie{subject, magnet})
        fmt.Printf("magnet-->  %s    %t\n", magnet, exist)
      }
    })
  })             
}                             
