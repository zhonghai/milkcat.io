package main

import (
  "fmt"
  "log"
  "encoding/json"
  "net/http"
  "milkcat"
  "deptree2svg"
)

var parser *milkcat.Parser = nil

func treeSVGHandler(w http.ResponseWriter, r *http.Request) {
  query := r.FormValue("q")
  contentType := r.FormValue("ct")
  sentence := parser.Predict(query)
  tree := deptree2svg.NewTree()
  for idx, item := range sentence {
    // Ignores the other sentences
    if item.IsBeginOfSentence && idx != 0 {
      break
    }
    tree.Add(item.Word,
             item.PartOfSpeechTag,
             item.Head,
             item.DependencyLabel)
  }

  if contentType == "svg" {
    w.Header().Set("Content-Type", "image/svg+xml")
  }

  fmt.Fprintln(w, deptree2svg.TreeToSVG(tree))
}

func parserHandler(w http.ResponseWriter, r *http.Request) {
  query := r.FormValue("q")
  sentence := parser.Predict(query)
  b, err := json.Marshal(sentence)
  if err != nil {
    http.Error(w, err.Error(), 500)
    return
  }

  w.Write(b)
}

func main() {
  parserOpt := milkcat.NewParserOptions()
  parserOpt.UseMixedSegmenter()
  parserOpt.UseFastCrfPOSTagger()
  parserOpt.UseBeamDependencyParser()

  var err error
  parser, err = milkcat.NewParser(parserOpt);
  if err != nil {
    log.Fatal(err)
  }

  fs := http.FileServer(http.Dir("static"))
  http.Handle("/", fs)

  http.HandleFunc("/tree2svg", treeSVGHandler)
  http.HandleFunc("/predict", parserHandler)
  http.ListenAndServe(":8080", nil)
}
