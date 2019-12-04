package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	/*TODO: add code and additional functions to do the following:
	- Add an HTTP header to the response with the name
	 `Access-Control-Allow-Origin` and a value of `*`. This will
	  allow cross-origin AJAX requests to your server.
	- Get the `url` query string parameter value from the request.
	  If not supplied, respond with an http.StatusBadRequest error.
	- Call fetchHTML() to fetch the requested URL. See comments in that
	  function for more details.
	- Call extractSummary() to extract the page summary meta-data,
	  as directed in the assignment. See comments in that function
	  for more details
	- Close the response HTML stream so that you don't leak resources.
	- Finally, respond with a JSON-encoded version of the PageSummary
	  struct. That way the client can easily parse the JSON back into
	  an object. Remember to tell the client that the response content
	  type is JSON.
	Helpful Links:
	https://golang.org/pkg/net/http/#Request.FormValue
	https://golang.org/pkg/net/http/#Error
	https://golang.org/pkg/encoding/json/#NewEncoder
	*/

	/*
		// hard-coded test page
		testPage := "https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/"
	*/

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "URL not supplied", 400)
		return
	}

	// fetch
	body, err := fetchHTML(url)
	defer body.Close()
	if err != nil {
		http.Error(w, "error fetching URL", 400)
	}

	// extract
	pageSummary, err := extractSummary(url, body)
	if err != nil {
		http.Error(w, "error extracting summary", 400)
	}

	json.NewEncoder(w).Encode(pageSummary)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {
	/*TODO: Do an HTTP GET for the page URL. If the response status
	code is >= 400, return a nil stream and an error. If the response
	content type does not indicate that the content is a web page, return
	a nil stream and an error. Otherwise return the response body and
	no (nil) error.
	To test your implementation of this function, run the TestFetchHTML
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestFetchHTML
	Helpful Links:
	https://golang.org/pkg/net/http/#Get
	*/

	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New("unsuccessful response")
	}

	contentType := resp.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "text/html") {
		return nil, errors.New("not a webpage")
	}

	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {
	/*TODO: tokenize the `htmlStream` and extract the page summary meta-data
	according to the assignment description.
	To test your implementation of this function, run the TestExtractSummary
	test in summary_test.go. You can do that directly in Visual Studio Code,
	or at the command line by running:
		go test -run TestExtractSummary
	Helpful Links:
	https://drstearns.github.io/tutorials/tokenizing/
	http://ogp.me/
	https://developers.facebook.com/docs/reference/opengraph/
	https://golang.org/pkg/net/url/#URL.ResolveReference
	*/

	summary := PageSummary{}
	tokenizer := html.NewTokenizer(htmlStream)
	for {
		tokenType := tokenizer.Next()

		switch tokenType {
		case html.ErrorToken:
			err := tokenizer.Err()
			if err == io.EOF {
				// end of file, something went wrong if we haven't returned yet
				return &summary, err
			}
			// otherwise it's error tokenizing
			return &summary, err

		// deal with start and self-closing tag in one case
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			switch token.Data {
			case "title":
				// don't override open graph title
				if len(summary.Title) == 0 {
					tokenType = tokenizer.Next()
					if tokenType == html.TextToken {
						summary.Title = tokenizer.Token().Data
					}
				}
			case "meta":
				for _, attr := range token.Attr {
					switch attr.Key {
					// deal with open graph property
					case "property":
						for _, contentAttr := range token.Attr {
							if contentAttr.Key == "content" {
								switch attr.Val {
								case "og:type":
									summary.Type = contentAttr.Val
								case "og:url":
									url, err := resolveURL(pageURL, contentAttr.Val)
									if err == nil {
										summary.URL = url
									}
								case "og:title":
									summary.Title = contentAttr.Val
								case "og:site_name":
									summary.SiteName = contentAttr.Val
								case "og:description":
									summary.Description = contentAttr.Val
								case "og:image":
									image := PreviewImage{}
									url, err := resolveURL(pageURL, contentAttr.Val)
									if err == nil {
										image.URL = url
									}
									summary.Images = append(summary.Images, &image)
								case "og:image:url":
									url, err := resolveURL(pageURL, contentAttr.Val)
									if err == nil {
										summary.Images[len(summary.Images)-1].URL = url
									}
								case "og:image:secure_url":
									url, err := resolveURL(pageURL, contentAttr.Val)
									if err == nil {
										summary.Images[len(summary.Images)-1].SecureURL = url
									}
								case "og:image:type":
									summary.Images[len(summary.Images)-1].Type = contentAttr.Val
								case "og:image:width":
									summary.Images[len(summary.Images)-1].Width, _ = strconv.Atoi(contentAttr.Val)
								case "og:image:height":
									summary.Images[len(summary.Images)-1].Height, _ = strconv.Atoi(contentAttr.Val)
								case "og:image:alt":
									summary.Images[len(summary.Images)-1].Alt = contentAttr.Val
								}
							}
						}
					case "name":
						for _, contentAttr := range token.Attr {
							if contentAttr.Key == "content" {
								switch attr.Val {
								case "description":
									// don't override open graph description
									if len(summary.Description) == 0 {
										summary.Description = contentAttr.Val
									}
								case "author":
									summary.Author = contentAttr.Val
								case "keywords":
									arr := strings.Split(contentAttr.Val, ",")
									for i := range arr {
										arr[i] = strings.TrimSpace(arr[i])
									}
									summary.Keywords = arr
								}
							}
						}
					}
				}
			// icon
			case "link":
				for _, relAttr := range token.Attr {
					// look for rel attribute
					if relAttr.Key == "rel" && relAttr.Val == "icon" {
						// yes it's an icon
						icon := PreviewImage{}
						for _, hrefAttr := range token.Attr {
							if hrefAttr.Key == "href" {
								url, err := resolveURL(pageURL, hrefAttr.Val)
								if err == nil {
									icon.URL = url
								}
							}
						}
						for _, sizesAttr := range token.Attr {
							if sizesAttr.Key == "sizes" {
								sizeArr := strings.Split(sizesAttr.Val, "x")
								for i := range sizeArr {
									sizeArr[i] = strings.TrimSpace(sizeArr[i])
								}
								if len(sizeArr) == 2 {
									icon.Height, _ = strconv.Atoi(sizeArr[0])
									icon.Width, _ = strconv.Atoi(sizeArr[1])
								}
							}
						}
						for _, typeAttr := range token.Attr {
							if typeAttr.Key == "type" {
								icon.Type = typeAttr.Val
							}
						}
						summary.Icon = &icon
					}
				}
			}

		case html.EndTagToken:
			token := tokenizer.Token()
			if "head" == token.Data {
				// we're at the end of head, return
				return &summary, nil
			}
		}
	}
}

func resolveURL(pageURL string, relativeURL string) (string, error) {
	u, err := url.Parse(relativeURL)
	if err != nil {
		return "", err
	}
	base, err := url.Parse(pageURL)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(u).String(), nil
}