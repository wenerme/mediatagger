package mediatagger

import (
	"regexp"
	"strconv"
)

type Episode struct {
	Season       int
	Episode      int
	EpisodeStart int
	EpisodeEnd   int
}

var (
	seasonStd      = regexp.MustCompile(`(?i)(\b|_)s(\d+)(e(\d+))(\b|_)`)
	seasonCn       = regexp.MustCompile(`第([0-9零〇一二三四五六七八九十百千]+)[季]`)
	episodeCn      = regexp.MustCompile(`第([0-9零〇一二三四五六七八九十百千]+)[集话卷回]`)
	episodeSimple  = regexp.MustCompile(`\[(\d+)]`)
	episodeSimple2 = regexp.MustCompile(`[.](\d+)[.]`)
)

func ExtractEpisode(in string) (out string, episode Episode) {
	var submatch []string
	out, submatch = extractStringSubmatch(seasonStd, in)
	if len(submatch) > 0 {
		s, _ := strconv.ParseInt(submatch[2], 10, 64)
		e, _ := strconv.ParseInt(submatch[4], 10, 64)
		episode.Season = int(s)
		episode.Episode = int(e)
		return
	}
	out = in
	out, submatch = extractStringSubmatch(seasonCn, out)
	if len(submatch) > 0 {
		if s, ok := parseChineseNumber(submatch[1]); ok {
			episode.Season = s
		}
	}

	out, submatch = extractStringSubmatch(episodeCn, out)
	if len(submatch) > 0 {
		if s, ok := parseChineseNumber(submatch[1]); ok {
			episode.Episode = s
		}
	}
	if episode.Episode == 0 {
		out, submatch = extractStringSubmatch(episodeSimple, out)
		if len(submatch) > 0 {
			if s, ok := parseChineseNumber(submatch[1]); ok {
				episode.Episode = s
			}
		}
	}
	if episode.Episode == 0 {
		out, submatch = extractStringSubmatch(episodeSimple2, out)
		if len(submatch) > 0 {
			if s, ok := parseChineseNumber(submatch[1]); ok {
				episode.Episode = s
			}
		}
	}
	return
}

var cnUnit = map[rune]int{
	'十': 10,
	'拾': 10,
	'百': 100,
	'佰': 100,
	'千': 1000,
	'仟': 1000,
	'万': 10000,
	'萬': 10000,
	'亿': 100000000,
	'億': 100000000,
}

var cn = map[rune]int{
	'〇': 0,
	'零': 0,
	'一': 1,
	'壹': 1,
	'二': 2,
	'贰': 2,
	'三': 3,
	'叁': 3,
	'四': 4,
	'肆': 4,
	'五': 5,
	'伍': 5,
	'六': 6,
	'陆': 6,
	'七': 7,
	'柒': 7,
	'八': 8,
	'捌': 8,
	'玖': 9,
	'九': 9,
}

func parseChineseNumber(in string) (out int, ok bool) {
	if i, err := strconv.ParseInt(in, 10, 64); err == nil {
		return int(i), true
	}
	n := []rune(in)
	for i := 0; i < len(n); i++ {
		if v, ok := cn[n[i]]; ok {
			out += v
		} else if v, ok := cnUnit[n[i]]; ok {
			out *= v
		} else {
			return 0, false
		}
	}
	ok = true
	return
}
