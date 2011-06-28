package webmachine

import (
  "mime"
  "strconv"
  "strings"
)

type acceptMatch struct {
  thetype string
  subtype string
  parameters map[string]string
}

type standardMatch struct {
  strMatch string
  parameters map[string]string
}

func guessMime(filename string) string {
  if index := strings.LastIndex(filename, "."); index >= 0 {
    return mime.TypeByExtension(filename[index:])
  }
  return mime.TypeByExtension("." + filename)
}

func splitAcceptString(accept string) []acceptMatch {
  return splitAcceptArray(strings.Split(accept, ",", -1))
}

func splitAcceptArray(accept []string) []acceptMatch {
  retval := make([]acceptMatch, len(accept))
  for i, m := range accept {
    retval[i].parameters = make(map[string]string)
    parts := strings.Split(m, ";", -1)
    if index := strings.Index(parts[0], "/"); index >= 0 {
      retval[i].thetype = strings.TrimSpace(parts[0][0:index])
      retval[i].subtype = strings.TrimSpace(parts[0][index+1:])
    } else {
      retval[i].thetype = strings.TrimSpace(parts[0])
    }
    for j := 1; j < len(parts); j++ {
      if index := strings.Index(parts[j], "="); index >= 0 {
        k := strings.TrimSpace(parts[j][0:index])
        v := strings.TrimSpace(parts[j][index+1:])
        if len(k) > 0 && len(v) > 0 {
          retval[i].parameters[k] = v
        }
      }
    }
  }
  return retval
}

func chooseMediaType(providedStrings []string, accept string) string {
  bestScore := -1.0
  bestMatch := ""
  accepts := splitAcceptString(accept)
  providedSets := splitAcceptArray(providedStrings)
  for _, acceptMatch := range accepts {
    for i, provided := range providedSets {
      if (acceptMatch.thetype == provided.thetype || acceptMatch.thetype == "*") && (acceptMatch.subtype == provided.subtype || acceptMatch.subtype == "*") {
        score := 100.0
        if len(provided.subtype) > 0 && len(acceptMatch.subtype) > 0 && (provided.subtype == acceptMatch.subtype || acceptMatch.subtype == "*") {
          score += 10.0
        }
        for k,v := range acceptMatch.parameters {
          if k != "q" {
            if v2, ok := provided.parameters[k]; ok && v == v2 {
              score += 1.0
            }
          }
        }
        if q, ok := acceptMatch.parameters["q"]; ok {
          if qf, err := strconv.Atof64(q); err != nil {
            score *= qf
          }
        }
        if score > bestScore {
          bestScore = score
          bestMatch = providedStrings[i]
        }
      }
    }
  }
  return bestMatch
}


func chooseMediaTypeDefault(providedStrings []string, matchString string, defaultString string) string {
  s := chooseMediaType(providedStrings, matchString)
  if len(s) <= 0 {
    return defaultString
  }
  return s
}

func splitStandardMatchString(matchString string) []standardMatch {
  return splitStandardMatchArray(strings.Split(matchString, ",", -1))
}

func splitStandardMatchArray(matchStrings []string) []standardMatch {
  retval := make([]standardMatch, len(matchStrings))
  for i, m := range matchStrings {
    retval[i].parameters = make(map[string]string)
    parts := strings.Split(m, ";", -1)
    retval[i].strMatch = strings.TrimSpace(parts[0])
    for j := 1; j < len(parts); j++ {
      if index := strings.Index(parts[j], "="); index >= 0 {
        k := strings.TrimSpace(parts[j][0:index])
        v := strings.TrimSpace(parts[j][index+1:])
        if len(k) > 0 && len(v) > 0 {
          retval[i].parameters[k] = v
        }
      }
    }
  }
  return retval
}

func chooseStandardMatchString(providedStrings []string, matchString string) string {
  matchStrings := splitStandardMatchString(matchString)
  return chooseStandardMatch(providedStrings, matchStrings)
}

