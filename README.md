# Find Closest File Names

Helps finding the same or very similar file names in a directory. For every combination of names, it will print:

1. First name and file size in parentheses
2. Second name and file size in parentheses
3. Levenshtein distance between the normalised names
4. First normalised name
5. Second normalised name
6. Empty line

The name combinations will be sorted by the levenshtein distance. The exactly same normalised strings will have the distance set to zero. The string normalisation converts the names to a sequence of words and removes special characters and particles.

## Example

```txt
❯ ./closest-file-names .
7 files
21 combinations

go.mod (54)
main.go (3,833)
3
go mod
go main

.git (288)
main.go (3,833)
5
git
go main

.vscode (96)
go.mod (54)
5
vscode
go mod
```
