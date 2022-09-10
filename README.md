# depslo
Psuedo Localization Tool

## Why?
A software that is available in multiple languages is usually designed and developed in one language. One the developement is done of that iteration, all the finalised texts are translated by translators who are unaware of the layout constraints in the design. This can lead to broken layout, or too much white space in some languages.

A software product that supports multiple languages has the vulnerability of having its layout broken during each release as new text is added. This tool attempts aid developers find these issues during the developement process by trying to suggesting a alternative lengthened strings in psudeo localized format.

## How does this solve that problem?
The tool tries to mitigate those problems by performing psuedo localization of english sentences. The developers can hit the endpoint with the JSON of translations in english and get back the psuedo localized JSON of all the strings.

## Features
* Endpoint to pass a JSON with english strings and get back a JSON with psuedo localized strings

## To Implement
* Be able to send a JSON file with english translation and get back a JSON file with psudeo localized strings
* Run on change to JSON file with english translation.
* Psuedo localization based on language locale, ideally based on locale in file name.

## To Research
* Research into optimum string elongation and contraction overall with respect to English in global languages
* Research into optimum string elongation and contraction in specific languages with respect to English

## For the curious
If you wish to try it out today,

1. Clone the repo, navigate to the `depslo` directory
2. Run `./run.sh`, the server is running now.
  - Can be verified by `curl http://localhost:1234/ping`, should output `pong`
3. While the server is running, make a POST request with a JSON with english strings that will be translated.
  - Sample POST request
  ```bash
  curl http://localhost:1234/translate --include --header "Content-Type: application/json" --data '{"HELLO": "Hello, this is rakshith", "tITLE": "the coolest developer tool"}'
  ```

  - Sample JSON data
  ```json
  {
    "KEY1": "string1",
    "KEY2": "string2",
    .....
    ....
  }
  ```
