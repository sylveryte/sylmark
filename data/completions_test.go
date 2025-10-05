package data

import (
	"fmt"
	"testing"
)

func TestAnalyzeCompletionTrigger(t *testing.T) {
	linenocompl := "Cool new stuff n"
	t.Run("1 No completions at 0", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(1, linenocompl)
		assert(t, CompletionNone, kind, "", arg, 0, 0, cstart, cend)
	})
	t.Run("2 No completions at end of word 4", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(4, linenocompl)
		assert(t, CompletionNone, kind, "", arg, 0, 0, cstart, cend)
	})
	t.Run("3 No completions at end", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(len(linenocompl)-1, linenocompl)
		assert(t, CompletionNone, kind, "", arg, 0, 0, cstart, cend)
	})

	// tags
	t.Run("4 Tag at start only hash", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(1, "#")
		assert(t, CompletionTag, kind, "", arg, 0, 1, cstart, cend)
	})
	t.Run("5 Tag at start some text ", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(3, "#sup")
		assert(t, CompletionTag, kind, "su", arg, 0, 3, cstart, cend)
	})
	t.Run("6 Tag at [[ d some text ", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(16, "cool new wo #sup")
		assert(t, CompletionTag, kind, "sup", arg, 12, 16, cstart, cend)
	})
	t.Run("7 # Tag  mid hash", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(13, "cool new wo # jjo")
		assert(t, CompletionTag, kind, "", arg, 12, 13, cstart, cend)
	})

	// wikilink
	t.Run("8 Wiki at  start simple", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(2, "[[ ")
		assert(t, CompletionWiki, kind, "", arg, 0, 2, cstart, cend)
	})
	t.Run("9 Wiki at  start pre text", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(17, "cool some text [[")
		assert(t, CompletionWiki, kind, "", arg, 15, 17, cstart, cend)
	})
	t.Run("10 Wiki at end pre heading but cursor before trigger", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(15, "# cool some text [[")
		assert(t, CompletionNone, kind, "", arg, 0, 0, cstart, cend)
	})
	t.Run("11 Wiki at end with pre heading", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(19, "# cool some text [[")
		assert(t, CompletionWiki, kind, "", arg, 17, 19, cstart, cend)
	})
	t.Run("12 Wiki at mid with some text", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(16, "cool text [[some ")
		assert(t, CompletionWiki, kind, "some", arg, 10, 16, cstart, cend)
	})
	t.Run("13 Wiki at mid with some text with space", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(17, "cool text [[some ")
		assert(t, CompletionWiki, kind, "some ", arg, 10, 17, cstart, cend)
	})
	t.Run("14 Wiki at mid with some text before", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(15, "cool text [[some ")
		assert(t, CompletionWiki, kind, "some", arg, 10, 15, cstart, cend)
	})
	t.Run("15 Wiki at mid with some text before2", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(13, "cool text [[some  cool")
		assert(t, CompletionWiki, kind, "some", arg, 10, 13, cstart, cend)
	})

	// with subheading
	t.Run("16 Wiki at  start with subheading", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(20, "cool some text [[some# ")
		assert(t, CompletionWiki, kind, "some#", arg, 15, 20, cstart, cend)
	})

	// with endings
	t.Run("17.1 Wiki at  start with ending", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(3, "[[k ]]")
		assert(t, CompletionWikiWithEnd, kind, "k ", arg, 0, 4, cstart, cend)
	})
	t.Run("17.1.1 Wiki at  start with ending zero", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(2, "[[]]")
		assert(t, CompletionWikiWithEnd, kind, "", arg, 0, 2, cstart, cend)
	})
	t.Run("17.2 Wiki at  start with ending", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(3, "[[k]]")
		assert(t, CompletionWikiWithEnd, kind, "k", arg, 0, 3, cstart, cend)
	})
	t.Run("17.3 Wiki at  start with ending", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(20, "# cool some text [[ ]]")
		assert(t, CompletionWikiWithEnd, kind, " ", arg, 17, 20, cstart, cend)
	})
	t.Run("18 Wiki at  start with ending and subheading", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(20, "cool some text [[some# ]]")
		assert(t, CompletionWikiWithEnd, kind, "some# ", arg, 15, 23, cstart, cend)
	})

	t.Run("19 Wiki at  start with ending and subheading", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(20, "cool some text [[some text ]]")
		assert(t, CompletionWikiWithEnd, kind, "some text ", arg, 15, 27, cstart, cend)
	})

	t.Run("20 Wiki end mid with wiki and tag at and at ending", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(32, " [[cool]]some #k text [[ anoterh ]]")
		assert(t, CompletionWikiWithEnd, kind, " anoterh ", arg, 22, 33, cstart, cend)
	})

	// complex ones
	t.Run("21 Wiki start with wiki and tag at and at ending", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(4, " [[cool]]some #k text [[ anoterh ]]")
		assert(t, CompletionWikiWithEnd, kind, "cool", arg, 1, 7, cstart, cend)
	})

	t.Run("22 Wiki start with wiki and tag at and at ending cursor at [[", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(4, " [[cool]]some #k text [[ anoterh ]]")
		assert(t, CompletionWikiWithEnd, kind, "cool", arg, 1, 7, cstart, cend)
	})

	// tag
	t.Run("23 Tag in mid with wiki at and tag ending", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(16, " [[cool]]some #k text [[ anoterh ]]")
		assert(t, CompletionTag, kind, "k", arg, 14, 16, cstart, cend)
	})

	t.Run("24 Wiki end with", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(13, "[[Black wiki#]]")
		assert(t, CompletionWikiWithEnd, kind, "Black wiki#", arg, 0, 13, cstart, cend)
	})

	// markdown link
	t.Run("25 markdown link basic blank", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(8, "some [](")
		assert(t, CompletionInlineLink, kind, "", arg, 5, 8, cstart, cend)
	})

	t.Run("26 markdown link basic blank", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(8, "some []()")
		assert(t, CompletionInlineLinkEnd, kind, "", arg, 5, 9, cstart, cend)
	})

	t.Run("27 markdown link basic blank", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(12, "some [](cool)")
		assert(t, CompletionInlineLinkEnd, kind, "cool", arg, 5, 13, cstart, cend)
	})
	t.Run("28 markdown link basic blank at start", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(4, "[](k)")
		assert(t, CompletionInlineLinkEnd, kind, "k", arg, 0, 5, cstart, cend)
	})
	t.Run("29 markdown link basic blank, non empty", func(t *testing.T) {
		kind, arg, arg2, cstart, cend := analyzeTriggerKind(19, "some [jhsu kd](cool)")
		assertMarkdown(t, CompletionInlineLinkEndHasText, kind, "cool", "jhsu kd", arg, arg2, 5, 20, cstart, cend)
	})
	t.Run("30 markdown link basic blank at start,non empty", func(t *testing.T) {
		kind, arg, arg2, cstart, cend := analyzeTriggerKind(12, "[okay bro](k)")
		assertMarkdown(t, CompletionInlineLinkEndHasText, kind, "k", "okay bro", arg, arg2, 0, 13, cstart, cend)
	})
	t.Run("31 markdown link basic blank at start,non empty", func(t *testing.T) {
		kind, arg, arg2, cstart, cend := analyzeTriggerKind(8, "[nice](k)")
		assertMarkdown(t, CompletionInlineLinkEndHasText, kind, "k", "nice", arg, arg2, 0, 9, cstart, cend)
	})

	// foot note
	t.Run("32 Footnote with end ", func(t *testing.T) {
		kind, arg, _, cstart, cend := analyzeTriggerKind(2, "[k]")
		assert(t, CompletionFootNoteEnd, kind, "k", arg, 0, 3, cstart, cend)
	})
}

func assertMarkdown(t *testing.T, want, got CompletionTriggerKind, argwant, arg2want, arggot, arg2got string, wantcstart, wantcend, cstart, cend int) {
	if arg2got != arg2want {
		t.Errorf(fmt.Sprintf("Arg2 >>> [%s] got [%s]", arg2want, arg2got))
	}
	assert(t, want, got, argwant, arggot, wantcstart, wantcend, cstart, cend)
}

func assert(t *testing.T, want, got CompletionTriggerKind, argwant, arggot string, wantcstart, wantcend, cstart, cend int) {
	if want != got {
		t.Errorf(fmt.Sprintf("Kind >>> [%d] got [%d]", want, got))
	}
	if argwant != arggot {
		t.Errorf(fmt.Sprintf("Arg >>> [%s] got [%s]", argwant, arggot))
	}
	if wantcstart != cstart {
		t.Errorf(fmt.Sprintf("Start >>> [%d] got [%d]", wantcstart, cstart))
	}

	if wantcend != cend {
		t.Errorf(fmt.Sprintf("End >>> [%d] got [%d]", wantcend, cend))
	}
}
