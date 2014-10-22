package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	_ "gopkg.in/cq.v1"
)

type Response struct {
	Tracks TrackResponse `json:"tracks"`
}

type TrackResponse struct {
	Items []Track `json:"items"`
}

type Track struct {
	Duration   int    `json:"duration_ms"`
	Id         string `json:"id"`
	PreviewURL string `json:"preview_url"`
	Album
}

type Album struct {
	Name   string  `json:"name"`
	Images []Image `json:"images"`
}

type Image struct {
	Height int    `json:"width"`
	Width  int    `json:"width"`
	Url    string `json:"url"`
}

func main() {
	db, err := sql.Open("neo4j-cypher", "http://localhost:7474")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare(`
	  MATCH (t:Track)
	  return t.name, t.artist
   `)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var artist string
		err := rows.Scan(&name, &artist)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(fmt.Sprintf("https://api.spotify.com/v1/search?q=track:%s+artist:%s&type=track", url.QueryEscape(name), url.QueryEscape(artist)))
		defer resp.Body.Close()

		r := Response{}
		err = json.NewDecoder(resp.Body).Decode(&r)
		if err != nil {
			log.Fatal(err)
		}

		if len(r.Tracks.Items) > 0 {
			track = r.Tracks.Items[0]
			// update the db
			imgUrl := ""
			if len(r.Tracks.Items[0].Album.Images) > 0 {
				r.Tracks.Items[0].Album.Images[0].Url
			}
			_, err := db.Exec(`
			MATCH (t:Track {name:{0}, artist:{1}}) 
			  SET t.duration = {2}, 
			      t.trackId = {3},
			      t.previewUrl = {4},
			      t.album = {5},
			      t.albumImageUrl = {6}
			  `, name, artist, track.Duration, track.Id, track.PreviewURL, track.Album.Name, imgUrl)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
