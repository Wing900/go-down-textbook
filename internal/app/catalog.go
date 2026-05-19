package app

import (
	"fmt"
	"sort"

	"github.com/chenyb-go/go-down-textbook/internal/api"
	"github.com/chenyb-go/go-down-textbook/internal/auth"
	"github.com/chenyb-go/go-down-textbook/internal/download"
	"github.com/chenyb-go/go-down-textbook/internal/models"
)

func (s *Service) EnsureCatalog() (*CatalogData, error) {
	session := auth.NewSessionManager(auth.LoginViaBrowserQuiet)
	token, err := session.EnsureToken()
	if err != nil {
		return nil, err
	}

	client := api.NewClient(token)
	client.SetUnauthorizedHandler(session.RefreshToken)
	entries, err := api.FetchCatalog(client)
	if err != nil {
		return nil, err
	}

	grades := api.GetGrades(entries)
	if len(grades) == 0 {
		return nil, fmt.Errorf("未找到年级信息")
	}
	return &CatalogData{Token: token, Entries: entries, Grades: grades}, nil
}

func (s *Service) Subjects(entries []models.CatalogEntry, grade string) []string {
	subjects := api.GetSubjects(entries, grade)
	sort.Strings(subjects)
	return subjects
}

func (s *Service) Books(entries []models.CatalogEntry, grade, subject string) []SelectableBook {
	filtered := api.FilterBySubject(api.FilterByGrade(entries, grade), subject)
	books := make([]SelectableBook, 0, len(filtered))
	for _, entry := range filtered {
		books = append(books, SelectableBook{
			BookItem:   models.BookItem{ID: entry.ID, Title: entry.Title, Grade: grade, Subject: subject},
			Downloaded: download.HasBook(s.OutputDir, entry.ID),
		})
	}
	sort.SliceStable(books, func(i, j int) bool {
		if books[i].Downloaded != books[j].Downloaded {
			return !books[i].Downloaded
		}
		return books[i].Title < books[j].Title
	})
	return books
}
