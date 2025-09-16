![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)
# Synodict-Go

**Synodict-Go** is a console-based synonym dictionary written in Go.  
It stores words and their relations as an undirected graph, allowing you to add and check synonyms, as well as import/export dictionaries in multiple formats.

## Features

- Add words and create synonym relations
- Check:
  - whether a word exists in the dictionary
  - whether two words are synonyms (direct or transitive)
  - whether a direct link exists between two words
- List all direct-linked synonyms of a word
- Import/export dictionaries in:
  - **GOB** (Go serialization format)
  - **CSV** (Saves the original word order)
  - **CSV condensed** (Does not save the original word order but uses less memory)
- Import conflict modes:
  - **Overwrite (o)** — fully replace the current dictionary
  - **Merge (m)** — merge the dictionaries
  - **Cancel (c)** — cancel the import
- Safe confirmation prompts before overwriting data
- Supports words with Latin, Cyrillic, diacritics, spaces, and hyphens

## Usage

```bash
git clone git@github.com:neimee8/synodict-go.git
cd synodict-go
go run main.go
```

### Example session:
```
type "help" for instructions
 > add "fast" "quick"
 > check "fast" "quick"
synonims: yes
 > exists "slow"
exists: no
 > export
please choose the export format:
gob  - GOB
csv  - CSV
csvc - CSV condensed
c    - go back without export
done - stop execution
 > csv
please specify the file location:
c    - go back to the previous step
done - stop execution
 > dict.csv
exported successfully
 > done
goodbye!
```

## Commands

```
add "word1"...
```
Adds each word to the dictionary if not already present, and links them as synonyms

```
add-words "word1"...
```
Adds each word to the dictionary if not already present (does not link them as synonyms)

```
remove "word1"...
```
Removes each word from the dictionary if already present

```
unlink "word1" "word2"
```
Removes synonym link between words (does not delete the words themselves)

```
unlink-clean "word1" "word2"
```
Removes synonym link between words (deletes words if they have no other synonyms)

```
check "word1" "word2"
```
Checks if the words are synonyms (directly or transitively)

```
check-direct "word1" "word2"
```
Checks if the words are directly linked as synonyms

```
exists "word"
```
Checks if the word exists in the dictionary

```
count "word"
```
Prints the number of synonyms of the word

```
synonyms "word"
```
Prints all synonyms of the word

```
direct-synonyms "word"
```
Prints only directly linked synonyms (words that were explicitly connected)

```
count-groups
```
Prints the number of synonym groups

```
groups
```
Prints all synonym groups

```
count-words
```
Prints the total number of words in the dictionary

```
words
```
Prints all words

```
cleanup
```
Removes words that have no synonyms from the dictionary

```
clear
```
Clears the dictionary (warning: cannot be undone)

```
import
```
Import dictionary (supports gob/csv); if current dictionary is not empty, you will be prompted to save, merge, or overwrite

```
export
```
Export dictionary (supports gob/csv)

```
help
```
Prints this help message

```
done
```
Stops execution
