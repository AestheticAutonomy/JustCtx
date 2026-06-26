package scanner

import (
	"math"
	"regexp"
	"strings"

	"github.com/AestheticAutonomy/justctx/pkg/schema"
)

var (
	headingRegex = regexp.MustCompile(`^(#+)\s+(.+)$`)
	wordRegex    = regexp.MustCompile(`[a-zA-Z0-9]+`)
)

func DetectConflicts(chunks []schema.Chunk, sources []schema.Source) []schema.Conflict {
	var conflicts []schema.Conflict

	// 1. Detect duplicate headings
	duplicateHeadings := detectDuplicateHeadings(chunks)
	conflicts = append(conflicts, duplicateHeadings...)

	// 2. Detect near duplicate paragraphs
	nearDuplicates := detectNearDuplicateParagraphs(chunks)
	conflicts = append(conflicts, nearDuplicates...)

	// 3. Detect contradicting imperatives
	contradictions := detectContradictingImperatives(chunks)
	conflicts = append(conflicts, contradictions...)

	return conflicts
}

func detectDuplicateHeadings(chunks []schema.Chunk) []schema.Conflict {
	// Map heading text to list of SourceIDs that contain it
	headingSources := make(map[string][]string)

	for _, chunk := range chunks {
		lines := strings.Split(chunk.Content, "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if match := headingRegex.FindStringSubmatch(trimmed); match != nil {
				// Normalize heading text (lowercase, trimmed)
				heading := strings.ToLower(strings.TrimSpace(match[2]))
				if heading == "" {
					continue
				}

				// Check if this source already listed for this heading
				alreadyAdded := false
				for _, sid := range headingSources[heading] {
					if sid == chunk.SourceID {
						alreadyAdded = true
						break
					}
				}
				if !alreadyAdded {
					headingSources[heading] = append(headingSources[heading], chunk.SourceID)
				}
			}
		}
	}

	var conflicts []schema.Conflict
	for heading, sids := range headingSources {
		if len(sids) > 1 {
			// Find original case heading to display nicely
			conflicts = append(conflicts, schema.Conflict{
				Type:      "duplicate_heading",
				Heading:   heading,
				SourceIDs: sids,
			})
		}
	}

	return conflicts
}

type paragraphInfo struct {
	content  string
	sourceID string
	words    map[string]bool
}

func detectNearDuplicateParagraphs(chunks []schema.Chunk) []schema.Conflict {
	var paragraphs []paragraphInfo

	for _, chunk := range chunks {
		// Split by double newline to get paragraphs
		paras := strings.Split(chunk.Content, "\n\n")
		for _, para := range paras {
			trimmed := strings.TrimSpace(para)
			// Ignore headings and very short text
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}

			words := tokenize(trimmed)
			if len(words) < 8 { // Ignore short lists/phrases
				continue
			}

			paragraphs = append(paragraphs, paragraphInfo{
				content:  trimmed,
				sourceID: chunk.SourceID,
				words:    words,
			})
		}
	}

	var conflicts []schema.Conflict
	seenPairs := make(map[string]bool)

	for i := 0; i < len(paragraphs); i++ {
		for j := i + 1; j < len(paragraphs); j++ {
			p1 := paragraphs[i]
			p2 := paragraphs[j]

			// Only compare different sources
			if p1.sourceID == p2.sourceID {
				continue
			}

			overlap := computeOverlap(p1.words, p2.words)
			if overlap > 0.85 {
				pairKey := p1.sourceID + "-" + p2.sourceID
				if p2.sourceID < p1.sourceID {
					pairKey = p2.sourceID + "-" + p1.sourceID
				}

				if !seenPairs[pairKey] {
					seenPairs[pairKey] = true
					conflicts = append(conflicts, schema.Conflict{
						Type:      "near_duplicate_paragraph",
						SourceIDs: []string{p1.sourceID, p2.sourceID},
					})
				}
			}
		}
	}

	return conflicts
}

