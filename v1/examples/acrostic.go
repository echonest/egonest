// adapted from pyen's example/acrostic.py
package main

import (
	"fmt"
	"github.com/echonest/egonest/v1"
	"net/url"
	"os"
	"strings"
	"unicode"
)

var en egonest.Host

func init() {
	en.Throttle = true
}

func get_songs(genre string, c chan string) {
	args := make(url.Values)
	args.Set("type", "genre-radio")
	args.Set("genre", genre)
	resp, err := en.GetCall("playlist/dynamic/create", args)
	if err != nil {
		panic(err)
	}
	decoded, err := egonest.GenericUnmarshal(resp, false)
	if err != nil {
		panic(err)
	}
	session_id, ok := egonest.Dig(decoded, "response", "session_id").(string)
	if !ok {
		panic("no session id!")
	}

	args = make(url.Values)
	args.Set("results", "5")
	args.Set("session_id", session_id)
	for {
		resp, err := en.GetCall("playlist/dynamic/next", args)
		if err != nil {
			panic(err)
		}
		decoded, err := egonest.GenericUnmarshal(resp, false)
		if err != nil {
			panic(err)
		}
		for _, song := range egonest.Dig(decoded, "response", "songs").([]interface{}) {
			c <- fmt.Sprint(egonest.Dig(song, "title"), " by ", egonest.Dig(song, "artist_name"))
		}
	}
}

func checksong(r rune, song string) bool {
	return []rune(strings.ToLower(song))[0] == r
}

func build_acrostic(message string, c chan string) {
	songbuf := make([]string, 0)
MESSAGE:
	for _, r := range []rune(message) {
		if unicode.IsLetter(r) {
			for tries := 200; tries > 0; {
				for i, song := range songbuf {
					if checksong(r, song) {
						fmt.Println(song)
						songbuf[i], songbuf[len(songbuf)-1] = songbuf[len(songbuf)-1], songbuf[i]
						songbuf = songbuf[:len(songbuf)-2]
						continue MESSAGE
					}
					tries--
				}
				for tries >= 0 {
					song := <-c
					if checksong(r, song) {
						fmt.Println(song)
						continue MESSAGE
					}
					songbuf = append(songbuf, song)
				}

			}

		} else {
			fmt.Println(string([]rune{r}))
		}
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("acrostic genre 'secret message'")
		return
	}
	c := make(chan string, len(os.Args[2]))
	go get_songs(strings.ToLower(os.Args[1]), c)
	build_acrostic(os.Args[2], c)
}
