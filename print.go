package redblack

import (
	"bytes"
	"fmt"
	"strings"
)

type RuneSet struct {
	v, h               rune
	ctl, ctr, cbl, cbr rune
	bt, bb, bl, br     rune
}

var (
	regularRuneSet = RuneSet{
		'│', '─',
		'┌', '┐', '└', '┘',
		'┴', '┬', '┤', '├',
	}
	boldRuneSet = RuneSet{
		'║', '═',
		'╔', '╗', '╚', '╝',
		'╧', '╤', '╢', '╟',
	}
)

func boxed(s string, bold, hasParent, hasLeft, hasRight bool) [3]string {
	strLen := len(s)
	if strLen%2 != 1 {
		s = s + " "
		strLen++
	}

	rs := &regularRuneSet
	if bold {
		rs = &boldRuneSet
	}
	ret := [3]string{}

	// Build top string
	topStringBuilder := strings.Builder{}
	topStringBuilder.WriteRune(rs.ctl)
	for i := 0; i < (strLen+2)/2; i++ {
		topStringBuilder.WriteRune(rs.h)
	}
	if hasParent {
		topStringBuilder.WriteRune(rs.bt)
	} else {
		topStringBuilder.WriteRune(rs.h)
	}
	for i := 0; i < (strLen+2)/2; i++ {
		topStringBuilder.WriteRune(rs.h)
	}
	topStringBuilder.WriteRune(rs.ctr)
	ret[0] = topStringBuilder.String()

	// Build center string
	cenStringBuilder := strings.Builder{}
	cenStringBuilder.WriteRune(rs.v)
	cenStringBuilder.WriteRune(' ')
	cenStringBuilder.WriteString(s)
	cenStringBuilder.WriteRune(' ')
	cenStringBuilder.WriteRune(rs.v)
	ret[1] = cenStringBuilder.String()

	//	Build bottom string
	botStringBuilder := strings.Builder{}
	botStringBuilder.WriteRune(rs.cbl)
	botStringBuilder.WriteString(bottomString(rs, strLen/2+1, hasLeft))
	botStringBuilder.WriteRune(rs.h)
	botStringBuilder.WriteString(bottomString(rs, strLen/2+1, hasRight))
	botStringBuilder.WriteRune(rs.cbr)
	ret[2] = botStringBuilder.String()

	return ret
}

func bottomString(rs *RuneSet, len int, withSibling bool) string {
	sb := strings.Builder{}
	if withSibling {
		p1 := len / 2
		p2 := len - p1
		for i := 0; i < p1; i++ {
			sb.WriteRune(rs.h)
		}
		sb.WriteRune(rs.bb)
		for i := 0; i < p2-1; i++ {
			sb.WriteRune(rs.h)
		}
	} else {
		for i := 0; i < len; i++ {
			sb.WriteRune(rs.h)
		}
	}
	return sb.String()
}

func buildConnector(len int, left bool) string {
	sb := strings.Builder{}
	if left {
		sb.WriteRune(regularRuneSet.ctl)
	} else {
		sb.WriteRune(regularRuneSet.cbl)
	}
	for i := 0; i < len-2; i++ {
		sb.WriteRune(regularRuneSet.h)
	}
	if left {
		sb.WriteRune(regularRuneSet.cbr)
	} else {
		sb.WriteRune(regularRuneSet.ctr)
	}
	return sb.String()
}

// String returns a string representation of the tree.
func (t Tree[V, T]) String() string {
	level := t.GetTreeLevels()
	h := len(level)
	strLen := 9
	w := (1 << (h - 1)) * (strLen + 3)
	buffer := bytes.NewBuffer(make([]byte, 0, 1000))
	for i, l := range level {
		div := 1 << i
		whitespace := (w/div - strLen) / 2
		format := fmt.Sprintf("%%%ds%%s%%%ds", whitespace, whitespace)
		stringBuilders := [4]strings.Builder{}
		for _, n := range l {
			if n != nil {
				strings := boxed(fmt.Sprint(n.value), n.red, i > 0, n.left != nil, n.right != nil)
				for k, s := range strings {
					stringBuilders[k].WriteString(fmt.Sprintf(format, "", s, ""))
					stringBuilders[k].WriteRune(' ')
				}

				// print connectors
				lenCon := whitespace/2 + strLen/3 - 1
				formatCon := fmt.Sprintf("%%%ds%%%ds", (whitespace)/2+2, (whitespace)/2+2)

				if n.left != nil {
					stringBuilders[3].WriteString(fmt.Sprintf(formatCon, "", buildConnector(lenCon, true)))
				} else {
					stringBuilders[3].WriteString(fmt.Sprintf(formatCon, "", ""))
				}
				for k := 0; k < strLen/3; k++ {
					stringBuilders[3].WriteRune(' ')
				}
				if n.right != nil {
					stringBuilders[3].WriteString(fmt.Sprintf(formatCon, buildConnector(lenCon, false), ""))
				} else {
					stringBuilders[3].WriteString(fmt.Sprintf(formatCon, "", ""))
				}
				stringBuilders[3].WriteRune(' ')
			} else {
				for k := range stringBuilders {
					stringBuilders[k].WriteString(fmt.Sprintf(format, "", "          ", ""))
				}
			}
		}
		for _, sb := range stringBuilders {
			buffer.WriteString(sb.String() + "\n")
			sb.Reset()
		}
	}
	return buffer.String()
}
