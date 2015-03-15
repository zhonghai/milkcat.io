package deptree2svg

import (
  "math"
  "fmt"
  "strings"
)

type node struct {
  word string
  postag string
  head int
  label string
}

type Tree []node

func NewTree() *Tree {
  tree := Tree{}
  tree.Add("ROOT", "ROOT", 0, "ROOT")
  return &tree
}

func (self *Tree) Add(word string, postag string, head int, label string) {
  *self = append(*self, node{word, postag, head, label})
}

const header = `<svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="%d" height="%d">`
const arrow = `<defs><marker id="head" orient="auto" markerWidth="6" markerHeight="12" refX="0.1" refY="6"><path d="M0,0 V12 L6,6 Z" fill="blue" /></marker></defs>`
const footer = `</svg>`
const textFmt = `<text x="%.1f" y="%.1f" fill="red" style="font-family: arial,sans-serif; font-size: 14px; text-anchor: middle;">%s</text>`
const wordAndTagTextFmt = `<text x="%.1f" y="%.1f" style="text-anchor: middle; font-size: 16px; font-family: arial,sans-serif">%s</text>`
const arcFmt = `<path marker-mid="url(#head)" stroke-width="1" stroke="blue" fill="none" d="%s"></path>`

const arcBottom = 60.0
const termBottom = 40.0
const postagBottom = 20.0
const fontSize = 16.0
const textMargin = 24.0

func heightOfWidth(width float64) float64 {
  return math.Atan(math.Abs(width) / 500 + 0.05) / 1.57 * 200;
}

func arcSVG(from, to int,
            label string,
            tokenPos []float64,
            svgPossibleHeight float64) (string, string) {
  y := svgPossibleHeight - arcBottom
  fromX := tokenPos[from]
  toX := tokenPos[to]
  middleY := y - heightOfWidth(math.Abs(toX - fromX))
  middleX := fromX + (toX - fromX) / 2
  d := fmt.Sprintf("M%.1f,%.1f Q%.1f,%.1f %.1f,%.1f Q%.1f,%.1f %.1f,%.1f",
                   fromX, y, fromX, middleY, middleX,
                   middleY, toX, middleY, toX, y);
  return fmt.Sprintf(arcFmt, d),
         fmt.Sprintf(textFmt, middleX, middleY - 5.0, label);
}

func wordTagSVG(index int,
                word, postag string,
                svgPossibleHeight, startX float64) (string, float64) {
  x := startX + textMargin + possibleWidth(word) / 2.0
  wordSvg := fmt.Sprintf(wordAndTagTextFmt,
                         x,
                         svgPossibleHeight - termBottom,
                         word)
  postagSvg := fmt.Sprintf(wordAndTagTextFmt,
                           x,
                           svgPossibleHeight - postagBottom,
                           postag)
  return wordSvg + postagSvg, x
}

func possibleWidth(word string) float64 {
  characters := strings.Split(word, "")
  possibleWidth := 0.0
  for _, ch := range characters {
    if len(ch) >= 3 {
      possibleWidth += fontSize
    } else {
      possibleWidth += fontSize / 2.0
    }
  }
  return possibleWidth
}

func TreeToSVG(tree *Tree) string {
  svgText := make([]string, 0)
  tokenPos := make([]float64, len(*tree))
  svgText = append(svgText, arrow)
  
  // Calculates the possible width of svg image
  svgPossibleWidth := textMargin
  for _, n := range *tree {
    svgPossibleWidth += textMargin + possibleWidth(n.word)
  }
  svgPossibleHeight := heightOfWidth(svgPossibleWidth) + 80.0

  nextX := 0.0
  for idx, n := range *tree {
    svg, x := wordTagSVG(idx, n.word, n.postag, svgPossibleHeight, nextX)
    tokenPos[idx] = x
    svgText = append(svgText, svg)
    nextX = x + possibleWidth(n.word) / 2.0
  }
  
  labelSVG := make([]string, 0)
  for idx, n := range *tree {
    if idx != 0 {
      arc, label := arcSVG(n.head, idx, n.label, tokenPos, svgPossibleHeight)
      svgText = append(svgText, arc)
      labelSVG = append(labelSVG, label)
    }
  }

  for _, label := range labelSVG {
    svgText = append(svgText, label)
  }

  svgText = append(svgText, footer)
  return fmt.Sprintf(header, int(nextX + textMargin), int(svgPossibleHeight)) + 
         strings.Join(svgText, "\n")
}





