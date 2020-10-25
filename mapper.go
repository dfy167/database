package database

import (
	"strings"
	"unsafe"
)

// Mapper represents a name convertation between struct's fields name and table's column name
type Mapper interface {
	Obj2Table(string) string
	Table2Obj(string) string
}

// SameMapper implements IMapper and provides same name between struct and
// database table
type SameMapper struct {
}

// Obj2Table Obj2Table
func (m SameMapper) Obj2Table(o string) string {
	return o
}

// Table2Obj Table2Obj
func (m SameMapper) Table2Obj(t string) string {
	return t
}

// SnakeMapper implements IMapper and provides name transaltion between
// struct and database table
type SnakeMapper struct {
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func snakeCasedName(name string) string {
	newstr := make([]byte, 0, len(name)+1)
	for i := 0; i < len(name); i++ {
		c := name[i]
		if isUpper := 'A' <= c && c <= 'Z'; isUpper {
			if i > 0 {
				newstr = append(newstr, '_')
			}
			c += 'a' - 'A'
		}
		newstr = append(newstr, c)
	}

	return b2s(newstr)
}

// Obj2Table Obj2Table
func (mapper SnakeMapper) Obj2Table(name string) string {
	return snakeCasedName(name)
}

func titleCasedName(name string) string {
	newstr := make([]byte, 0, len(name))
	upNextChar := true

	name = strings.ToLower(name)

	for i := 0; i < len(name); i++ {
		c := name[i]
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
		case c == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, c)
	}

	return b2s(newstr)
}

// Table2Obj Table2Obj
func (mapper SnakeMapper) Table2Obj(name string) string {
	return titleCasedName(name)
}

// GonicMapper implements IMapper. It will consider initialisms when mapping names.
// E.g. id -> ID, user -> User and to table names: UserID -> user_id, MyUID -> my_uid
type GonicMapper map[string]bool

func isASCIIUpper(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

func toASCIIUpper(r rune) rune {
	if 'a' <= r && r <= 'z' {
		r -= ('a' - 'A')
	}
	return r
}

func gonicCasedName(name string) string {
	newstr := make([]rune, 0, len(name)+3)
	for idx, chr := range name {
		if isASCIIUpper(chr) && idx > 0 {
			if !isASCIIUpper(newstr[len(newstr)-1]) {
				newstr = append(newstr, '_')
			}
		}

		if !isASCIIUpper(chr) && idx > 1 {
			l := len(newstr)
			if isASCIIUpper(newstr[l-1]) && isASCIIUpper(newstr[l-2]) {
				newstr = append(newstr, newstr[l-1])
				newstr[l-1] = '_'
			}
		}

		newstr = append(newstr, chr)
	}
	return strings.ToLower(string(newstr))
}

// Obj2Table Obj2Table
func (mapper GonicMapper) Obj2Table(name string) string {
	return gonicCasedName(name)
}

// Table2Obj Table2Obj
func (mapper GonicMapper) Table2Obj(name string) string {
	newstr := make([]rune, 0)

	name = strings.ToLower(name)
	parts := strings.Split(name, "_")

	for _, p := range parts {
		_, isInitialism := mapper[strings.ToUpper(p)]
		for i, r := range p {
			if i == 0 || isInitialism {
				r = toASCIIUpper(r)
			}
			newstr = append(newstr, r)
		}
	}

	return string(newstr)
}

// LintGonicMapper is A GonicMapper that contains a list of common initialisms taken from golang/lint
var LintGonicMapper = GonicMapper{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

// PrefixMapper provides prefix table name support
type PrefixMapper struct {
	Mapper Mapper
	Prefix string
}

// Obj2Table Obj2Table
func (mapper PrefixMapper) Obj2Table(name string) string {
	return mapper.Prefix + mapper.Mapper.Obj2Table(name)
}

// Table2Obj Table2Obj
func (mapper PrefixMapper) Table2Obj(name string) string {
	return mapper.Mapper.Table2Obj(name[len(mapper.Prefix):])
}

//NewPrefixMapper NewPrefixMapper
func NewPrefixMapper(mapper Mapper, prefix string) PrefixMapper {
	return PrefixMapper{mapper, prefix}
}

// SuffixMapper provides suffix table name support
type SuffixMapper struct {
	Mapper Mapper
	Suffix string
}

// Obj2Table Obj2Table
func (mapper SuffixMapper) Obj2Table(name string) string {
	return mapper.Mapper.Obj2Table(name) + mapper.Suffix
}

// Table2Obj Table2Obj
func (mapper SuffixMapper) Table2Obj(name string) string {
	return mapper.Mapper.Table2Obj(name[:len(name)-len(mapper.Suffix)])
}

// NewSuffixMapper NewSuffixMapper
func NewSuffixMapper(mapper Mapper, suffix string) SuffixMapper {
	return SuffixMapper{mapper, suffix}
}
