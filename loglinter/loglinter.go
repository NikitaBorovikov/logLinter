package loglinter

import (
	"fmt"
	"go/ast"
	"loglinter/config"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglint",
	Doc:  "checks if log messages follow project conventions",
	Run:  run,
}

var (
	cfg    *config.Config
	cfgErr error
	once   sync.Once
)

var logMethods = map[string]bool{
	"Info":  true,
	"Error": true,
	"Warn":  true,
	"Debug": true,
}

func run(pass *analysis.Pass) (interface{}, error) {
	once.Do(func() {
		cfg, cfgErr = config.InitConfig()
	})
	if cfgErr != nil {
		return nil, fmt.Errorf("failed to init config: %v", cfgErr)
	}

	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			callExpr, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			selExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				return true
			}

			if !isLogMethod(pass, selExpr) {
				return true
			}

			if len(callExpr.Args) == 0 {
				return true
			}

			arg := callExpr.Args[0]
			lit, ok := arg.(*ast.BasicLit)
			if !ok {
				checkSensitiveDataInExpr(pass, arg)
				return true
			}

			logMsg, err := strconv.Unquote(lit.Value)
			if err != nil {
				return true
			}

			checkLogMessage(pass, lit, logMsg)
			return true
		})
	}
	return nil, nil
}

func isLogMethod(pass *analysis.Pass, sel *ast.SelectorExpr) bool {
	if !logMethods[sel.Sel.Name] {
		return false
	}

	// try to determine the type using the type analyzer
	if tv, ok := pass.TypesInfo.Types[sel.X]; ok && tv.Type != nil {
		typStr := tv.Type.String()
		for _, logPkg := range cfg.LogPkgNames {
			if strings.Contains(typStr, logPkg) {
				return true
			}
		}
	}

	// check by identifier name
	if ident, ok := sel.X.(*ast.Ident); ok {
		for _, logIdentName := range cfg.LogIndentNames {
			if ident.Name == logIdentName {
				return true
			}
		}
	}
	return false
}

func checkSensitiveDataInExpr(pass *analysis.Pass, expr ast.Expr) {
	ast.Inspect(expr, func(n ast.Node) bool {
		if ident, ok := n.(*ast.Ident); ok {
			if isContainSensitiveData(ident.Name) {
				pass.Reportf(ident.Pos(), "avoid logging sensitive variable: %s", ident.Name)
			}
		}
		return true
	})
}

func checkLogMessage(pass *analysis.Pass, node ast.Node, msg string) {
	if len(msg) == 0 {
		return
	}
	if !isFirstLetterLowerCase(msg) {
		pass.Reportf(node.Pos(), "log message should start with a lowercase letter")
		return
	}
	if !isOnlyEnglishLetters(msg) {
		pass.Reportf(node.Pos(), "log message should be in English only")
		return
	}
	if isMsgContainSpecialChars(msg) {
		pass.Reportf(node.Pos(), "log message contains forbidden characters or emojis")
		return
	}
	if isContainSensitiveData(msg) {
		pass.Reportf(node.Pos(), "log message might contain sensitive data")
		return
	}
}

func isFirstLetterLowerCase(msg string) bool {
	firstLetter := []rune(msg)[0]
	return unicode.IsLower(firstLetter)
}

func isOnlyEnglishLetters(msg string) bool {
	for _, r := range msg {
		if unicode.IsLetter(r) && r > unicode.MaxASCII {
			return false
		}
	}
	return true
}

func isMsgContainSpecialChars(msg string) bool {
	lowerMsg := strings.ToLower(msg)

	for _, r := range lowerMsg {
		isLetter := (r >= 'a' && r <= 'z')
		isDigit := (r >= '0' && r <= '9')
		isSpace := unicode.IsSpace(r)

		if !isLetter && !isDigit && !isSpace {
			return true
		}
	}
	return false
}

func isContainSensitiveData(msg string) bool {
	lowerMsg := strings.ToLower(msg)

	for _, word := range cfg.SensitiveWords {
		if strings.Contains(lowerMsg, word) {
			return true
		}
	}
	return false
}
