package dscore

type keep struct {
	repo, hidden bool
}
type storedata struct {
	sourcedir, central bool
}
type prefs struct {
	keep
	storedata
	//useGitIgnore     bool //will not move files marked by gitignore (eventually)
	//ignorefile/internalignore; avoid this for now to simplify
}

// hold global settings
type GlobalData struct {
	prefs                  //these will act as default, to be overridden with flags
	mainStoragePath string //primary location to use as destination.
	//any config using the primaryStore will default to a subdirectory ./[ Config.name ]
}

// global defaults for settings fields
var globalDefault = GlobalData{
	prefs:           prefs{keep: keep{repo: true, hidden: true}, storedata: storedata{sourcedir: true, central: true}},
	mainStoragePath: "~/.config/dotstrike/store",
}
