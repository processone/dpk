package dpk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

type Variant struct {
	Bitrate     string
	ContentType string `json:"content_type"`
	Url         string
}

// Sort variants by reverse bitrate. highest quality video will be first.
type Variants []Variant

func (v Variants) Len() int           { return len(v) }
func (v Variants) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v Variants) Less(i, j int) bool { return v[i].Bitrate > v[j].Bitrate }

type VideoInfo struct {
	AspectRatio []string `json:"aspect_ratio"`
	Variants    Variants
}

// Media filename is TweetID-FilePartOfTheMediaUrl
type Media struct {
	// Type can be "photo", "animated_gif", "video"
	Type      string
	Url       string
	MediaUrl  string    `json:"media_url"`
	VideoInfo VideoInfo `json:"video_info"`
}

type ExtendedEntities struct {
	Media []Media
}

type Tweet struct {
	Id               string `json:"id_str"`
	FullText         string `json:"full_text"`
	Lang             string
	Retweeted        bool
	FavoriteCount    string `json:"favorite_count"`
	RetweetCount     string `json:"retweet_count"`
	CreatedAt        string `json:"created_at"`
	ReplyToTweetId   string `json:"in_reply_to_status_id_str"`
	ReplyToUser      string `json:"in_reply_to_screen_name"`
	ClientLink       string `json:"source"`
	Entities         Entities
	ExtendedEntities ExtendedEntities `json:"extended_entities"`
	Truncated        bool
	Timestamp        time.Time
}

// Implements Sorter interface
type Tweets []Tweet

func (t Tweets) Len() int           { return len(t) }
func (t Tweets) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t Tweets) Less(i, j int) bool { return t[i].Timestamp.Before(t[j].Timestamp) }

//=============================================================================
// Post metadata struct for marshaling

type Metadata struct {
	Type      string
	Lang      string
	HashTags  []HashTag `json:",omitempty"`
	CreatedAt time.Time
}

type localMedia struct {
	mediaType    string
	filename     string
	originalLink string
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

		// Create directory for post
		targetDir := filepath.Join(OutputDir, newDir, fmt.Sprintf("%03d", index))
		if err = os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}
		// Generate markdown for post
		if err = ioutil.WriteFile(filepath.Join(targetDir, "post.md"), []byte(tweetToMd(tweet)), 0644); err != nil {
			return err
		}
		// Copy media
		for _, mediafile := range getMedia(tweet) {
			err = copyFile(
				filepath.Join(archiveDir, "tweet_media", mediafile.filename),
				filepath.Join(targetDir, mediafile.filename))
			if err != nil {
				fmt.Println("Error copying", mediafile.filename)
			}
		}
		// Generate Metadata file
		metadata := Metadata{
			Type:      "microblog",
			Lang:      tweet.Lang,
			CreatedAt: tweet.Timestamp,
		}
		meta, err := json.Marshal(metadata)
		if err != nil {
			return err
		}
		if err = ioutil.WriteFile(filepath.Join(targetDir, "metadata.json"), meta, 0644); err != nil {
			return err
		}
	}

	return nil
}

func getMedia(tweet Tweet) []localMedia {
	var files []localMedia
	for _, media := range tweet.ExtendedEntities.Media {
		mediafile := ""
		switch media.Type {
		case "photo":
			mediafile = baseName(media.MediaUrl)
		case "animated_gif":
			variants := media.VideoInfo.Variants
			if len(variants) > 0 {
				sort.Sort(variants)
				mediafile = baseName(variants[0].Url)
			}
		case "video":
			variants := media.VideoInfo.Variants
			if len(variants) > 0 {
				sort.Sort(variants)
				mediafile = baseName(variants[0].Url)
			}
		}
		if mediafile != "" {
			filename := fmt.Sprintf("%s-%s", tweet.Id, mediafile)
			files = append(files, localMedia{
				mediaType:    media.Type,
				filename:     filename,
				originalLink: media.Url,
			})
		}
	}
	return files
}

func baseName(url string) string {
	name := filepath.Base(url)
	nameWithoutParams := strings.Split(name, "?")
	return nameWithoutParams[0]
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

// TODO: Render links to mentioned people to Twitter accounts.
// TODO: Replace other shortened URL buff.ly, tinyurl, etc, to remove dependency to third-party service.
func tweetToMd(tweet Tweet) string {
	// Insert two spaces at end of line to generate Markdown line break
	markdown := strings.Replace(tweet.FullText, "\n", "  \n", 0)
	// Replace Twitter URLs with original URLs
	for _, u := range tweet.Entities.Urls {
		mdURL := fmt.Sprintf("[%s](%s)", u.DisplayUrl, u.ExpandedUrl)
		markdown = strings.Replace(markdown, u.Url, mdURL, 1)
	}
	return markdown
}

// copyFile the src file to dst. Any existing file will be overwritten and will not
// copy file attributes.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
