package match

var ( // these are just here to think about how I want to handle all this
	exPathLocal          = ".\\foo.txt"
	exPathLocalWild      = ".\\bar*"
	exPathAnyLSubdir     = "*\\secret.hex"
	exPathAnyLSubdirWild = "*\\mix*"
	exPathAnySubdirRW    = "**\\.*"
	exPathFilename       = "tfoot.footxt"
)
