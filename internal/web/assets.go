package web

import "embed"

//go:embed templates/*.html
var templates embed.FS

//go:embed static/*
var staticFiles embed.FS
