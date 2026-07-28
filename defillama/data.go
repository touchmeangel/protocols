package defillama

import (
	"errors"
	"io"

	"golang.org/x/net/html"
)

var ErrNotFound = errors.New("__NEXT_DATA__ script tag not found")

func Extract(r io.Reader) ([]byte, error) {
	z := html.NewTokenizer(r)

	for {
		switch z.Next() {
		case html.ErrorToken:
			return nil, ErrNotFound

		case html.StartTagToken:
			tok := z.Token()
			if tok.Data != "script" || !hasID(tok, "__NEXT_DATA__") {
				continue
			}
			if z.Next() != html.TextToken {
				return nil, ErrNotFound
			}
			return append([]byte(nil), z.Text()...), nil
		}
	}
}

func hasID(tok html.Token, id string) bool {
	for _, a := range tok.Attr {
		if a.Key == "id" && a.Val == id {
			return true
		}
	}
	return false
}
