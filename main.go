package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/ChimeraCoder/anaconda"
)

func main() {

	var (
		consumerKey       = os.Getenv("CONSUMER_KEY")
		consumerSecret    = os.Getenv("CONSUMER_SECRET")
		accessToken       = os.Getenv("ACCESS_TOKEN")
		accessTokenSecret = os.Getenv("ACCESS_TOKEN_SECRET")
	)

	api := anaconda.NewTwitterApiWithCredentials(accessToken, accessTokenSecret, consumerKey, consumerSecret)

	// tweet, _ := api.GetTweet(1246007863895449600, nil)
	// m := tweet.InReplyToUserID
	// fmt.Println(m)

	val := url.Values{
		"track": []string{"@bulundindia1337"},
	}
	s := api.PublicStreamFilter(val)

	for t := range s.C {
		switch v := t.(type) {
		case anaconda.Tweet:
			fmt.Printf("New mention %d \n", v.Id)
			fmt.Printf("Parent tweet id %d\n", v.InReplyToStatusID)
			parentID := v.InReplyToStatusID
			a := praseURL(parentID, api)
			err := downloadFile(a)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(a)
		}
	}
}

func praseURL(tweetID int64, api *anaconda.TwitterApi) string {
	tweet, err := api.GetTweet(tweetID, nil)
	if err != nil {
		fmt.Println(err)
	}
	m := tweet.ExtendedEntities.Media[0].VideoInfo.Variants[2].Url
	v := strings.SplitAfter(m, "mp4")
	return v[0]
}

func downloadFile(url string) error {
	// f := strings.Split(url, "/")
	// filename := f[len(f)-1]
	filename := path.Base(url)
	out, err := os.Create("/home/nish/Documents/" + filename)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(n, "bytes downloaded")

	return nil
}
