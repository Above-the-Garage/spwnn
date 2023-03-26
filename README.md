# spwnn
* **s**hallow **p**air-**w**ise **n**eural **n**etwork
* spelling correction / suggestion version
* 40 years in the making!

## Plea
This has the MIT License; but honestly the algorithm is so simple you could just code it up.  I would like to ask
that if you find that a reverse index like this one with pair-wise associations helpful because you learned about
it here that you give me a bit of credit.

## Examples

Various correctly spelled and misspelled words; some in the dictionary, some not.

```
PS C:\Users\Stephen\Go\src\github.com\above-the-garage\spwnncli> .\spwnncli.exe
.................................................
49167 words in dictionary

Command or word: contention
  100%  0       _contention_
Words Touched = 43%

Command or word: contension
  81%   0       _consistent_
  81%   0       _contention_
  81%   1       _contentions_
  81%   2       _consistently_
  81%   3       _comprehension_
  81%   5       _intensification_
Words Touched = 42%

Command or word: contemplation
  100%  0       _contemplation_
Words Touched = 47%

Command or word: ontemplation
  92%   1       _contemplation_
Words Touched = 43%

Command or word: templation
  90%   3       _contemplation_
Words Touched = 42%

Command or word: stepheeen
  80%   1       _stephenson_
  80%   2       _stephen_
Words Touched = 41%

Command or word: halleilujah
  41%   0       _hallucinate_
  41%   2       _hillsdale_
  41%   5       _halley_
Words Touched = 26%

Command or word: hallelulah
  45%   1       _hallucinate_
  45%   1       _tallahassee_
  45%   2       _callahan_
  45%   2       _cellular_
  45%   2       _tallahatchie_
  45%   3       _multicellular_
  45%   4       _halley_
  45%   4       _mullah_
  45%   4       _gullah_
  45%   5       _allah_
  45%   11      _electroencephalograph_
Words Touched = 31%

Command or word: hallelujah
  45%   1       _hallucinate_
  45%   4       _halley_
  45%   11      _electroencephalograph_
Words Touched = 26%

Command or word: -e
Bye!
PS C:\Users\Stephen\Go\src\github.com\above-the-garage\spwnncli>
```

## Try it out here
Ideally this is up and running.  :)

http://ec2-54-221-105-181.compute-1.amazonaws.com/


## Algorithm:

### Build dictionary/neural net:
* ForEach word in dictionary:
  * Add "_" to beginning and end of word
  * Add the word to a list rooted at each character pair

### So, for example:
  * "bear" -> "\_bear\_"
  * add "\_bear\_" to lists anchored at \_b, be, ea, ar, r\_

### Using neural net:
  * Add "_" to beginng and end of input word
  * For each letter pair
    * For each word in list rooted at each letter pair
      * Increase the score for that word by 1.0 / float64(len(input word)-1)
  * Sort all the words by their score and keep all of the words with a matching highest score
  * Sort the highest scoring words by the amount their length differs from the input word

## Display results:

### Print winning words
* Clients
  * spwnncli - on a command line
  * spwnnlambda - via an AWS  lambda
  * spwnnweb - via a web browser

### Verify dictionary, to show that most words correct to themselves:
* Run the algo on every word in the dictionary
* Print out any results where there is more than one word with the highest score that is the same length as the input word 
* Clients
  * spwnncli - on a command line (command -g)
  * spwnnmark - benchmark, also shows validation results; see README.md for spwnnmark for some results

## Effectiveness:
I have three dictionaries.
* ispell.words - about 49,000 words (of somewhat unknown origin, circa 2000)
  * 12 failures / 49167 total words = 0.0244 % error rate
* american.sml+.mwl - about 23,000 words (from ispell, 2022)
  * 8 failures / 23115 total words = 0.378% error rate
* dwyl-english-words-words_alpha.txt - about 370,000 words (from aspell, 2022)
  * 408 failures / 370104 total words = 0.0011% error rate

An interesting test and a good way to validate the code is to see if each word in the dictionary will produce just itself as a result.

