package data

import (
	"fmt"
	"strings"
	"sylmark/lsp"
	"sylmark/utils"
	"time"

	"github.com/tj/go-naturaldate"
)

type DateStore map[string]string

func NewDateStore() DateStore {

	store := DateStore{}

	commons := []string{
		"today",
		"tomorrow",
		"yesterday",
	}
	t := time.Now()
	for _, c := range commons {
		d, err := naturaldate.Parse(c, t, naturaldate.WithDirection(naturaldate.Future))
		if err == nil {
			ds := utils.FormatDate(d)
			store[c] = ds
		}
	}

	day_prefix := map[string]time.Time{}
	day_prefix["next"] = t.Add(time.Hour * 24 * 7)
	day_prefix["last"] = t
	day_prefix["previous"] = t
	day_prefix["this"] = t
	day_prefix[""] = t

	days := []string{
		"monday",
		"tuesday",
		"wednesday",
		"thursday",
		"friday",
		"saturday",
		"sunday",
	}

	for k, v := range day_prefix {
		for _, d := range days {
			c := strings.TrimSpace(k + " " + d)
			d, err := naturaldate.Parse(c, v, naturaldate.WithDirection(naturaldate.Future))
			if err == nil {
				ds := utils.FormatDate(d)
				store[c] = ds
			}
		}
	}

	return store
}

func (s *Store) getDateCompletions(arg string, needEnd bool, rng lsp.Range) (items []lsp.CompletionItem) {

	time, err := naturaldate.Parse(arg, time.Now())
	if err == nil {
		date := utils.FormatDate(time)
		item, ok := getDateCompletion(arg, date, needEnd, rng)
		if ok {
			items = append(items, item)
		}
	}

	for k, v := range s.DateStore {
		item, ok := getDateCompletion(k, v, needEnd, rng)
		if ok {
			items = append(items, item)
		}
	}

	return items
}

func getDateCompletion(arg string, date string, needEnd bool, rng lsp.Range) (comp lsp.CompletionItem, ok bool) {
	var link string
	if needEnd {
		link = "[[" + date + "]]"
	} else {
		link = "[[" + date
	}
	label := fmt.Sprintf("[[%s (%s)]]", arg, date)
	comp = lsp.CompletionItem{
		Label:    label,
		Kind:     lsp.ValueCompletion,
		SortText: "a",
		TextEdit: &lsp.TextEdit{
			Range:   rng,
			NewText: link,
		},
	}
	return comp, true
}
