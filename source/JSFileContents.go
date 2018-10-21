package source

import (
	"strings"
)

// JSFileContents describes a systemjs file
type JSFileContents struct {
	preamble         []string
	imports          []string
	body             []string
	sourceMappingURL string
	lineCount        int
	isSystemJS       bool
}

// BundleLines returns a list of lines from the body ready to include in a SystemJSBundle
func (jsfc *JSFileContents) BundleLines() []string {
	return jsfc.body
}

// ParseJSFileContents parses the contents of a JS file
func ParseJSFileContents(name string, fileContents string) (*JSFileContents, error) {
	lines := stringToLines(fileContents)

	numLines := len(lines)
	var imports []string
	var foundRegister = false
	var sourceMappingURL = ""
	var foundSourceMap = false
	var body []string
	var preamble []string
	var numPreambleLines int
	if numLines > 0 {
		preamble, numPreambleLines = skipPreamble(lines)
		registerLine := lines[numPreambleLines]
		imports, foundRegister = ParseRegisterDependencies(registerLine, false)

		sourceMapLine := lines[numLines-1]
		sourceMappingURL, foundSourceMap = parseSourceMappingURL(sourceMapLine)

		body = chooseBodyLines(lines, numPreambleLines, foundSourceMap)
	}

	bodyCopy := []string(nil)
	bodyCopy = append(bodyCopy, preamble...)
	bodyCopy = append(bodyCopy, body...)
	bodyCopy[len(preamble)] = getRegisterLineForBundle(name, imports)

	return &JSFileContents{
		preamble:         preamble,
		imports:          imports,
		body:             bodyCopy,
		sourceMappingURL: sourceMappingURL,
		lineCount:        numLines,
		isSystemJS:       foundRegister,
	}, nil
}

// getRegisterLineForBundle outputs the System.register line with a name
func getRegisterLineForBundle(name string, imports []string) string {
	importsJoined := strings.Join(imports, ", ")
	return "System.register(\"" + name + ".js\", [" + importsJoined + "], function (exports_1, context_1) {"
}
