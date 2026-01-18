# CLAUDE.md - Project Context for Claude Code

## Project Overview
**spwnn** - A neural-network inspired spelling corrector using bigram (letter-pair) indexing. Created by Stephen Clarke-Willson, Ph.D.

## Repository Structure
All repos are under `github.com/Above-the-Garage` (public):

| Repo | Description |
|------|-------------|
| `spwnn` | Core library - the spelling correction algorithm |
| `spwnncli` | Command-line interface with interactive mode |
| `spwnnmark` | Benchmark/test tool - verifies each word corrects to itself |
| `spwnnweb` | Web server interface |
| `spwnnlambda` | AWS Lambda deployment |

## Build & Test
```bash
# Build any package
cd spwnn && go build

# Run tests (each word should correct to itself)
cd spwnnmark && ./spwnnmark.exe
# or
cd spwnncli && ./spwnncli.exe -test

# With custom dictionary
./spwnncli.exe -test -dict mywords.txt
```

## Key APIs
```go
// Load dictionary
dict := spwnn.ReadDictionary("knownWords.txt", true)  // filename, noisy

// Correct a word
results, wordsTouched := spwnn.CorrectSpelling(dict, "wrold", false)  // word, strictLen

// Dictionary methods
dict.WordCount()  // total words
dict.Words()      // slice of all words
```

## AWS Lambda Deployment (spwnnlambda)
```bash
# Scripts in spwnnlambda directory
./CreateGateway   # Creates Lambda function (uses dynamic AWS account lookup)
./UpdateGateway   # Updates function code
./DeleteGateway   # Deletes function
```

## Coding Style Preferences
- Go idiomatic style (methods over getter functions)
- Comments should match code (e.g., sAlphabetSize = 28 is 26 letters + underscore + non-alpha)
- Use `sort.Stable` for secondary sorts to preserve primary ordering
- Prefer parameterized functions over hardcoded values (e.g., dictionary filename)

## Git Workflow
- Commits should be pushed to GitHub after committing
- Sensitive data (AWS account IDs, tokens) must not appear in git history
- Use `git filter-branch` or similar to scrub history if needed before making repos public
