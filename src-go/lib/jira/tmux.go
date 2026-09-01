package jira

import (
	"os/exec"
	"strings"
	"unicode"
)

const maxTmuxIDLen = 38

func sanitizeTmuxRune(r rune) rune {
	if r <= unicode.MaxASCII {
		if r == '-' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' {
			return r
		}
		return '-'
	}
	r = unicode.ToLower(r)
	switch r {
	case 'à', 'á', 'â', 'ã', 'ä', 'å', 'ā', 'ă', 'ą', 'ǎ', 'ǟ', 'ǡ', 'ǻ', 'ȁ', 'ȃ', 'ȧ', 'ḁ', 'ạ', 'ả', 'ấ', 'ầ', 'ẩ', 'ẫ', 'ậ', 'ắ', 'ằ', 'ẳ', 'ẵ', 'ặ':
		r = 'a'
	case 'ḃ', 'ḅ', 'ḇ':
		r = 'b'
	case 'ç', 'ć', 'ĉ', 'ċ', 'č':
		r = 'c'
	case 'ď', 'ḋ', 'ḍ', 'ḏ', 'ḑ', 'ḓ':
		r = 'd'
	case 'è', 'é', 'ê', 'ë', 'ē', 'ĕ', 'ė', 'ę', 'ě', 'ȅ', 'ȇ', 'ȩ', 'ḕ', 'ḗ', 'ḙ', 'ḛ', 'ḝ', 'ẹ', 'ẻ', 'ẽ', 'ế', 'ề', 'ể', 'ễ', 'ệ':
		r = 'e'
	case 'ḟ':
		r = 'f'
	case 'ĝ', 'ğ', 'ġ', 'ģ', 'ǧ', 'ǵ', 'ḡ':
		r = 'g'
	case 'ĥ', 'ȟ', 'ḣ', 'ḥ', 'ḧ', 'ḩ', 'ḫ':
		r = 'h'
	case 'ì', 'í', 'î', 'ï', 'ĩ', 'ī', 'ĭ', 'į', 'ǐ', 'ȉ', 'ȋ', 'ḭ', 'ḯ', 'ỉ', 'ị':
		r = 'i'
	case 'ĵ', 'ǰ':
		r = 'j'
	case 'ķ', 'ǩ', 'ḱ', 'ḳ', 'ḵ':
		r = 'k'
	case 'ĺ', 'ļ', 'ľ', 'ŀ', 'ḷ', 'ḹ', 'ḻ', 'ḽ':
		r = 'l'
	case 'ḿ', 'ṁ', 'ṃ':
		r = 'm'
	case 'ñ', 'ń', 'ņ', 'ň', 'ǹ', 'ṅ', 'ṇ', 'ṉ', 'ṋ':
		r = 'n'
	case 'ò', 'ó', 'ô', 'õ', 'ö', 'ø', 'ō', 'ŏ', 'ő', 'ǒ', 'ǫ', 'ǭ', 'ǿ', 'ȍ', 'ȏ', 'ȫ', 'ȭ', 'ȯ', 'ȱ', 'ṍ', 'ṏ', 'ṑ', 'ṓ', 'ọ', 'ỏ', 'ố', 'ồ', 'ổ', 'ỗ', 'ộ', 'ớ', 'ờ', 'ở', 'ỡ', 'ợ':
		r = 'o'
	case 'ṕ', 'ṗ':
		r = 'p'
	case 'ŕ', 'ŗ', 'ř', 'ȑ', 'ȓ', 'ṙ', 'ṛ', 'ṝ', 'ṟ':
		r = 'r'
	case 'ś', 'ŝ', 'ş', 'š', 'ș', 'ṡ', 'ṣ', 'ṥ', 'ṧ', 'ṩ':
		r = 's'
	case 'ţ', 'ť', 'ț', 'ṫ', 'ṭ', 'ṯ', 'ṱ':
		r = 't'
	case 'ù', 'ú', 'û', 'ü', 'ũ', 'ū', 'ŭ', 'ů', 'ű', 'ų', 'ǔ', 'ǖ', 'ǘ', 'ǚ', 'ǜ', 'ȕ', 'ȗ', 'ṳ', 'ṵ', 'ṷ', 'ṹ', 'ṻ', 'ụ', 'ủ', 'ứ', 'ừ', 'ử', 'ữ', 'ự':
		r = 'u'
	case 'ṽ', 'ṿ':
		r = 'v'
	case 'ŵ', 'ẁ', 'ẃ', 'ẅ', 'ẇ', 'ẉ':
		r = 'w'
	case 'ẋ', 'ẍ':
		r = 'x'
	case 'ý', 'ÿ', 'ŷ', 'ȳ', 'ẏ', 'ỳ', 'ỵ', 'ỷ', 'ỹ':
		r = 'y'
	case 'ź', 'ż', 'ž', 'ẑ', 'ẓ', 'ẕ':
		r = 'z'
	default:
		return '-'
	}
	return r
}

