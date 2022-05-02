package filenamify

import (
	"errors"
	"math"
	"path/filepath"
	"regexp"
)

type Options struct {
	// String for substitution
	Replacement string
	// maxlength
	MaxLength int
}

const MAX_FILENAME_LENGTH = 100

func Filenamify(str string, options Options) (string, error) {
	var replacement string

	reControlCharsRegex := regexp.MustCompile("[\u0000-\u001f\u0080-\u009f]")

	reRelativePathRegex := regexp.MustCompile(`^\.+`)

	// https://github.com/sindresorhus/filename-reserved-regex/blob/master/index.js
	filenameReservedRegex := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	filenameReservedWindowsNamesRegex := regexp.MustCompile(`(?i)^(con|prn|aux|nul|com[0-9]|lpt[0-9])$`)

	if options.Replacement == "" {
		replacement = "!"
	} else {
		replacement = options.Replacement
	}

	if filenameReservedRegex.MatchString(replacement) && reControlCharsRegex.MatchString(replacement) {
		return "", errors.New("Replacement string cannot contain reserved filename characters")
	}

	// reserved word
	str = filenameReservedRegex.ReplaceAllString(str, replacement)

	// continue
	str = reControlCharsRegex.ReplaceAllString(str, replacement)
	str = reRelativePathRegex.ReplaceAllString(str, replacement)

	// for repeat
	if len(replacement) > 0 {
		str = trimRepeated(str, replacement)

		if len(str) > 1 {
			str = stripOuter(str, replacement)
		}
	}

	// for windows names
	if filenameReservedWindowsNamesRegex.MatchString(str) {
		str = str + replacement
	}

	// limit length
	var limitLength int
	if options.MaxLength > 0 {
		limitLength = options.MaxLength
	} else {
		limitLength = MAX_FILENAME_LENGTH
	}
	strBuf := []byte(str)
	strBuf = strBuf[0:int(math.Min(float64(limitLength), float64(len(strBuf))))]

	return string(strBuf), nil
}

func FilenamifyV2(str string, optFuns ...func(options *Options)) (string, error) {
	options := Options{
		Replacement: "!",
		MaxLength:   MAX_FILENAME_LENGTH,
	}
	for _, fn := range optFuns {
		fn(&options)
	}
	return Filenamify(str, options)
}

func Path(filePath string, options Options) (string, error) {
	p, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	p, err = Filenamify(filepath.Base(p), options)
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(p), p), nil
}

func PathV2(str string, optFuns ...func(options *Options)) (string, error) {
	options := Options{
		Replacement: "!",
		MaxLength:   MAX_FILENAME_LENGTH,
	}
	for _, fn := range optFuns {
		fn(&options)
	}
	return Path(str, options)
}

func escapeStringRegexp(str string) string {
	// https://github.com/sindresorhus/escape-string-regexp/blob/master/index.js
	reg := regexp.MustCompile(`[|\\{}()[\]^$+*?.-]`)
	str = reg.ReplaceAllStringFunc(str, func(s string) string {
		return `\` + s
	})
	return str
}

func trimRepeated(str string, replacement string) string {
	reg := regexp.MustCompile(`(?:` + escapeStringRegexp(replacement) + `){2,}`)
	return reg.ReplaceAllString(str, replacement)
}

func stripOuter(input string, substring string) string {
	// https://github.com/sindresorhus/strip-outer/blob/master/index.js
	substring = escapeStringRegexp(substring)
	reg := regexp.MustCompile(`^` + substring + `|` + substring + `$`)
	return reg.ReplaceAllString(input, "")
}
