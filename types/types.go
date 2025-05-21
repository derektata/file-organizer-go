package types

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileExtensions maps a category name (e.g. "image") to the list of accepted
// file‑name extensions (".jpg", ".png", …).  The custom Unmarshal below lets
// the YAML config use any of these forms per category:
//   programming: ".go, .rs"            # one scalar, comma‑separated
//   archive: [".zip", ".tar.gz"]        # flow sequence
//   audio:
//     - ".flac"                          # block sequence
//     - ".mp3"
// Keeping the parser flexible means the rest of the application can always rely
// on a uniform in‑memory structure: map[string][]string.

type FileExtensions map[string][]string

// UnmarshalYAML implements yaml.v3 custom decoding so the config *author* can
// mix scalar / block / flow styles while the *caller* sees a clean []string.
func (fe *FileExtensions) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("file‑extensions root must be a mapping, got %v", node.Kind)
	}

	out := make(FileExtensions)

	for i := 0; i < len(node.Content); i += 2 {
		keyNode := node.Content[i]
		valNode := node.Content[i+1]

		var category string
		if err := keyNode.Decode(&category); err != nil {
			return err
		}

		switch valNode.Kind {
		case yaml.ScalarNode:
			out[category] = splitAndClean(valNode.Value)

		case yaml.SequenceNode:
			var seq []string
			if err := valNode.Decode(&seq); err != nil {
				return err
			}
			out[category] = clean(seq)

		default:
			return fmt.Errorf("invalid YAML kind for %q: %v", category, valNode.Kind)
		}
	}

	*fe = out
	return nil
}

// splitAndClean breaks a comma‑separated scalar and normalises each token.
func splitAndClean(s string) []string {
	return clean(strings.Split(s, ","))
}

// clean trims whitespace / quotes, lower‑cases, dedups and drops empties.
func clean(raw []string) []string {
	seen := map[string]struct{}{}
	var out []string
	for _, p := range raw {
		p = strings.TrimSpace(strings.Trim(p, `"'`))
		if p == "" {
			continue
		}
		p = strings.ToLower(p) // case‑folding makes look‑ups predictable
		if !strings.HasPrefix(p, ".") {
			p = "." + p
		}
		if _, dup := seen[p]; dup {
			continue
		}
		seen[p] = struct{}{}
		out = append(out, p)
	}
	return out
}
