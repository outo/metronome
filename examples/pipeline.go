package examples

import (
   "fmt"
)

// function reversing letters in input string and returning it
// for brevity, ASCII characters only (https://en.wikipedia.org/wiki/ASCII)
func reverseLetters(input string) (output string) {
   for i := len(input) - 1; i >= 0; i -- {
      output += string(input[i])
   }
   return
}

// function wrapping input string after wrapAtLength characters and
// returning all wrapped segments as a slice of strings
// for brevity, ASCII characters only (https://en.wikipedia.org/wiki/ASCII)
func wrapString(wrapAtLength int, input string) (output []string) {
   i := 0
   for ; i < len(input)-wrapAtLength; i += wrapAtLength {
      output = append(output, input[i:i+wrapAtLength])
   }
   if i < len(input)-1 {
      output = append(output, input[i:])
   }
   return
}

// service reversing letters (ASCII only, oversimplifying it - use English alphabet)
func letterReversingService(inputStrings chan string) (outputStrings chan string) {
   outputStrings = make(chan string)
   go func() {
      for inputString := range inputStrings {
         outputStrings <- reverseLetters(inputString)
      }
      close(outputStrings)
   }()
   return
}

// service wrapping words at given length (ASCII only, oversimplifying it - use English alphabet)
func textWrappingService(wrapAt int, inputStrings chan string) (outputStrings chan string) {
   outputStrings = make(chan string)
   go func() {
      for inputString := range inputStrings {
         wrapped := wrapString(wrapAt, inputString)
         for _, segment := range wrapped {
            outputStrings <- segment
         }
      }
      close(outputStrings)
   }()
   return
}

var words = []string{
   "apple",
   "cucumber",
   "onion",
   "cabbage",
   "aubergine",
}

func main() {
   textChan := make(chan string)
   wrappedChan := textWrappingService(5, textChan)
   reversedChan := letterReversingService(wrappedChan)

   go func() {
      for _, text := range words {
         textChan <- text
      }
      close(textChan)
   }()

   for word := range reversedChan {
      fmt.Println(word)
   }

}