func detectContradictingImperatives(chunks []schema.Chunk) []schema.Conflict {
	type lineInfo struct {
		original string
		stripped string
		modal    string // always, never, do, dont, should, shouldnt
		sourceID string
	}

	var lines []lineInfo

	modalPairs := map[string]string{
		"always":   "never",
		"never":    "always",
		"do":       "dont",
		"dont":     "do",
		"should":   "shouldnt",
		"shouldnt": "should",
		"must":     "mustnot",
		"mustnot":  "must",
	}

	for _, chunk := range chunks {
		linesRaw := strings.Split(chunk.Content, "\n")
		for _, line := range linesRaw {
			trimmed := strings.TrimSpace(line)
			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				continue
			}

			lower := strings.ToLower(trimmed)
			modal := ""
			stripped := lower

			if strings.Contains(lower, "always") {
				modal = "always"
				stripped = strings.ReplaceAll(stripped, "always", "")
			} else if strings.Contains(lower, "never") {
				modal = "never"
				stripped = strings.ReplaceAll(stripped, "never", "")
			} else if strings.Contains(lower, "don't") || strings.Contains(lower, "do not") {
				modal = "dont"
				stripped = strings.ReplaceAll(stripped, "don't", "")
				stripped = strings.ReplaceAll(stripped, "do not", "")
			} else if strings.Contains(lower, "should not") || strings.Contains(lower, "shouldn't") {
				modal = "shouldnt"
				stripped = strings.ReplaceAll(stripped, "should not", "")
				stripped = strings.ReplaceAll(stripped, "shouldn't", "")
			} else if strings.Contains(lower, "should") {
				modal = "should"
				stripped = strings.ReplaceAll(stripped, "should", "")
			} else if strings.Contains(lower, "must not") || strings.Contains(lower, "mustn't") {
				modal = "mustnot"
				stripped = strings.ReplaceAll(stripped, "must not", "")
				stripped = strings.ReplaceAll(stripped, "mustn't", "")
			} else if strings.Contains(lower, "must") {
				modal = "must"
				stripped = strings.ReplaceAll(stripped, "must", "")
			} else if strings.HasPrefix(lower, "do ") {
				modal = "do"
				stripped = stripped[3:]
			}

			if modal != "" {
				// Clean up stripped string
				strippedWords := tokenize(stripped)
				if len(strippedWords) >= 3 { // Need a minimum phrase size
					lines = append(lines, lineInfo{
						original: trimmed,
						stripped: rebuildFromMap(strippedWords),
						modal:    modal,
						sourceID: chunk.SourceID,
					})
				}
			}
		}
	}

	var conflicts []schema.Conflict
	seenPairs := make(map[string]bool)

	for i := 0; i < len(lines); i++ {
		for j := i + 1; j < len(lines); j++ {
			l1 := lines[i]
			l2 := lines[j]

			if l1.sourceID == l2.sourceID {
				continue
			}

			// Check if modals are opposite pairs
			if expectedOpposite, ok := modalPairs[l1.modal]; ok && l2.modal == expectedOpposite {
				// Compare stripped phrases
				w1 := tokenize(l1.stripped)
				w2 := tokenize(l2.stripped)
				overlap := computeOverlap(w1, w2)

				if overlap > 0.80 {
					pairKey := l1.sourceID + "-" + l2.sourceID + "-" + l1.stripped
					if l2.sourceID < l1.sourceID {
						pairKey = l2.sourceID + "-" + l1.sourceID + "-" + l1.stripped
					}

					if !seenPairs[pairKey] {
						seenPairs[pairKey] = true
						conflicts = append(conflicts, schema.Conflict{
							Type:      "contradicting_imperative",
							Heading:   l1.original + " vs " + l2.original,
							SourceIDs: []string{l1.sourceID, l2.sourceID},
						})
					}
				}
			}
		}
	}

	return conflicts
}

func tokenize(text string) map[string]bool {
	words := make(map[string]bool)
	matches := wordRegex.FindAllString(strings.ToLower(text), -1)
	for _, m := range matches {
		if len(m) > 1 { // Skip single letters
			words[m] = true
		}
	}
	return words
}

func computeOverlap(w1, w2 map[string]bool) float64 {
	if len(w1) == 0 || len(w2) == 0 {
		return 0.0
	}

	intersection := 0
	for k := range w1 {
		if w2[k] {
			intersection++
		}
	}

	maxLen := math.Max(float64(len(w1)), float64(len(w2)))
	return float64(intersection) / maxLen
}

func rebuildFromMap(m map[string]bool) string {
	var words []string
	for k := range m {
		words = append(words, k)
	}
	return strings.Join(words, " ")
}
