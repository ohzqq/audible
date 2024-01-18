package audible

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/ohzqq/avtools"
)

type ContentRequest struct {
	*Request
}

type ContentResponse struct {
	Meta struct {
		ChapterInfo      `json:"chapter_info"`
		ContentReference struct {
			Asin string `json:"asin"`
		} `json:"content_reference"`
	} `json:"content_metadata"`
}

type ChapterInfo struct {
	IntroDurationMs int       `json:"brandIntroDurationMs"`
	OutroDurationMs int       `json:"brandOutroDurationMs"`
	IsAccurate      bool      `json:"is_accurate"`
	DurationMs      int       `json:"runtime_length_ms"`
	DurationSec     int       `json:"runtime_length_sec"`
	Chapters        []Chapter `json:"chapters"`
}

type Chapter struct {
	Length    int    `json:"length_ms"`
	StartMs   int    `json:"start_offset_ms"`
	ChapTitle string `json:"title"`
}

func Content() *ContentRequest {
	req := &ContentRequest{Request: NewRequest()}
	req.SetParam("response_groups", responseGroups[content])
	req.AppendPath(content)
	return req
}

func (c *ContentRequest) URL(u string) (*ContentResponse, error) {
	_, err := c.ParseURL(u)
	if err != nil {
		return &ContentResponse{}, err
	}
	c.AppendPath("metadata")

	return c.Get(c.Query.String())
}

func (c *ContentRequest) Asin(id string) (*ContentResponse, error) {
	c.AppendPath(id)
	c.AppendPath("metadata")
	c.Request.Asin = id

	return c.Get(c.Query.String())
}

func (c *ContentRequest) Get(u string) (*ContentResponse, error) {
	d, err := get(u)
	if err != nil {
		return &ContentResponse{}, err
	}

	res := &ContentResponse{}
	err = json.Unmarshal(d, res)
	if err != nil {
		return &ContentResponse{}, err
	}

	return res, nil
}

func (c *ContentResponse) Chapters() []*avtools.Chapter {
	var chaps []*avtools.Chapter
	for _, ch := range c.Meta.ChapterInfo.Chapters {
		chaps = append(chaps, avtools.NewChapter(ch))
	}
	return chaps
}

func (c Chapter) Start() time.Duration {
	ss := strconv.Itoa(c.StartMs)
	dur, err := time.ParseDuration(ss + "ms")
	if err != nil {
		log.Fatal(err)
	}
	return dur
}

func (c Chapter) End() time.Duration {
	to := strconv.Itoa(c.StartMs + c.Length)
	dur, err := time.ParseDuration(to + "ms")
	if err != nil {
		log.Fatal(err)
	}
	return dur
}

func (c Chapter) Title() string {
	return c.ChapTitle
}
