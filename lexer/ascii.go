package lexer

// Checks if byte is ASCII alphabetical
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// Checks if it is ASCII numeric
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Checks if byte is ASCII alphanumeric
func isLegalIdentChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch)
}
