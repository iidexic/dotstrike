package cfg

type componentType int

const (
	Ccore      componentType = iota
	CSet                     // "group" of apps
	CApp                     // bg info/app info
	CSource                  //path details
	CDest                    //path details
	CStructure               //config
)

type CMP interface {
	generate(data []string)
}

// cmpConfig stores settings for a given component
type cmpConfig struct {
	deep, use_gitignore, keep_historic bool
	ignore_style                       string
	relpath_ignore                     string
}

// Config Mananger, stores all other config components
type Manager struct {
	g          Global
	components []CMP
}

// Global contains the whole set of components for an instance of dotstrike
type Global struct {
	path_primary_store string
	path_drive2_store  string
}

// Set
type Set struct {
}

type App struct {
	name    string
	opts    cmpConfig
	appPath string
}

func DataManager() *Manager {
	mgr := Manager{}
	return &mgr
}

func (m *Manager) NewComponent(ctype componentType, data CMP) {

}
