package dpk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

//=============================================================================
// Structs for JSON parsing data from Twitter

type Url struct {
	Url         string
	ExpandedUrl string `json:"expanded_url"`
	DisplayUrl  string `json:"display_url"`
	Indices     []string
}

type Mention struct {
	Name       string
	ScreenName string `json:"screen_name"`
}

type HashTag struct {
	Text    string
	Indices []string
}

type Symbol struct {
	Text    string
	Indices []string
}

type Entities struct {
	HashTags     []HashTag
	Symbols      []Symbol
	UserMentions []Mention
	Urls         []Url
}

type Tweet struct {
	Id             string `json:"id_str"`
	FullText       string `json:"full_text"`
	Lang           string
	Retweeted      bool
	FavoriteCount  string `json:"favorite_count"`
	RetweetCount   string `json:"retweet_count"`
	CreatedAt      string `json:"created_at"`
	ReplyToTweetId string `json:"in_reply_to_status_id_str"`
	ReplyToUser    string `json:"in_reply_to_screen_name"`
	ClientLink     string `json:"source"`
	Entities       Entities
	// TODO Media support
	Truncated bool
	Timestamp time.Time
}

// Implements Sorter interface
type Tweets []Tweet

func (t Tweets) Len() int           { return len(t) }
func (t Tweets) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tweets) Less(i, j int) bool { return t[i].Timestamp.Before(t[j].Timestamp) }

//=============================================================================
// Post metadata struct for marshaling

type Metadata struct {
	CreatedAt time.Time
}

//=============================================================================
// Data conversion

func TwitterToMD(archiveDir, OutputDir string) error {
	// =================================
	// Read Tweets
	data, err := ioutil.ReadFile(filepath.Join(archiveDir, "tweet.js"))
	if err != nil {
		return err
	}

	var tweets Tweets
	jsonData := bytes.TrimPrefix(data, []byte("window.YTD.tweet.part0 = "))
	if err = json.Unmarshal(jsonData, &tweets); err != nil {
		return err
	}

	// =================================
	// Parse the date for all tweets
	for i, tweet := range tweets {
		tweets[i].Timestamp, err = rubyDateToTime(tweet.CreatedAt)
		if err != nil {
			return err
		}
	}

	// =================================
	// Sort tweets by creation date
	sort.Sort(tweets)

	// =================================
	// Convert each tweet to Markdown and prepare a directory structure for it
	index := 1
	currentDir := ""
	for _, tweet := range tweets {
		if isReply(tweet) {
			continue
		}

		if isTruncated(tweet) {
			continue
		}

		year := tweet.Timestamp.Year()
		month := tweet.Timestamp.Month()
		day := tweet.Timestamp.Day()
		newDir := filepath.Join(
			fmt.Sprintf("%02d", year),
			fmt.Sprintf("%02d", month),
			fmt.Sprintf("%02d", day))
		if newDir == currentDir {
			index++
		} else {
			currentDir = newDir
			index = 1
		}

		// Real stuff
		targetDir := filepath.Join(OutputDir, newDir, fmt.Sprintf("%03d", index))
		if err = os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}
		err = ioutil.WriteFile(filepath.Join(targetDir, "post.md"), []byte(tweetToMd(tweet)), 0644)
	}

	return nil
}

func rubyDateToTime(timeString string) (time.Time, error) {
	rubyDateFormat := "Mon Jan 02 15:04:05 -0700 2006"
	timestamp, err := time.Parse(rubyDateFormat, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// Check if this is a reply (in_reply_to flag or starting by @ or .
func isReply(tweet Tweet) bool {
	if tweet.ReplyToTweetId != "" {
		return true
	}
	if strings.HasPrefix(tweet.FullText, "@") {
		return true
	}
	// Usually this is a reply / complain to customer service
	if strings.HasPrefix(tweet.FullText, ". ") {
		return true
	}

	return false
}

func isTruncated(tweet Tweet) bool {
	if strings.HasSuffix(tweet.FullText, "…") {
		return true
	}
	return false
}

// TODO: Replace shorten URL with real URLs
func tweetToMd(tweet Tweet) string {
	// Insert two spaces at end of line to generate Markdown line break
	return strings.Replace(tweet.FullText, "\n", "  \n", 0)
}
