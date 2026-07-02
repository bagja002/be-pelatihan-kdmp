// Command generate scaffolds a full clean-architecture slice for a new
// entity: entity, dto, repository, service, handler and route files.
//
// Usage:
//
//	go run ./cmd/generate -name Product
//	go run ./cmd/generate -name OrderItem -plural order_items
//	go run ./cmd/generate -name Product -force   # overwrite existing files
package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*.tmpl
var templatesFS embed.FS

// moduleName must match the module path in go.mod.
const moduleName = "knmp-backend"

type target struct {
	tmpl string
	out  string
}

type data struct {
	Name    string // PascalCase singular, e.g. Product
	VarName string // camelCase singular, e.g. product
	Plural  string // route plural, e.g. products
	File    string // file base name, e.g. product
	Module  string // go module path
}

func main() {
	name := flag.String("name", "", "entity name, singular (e.g. Product)")
	plural := flag.String("plural", "", "route plural (optional, e.g. products)")
	force := flag.Bool("force", false, "overwrite files that already exist")
	flag.Parse()

	if *name == "" {
		fmt.Println("usage: go run ./cmd/generate -name Product [-plural products] [-force]")
		os.Exit(1)
	}

	entity := pascal(*name)
	file := strings.ToLower(entity)
	pl := *plural
	if pl == "" {
		pl = pluralize(file)
	}

	d := data{
		Name:    entity,
		VarName: lowerFirst(entity),
		Plural:  pl,
		File:    file,
		Module:  moduleName,
	}

	targets := []target{
		{"templates/entity.tmpl", filepath.Join("internal", "entity", file+".go")},
		{"templates/dto.tmpl", filepath.Join("internal", "dto", file+"_dto.go")},
		{"templates/repository.tmpl", filepath.Join("internal", "repository", file+"_repository.go")},
		{"templates/service.tmpl", filepath.Join("internal", "service", file+"_service.go")},
		{"templates/handler.tmpl", filepath.Join("internal", "handler", file+"_handler.go")},
		{"templates/route.tmpl", filepath.Join("internal", "router", file+"_route.go")},
	}

	for _, t := range targets {
		if err := render(t, d, *force); err != nil {
			fmt.Printf("error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Printf("\n✅ generated %q (route group: /%s)\n\n", entity, pl)
	fmt.Println("Last step — register the route (1 line) in internal/router/router.go:")
	fmt.Printf("\n\tRegister%sRoutes(api, db)\n\n", entity)
	fmt.Println("Migration is automatic: the entity self-registers via init().")
}

func render(t target, d data, force bool) error {
	if _, err := os.Stat(t.out); err == nil && !force {
		fmt.Printf("skip (exists): %s\n", t.out)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(t.out), 0o755); err != nil {
		return err
	}

	tmpl, err := template.ParseFS(templatesFS, t.tmpl)
	if err != nil {
		return err
	}

	f, err := os.Create(t.out)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := tmpl.Execute(f, d); err != nil {
		return err
	}
	fmt.Printf("created: %s\n", t.out)
	return nil
}

// pascal converts "product", "order_item" or "order-item" to "Product",
// "OrderItem".
func pascal(s string) string {
	parts := strings.FieldsFunc(strings.TrimSpace(s), func(r rune) bool {
		return r == '_' || r == '-' || r == ' '
	})
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + strings.ToLower(p[1:]))
	}
	return b.String()
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// pluralize is a naive English pluralizer; pass -plural to override.
func pluralize(s string) string {
	switch {
	case strings.HasSuffix(s, "y") && !endsWithVowelBeforeY(s):
		return s[:len(s)-1] + "ies"
	case strings.HasSuffix(s, "s"),
		strings.HasSuffix(s, "x"),
		strings.HasSuffix(s, "z"),
		strings.HasSuffix(s, "ch"),
		strings.HasSuffix(s, "sh"):
		return s + "es"
	default:
		return s + "s"
	}
}

func endsWithVowelBeforeY(s string) bool {
	if len(s) < 2 {
		return false
	}
	switch s[len(s)-2] {
	case 'a', 'e', 'i', 'o', 'u':
		return true
	}
	return false
}
