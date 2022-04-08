package mediatagger

import "regexp"

func extractStringSubmatch(re *regexp.Regexp, in string) (out string, matches []string) {
	idx := re.FindStringSubmatchIndex(in)
	if len(idx) == 0 {
		return in, nil
	}
	rn := []byte(in)
	n := len(idx)
	for i := 0; i < n; i += 2 {
		if i < 0 {
			continue
		}
		a := idx[i]
		b := idx[i+1]
		matches = append(matches, in[a:b])
	}

	for a := idx[0]; a < idx[1]; a++ {
		rn[a] = ' '
	}

	out = string(rn)
	return
}

func extractAllString(re *regexp.Regexp, in string) (out string, matches [][]string) {
	all := re.FindAllStringIndex(in, -1)
	if len(all) == 0 {
		return in, nil
	}
	rn := []byte(in)
	for _, idx := range all {
		n := len(idx)
		var m []string
		for i := 0; i < n; i += 2 {
			if i < 0 {
				continue
			}
			a := idx[i]
			b := idx[i+1]
			m = append(m, in[a:b])
		}

		matches = append(matches, m)
		for a := idx[0]; a < idx[1]; a++ {
			rn[a] = ' '
		}
	}
	out = string(rn)
	return
}

func extract(re *regexp.Regexp, in string) (out string, match bool) {
	index := re.FindAllStringIndex(in, -1)
	if len(index) == 0 {
		return in, false
	}
	rn := []byte(in)
	for _, i := range index {
		for a := i[0]; a < i[1]; a++ {
			rn[a] = ' '
		}
	}
	out = string(rn)
	match = true
	return
}