func chooseStandardMatch(providedStrings []string, matchStrings []standardMatch) string {
  bestScore := -1.0
  bestMatch := ""
  providedSets := splitStandardMatchArray(providedStrings)
  for _, match := range matchStrings {
    for i, provided := range providedSets {
      if match.strMatch == provided.strMatch || match.strMatch == "*" {
        score := 100.0
        for k,v := range match.parameters {
          if k != "q" {
            if v2, ok := provided.parameters[k]; ok && v == v2 {
              score += 1.0
            }
          }
        }
        if q, ok := match.parameters["q"]; ok {
          if qf, err := strconv.Atof64(q); err != nil {
            score *= qf
          }
        }
        if score > bestScore {
          bestScore = score
          bestMatch = providedStrings[i]
        }
      }
    }
  }
  return bestMatch
}

func chooseStandardMatchDefaultString(providedStrings []string, matchString string, defaultString string) string {
  matchStrings := splitStandardMatchString(matchString)
  return chooseStandardMatchDefault(providedStrings, matchStrings, defaultString)
}

func chooseCharsetWithDefaultString(providedStrings []string, matchString string) string {
  return chooseStandardMatchDefaultString(providedStrings, matchString, "UTF-8")
}

func chooseEncodingWithDefaultString(providedStrings []string, matchString string) string {
  return chooseStandardMatchDefaultString(providedStrings, matchString, "identity")
}

func chooseStandardMatchDefault(providedStrings []string, matchStrings []standardMatch, defaultString string) string {
  s := chooseStandardMatch(providedStrings, matchStrings)
  if len(s) <= 0 {
    return defaultString
  }
  return s
}

func chooseCharsetWithDefault(providedStrings []string, matchStrings []standardMatch) string {
  return chooseStandardMatchDefault(providedStrings, matchStrings, "UTF-8")
}

func chooseEncodingWithDefault(providedStrings []string, matchStrings []standardMatch) string {
  return chooseStandardMatchDefault(providedStrings, matchStrings, "identity")
}



func (p *PassThroughMediaTypeHandler) splitRangeHeaderString(rangeHeader string) ([][2]int64) {
  if len(rangeHeader) > 6 && rangeHeader[0:6] == "bytes=" {
    rangeStrings := strings.Split(rangeHeader[6:], ",", -1)
    ranges := make([][2]int64, len(rangeStrings))
    for i, rangeString := range rangeStrings {
      trimmedRangeString := strings.TrimSpace(rangeString)
      dashIndex := strings.Index(rangeString, "-")
      switch {
      case dashIndex < 0:
        // single byte, e.g. 507
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString)
        ranges[i][1] = ranges[i][0] + 1
      case dashIndex == 0:
        // start from end, e.g -51
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString)
        ranges[i][0] += p.numberOfBytes
        ranges[i][1] = p.numberOfBytes
      case dashIndex == len(trimmedRangeString):
        // byte to end, e.g. 9500-
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString)
        ranges[i][1] = p.numberOfBytes
      default:
        // range, e.g. 400-500
        ranges[i][0], _ = strconv.Atoi64(trimmedRangeString[0:dashIndex])
        ranges[i][1], _ = strconv.Atoi64(trimmedRangeString[dashIndex:])
        ranges[i][1] += 1
      }
      if ranges[i][0] >= p.numberOfBytes {
        continue
      }
      if ranges[i][1] > p.numberOfBytes {
        ranges[i][1] = p.numberOfBytes
      }
      // TODO sorting and compression of byte ranges
    }
    // sort ranges in ascending order
    for i, arange := range ranges {
      for ; i > 0 && ranges[i-1][0] > arange[0]; i-- {
        ranges[i-1][0], ranges[i-1][1], ranges[i][0], ranges[i][1] = ranges[i][0], ranges[i][1], ranges[i-1][0], ranges[i-1][1]
      }
    }
    // perform range compression for non-canonical ranges
    l := list.New()
    lastRange := ranges[0]
    for i, arange := range ranges {
      if i == 0 || lastRange[1] < arange[0] {
        l.PushBack(arange)
        lastRange = arange
      } else if lastRange[1] >= arange[0] {
        if lastRange[1] < arange[1] {
          lastRange[1] = arange[1]
        }
      } else {
        l.PushBack(arange)
        lastRange = arange
      }
    }
    theranges := make([][2]int64, l.Len())
    for i, elem :=0, l.Front(); elem != nil; elem, i = elem.Next(), i + 1 {
      theranges[i] = elem.Value.([2]int64)
    }
    return theranges
  }
  theranges := make([][2]int64, 1)
  theranges[0][0] = 0
  theranges[0][1] = p.numberOfBytes
  return theranges
}
