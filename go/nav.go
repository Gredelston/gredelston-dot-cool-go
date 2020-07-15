package main

// Nav contains the data necessary to render a single element of the navbar.
type Nav struct {
	Text  string
	HRef  string
	Active bool
}

// These are the Nav items we'll typically render.
// They are individually named so that they can be individually modified, such as setting as Active.
var (
	NavHome  = Nav{Text: "Home", HRef: "/"}
	NavAbout = Nav{Text: "About Me", HRef: "/about"}
)

// NewNavs returns the standard array of Navs.
func NewNavs() []Nav {
	return []Nav{NavHome, NavAbout}
}

// NewNavsWithActive returns the standard array of Navs, with one set as Active.
func NewNavsWithActive(active Nav) []Nav {
	ns := NewNavs()
	for i, n := range ns {
		ns[i].Active = (n == active)
	}
	return ns
}
