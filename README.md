Build dictionary/neural net:
  Add "_" to beginning and end of word
  Add the word to a list rooted at each character pair
  So, for example:
     "bear" -> "_bear_"
     add "_bear_" to lists anchored at _b, be, ea, ar, r_

Use neural net:
  Add "_" to beginng and end of word
  For each letter pair
    For each word in list rooted at each letter pair
       Increase the score for that word by 1.0 / float64(len(word)-1)
  Sort all the words by their score and keep all of the words with a matching highest score
  Sort the highest scoring words by the amount their length differs from the input word

Show results:
  Print those words

Verify dictionary, to show that most words correct to themselves:
  Run the algo on every word in the dictionary
  Print out any results where there is more than one word with the highest score that is the same length as the input word 

 
 