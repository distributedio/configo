package configo

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

type Rule interface {
	Parse() []Validator
}
type rule string

const (
	rule_start = iota
	vldt_start //validator start
	vldt_end   //validator end
	rule_end
)

func (r rule) Parse(rt reflect.Type) ([]Validator, error) {
	var vlds []Validator

	st := rule_start
	for i := 0; i < len(r); i++ {
		c := r[i]
		switch c {
		case '(', '[':
			switch st {
			case rule_start, vldt_end:
				st = vldt_start
				v, pos, err := parseRange(string(r), i)
				if err != nil {
					return nil, err
				}
				i = pos
				vlds = append(vlds, v)
				st = vldt_end
			default:
				return nil, fmt.Errorf("parse error at position %d of rule %s", i, r)
			}
		case '<', '>':
			switch st {
			case rule_start, vldt_end:
				st = vldt_start
				v, pos, err := parseCompExp(string(r), i)
				if err != nil {
					return nil, err
				}
				i = pos
				vlds = append(vlds, v)
				st = vldt_end
			}
		case 'r':
		case ' ', '\t':
			continue
		}
	}
	return vlds, nil
}

func parseInt(val string) (int64, error) {
	num, err := strconv.ParseInt(val, 0, 64)
	if err != nil {
		d, err := time.ParseDuration(val)
		if err != nil {
			return 0, err
		}
		num = int64(d)
	}
	return num, nil
}

func parseRange(r string, pos int) (*nrange, int, error) {
	i := pos
	v := &nrange{}
	vstart := 0
	vend := 0

	const (
		range_start = iota
		min_start
		min_val
		min_end
		max_start
		max_val
		max_end
		range_end
	)
	st := range_start

LOOP:
	for ; i < len(r); i++ {
		c := r[i]
		switch c {
		case '(', '[':
			if st == range_start {
				st = min_start
				vstart = i + 1
			} else {
				return nil, i, fmt.Errorf("range parse error at position %d of rule %s", i, r)
			}
			if c == '[' {
				v.left = true
			}
		case ',':
			switch st {
			case min_val:
				st = min_end
				vend = i
				fallthrough
			case min_end:
				val := r[vstart:vend]
				if val == "" {
					v.min = math.MinInt64
				} else {
					num, err := parseInt(val)
					if err != nil {
						return nil, i, err
					}
					v.min = num
				}

				st = max_start
				vstart = i + 1
			default:
				return nil, i, fmt.Errorf("range parse error at position %d of rule %s", i, r)
			}
		case ')', ']':
			switch st {
			case max_val:
				st = max_end
				vend = i
				fallthrough
			case max_end:
				val := r[vstart:vend]
				if val == "" {
					v.max = math.MaxInt64
				} else {
					num, err := parseInt(val)
					if err != nil {
						return nil, i, err
					}
					v.max = num
				}

				if c == ']' {
					v.right = true
				}

				st = range_end
				break LOOP
			default:
				return nil, i, fmt.Errorf("range parse error at position %d of rule %s", i, r)
			}
		case ' ', '\t':
			switch st {
			case min_start, max_start:
				vstart = i + 1
			case min_val:
				vend = i
				st = min_end
			case max_val:
				vend = i
				st = max_end
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '.', //positive, negative, or float digital
			'A', 'B', 'C', 'D', 'E', 'F', 'a', 'b', 'c', 'd', 'e', 'f', 'X', 'x', //hex
			'h', 'm', 's', 'n', 'u', 'µ': //time.Duration
			switch st {
			case min_start:
				st = min_val
			case max_start:
				st = max_val
			}
		default:
			return nil, i, fmt.Errorf("unknown char at position %d of rule %s", i, r)
		}
	}

	if st != range_end {
		return nil, i, fmt.Errorf("rule is incompleted: %s", r)
	}

	return v, i, nil
}

//parse comparison expression
func parseCompExp(r string, pos int) (*nrange, int, error) {
	i := pos
	v := &nrange{}

	const (
		ce_start = iota
		val_start
		val_in
		val_end
		ce_end
	)
	st := ce_start
	vstart := 0
	vend := 0

	lt := true

LOOP:
	for ; i < len(r); i++ {
		c := r[i]
		switch c {
		case '<', '>':
			if st != ce_start {
				return nil, i, fmt.Errorf("rule parse error at position %d of rule %s", i, r)
			}
			if c == '<' {
				lt = true
			} else {
				lt = false
			}

			vstart = i + 1
			st = val_start
		case '=':
			if st != val_start {
				return nil, i, fmt.Errorf("rule parse error at position %d of rule %s, '=' should comes after '>' or '<'", i, r)
			}
			if lt {
				v.left = true
			} else {
				v.right = true
			}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-', '.', //positive, negative, or float digital
			'A', 'B', 'C', 'D', 'E', 'F', 'a', 'b', 'c', 'd', 'e', 'f', 'X', 'x', //hex
			'h', 'm', 's', 'n', 'u', 'µ': //time.Duration
			if st == val_start {
				st = val_in
			}
		case ' ', '\t':
			switch st {
			case val_start:
				vstart = i + 1
			case val_in:
				st = val_end
				vend = i
				val := r[vstart:vend]
				num, err := parseInt(val)
				if err != nil {
					return nil, i, err
				}
				if lt {
					v.max = num
				} else {
					v.min = num
				}
				st = ce_end
				break LOOP
			default:
			}
		default:
			return nil, i, fmt.Errorf("unknown char at position %d of rule %s", i, r)
		}
	}
	if i == len(r) && st == val_in {
		st = val_end
		vend = i
		val := r[vstart:vend]
		num, err := parseInt(val)
		if err != nil {
			return nil, i, err
		}
		if lt {
			v.max = num
		} else {
			v.min = num
		}
		st = ce_end

	}
	if st != ce_end {
		return nil, i, fmt.Errorf("rule is incompleted %s", r)
	}

	if lt {
		v.min = math.MinInt64
	} else {
		v.max = math.MaxInt64
	}
	return v, i, nil
}

func parseRegex(r string, pos int) (*regex, int, error) {
	const (
		re_start = iota
		exp_start
		exp_val
		exp_end
		re_end
	)

	v := &regex{}

	vstart := 0
	vend := 0
	i := pos
	st := re_start
LOOP:
	for ; i < len(r); i++ {
		c := r[i]
		switch c {
		case '/':
			switch st {
			case re_start:
				vstart = i + 1
				st = exp_start
			case exp_start, exp_val:
				vend = i
				st = exp_end
				val := r[vstart:vend]
				v.exp = val
				st = re_end
				break LOOP
			}
		case ' ', '\t':
			if st == re_start {
				continue
			}
			fallthrough
		default:
			switch st {
			case re_start:
				return nil, i, fmt.Errorf("regex parse error at position %d, rule %q", i, r)
			case exp_start:
				st = exp_val
			}
		}
	}
	if st != re_end {
		return nil, i, fmt.Errorf("rule is incompleted: %s", r)
	}
	return v, i, nil
}
