package lexer

// Checks if byte is ASCII lower character
func isLower(ch byte) bool {
	return 'a' <= ch && ch <= 'z'
}

// Checks if byte is ASCII upper character
func isUpper(ch byte) bool {
	return 'A' <= ch && ch <= 'Z'
}

// Checks if byte is ASCII alphabetical or underscore
func isLetter(ch byte) bool {
	return isLower(ch) || isUpper(ch) || ch == '_'
}

// Checks if it is ASCII numeric
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Checks if byte is ASCII alphanumeric
func isLegalIdentChar(ch byte) bool {
	return isLetter(ch) || isDigit(ch) || ch == '_'
}
