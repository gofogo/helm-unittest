package common

import (
	"bytes"
	"regexp"
	"strings"
)

const (
	bsCode                           byte = byte('\\')
	metaCharactersNeedsEscapePattern      = `.*[.+*?()|[\]{}^$].*`
)

var metaRegex = regexp.MustCompile(metaCharactersNeedsEscapePattern)

type YmlEscapeHandlers struct{}

// Escape function is required, as yaml library no longer maintained
// yaml unmaintained library issue https://github.com/go-yaml/yaml/pull/862
func (y *YmlEscapeHandlers) Escape(content string) []byte {
	if !strings.Contains(content, `\`) && !metaRegex.MatchString(content) {
		return nil
	}
	return escapeBackslashes([]byte(content))
}

// escapeBackslashes escapes backslashes in the given byte slice.
// It ensures that an even number of backslashes are present by doubling any single backslash found.
func escapeBackslashes(content []byte) []byte {
	var result bytes.Buffer
	i := 0
	for i < len(content) {
		if content[i] != bsCode {
			result.WriteByte(content[i])
			i++
			continue
		}

		count := 1
		for i+count < len(content) && content[i+count] == bsCode {
			count++
		}

		times := count
		if count%2 == 1 {
			times++
		}

		for j := 0; j < times; j++ {
			result.WriteByte(bsCode)
		}

		i += count
	}
	return result.Bytes()
}