func (t Ticket) TmuxSessionID() string {
	var result [maxTmuxIDLen]byte
	n := 0
	prevDash := false
	writeRune := func(r rune) bool {
		r = sanitizeTmuxRune(r)
		if r == '-' {
			if prevDash {
				return false
			}
			prevDash = true
		} else {
			prevDash = false
		}
		result[n] = byte(r)
		n++
		return n == len(result)
	}
	write := func(s string, lower bool) bool {
		for _, r := range s {
			if lower && r <= unicode.MaxASCII {
				r = unicode.ToLower(r)
			}
			if writeRune(r) {
				return true
			}
		}
		return false
	}

	if !write(t.Current, false) && !writeRune('-') {
		label := t.Label
		for {
			start := strings.IndexByte(label, '[')
			if start < 0 {
				write(label, true)
				break
			}
			end := strings.IndexByte(label[start:], ']')
			if end < 0 {
				write(label, true)
				break
			}
			if write(label[:start], true) {
				break
			}
			label = label[start+end+1:]
		}
	}
	for n > 0 && result[n-1] == '-' {
		n--
	}
	return string(result[:n])
}

func TmuxSessionPrefixes(currentIDs map[string]bool) []string {
	prefixes := make([]string, 0, len(currentIDs))
	for id := range currentIDs {
		prefixes = append(prefixes, SanitizeTmuxID(id)+"-")
	}
	return prefixes
}

func MatchesTmuxSession(session string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(session, p) {
			return true
		}
	}
	return false
}

// SanitizeTmuxID transliterates Latin diacritics, replaces other disallowed
// chars with '-', and collapses runs of dashes into one, in a single pass.
func SanitizeTmuxID(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		r = sanitizeTmuxRune(r)
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
		} else {
			prevDash = false
		}
		b.WriteByte(byte(r)) // ponytail: all valid output chars are ASCII [a-zA-Z0-9-]
	}
	return b.String()
}

// NotifCounts is the per-session tally of agent notification states:
// panes finished (value 1, red badge), still working (value 2, blue badge),
// or manually flagged (value 3, yellow badge).
type NotifCounts struct {
	Working  int
	Finished int
	Manual   int
}

func LoadNotificationCounts() map[string]NotifCounts {
	out, err := exec.Command("tmux", "list-panes", "-a", "-F",
		"#{session_name}\t#{@zsb_agent_notif}").Output()
	if err != nil {
		return map[string]NotifCounts{}
	}
	return ParseNotificationCounts(string(out))
}

// ParseNotificationCounts tallies panes per session by state: value "1" → Finished,
// "2" → Working. Everything else (unset/empty/"0") is ignored.
func ParseNotificationCounts(out string) map[string]NotifCounts {
	counts := map[string]NotifCounts{}
	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		c := counts[parts[0]]
		switch strings.TrimSpace(parts[1]) {
		case "1":
			c.Finished++
		case "2":
			c.Working++
		case "3":
			c.Manual++
		default:
			continue
		}
		counts[parts[0]] = c
	}
	return counts
}
