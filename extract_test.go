package mediatagger

import (
	"reflect"
	"sort"
	"testing"
)

func TestExtract(t *testing.T) {
	for _, test := range []struct {
		in string
		e  ExtractInfo
	}{
		{in: "test-s03e3.[1080p]", e: ExtractInfo{Rest: "test-     .[     ]", Tags: []string{"1080P"}, Episode: Episode{Season: 3, Episode: 3}}},
		{
			in: "[天空树字幕组][ONE PIECE 海贼王].第一季.[971][X264][720P][GB_JP][MP4][CRRIP][中日双语字幕]",
			e: ExtractInfo{
				Rest:    "[                  ][ONE PIECE 海贼王].         .     [    ][    ][     ][   ][CRRIP][                  ]",
				Tags:    []string{"720P", "天空树字幕组", "H264", "MP4", "日语字幕", "简体中文字幕"},
				Episode: Episode{Season: 1, Episode: 971},
			},
		},
	} {
		out := Extract(test.in)
		sort.Strings(test.e.Tags)
		sort.Strings(out.Tags)
		if !reflect.DeepEqual(test.e, out) {
			t.Error("Expected", test.e, "got", out)
		}

	}
}
