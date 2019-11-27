package ini

// The Options type is used to configure the behavior during unmarshalling.
type Options struct {
	// AllowMultilineValues enables a property value to contain multiple lines.
	// Currently supported methods:
	//
	// - Escaped newlines: A newline character preceded by a single backslash
	// - Space-prefixed: A line beginning with one or more spaces
	AllowMultilineValues bool

	// AllowNumberSignComments treats lines beginning with the number sign (#)
	// as a comment.
	AllowNumberSignComments bool
}