### ispell.words failures
```
'_bringing_' could be '[{1 0 _bridging_} {1 0 _bringing_}]'
'_contented_' could be '[{1 0 _consented_} {1 0 _contended_} {1 0 _contented_}]'
'_deeded_' could be '[{1 0 _deeded_} {1 0 _deemed_}]'
'_descendent_' could be '[{1 0 _descendant_} {1 0 _descendent_}]'
'_indented_' could be '[{1 0 _indented_} {1 0 _intended_}]'
'_intended_' could be '[{1 0 _indented_} {1 0 _intended_}]'
'_microeconomics_' could be '[{1 0 _macroeconomics_} {1 0 _microeconomics_}]'
'_ratification_' could be '[{1 0 _ramification_} {1 0 _ratification_}]'
'_ringing_' could be '[{1 0 _rigging_} {1 0 _ringing_}]'
'_tinting_' could be '[{1 0 _tenting_} {1 0 _tinting_}]'
'_unindented_' could be '[{1 0 _unindented_} {1 0 _unintended_}]'
'_unintended_' could be '[{1 0 _unindented_} {1 0 _unintended_}]'
```
Out of 49,167 words, the algorithm produces 12 words that have the same score and are the same length.

This suggests that while the system isn't perfect, it is very good, which further suggests to me that human brains tend to want to make words that are highly differentiated based on pair-wise combinations.

### Interesting words
* "Intended" and "indented" are two words that show where the system breaks down:  both have the exact same set of letter pairs!  In my 49,000 word dictionary this only happens a couple of times: intended/indented; contented/contended.  (Unintended/Unindented are the nearly the same words as intended/indented.)
* The other words that don't validate _only_ to themselves have a pattern where there is a letter-pair that is duplicated:  "ringing" has 'in' and 'ng' duplicated, and produces the same score as "rigging", which also contains 'in and 'ng'.  Ah, but you say, shouldn't the repeated 'ng' in "ringing" make it a clear winner?  Yes, except, it doesn't quite work that way.

 ## Theory and history
 I started thinking about this approximately 1980 when I was a student at UCI.  
 
 This started when I wondered why, if people can read whole words at once - i.e., parse the whole word in parallel, and yet there needs to be an ordering to the letters in a word. And why is there such a thing as alphabetical order?  This led to me thinking about how little it takes before one identifies a series of letters as part of the alphabetical order:  ab  ... pq ...  xy; these letter orders are ingrained in us from an early age.

 So I decided to see if I could make a spelling corrector based solely on letter-pairs, and the ordering of the letters in the word would be implied by the way the letter-pairs interacted.
 
 My first version was in Lisp on a PDP-10 which had limited memory so I could only load in the words that started with "a".  I'm not sure what dictionary I used back then.  It seemed to work but the test cases were very limited.

 Around 2000 I realized my home computer was probably enough to run the algorithm, and I made a version in C.

 Around 2016 I made a version in Go, which mostly worked.

 Around 2018 I went on a programming retreat, where I fixed a key bug that had been in the original, and updated the Go version to support a command line client, a web client, and an AWS lambeda.  This is when I found a key bug:  when adding a word to a list of words containing a letter-pair, the code would add (for example) 'ng' in "ringing" twice, because it just interated through all of the letter-pairs.  But then when checking a word, it would double count the contribution of 'ng'.  So I put in a little test to only use each letter-pair once when building the dictionary.

 End of year 2022 I cleaned up all of my Go code, relearned how to launch and run a lambda (things were greatly improved!), made the standalone benchmark version, and started to document how it works.

 ## Summary
 I mean, it mostly works!  Why is that?  I don't really know, other than I had a hypothesis, I thought up a way to test my hypothesis, and the results are pretty good.  Is there some deeper meaning to this?  Maybe!  In college I did some tests with word-pairs. I typed in the first paragraph of an article from the school paper "The New U" and assigned an atom (this was lisp!) to the word-pairs.  Atoms would be things like "football" or "politics".  It turns out there was very little variety in New U articles so it wasn't much of a test, but it sort of worked.  Later, when I read about how to make a reverse-index for web pages, I recognized it as similar, but using individual words rather than word pairs.  I suspect word-pairs would work better.  (And letter-triples and word-triples would probably work even better, but why complicate things?)

 
