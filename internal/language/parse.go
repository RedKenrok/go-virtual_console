package language

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Token struct {
	Content string

	Row    int
	Column int
}

type ImportResolver func(
	importPath string,
	baseFile string,
) (
	[]byte,
	error,
)

func defaultImportResolver(
	importPath string,
	baseFile string,
) (
	[]byte,
	error,
) {
	return os.ReadFile(importPath)
}

var validSymbol = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
var validInt = regexp.MustCompile(`^-?[0-9]+$`)
var validFloat = regexp.MustCompile(`^-?[0-9]*\.[0-9]+([eE]-?[0-9]+)?$`)

func tokenize(
	input string,
) []Token {
	var tokens []Token
	var currentToken strings.Builder
	var tokenStartRow, tokenStartColumn int
	inString := false
	row, column := 1, 1

	flushToken := func() {
		if currentToken.Len() > 0 {
			tokens = append(tokens, Token{
				Content: currentToken.String(),

				Row:    tokenStartRow,
				Column: tokenStartColumn,
			})
			currentToken.Reset()
		}
	}

	for i := 0; i < len(input); i++ {
		c := input[i]
		if currentToken.Len() == 0 {
			tokenStartRow = row
			tokenStartColumn = column
		}
		if inString {
			if c == CHR_ESCAPE {
				if i+1 < len(input) {
					i++
					column++
					currentToken.WriteByte(input[i])
				}
			} else if c == CHR_STRING {
				inString = false
				currentToken.WriteByte(c)
				tokens = append(tokens, Token{
					Content: STR_STRING + currentToken.String()[1:len(currentToken.String())-1] + STR_STRING,
					Row:     tokenStartRow,
					Column:  tokenStartColumn,
				})
				currentToken.Reset()
			} else {
				currentToken.WriteByte(c)
			}
		} else {
			if c == CHR_STRING {
				inString = true
				currentToken.Reset()
				tokenStartRow = row
				tokenStartColumn = column
				currentToken.WriteByte(c)
			} else if c == CHR_LIST_START || c == CHR_LIST_END {
				flushToken()
				tokens = append(tokens, Token{
					Content: string(c),
					Row:     row,
					Column:  column,
				})
			} else if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
				flushToken()
			} else {
				currentToken.WriteByte(c)
			}
		}

		if c == '\n' || c == '\r' {
			if i+1 < len(input) && (input[i+1] == '\n' || input[i+1] == '\r') {
				i++
			}
			row++
			column = 1
		} else {
			column++
		}
	}
	flushToken()
	return tokens
}

func parseTokens(
	tokens []Token,
	fileName string,
	resolver ImportResolver,
) (
	Value,
	[]Token,
	error,
) {
	if len(tokens) == 0 {
		return Value{}, tokens, errors.New(fileName + ": unexpected EOF")
	}

	token := tokens[0]
	tokens = tokens[1:]
	if token.Content == STR_LIST_START {
		// Check for import statements.
		if len(tokens) >= 2 && tokens[0].Content == "import" {
			importToken := tokens[0]
			fileToken := tokens[1]
			tokens = tokens[2:]
			if len(tokens) == 0 || tokens[0].Content != STR_LIST_END {
				errMsg := fmt.Sprintf("%s:%d:%d missing %s after import statement", fileName, importToken.Row, importToken.Column, STR_LIST_END)
				return Value{}, tokens, errors.New(errMsg)
			}

			tokens = tokens[1:]
			importPath := fileToken.Content
			if resolver == nil {
				resolver = defaultImportResolver
			}
			data, err := resolver(importPath, fileName)

			if err != nil {
				errMsg := fmt.Sprintf("%s:%d:%d failed to import file %s: %v", fileName, importToken.Row, importToken.Column, importPath, err)
				return Value{}, tokens, errors.New(errMsg)
			}
			importValue, err := Parse(string(data), importPath, resolver)

			if err != nil {
				errMsg := fmt.Sprintf("%s:%d:%d error in imported file %s: %v", fileName, importToken.Row, importToken.Column, importPath, err)
				return Value{}, tokens, errors.New(errMsg)
			}
			return importValue, tokens, nil
		}

		var listValue []Value
		for len(tokens) > 0 && tokens[0].Content != STR_LIST_END {
			var expression Value
			var err error
			expression, tokens, err = parseTokens(
				tokens,
				fileName,
				resolver,
			)
			if err != nil {
				return Value{}, tokens, err
			}

			listValue = append(listValue, expression)
		}

		if len(tokens) == 0 {
			errMsg := fmt.Sprintf("%s:%d:%d missing %s", fileName, token.Row, token.Column, STR_LIST_END)
			return Value{}, tokens, errors.New(errMsg)
		}
		tokens = tokens[1:]

		return Value{
			Type: List,
			Data: listValue,
		}, tokens, nil
	} else if token.Content == STR_LIST_END {
		errMsg := fmt.Sprintf("%s:%d:%d unexpected %s", fileName, token.Row, token.Column, token.Content)
		return Value{}, tokens, errors.New(errMsg)
	}

	value, err := atom(token, fileName)
	if err != nil {
		return Value{}, tokens, err
	}

	return value, tokens, nil
}

func atom(
	token Token,
	fileName string,
) (
	Value,
	error,
) {
	if len(token.Content) > 0 && token.Content[0] == CHR_STRING && token.Content[len(token.Content)-1] == CHR_STRING {
		content := token.Content[1 : len(token.Content)-1]
		return Value{
			Type: String,
			Data: content,
		}, nil
	}

	if token.Content == "true" {
		return Value{
			Type: Bool,
			Data: true,
		}, nil
	}
	if token.Content == "false" {
		return Value{
			Type: Bool,
			Data: false,
		}, nil
	}

	if validFloat.MatchString(token.Content) {
		if floatValue, err := strconv.ParseFloat(token.Content, 64); err == nil {
			return Value{
				Type: Float,
				Data: floatValue,
			}, nil
		}
	}

	if validInt.MatchString(token.Content) {
		if intValue, err := strconv.ParseInt(token.Content, 10, 64); err == nil {
			return Value{
				Type: Int,
				Data: intValue,
			}, nil
		}
	}

	if validSymbol.MatchString(token.Content) {
		return Value{
			Type: Symbol,
			Data: token.Content,
		}, nil
	}

	errMsg := fmt.Sprintf("%s:%d:%d invalid token: %s", fileName, token.Row, token.Column, token.Content)
	return Value{}, errors.New(errMsg)
}

func joinTokenContents(
	tokens []Token,
) string {
	var parts []string
	for _, t := range tokens {
		parts = append(parts, t.Content)
	}
	return strings.Join(parts, " ")
}

func Parse(
	input string,
	fileName string,
	resolver ImportResolver,
) (
	Value,
	error,
) {
	tokens := tokenize(input)
	expression, remaining, err := parseTokens(
		tokens,
		fileName,
		resolver,
	)
	if err != nil {
		return Value{}, err
	}
	if len(remaining) != 0 {
		unexpected := remaining[0]
		return Value{}, errors.New(fmt.Sprintf("%s:%d:%d unexpected tokens after parsing: %s", fileName, unexpected.Row, unexpected.Column, joinTokenContents(remaining)))
	}
	return expression, nil
}
