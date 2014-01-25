package bwt

import (
    "fmt"
    "sort"
    "strings"
)

// Bstring holds one ciclic permutation of the input strig to BWT
type Bstring struct {
    Pos int
    Ch  string
}

// Bstrings holds a list of Bstring
type Bstrings []Bstring

// FIXME: This sucks big time. Was the quickest thing to get sorting work though.
var strX string

// LessThan compares two Bstring's and returns which should have a "lower" position in the BWT
func (bs Bstring) LessThan(other Bstring) bool {
    a, b, n := bs.Pos, other.Pos, len(strX)

    lim := n - 1
    if a > b {
        lim -= a
    } else {
        lim -= b
    }

    if strX[a:a+lim] != strX[b:b+lim] {
        return strX[a:a+lim] < strX[b:b+lim]
    }

    a += lim
    b += lim

    for {
        if strX[a] != strX[b] {
            return strX[a] < strX[b]
        }

        a, b = (a+1)%n, (b+1)%n
    }

    return false
}

func (bs Bstrings) Swap(i, j int)      { bs[i], bs[j] = bs[j], bs[i] }
func (bs Bstrings) Len() int           { return len(bs) }
func (bs Bstrings) Less(i, j int) bool { return bs[i].LessThan(bs[j]) }

// BurrowsWheelerTransform implments the operation with the same name upon str
func BurrowsWheelerTransform(str string) (out string) {
    strX = str
    n := len(str)
    sx := make(Bstrings, n)
    for k := range str {
        end := (k + n - 1) % n
        sx[k] = Bstring{k, string(str[end])}
    }

    sort.Sort(sx)

    for _, bstr := range sx {
        out += bstr.Ch
    }

    return
}

// BurrowsWheelerTransformNaive implements the operation with the same name upon str,
// using a naive (and thus suboptimal) approach.
func BurrowsWheelerTransformNaive(str string) (out string) {
    strX = str
    sx := make([]string, len(str))
    for k := range str {
        curr := str[k:] + str[:k]
        sx[k] = curr
    }

    sort.Strings(sx)

    out = ""
    for _, line := range sx {
        out += fmt.Sprintf("%s", string(line[len(line)-1]))
    }

    return
}

// Inverse BWT

// Pos seems to be pretty much the same as Bstring.
// FIXME: reconcile them BUT only after adding tests.
type Pos struct {
    Ch  uint8
    Pos int
}

func sortedString(str string) string {
    str1x := strings.Split(str, "")
    sort.Strings(str1x)
    return strings.Join(str1x, "")
}

func countOccurences(str string) (occ []int) {
    occ = make([]int, len(str))
    counts := map[rune]int{}

    for k, v := range str {
        counts[v]++
        occ[k] = counts[v]
    }

    return
}

// InvertedBurrowsWheelerTransform implements the operation with the same name upon str
func InvertedBurrowsWheelerTransform(str string) string {
    str1 := sortedString(str)
    left, right, rightToLeft := countOccurences(str1), countOccurences(str), map[Pos]Pos{}
    for k := range str {
        posR, posL := Pos{str[k], right[k]}, Pos{str1[k], left[k]}
        rightToLeft[posR] = posL
    }

    out, prev := "", Pos{'$', 1}
    for {
        curr := rightToLeft[prev]
        out += string(curr.Ch)
        prev = curr
        if curr.Ch == '$' {
            break
        }
    }

    return out
}
