package pops

type operationType int

type PathOperator struct {
	cursor string
	cwd    string
}

type oplog struct {
}

const (
	internal operationType = iota
	ls
	cd
	sel
	selAdd
	unsel
	dirMake
	dirMakeTemp
	dirDelete
	dirRename
	dirMove
	dirGlob
	dirGlobRec
	fileMake
	fileMakeTemp
	fileDelete
	fileRename
	fileMove
	fileRead
	fileReadPart
	fileWrite
	fileWritePart
)
