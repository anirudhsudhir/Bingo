package main

import "github.com/anirudhsudhir/Bingo/internal/models"

type templateData struct {
	Snip  *models.Snip
	Snips []*models.Snip
}
