// Copyright 2019 Sidhartha Mani <sidharthamn@gmail.com>                                                                                            
//                                                                                                                                                  
// Licensed under the Apache License, Version 2.0 (the "License");                                                                                  
// you may not use this file except in compliance with the License.                                                                                 
// You may obtain a copy of the License at                                                                                                          
//                                                                                                                                                  
//     http://www.apache.org/licenses/LICENSE-2.0                                                                                                   
//                                                                                                                                                  
// Unless required by applicable law or agreed to in writing, software                                                                              
// distributed under the License is distributed on an "AS IS" BASIS,                                                                                
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.                                                                         
// See the License for the specific language governing permissions and                                                                              
// limitations under the License.                                                                                                                   

package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"
)

type state int

const (
	def state = iota
	nextArg
	capture
	done
)

type parseFn func(byte, io.Writer, string) (interface{}, state, error)

var (
	captured = false
	scratch  = ""
)

func printf(out io.Writer, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("printf: atleast 1 argument expected")
	}
	format := args[0]
	if len(args) == 1 {
		_, err := fmt.Fprintf(out, format)
		return fmt.Errorf("printf: %v", err)
	}
	var parseFunc parseFn
	parseFunc = parseFormat
	var ps state
	ok := false
	j := 1

	for i := range format {
		var arg string
		if j < len(args) {
			arg = args[j]
		} else {
			arg = "MISSING"
		}
		var err error
		var parseFnInterface interface{}
		parseFnInterface, ps, err = parseFunc(format[i], out, arg)
		if err != nil {
			return fmt.Errorf("printf: %v", err)
		}
		if ps == nextArg {
			j++
		}
		if ps == capture {
			captured = true
			scratch = scratch + string(format[i])
		}
		if parseFunc, ok = parseFnInterface.(func(byte, io.Writer, string) (interface{}, state, error)); !ok {
			return fmt.Errorf("printf: invalid parseFn type")
		}
	}
	if ps != done && ps != nextArg {
		return fmt.Errorf("printf: incomplete control sequence")
	}
	return nil
}

func parseFormat(r byte, out io.Writer, arg string) (interface{}, state, error) {
	if r == '%' {
		return parsePercent, def, nil
	}
	fmt.Fprintf(out, "%s", string(r))
	return parseFormat, done, nil
}

func parsePercent(r byte, out io.Writer, arg string) (interface{}, state, error) {
	var size int
	defer func() {
		if captured && size != 0 {
			scratches := strings.Split(scratch, ".")
			if len(scratches) > 2 {
				return
			}
			if len(scratches[0]) > 0 {
				n, err := atoi(scratches[0])
				if err != nil {
					fmt.Fprintf(out, "%v", err)
				}
				if n > int64(size) {
					fmt.Fprintf(out, "%*s", n-int64(size), " ")
				}
			}
			captured = false
			scratch = ""
		}
	}()
	if r == '%' {
		_, err := fmt.Fprintf(out, "%%")
		size = 1
		return parseFormat, def, err
	}
	if r == 's' {
			size = len(arg)
			_, err := fmt.Fprintf(out, arg)
			return parseFormat, nextArg, err
	}
	if r == 'c' {
		if argRune, ok := utf8.DecodeRuneInString(arg); ok > 0 {
			return parseFormat, def, fmt.Errorf("expected rune, got %v", arg)
		} else {
			size = utf8.RuneLen(argRune)
			_, err := fmt.Fprintf(out, "%c", argRune)
			return parseFormat, nextArg, err
		}
	}
	if r == 'o' {
			n, err := atoi(arg)
			if err != nil {
				return parseFormat, def, fmt.Errorf("expected integer, got %v", arg)
			}
			size = len(arg)
			_, err = fmt.Fprintf(out, "%o", n)
			return parseFormat, nextArg, err
	}
	if r == 'x' {
			n, err := atoi(arg)
			if err != nil {
				return parseFormat, def, fmt.Errorf("expected integer, got %v", arg)
			}
			size = len(arg)
			_, err = fmt.Fprintf(out, "%x", n)
			return parseFormat, nextArg, err
	}
	if r == 'X' {
			n, err := atoi(arg)
			if err != nil {
				return parseFormat, def, fmt.Errorf("expected integer, got %v", arg)
			}
			size = len(arg)
			_, err = fmt.Fprintf(out, "%X", n)
			return parseFormat, nextArg, err
	}
	if r == 'd' {
			n, err := atoi(arg)
			if err != nil {
				return parseFormat, def, fmt.Errorf("expected integer, got %v", arg)
			}
			size = len(arg)
			_, err = fmt.Fprintf(out, "%d", n)
			return parseFormat, nextArg, err
	}
	if r == 'b' {
			n, err := atoi(arg)
			if err != nil {
				return parseFormat, def, fmt.Errorf("expected integer, got %v", arg)
			}
			size = len(arg)
			_, err = fmt.Fprintf(out, "%b", n)
			return parseFormat, nextArg, err
	}
	if r == 'f' {
			precision := int64(6)
			if captured {
				scratches := strings.Split(scratch, ".")
				if len(scratches) > 2 {
					return parseFormat, def, fmt.Errorf("illegal argument")
				}
				if len(scratches) == 2 {
					var err error
					precision, err = atoi(scratches[1])
					if err != nil {
						return parseFormat, def, err
					}
				}
			}
			f, err := strconv.ParseFloat(arg, 64)
			if err != nil {
				return parseFormat, def, fmt.Errorf("expected float, got %v", arg)
			}
			size = len(arg)
			_, err = fmt.Fprintf(out, "%.*f", precision, f)
			return parseFormat, nextArg, err
	}
	if r >= '0' && r <= '9' {
		return parsePercent, capture, nil
	}
	if r == '.' {
		return parsePercent, capture, nil
	}
	return parseFormat, def, fmt.Errorf("invalid control sequence")
}

func atoi(d string) (int64, error) {
	return strconv.ParseInt(d, 0, 64)
}
