package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

type Sangam struct {
	url        string
	nameRegex  *regexp.Regexp
	menu       *Menu
	cachedHtml string
	ctx        context.Context
}

func InitSangam(ctx context.Context) *Sangam {
	return &Sangam{
		url:       "https://www.sangam-aalen.de/speisekarte",
		nameRegex: regexp.MustCompile("^((\\w*\\d+)\\s*[-–]{1}\\s*)?(([\\w\\.-äöüÄÖÜß]{2,} ?)+).*$"),
		ctx:       ctx,
	}
}

func (s *Sangam) updateHtmlCache() error {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(s.ctx, http.MethodGet, s.url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("User-Agent", "Golang")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http request to %s returned code %s", s.url, res.Status)
	}
	b, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	s.cachedHtml = string(b)
	return nil
}

func (s *Sangam) updateMenuFromCache() error {
	newMenu := &Menu{
		Items: []MenuItem{},
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s.cachedHtml))
	if err != nil {
		return err
	}
	doc.Find(".menuItemBox").Each(func(i int, selection *goquery.Selection) {
		nameElement := selection.Find(".menuItemName").Text()
		matches := s.nameRegex.FindStringSubmatch(nameElement)
		if len(matches) == 0 {
			return
		}
		id := matches[2]
		if id == "" {
			id = "<missing>"
		}
		name := strings.Trim(matches[3], " ")
		price, err := strconv.ParseFloat(strings.Replace(strings.TrimSuffix(selection.Find(".menuItemPrice").Text(), "€"), ",", ".", 1), 32)
		if err != nil {
			return
		}
		newMenu.Items = append(newMenu.Items, MenuItem{Id: id, Name: name, Price: float32(price)})
	})
	s.menu = newMenu
	return nil
}

func (s *Sangam) UpdateCache() error {
	log.Ctx(s.ctx).Debug().Msg("update HTML cache in sangam provider")
	if err := s.updateHtmlCache(); err != nil {
		if s.cachedHtml != "" {
			log.Ctx(s.ctx).Warn().Err(err).Msg("could not load sangam menu HTML. Will fall back to cached version")
		} else {
			log.Ctx(s.ctx).Error().Err(err).Msg("could not load sangam menu HTML. No cached version to fall back to")
			return fmt.Errorf("could not load sangam menu html: %s", err.Error())
		}
	}
	log.Ctx(s.ctx).Debug().Msg("update menu in sangam provider")
	if err := s.updateMenuFromCache(); err != nil {
		if s.menu != nil {
			log.Ctx(s.ctx).Warn().Err(err).Msg("could not parse sangam menu. Will fall back to cached version")
		} else {
			log.Ctx(s.ctx).Error().Err(err).Msg("could not parse sangam menu. No cached version to fall back to")
			return fmt.Errorf("could not parse sangam menu: %s", err.Error())
		}
	}
	return nil
}

func (s *Sangam) GetMenu() *Menu {
	return s.menu
}

func (s *Sangam) GetName() string {
	return "sangam"
}

func (s *Sangam) CheckItems(inItems []string) []string {
	for _, menuItem := range s.menu.Items {
		for j, inItem := range inItems {
			if menuItem.Id != inItem {
				continue
			}
			log.Ctx(s.ctx).Debug().Msgf("'%s' is in sangam menu", inItem)
			// remove element from the list
			inItems[j] = inItems[len(inItems)-1]
			inItems = inItems[:len(inItems)-1]
		}
	}
	log.Ctx(s.ctx).Debug().Msgf("items '%v' are not in sangam menu", inItems)
	return inItems
}
