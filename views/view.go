package views

import (
	"fmt"
	"html/template"
	"strings"
	"sync"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gin-gonic/gin/render"
)

var viewBox *rice.Box

// Render type
type Render struct {
	vbox  *rice.Box
	tpls  sync.Map
	funcs sync.Map
	ext   string
}

// LoadViews init load view box
func LoadViews() {
	viewBox = rice.MustFindBox("../views")
}

// NewRender return an render instance
func NewRender() *Render {
	r := &Render{
		vbox: viewBox,
		ext:  "html",
	}

	// set templates
	setTpls(r)

	// set funcs
	setFuncs(r)

	return r
}

func setTpls(r *Render) {
	// error
	r.tpls.Store("error", []string{"layouts/normal", "error/error"})
	// login
	r.tpls.Store("login", []string{"layouts/normal", "home/login"})
	// home
	r.tpls.Store("home", []string{"layouts/main", "layouts/nav", "home/index"})
	// user
	r.tpls.Store("user", []string{"layouts/main", "layouts/nav", "user/index", "user/search", "user/add", "user/edit"})
	// password
	r.tpls.Store("password", []string{"layouts/normal", "password/change"})
	// role
	r.tpls.Store("role", []string{"layouts/main", "layouts/nav", "role/index"})
}

func setFuncs(r *Render) {}

// Instance supply render string
func (r *Render) Instance(name string, data interface{}) render.Render {
	html := ""
	funcs := make(template.FuncMap)

	if v, ok := r.tpls.Load(name); ok {
		if tpls, ok := v.([]string); ok {
			if l := len(tpls); l > 0 {
				text := make([]string, 0, len(tpls))

				for _, name := range tpls {
					text = append(text, r.vbox.MustString(fmt.Sprintf("%s.%s", name, r.ext)))
				}

				html = strings.Join(text, "")
			}
		}
	}

	if v, ok := r.funcs.Load(name); ok {
		if fs, ok := v.(template.FuncMap); ok {
			funcs = fs
		}
	}

	return render.HTML{
		Template: template.Must(template.New(name).Funcs(funcs).Parse(html)),
		Data:     data,
	}
}
