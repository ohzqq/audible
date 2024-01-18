package audible

import (
	"net/url"
	"path"
	"testing"
)

const contentEndpoint = `https://api.audible.com/1.0/content?response_groups=chapter_info%2Ccontent_reference`
const testContentURL = `https://www.audible.com/pd/Red-Fish-Dead-Fish-Audiobook/B07B4LFT72`
const asinTest = `B07B4LFT72`
const testCall = `https://api.audible.com/1.0/content/B07B4LFT72/metadata?response_groups=chapter_info%2Ccontent_reference`

func TestContent(t *testing.T) {
	req := Content()
	if q := req.Query.String(); q != contentEndpoint {
		t.Errorf("req query %v != endpoint %v", q, contentEndpoint)
	}
}

func TestContentURL(t *testing.T) {
	req := Content()
	res, err := req.URL(testContentURL)
	if err != nil {
		t.Errorf("return error %v", err)
	}

	tu := curl()
	tu.Path = path.Join(tu.Path, asinTest, "metadata")
	o := tu.String()
	if q := req.String(); q != o {
		t.Errorf("req query %v != endpoint %v", q, o)
	}

	if req.Request.Asin != res.Meta.ContentReference.Asin {
		t.Errorf("req asin %v != res asin %v", req.Request.Asin, res.Meta.ContentReference.Asin)
	}
}

func TestChapters(t *testing.T) {
	req := Content()
	res, err := req.Asin(asinTest)
	if err != nil {
		t.Errorf("return error %v", err)
	}

	chaps := res.Chapters()
	if info := res.Meta.ChapterInfo.Chapters; len(info) != len(chaps) {
		t.Errorf("len chap info %d != len chapters %d", len(info), len(chaps))
	}
}

func TestContentAsin(t *testing.T) {
	req := Content()
	res, err := req.Asin(asinTest)
	if err != nil {
		t.Errorf("return error %v", err)
	}

	tu := curl()
	tu.Path = path.Join(tu.Path, asinTest, "metadata")
	o := tu.String()
	if q := req.Query.String(); q != o {
		t.Errorf("req query %v != endpoint %v", q, o)
	}

	if req.Request.Asin != res.Meta.ContentReference.Asin {
		t.Errorf("req asin %v != res asin %v", req.Request.Asin, res.Meta.ContentReference.Asin)
	}
}

func curl() *url.URL {
	u, _ := url.Parse(contentEndpoint)
	return u
}
