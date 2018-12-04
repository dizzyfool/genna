package model

import (
	"strings"

	"github.com/fatih/camelcase"
)

// Singular makes singular of plural english word
func Singular(input string) string {
	if !IsCountable(input) {
		return input
	}

	var singularDictionary = map[string]string{
		"are":      "is",
		"analyses": "analysis",
		"alumni":   "alumnus",
		"aliases":  "alias",
		"axes":     "axis",
		//"alumni": "alumnae", // for female - cannot have duplicate in map

		"genii":       "genius",
		"data":        "datum",
		"atlases":     "atlas",
		"appendices":  "appendix",
		"barracks":    "barrack",
		"beefs":       "beef",
		"buses":       "bus",
		"brothers":    "brother",
		"cafes":       "cafe",
		"corpuses":    "corpus",
		"campuses":    "campus",
		"cows":        "cow",
		"crises":      "crisis",
		"ganglions":   "ganglion",
		"genera":      "genus",
		"graffiti":    "graffito",
		"loaves":      "loaf",
		"matrices":    "matrix",
		"monies":      "money",
		"mongooses":   "mongoose",
		"moves":       "move",
		"movies":      "movie",
		"mythoi":      "mythos",
		"lice":        "louse",
		"niches":      "niche",
		"numina":      "numen",
		"octopuses":   "octopus",
		"opuses":      "opus",
		"oxen":        "ox",
		"penises":     "penis",
		"vaginas":     "vagina",
		"vertices":    "vertex",
		"viruses":     "virus",
		"shoes":       "shoe",
		"sexes":       "sex",
		"testes":      "testis",
		"turfs":       "turf",
		"teeth":       "tooth",
		"feet":        "foot",
		"cacti":       "cactus",
		"children":    "child",
		"criteria":    "criterion",
		"news":        "news",
		"deer":        "deer",
		"echoes":      "echo",
		"elves":       "elf",
		"embargoes":   "embargo",
		"foes":        "foe",
		"foci":        "focus",
		"fungi":       "fungus",
		"geese":       "goose",
		"heroes":      "hero",
		"hooves":      "hoof",
		"indices":     "index",
		"knifes":      "knife",
		"leaves":      "leaf",
		"lives":       "life",
		"men":         "man",
		"mice":        "mouse",
		"nuclei":      "nucleus",
		"people":      "person",
		"phenomena":   "phenomenon",
		"potatoes":    "potato",
		"selves":      "self",
		"syllabi":     "syllabus",
		"tomatoes":    "tomato",
		"torpedoes":   "torpedo",
		"vetoes":      "veto",
		"women":       "woman",
		"zeroes":      "zero",
		"natives":     "native",
		"hives":       "hive",
		"quizzes":     "quiz",
		"bases":       "basis",
		"diagnostic":  "diagnosis",
		"parentheses": "parenthesis",
		"prognoses":   "prognosis",
		"synopses":    "synopsis",
		"theses":      "thesis",
	}

	if result, ok := singularDictionary[strings.ToLower(input)]; !ok {
		// to handle words like apples, doors, cats
		if len(input) > 2 && input[len(input)-1] == 's' {
			return string(input[:len(input)-1])
		}
		return input
	} else {
		return result
	}
}

// IsCountable check if word can have plural form
func IsCountable(input string) bool {
	// dictionary of word that has no plural version
	var nonCountable = map[string]bool{
		"audio":        true,
		"bison":        true,
		"chassis":      true,
		"compensation": true,
		"coreopsis":    true,
		"data":         true,
		"deer":         true,
		"education":    true,
		"emoji":        true,
		"equipment":    true,
		"fish":         true,
		"furniture":    true,
		"gold":         true,
		"information":  true,
		"knowledge":    true,
		"love":         true,
		"rain":         true,
		"money":        true,
		"moose":        true,
		"nutrition":    true,
		"offspring":    true,
		"plankton":     true,
		"pokemon":      true,
		"police":       true,
		"rice":         true,
		"series":       true,
		"sheep":        true,
		"species":      true,
		"swine":        true,
		"traffic":      true,
		"wheat":        true,
	}

	_, ok := nonCountable[strings.ToLower(input)]
	return !ok
}

func IsUpper(c byte) bool {
	return c >= 'A' && c <= 'Z'
}

func IsLower(c byte) bool {
	return c >= 'a' && c <= 'z'
}

func ToUpper(c byte) byte {
	return c - 32
}

func ToLower(c byte) byte {
	return c + 32
}

// CamelCased converts string to camelCase
// from github.com/go-pg/pg/internal
func CamelCased(s string) string {
	r := make([]byte, 0, len(s))
	upperNext := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			if IsLower(c) {
				c = ToUpper(c)
			}
			upperNext = false
		}
		r = append(r, c)
	}
	return string(r)
}

func ModelName(input string) string {
	splitted := camelcase.Split(CamelCased(input))

	for i, split := range splitted {
		singular := Singular(split)
		if strings.ToLower(singular) != strings.ToLower(split) {
			splitted[i] = strings.Title(singular)
			break
		}
	}

	return strings.Join(splitted, "")
}

func ColumnName(input string) string {
	camelCased := CamelCased(input)

	if strings.HasSuffix(camelCased, "Id") {
		camelCased = camelCased[:len(camelCased)-1] + "D"
	}

	return strings.Title(camelCased)
}

func HasUpper(input string) bool {
	for i := 0; i < len(input); i++ {
		c := input[i]
		if IsUpper(c) {
			return true
		}
	}
	return false
}
