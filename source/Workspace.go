package source

import (
	"os"
	"path/filepath"
	"strings"
)

// Workspace is
type Workspace struct {
	rootPath string
}

var explicitSep = os.PathSeparator

func emulateUnix() {
	explicitSep = '/'
}

// NewWorkspace returns a new workspace for the given rootpath
func NewWorkspace(rootPath string) *Workspace {
	return &Workspace{
		rootPath: normaliseFilepath(rootPath, true),
	}
}

func normaliseFilepath(filepath string, requireSuffix bool) string {
	var normalisedFilepath string
	if explicitSep == '/' {
		normalisedFilepath = strings.Replace(filepath, "\\", "/", -1)
	} else {
		normalisedFilepath = strings.Replace(filepath, "/", "\\", -1)
	}

	if requireSuffix {
		if !strings.HasSuffix(normalisedFilepath, string(explicitSep)) {
			normalisedFilepath += string(explicitSep)
		}
	}

	if explicitSep != os.PathSeparator {
		explicitSep = os.PathSeparator
	}

	return normalisedFilepath
}

// RootPath returns the workspace's root path
func (ws *Workspace) RootPath() string {
	return ws.rootPath
}

// ReadSourceFile loads a source file
func (ws *Workspace) ReadSourceFile(imp *Import) (*File, error) {
	exists := false
	absoluteFilePath := ""
	for _, ext := range []string{"", ".js"} {
		absoluteFilePath = filepath.Join(ws.rootPath, (imp.Path() + ext))
		if _, err := os.Stat(absoluteFilePath); os.IsNotExist(err) {
			continue
		}
		exists = true
		break
	}

	if exists {
		return newFile(imp.Path(), absoluteFilePath), nil
	}

	return nil, os.ErrNotExist
}

// ToRelativePath converts an absolute filepath to a root-relative path
func (ws *Workspace) ToRelativePath(absoluteFilepath string) (string, bool) {
	if strings.HasPrefix(absoluteFilepath, ws.rootPath) {
		relativeFilepath := absoluteFilepath[len(ws.rootPath):]
		relativePath := strings.Replace(relativeFilepath, "\\", "/", -1)
		return relativePath, true
	}

	return "", false
}
