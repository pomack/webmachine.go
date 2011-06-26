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


