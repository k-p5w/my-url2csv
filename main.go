package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// WeatherObservations is ユーザ情報
type WeatherObservations struct {
	Lastupdate string
	Outputname string
	station    []WindInfo
}

// WindInfo is 地点ごとの情報
type WindInfo struct {
	Name    string
	Avgwind string
	Maxwind string
}

func main() {
	var Ebina WeatherObservations
	var item WindInfo

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}

	getdata := func(url string) {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
		}
		lastupdate := doc.Find("lastupdate").Text()

		Ebina.Outputname = "weather.csv"
		Ebina.Lastupdate = lastupdate

		doc.Find("data").Each(func(i int, s *goquery.Selection) {

			// For each item found, get the band and title
			pos, _ := s.Attr("station")

			if len(pos) > 0 {
				s.Find("obs").Each(func(i int, s *goquery.Selection) {
					avgwind, _ := s.Attr("ws")
					maxwind, _ := s.Attr("mom")
					fmt.Printf("%v:[%v] %v,%v \n", lastupdate, pos, avgwind, maxwind)
					item.Name = pos
					item.Avgwind = avgwind
					item.Maxwind = maxwind
					Ebina.station = append(Ebina.station, item)
				})
			}
		})
	}
	for _, arg := range flag.Args() {
		getdata(arg)
	}

	WriteLog(Ebina)
}

// WriteLog is ファイルへの書き出し用
func WriteLog(dt WeatherObservations) {

	// tsv出力
	//書き込みファイル作成
	file, err := os.OpenFile(dt.Outputname, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// 表示用関数
	w := bufio.NewWriter(file)
	defer w.Flush()

	for index := 0; index < len(dt.station); index++ {

		fmt.Fprintf(w, "%v,%v,%v,%v\n", dt.Lastupdate, dt.station[index].Name, dt.station[index].Avgwind, dt.station[index].Maxwind)

	}
	w.Flush() // ファイル出力
}
