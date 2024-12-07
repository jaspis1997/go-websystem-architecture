package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := initWeb()
	engine.Run(fmt.Sprintf("localhost:%s", os.Getenv("PORT")))
}

func routes(engine *gin.Engine) *gin.Engine {
	engine.StaticFile("/favicon.ico", os.Getenv("VIEWS_ROOT_PATH")+"/public/favicon.ico")
	engine.Static("/static", os.Getenv("VIEWS_ROOT_PATH")+"/public")
	engine.Static("/assets", os.Getenv("VIEWS_ROOT_PATH")+"/dist/assets")
	engine.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	engine.GET("/page/*path", func(c *gin.Context) {
		c.HTML(http.StatusOK, path.Join("pages", c.Param("path")), nil)
	})
	return engine
}

func parseTemplates(rootPath string) (*template.Template, error) {
	t, err := template.ParseFiles(rootPath + "/templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %v", err)
	}
	_, err = t.ParseGlob(rootPath + "/partials/*.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse partial templates: %v", err)
	}
	var parseTemplate func(parentPath string) error
	parseTemplate = func(parentPath string) error {
		rootPath := filepath.Join(rootPath, "templates", parentPath)
		pageDirs, err := os.ReadDir(rootPath)
		if err != nil {
			return fmt.Errorf("failed to read page directory: %v", err)
		}
		for _, dir := range pageDirs {
			// Skip non-directories
			if !dir.IsDir() {
				continue
			}
			// Skip hidden directories
			if strings.HasPrefix(dir.Name(), "__") {
				continue
			}
			var p *template.Template
			if _, err := os.Stat(filepath.Join(rootPath, dir.Name(), "index.html")); err == nil {
				name := path.Join(parentPath, dir.Name())
				p, err = template.New(name).ParseFiles(path.Join(rootPath, dir.Name(), "index.html"))
				if err != nil {
					log.Printf("failed to parse page template: %v", err)
					return err
				}
				t.AddParseTree(p.Name(), p.Lookup("index.html").Tree)
			}
			_, err := t.ParseGlob(rootPath + "/" + dir.Name() + "/_*.html")
			if err != nil {
				log.Printf("failed to parse page partial templates: %v", err)
			}
			if err = parseTemplate(path.Join(parentPath, dir.Name())); err != nil {
				return err
			}
		}
		return nil
	}
	if err := parseTemplate(""); err != nil {
		return nil, err
	}
	for _, v := range t.Templates() {
		log.Print(v.Name())
	}
	return t, nil
}

func initWeb() *gin.Engine {
	e := gin.Default()
	templates, err := parseTemplates(os.Getenv("TEMPLATE_ROOT_PATH"))
	if err != nil {
		panic(err)
	}
	e.SetHTMLTemplate(templates)

	return routes(e)
}
