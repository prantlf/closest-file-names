package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type pair struct {
	name string
	norm string
}

type item struct {
	name string
	norm string
	file fs.FileInfo
}

type comb struct {
	first  item
	second item
	dist   int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("missing directory")
		return
	}
	dir := os.Args[1]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d files\n", len(files))

	if len(files) < 3 {
		if len(files) == 2 {
			fmt.Printf("%s %s\n", files[0].Name(), files[1].Name())
		} else if len(files) == 1 {
			fmt.Printf("%s\n", files[0].Name())
		}
		return
	}

	combs := make([]comb, len(files)*(len(files)-1)/2)
	fmt.Printf("%d combinations\n", len(combs))
	k := 0
	for i, file1 := range files {
		firstf := prepare(file1, true)
		firstd := prepare(file1, false)
		for j := i + 1; j < len(files); j++ {
			file2 := files[j]
			keepext := !file1.IsDir() || file1.IsDir() != file2.IsDir()
			second := prepare(files[j], keepext)
			var first pair
			if keepext {
				first = firstf
			} else {
				first = firstd
			}
			combs[k] = comb{
				item{first.name, first.norm, file1},
				item{second.name, second.norm, file2},
				distance(first.norm, second.norm)}
			k++
		}
	}

	sort.Slice(combs, func(i, j int) bool {
		return combs[i].dist < combs[j].dist
	})

	fmt.Println()
	for i, comb := range combs {
		fmt.Printf("%s (%s)\n%s (%s)\n%d\n%s\n%s\n",
			comb.first.name, format(comb.first.file.Size()),
			comb.second.name, format(comb.second.file.Size()),
			comb.dist, comb.first.norm, comb.second.norm)
		if i < len(combs)-1 {
			fmt.Println()
		}
	}
}

func prepare(file fs.FileInfo, keepext bool) pair {
	name := file.Name()
	norm := name
	if !keepext {
		norm = cutext(norm)
	}
	norm = simplify(norm)
	return pair{name, norm}
}

func cutext(name string) string {
	if name[0] == '.' {
		return name
	}
	return name[:len(name)-len(filepath.Ext(name))]
}

var separators = regexp.MustCompile(`[-_.&,()\[\]{}]+`)
var sizes = regexp.MustCompile(`(\d+x\d+px\b)|(\d+px\b)|(\bx\d+)|(\d+p\b)|(\b720\b)|(\b1080\b)`)
var particles = regexp.MustCompile(`(\bthe\b)|(\ba\b)|(\bto\b)|(\bfrom\b)|(\bby\b)|(\bis\b)|(\bon\b)|(\bat\b)|(\bin\b)|(\bx\b)|(\band\b)|(\bfor\b)|(\bwith\b)`)
var whitespace = regexp.MustCompile(`\s+`)

func simplify(norm string) string {
	norm = strings.ToLower(norm)
	norm = separators.ReplaceAllString(norm, " ")
	norm = sizes.ReplaceAllString(norm, " ")
	norm = particles.ReplaceAllString(norm, " ")
	norm = whitespace.ReplaceAllString(norm, " ")
	norm = strings.TrimSpace(norm)
	parts := strings.Split(norm, " ")
	sort.Strings(parts)
	return strings.Join(parts, " ")
}

func distance(a, b string) int {
	if a == b {
		return 0
	}

	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func format(n int64) string {
	in := strconv.FormatInt(n, 10)
	numOfDigits := len(in)
	if n < 0 {
		numOfDigits-- // First character is the - sign (not a digit)
	}
	numOfCommas := (numOfDigits - 1) / 3

	out := make([]byte, len(in)+numOfCommas)
	if n < 0 {
		in, out[0] = in[1:], '-'
	}

	for i, j, k := len(in)-1, len(out)-1, 0; ; i, j = i-1, j-1 {
		out[j] = in[i]
		if i == 0 {
			return string(out)
		}
		if k++; k == 3 {
			j, k = j-1, 0
			out[j] = ','
		}
	}
}
