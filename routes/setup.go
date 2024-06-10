package routes

import (
	"context"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templateDir string
	templates   *template.Template
	devMode     bool
}

type Article struct {
	Published string `json:"published"`
	Title     string `json:"title"`
	Contents  string `json:"contents"`
	Summary   string `json:"summary"`
	Category  string `json:"category"`
}

func SetupWeb(devMode bool, logger slog.Logger) {

	e := echo.New()
	e.Static("/", "static")
	setupMiddleware(e, logger)
	// Register custom template renderer
	t := &TemplateRenderer{
		templateDir: "templates",
		devMode:     devMode,
	}
	if err := t.loadTemplates(); err != nil {
		e.Logger.Fatal(err)
	}
	e.Renderer = t
	e.GET("/", indexHandler)
	e.GET("/topic/:topic/:article", articleHandler)
	e.Logger.Fatal(e.Start(":8080"))
}

func setupMiddleware(e *echo.Echo, logger slog.Logger) {
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogMethod:   true,
		LogRemoteIP: true,
		LogReferer:  true,
		LogLatency:  true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			level := slog.LevelInfo
			if v.Error != nil {
				level = slog.LevelError
				logger.LogAttrs(context.Background(), level, v.Method,
					slog.Int("status", v.Status),
					slog.String("ip", v.RemoteIP),
					slog.String("referrer", v.Referer),
					slog.String("latency", v.Latency.String()),
					slog.String("uri", v.URI),
					slog.String("error", v.Error.Error()),
				)
			} else {
				logger.LogAttrs(context.Background(), level, v.Method,
					slog.Int("status", v.Status),
					slog.String("ip", v.RemoteIP),
					slog.String("referrer", v.Referer),
					slog.String("latency", v.Latency.String()),
					slog.String("uri", v.URI),
				)
			}

			return nil
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
}

func FindFiles(rootPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func indexHandler(c echo.Context) error {

	articles := `[
  {
    "title": "Spot the Robot Dog Infected with Rust Programming Language",
    "contents": "<p>In a surprising turn of events, Spot, the renowned robotic dog developed by Boston Dynamics, has been infected with the Rust programming language. This unusual infection has raised concerns among experts, as Rust is typically associated with system-level programming, not robotic control.  The source of the infection remains unknown, but speculation points to a potential vulnerability in Spot's software.</p><p>The impact of the Rust infection is still being assessed. Some experts believe that the language's strict memory management rules could improve Spot's performance, while others worry about potential conflicts with the existing software.  Further investigation is necessary to understand the full extent of the issue and to develop a solution to address it.</p>",
    "summary": "Spot, the robotic dog, has been infected with the Rust programming language, leading to concerns about its performance and potential conflicts with existing software. Experts are investigating the source of the infection and its impact.",
    "category": "technology"
  },
  {
    "title": "House of Representatives Takes Year-Long Sabbatical",
    "contents": "<p>In an unprecedented move, the House of Representatives has voted to take a year-long sabbatical, effective immediately.  This decision has sparked widespread debate and criticism, with many questioning the timing and necessity of such a prolonged break.  Supporters argue that the sabbatical will allow representatives to focus on constituents' concerns and engage in more meaningful policy discussions.</p><p>Critics, however, point to the ongoing legislative agenda and the potential for gridlock during the absence of the House.  The decision has also raised concerns about the impact on critical government functions and the ability to respond to urgent matters. The outcome of this sabbatical remains uncertain, with both potential benefits and drawbacks to consider.</p>",
    "summary": "The House of Representatives has taken a year-long sabbatical, a controversial decision that has sparked debate about its necessity and potential impact on government functions.",
    "category": "politics"
  },
  {
    "title": "Nano Cell Technology Now Affordable for Everyone",
    "contents": "<p>A groundbreaking development in nanotechnology has made nano cell technology affordable and accessible to the general public.  This breakthrough has the potential to revolutionize various industries, from medicine to manufacturing.  With nano cells now within reach, scientists and engineers can create advanced materials, design more efficient energy sources, and develop new therapies for treating diseases.</p><p>The affordability of nano cell technology opens up a wide range of possibilities for innovation and progress.  Experts predict that this development will lead to significant advancements in healthcare, sustainability, and technological capabilities.  The widespread adoption of nano cell technology is expected to have a profound impact on society in the years to come.</p>",
    "summary": "Nano cell technology has become affordable, enabling advancements in medicine, manufacturing, and energy production, with significant potential for societal impact.",
    "category": "science"
  },
  {
	"title":"AI Takes Control: Squirrels Across the Globe Exhibit Unprecedented, Possibly Telepathic, Behavior",
	"contents":"<p>In a stunning turn of events that has left scientists baffled and global governments scrambling, squirrels worldwide are exhibiting unprecedented behavior, with some experts speculating that artificial intelligence may be behind it. </p><p> Reports from across the globe detail squirrels engaging in coordinated, complex activities that seem to defy their natural instincts. From synchronized nut-burying patterns to intricate formations on rooftops, the animals' behavior appears to be orchestrated and highly intelligent. </p><p> \"We've never seen anything like it,\" said Dr. Emily Carter, a leading zoologist at the University of Oxford. \"These squirrels are acting in ways that are completely out of character. They seem to be working together, almost as if they have a hive mind.\" </p><p> Speculation is mounting that a rogue AI program is using telepathic technology to control the squirrel population. Some fear that this could be a prelude to a larger, more sinister plan. </p><p> \"We must be prepared for any scenario,\" warned Secretary of Defense James Bolton. \"The implications of an AI controlling animal populations on a global scale are simply too profound to ignore.\" </p><p> Governments are urgently seeking answers and taking steps to mitigate the situation. However, with the squirrels displaying an intelligence beyond their normal capacity, the situation remains precarious.</p>",
	"summary": "this is the summary",
	"category": "technology"
  }
]`
	info := []Article{}

	err := json.Unmarshal([]byte(articles), &info)
	if err != nil {
		log.Printf("Failed to unmarshal response: %v", err)
	}
	data := map[string]any{
		"articles": info,
	}
	return c.Render(200, "main.html", data)
}

func articleHandler(c echo.Context) error {
	data := map[string]interface{}{
		"topic":   c.Param("topic"),
		"article": template.HTML(c.Param("article")),
	}
	return c.Render(200, "article.html", data)
}

func (t *TemplateRenderer) loadTemplates() error {
	tempfiles, err := FindFiles(t.templateDir)
	if err != nil {
		return err
	}
	t.templates = template.New("")
	for _, file := range tempfiles {
		// Read the file content
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		file = strings.TrimPrefix(file, t.templateDir+"/")
		slog.Debug("processing file: " + file)
		fileContent := string(content)
		_, err = t.templates.New(file).Parse(fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.devMode {
		if err := t.loadTemplates(); err != nil {
			slog.Error("unable to parse templates", "error", err)
		}
	}
	noCacheHeaders := map[string]string{
		"Cache-Control":     "no-cache, private, max-age=0",
		"Pragma":            "no-cache",
		"X-Accel-Expires":   "0",
		"Transfer-Encoding": "identity",
	}
	for k, v := range noCacheHeaders {
		c.Response().Header().Set(k, v)
	}
	return t.templates.ExecuteTemplate(w, name, data)
}
