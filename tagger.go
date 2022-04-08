package mediatagger

import (
	"regexp"
	"strings"
)

type GroupRule struct {
	Name     string
	Continue bool
	Keep     bool
	Raw      bool
	Rules    []Rule
}

func (g *GroupRule) Match(m *Match) (out string, tags []string) {
	last := m.Input
	raw := m.Raw
	for _, r := range g.Rules {
		if g.Raw {
			m.Input = raw
		}
		if m.MatchFunc(r.Match) {
			// only modify when matched
			if g.Raw {
				m.Input = last
				last, _ = r.Match(m)
			}
			if !g.Continue {
				break
			}
		}
		if g.Keep {
			m.Input = last
		}
	}
	if g.Raw {
		m.Input = last
	}
	return
}

type Rule interface {
	Match(m *Match) (string, []string)
}

type RegExpRule struct {
	RegExp  *regexp.Regexp
	Tag     string
	TagFunc func([]string) []string
}

func (r *RegExpRule) Match(m *Match) (out string, tags []string) {
	if r.TagFunc != nil {
		out, matches := extractAllString(r.RegExp, m.Input)
		if len(matches) > 0 {
			for _, m := range matches {
				tags = append(tags, r.TagFunc(m)...)
			}
		}
		return out, tags
	}

	out, match := extract(r.RegExp, m.Input)
	if match {
		tags = append(tags, r.Tag)
	}

	return
}

var resolution = []*RegExpRule{
	{Tag: "4K", RegExp: regexp.MustCompile(`(?i)2160p|4k`)},
	{Tag: "1080P", RegExp: regexp.MustCompile(`(?i)1080p|1920.1080`)},
	{Tag: "1080i", RegExp: regexp.MustCompile(`(?i)1080i`)},
	{Tag: "720P", RegExp: regexp.MustCompile(`(?i)720p|1280.720`)},
	{Tag: "576P", RegExp: regexp.MustCompile(`(?i)576p|\d+.576`)},
}

var video = []*RegExpRule{
	{Tag: "H264", RegExp: regexp.MustCompile(`(?i)(\b|_)avc(\b|_)|[xh].?264`)},
	{Tag: "H265", RegExp: regexp.MustCompile(`(?i)(\b|_)hevc(\b|_)|[xh].?265`)},
}

var audio = []*RegExpRule{
	{Tag: "AAC", RegExp: regexp.MustCompile(`(?i)(\b|_)aac(\b|_)`)},
	{Tag: "DDP2.0", RegExp: regexp.MustCompile(`(?i)(\b|_)ddp2(.\d)?(\b|_)`)},
	{Tag: "DDP5.1", RegExp: regexp.MustCompile(`(?i)(\b|_)ddp5(.\d)?(\b|_)`)},
	{Tag: "FLAC", RegExp: regexp.MustCompile(`(?i)(\b|_)flac(\b|_)`)},
}

var quality = []*RegExpRule{
	{Tag: "Web", RegExp: regexp.MustCompile(`(?i)(\b|_)(webrip|web.?dl|web)(\b|_)`)},
	{Tag: "HDTV", RegExp: regexp.MustCompile(`(?i)(\b|_)(hdtvrip|hdtv)(\b|_)`)},
	{Tag: "TVRip", RegExp: regexp.MustCompile(`(?i)(\b|_)(tvrip)(\b|_)`)},
	{Tag: "BluRay", RegExp: regexp.MustCompile(`(?i)(\b|_)(bdrip|blu-?ray)(\b|_)`)},
}

var source = []*RegExpRule{
	{Tag: "Amazon", RegExp: regexp.MustCompile(`(?i)(\b|_)(amzn)(\b|_)`)},
	{Tag: "Netflix", RegExp: regexp.MustCompile(`(?i)(\b|_)(nf)(\b|_)`)},
	{Tag: "HBO Max", RegExp: regexp.MustCompile(`(?i)(\b|_)(hmax)(\b|_)`)},
}

var sub = []*RegExpRule{
	{Tag: "简体中文字幕", RegExp: regexp.MustCompile(`(?i)(\b|_)(chi|chs|gb)(\b|_)|[中简]\S*?((双语|外挂)(字幕)?|(双语|外挂)?(字幕))|简体|简繁`)},
	{Tag: "繁体中文字幕", RegExp: regexp.MustCompile(`(?i)(\b|_)(big5|cht)(\b|_)|[繁]\S*?(字幕|双语|外挂)|繁体|繁體|简繁`)},
	{Tag: "英语字幕", RegExp: regexp.MustCompile(`(?i)(\b|_)(eng)(\b|_)|[英]\S*?(字幕|双语|外挂)`)},
	{Tag: "日语字幕", RegExp: regexp.MustCompile(`(?i)(\b|_)(jap|jp)(\b|_)|[日]\S*?(字幕|双语|外挂)`)},
}

var info = []*RegExpRule{
	{Tag: "rarbg", RegExp: regexp.MustCompile(`(?i)(\b|_)(rarbg|rartv)(\b|_)`)},
	{Tag: "DUBBED", RegExp: regexp.MustCompile(`(?i)(\b|_)(dubbed)(\b|_)`)},
	{Tag: "EXTENDED", RegExp: regexp.MustCompile(`(?i)(\b|_)(extended)(\b|_)`)},
	{
		RegExp: regexp.MustCompile(`(?i)(\b|_)(mp4|mkv)(\b|_)`),
		TagFunc: func(in []string) []string {
			return []string{strings.ToUpper(in[0])}
		},
	},
}

func rulesOf[T Rule](a []T) (out []Rule) {
	out = make([]Rule, 0, len(a))
	for _, v := range a {
		out = append(out, Rule(v))
	}
	return
}

var rules = []*GroupRule{
	{Rules: rulesOf(resolution)},
	{Rules: rulesOf(quality)},
	{Rules: rulesOf(video)},
	{Rules: rulesOf(audio)},
	{Rules: rulesOf(sub), Continue: true, Raw: true},
	{Rules: rulesOf(info), Continue: true},
	{Rules: rulesOf(source)},
	{Rules: groups},
}

var groups = []Rule{
	&RegExpRule{
		RegExp: regexp.MustCompile(`(?i)[a-z0-9\x{4E00}-\x{9FA5}]+((字幕|搬运)[组社团]|工作室)`),
		TagFunc: func(in []string) []string {
			return []string{in[0]}
		},
	},
	&RegExpRule{
		RegExp: regexp.MustCompile(`(?i)(\b|_)[a-z0-9\x{4E00}-\x{9FA5}]+-?(Sub|Raws)(\b|_)`),
		TagFunc: func(in []string) []string {
			return []string{strings.ReplaceAll(in[0], "_", "")}
		},
	},
	&RegExpRule{RegExp: regexp.MustCompile(`(?i)F.?I.?X\s*字幕侠`), Tag: "F.I.X字幕侠"},
	&RegExpRule{RegExp: regexp.MustCompile(`(?i)人人影视|yyets`), Tag: "YYeTs字幕组"},
}

type Match struct {
	Input string
	Raw   string
	Tags  []string
}

func (m *Match) MatchFunc(f func(match *Match) (out string, tags []string)) bool {
	out, tags := f(m)
	if len(tags) > 0 {
		m.Input = out
		m.Tags = append(m.Tags, tags...)
		return true
	}
	return false
}

func ExtractTags(in string) (out string, tags []string) {
	m := &Match{
		Input: in,
		Raw:   in,
	}
	for _, r := range rules {
		m.MatchFunc(r.Match)
	}
	return m.Input, m.Tags
}
