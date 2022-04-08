package mediatagger

type ExtractInfo struct {
	Rest    string
	Tags    []string
	Episode Episode
}

func Extract(in string) (out ExtractInfo) {
	out.Rest, out.Episode = ExtractEpisode(in)
	out.Rest, out.Tags = ExtractTags(out.Rest)
	return
}
