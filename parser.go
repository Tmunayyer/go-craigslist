package gocraigslist

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Listing represents a craigslist listing as a struct
type Listing struct {
	DataPID      string
	DataRepostOf string
	Date         string
	Title        string
	Link         string
	Price        string
	Hood         string
}

var nilTime = time.Time{}

func parseSearchResults(data io.Reader) ([]Listing, int, error) {
	doc, err := html.Parse(data)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to parse data: %v", err)
	}

	// find the entrypoint to  the results section of the page
	resultSection, _ := findBy(doc, "id", "sortable-results")
	// find the resultList, everything in here will go into the listing slice
	resultList, _ := findBy(resultSection, "class", "rows")

	listings := extractListings(resultList, nilTime)

	totalCountSection, _ := findBy(doc, "class", "totalcount")
	totalCount, err := strconv.Atoi(findText(totalCountSection))
	if err != nil {
		return listings, 0, fmt.Errorf("unable to parse count: %v", err)
	}

	return listings, totalCount, nil
}

func parseSearchResultsAfter(data io.Reader, date time.Time) ([]Listing, int, error) {
	doc, err := html.Parse(data)
	if err != nil {
		return nil, 0, fmt.Errorf("unable to parse data: %v", err)
	}

	// find the entrypoint to  the results section of the page
	resultSection, _ := findBy(doc, "id", "sortable-results")
	// find the resultList, everything in here will go into the listing slice
	resultList, _ := findBy(resultSection, "class", "rows")

	listings := extractListings(resultList, date)

	totalCountSection, _ := findBy(doc, "class", "totalcount")
	totalCount, err := strconv.Atoi(findText(totalCountSection))
	if err != nil {
		return listings, 0, fmt.Errorf("unable to parse count: %v", err)
	}

	return listings, totalCount, nil
}

// findBy takes a parent node and iterates recursevly through the nodes
// to return the first target that matches the attribute key and name (class, nav-bar).
func findBy(n *html.Node, attrKey string, attrName string) (*html.Node, bool) {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == attrKey && attr.Val == attrName {
				return n, true
			}
		}
	}

	if n.FirstChild != nil {
		node, has := findBy(n.FirstChild, attrKey, attrName)
		if has {
			return node, has
		}
	}

	if n.NextSibling != nil {
		node, has := findBy(n.NextSibling, attrKey, attrName)
		if has {
			return node, has
		}
	}

	return nil, false
}

// findAttr will take in the list of attributes and pull the specific one needed
func findAttr(attributes []html.Attribute, targetName string) (name string, val string) {
	for _, attr := range attributes {
		if attr.Key == targetName {
			return attr.Key, attr.Val
		}
	}

	return "", ""
}

// Text is stored as data on a TextType node as a child of the parent
// that it is contained in. This function will go all the way down
// checking node types. When FirstChild and NextSibling are both nil
// as well as the correct type, return the data
func findText(n *html.Node) (text string) {
	if n == nil {
		return ""
	}

	if n.FirstChild != nil {
		return findText(n.FirstChild)
	}

	if n.NextSibling != nil {
		return findText(n.NextSibling)
	}

	if n.Type == html.TextNode {
		return n.Data
	}

	return ""
}

func extractListings(item *html.Node, cutoffDate time.Time) []Listing {
	listings := []Listing{}
	current := item.FirstChild

	for {
		if current == nil {
			break
		}

		if current.Type != html.ElementNode {
			current = current.NextSibling
			continue
		}

		// all data housed under this node
		info, _ := findBy(current, "class", "result-info")

		if info == nil {
			current = current.NextSibling
			continue
		}

		// pull some data off the parent node that is current
		_, dataPID := findAttr(current.Attr, "data-pid")
		_, dataRepostOf := findAttr(current.Attr, "data-repost-of")

		datetimeNode, _ := findBy(info, "class", "result-date")
		_, shortdate := findAttr(datetimeNode.Attr, "datetime")
		_, longdate := findAttr(datetimeNode.Attr, "title")

		longtime := strings.Split(longdate, " ")[3]
		longtimeseconds := strings.Split(longtime, ":")[2]

		shortdatetime := strings.Split(shortdate, " ")

		datetime := strings.Join([]string{shortdatetime[0], shortdatetime[1] + ":" + longtimeseconds}, " ")

		if cutoffDate != nilTime {
			layout := "2006-01-02 15:04:05"
			tz := cutoffDate.Location()
			t, err := time.ParseInLocation(layout, datetime, tz)
			if err != nil {
				panic(err)
			}

			if t.Before(cutoffDate) {
				break
			}
		}

		linkNode, _ := findBy(info, "class", "result-title hdrlnk")
		_, link := findAttr(linkNode.Attr, "href")
		title := findText(linkNode)

		priceNode, _ := findBy(info, "class", "result-price")
		price := findText(priceNode)

		hoodNode, _ := findBy(info, "class", "result-hood")
		hood := findText(hoodNode)

		newListing := Listing{
			DataPID:      dataPID,
			DataRepostOf: dataRepostOf,
			Date:         datetime,
			Title:        title,
			Link:         link,
			Price:        price,
			Hood:         hood,
		}

		listings = append(listings, newListing)
		current = current.NextSibling
	}

	return listings
}
