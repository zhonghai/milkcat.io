package milkcat

// #cgo LDFLAGS: -lmilkcat -lstdc++ -lm
// #include <stdlib.h>
// #include <milkcat.h>
import "C"

import (
  "unsafe"
  "errors"
  "runtime"
)

const (
  defaultPath = "-DEFAULT-PATH-"
)

type ParserOptions struct {
  options C.milkcat_parseroptions_t
}

type Parser struct {
	parserPtr *C.milkcat_parser_t
  iteratorPtr *C.milkcat_parseriterator_t
}

type Item struct {
  Word string `json:"word"`
  PartOfSpeechTag string `json:"postag"`
  Head int `json:"head"`
  DependencyLabel string `json:"deplabel"`
  IsBeginOfSentence bool `json:"bos"`
}

func NewParser(options *ParserOptions) (parser *Parser, err error) {
  parser = new(Parser)
  parser.parserPtr = C.milkcat_parser_new(&options.options)
  parser.iteratorPtr = C.milkcat_parseriterator_new()
  runtime.SetFinalizer(parser, func (parser *Parser) {
    C.milkcat_parser_destroy(parser.parserPtr)
    C.milkcat_parseriterator_destroy(parser.iteratorPtr)
  })
  if parser.parserPtr == nil {
    err = errors.New(C.GoString(C.milkcat_last_error()))
    parser = nil
    return
  }

  return
}

func (self *Parser) Predict(text string) []Item {
  textPtr := C.CString(text)
  defer C.free(unsafe.Pointer(textPtr))
  C.milkcat_parser_predict(self.parserPtr, self.iteratorPtr, textPtr)

  prediction := make([]Item, 0)
  for C.milkcat_parseriterator_next(self.iteratorPtr) {
    prediction = append(prediction, Item{
      C.GoString(self.iteratorPtr.word),
      C.GoString(self.iteratorPtr.part_of_speech_tag),
      int(self.iteratorPtr.head),
      C.GoString(self.iteratorPtr.dependency_label),
      bool(self.iteratorPtr.is_begin_of_sentence)})
  }

  return prediction
}


func NewParserOptions() *ParserOptions {
  parserOpt := new(ParserOptions)
  C.milkcat_parseroptions_init(&parserOpt.options)
  return parserOpt
}

func (self *ParserOptions) UseBigramSegmenter() {
  self.options.word_segmenter = C.MC_SEGMENTER_BIGRAM
}

func (self *ParserOptions) UseCrfSegmenter() {
  self.options.word_segmenter = C.MC_SEGMENTER_CRF
}

func (self *ParserOptions) UseMixedSegmenter() {
  self.options.word_segmenter = C.MC_SEGMENTER_MIXED
}

func (self *ParserOptions) UseFastCrfPOSTagger() {
  self.options.part_of_speech_tagger = C.MC_POSTAGGER_MIXED
}

func (self *ParserOptions) UseHmmPOSTagger() {
  self.options.part_of_speech_tagger = C.MC_POSTAGGER_HMM
}

func (self *ParserOptions) UseBeamDependencyParser() {
  self.options.dependency_parser = C.MC_DEPPARSER_BEAMYAMADA
}
