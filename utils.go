package main

import (
    "io"
    "os"
    "fmt"
    "bufio"
    "strings"
)


// Read lines from a file.  Lines will start being sent to the callback function between first
// and max
func LinesFromFile(filename string, firstResult int, maxResults int, callback func(line string) bool) error {
    var err error
    var file *os.File = nil

    if (filename == "-") {
        file = os.Stdin
    } else {
        file, err = os.Open(filename)
        if err != nil {
            return err
        }
        defer file.Close()
    }

    bufr := bufio.NewReader(file)
    resultCount := 0

    for line, err := bufr.ReadString('\n') ; err == nil ; line, err = bufr.ReadString('\n') {
        if (resultCount >= firstResult) {
            line = strings.TrimSpace(line)
            if (! callback(line)) {
                return nil
            }
        }
        resultCount++
        if ((resultCount >= firstResult + maxResults) && (maxResults != -1)) {
            fmt.Fprintf(os.Stderr, "Maximum number of lines encountered (%d).  Use -c to change.\n", maxResults)
            return nil
        }
    }

    if err != io.EOF {
        return err
    } else {
        return nil
    }
}

// Displays an error message and kills the program.
func Die(fmtstr string, args ...interface{}) {
    fmt.Fprintf(os.Stderr, fmtstr, args)
    os.Exit(1)
}


// Escapes the characters of the passed in string so they can safely be used as a filename.
// The characters allowed are all alphanumeric characters, ':', '-' and '.'.  Any other
// characters will be escaped in a form similar to QueryEscape (spaces will be escaped to %20).
//
// These can be unescaped using url.QueryUnescape
func EscapeIdForFilename(id string) string {
    escapedString := make([]byte, 0, len(id))

    for i := 0; i < len(id); i++ {
        c := id[i]
        if escapableCharacterForFilename(c) {
            d1 := "0123456789ABCDEF"[c>>4]
            d2 := "0123456789ABCDEF"[c&15]
            escapedString = append(escapedString, '%', d1, d2)
        } else {
            escapedString = append(escapedString, c)
        }
    }

    return string(escapedString)
}

func escapableCharacterForFilename(c byte) bool {
    if ((c >= 'a') && (c <= 'z')) || ((c >= 'A') && (c <= 'Z')) || ((c >= '0') && (c <= '9')) {
        return false
    } else if (c == ':') || (c == '-') || (c == '.') || (c == '_') {
        return false
    } else {
        return true
    }
}