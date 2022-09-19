package jsqt

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		// (transpose)
		{give: `3`, when: `(transpose)`, then: `3`},
		{give: `[{"a":3,"b":5,"c":7},{"a":4,"b":6}]`, when: `(transpose)`, then: `{"a":[3,4],"b":[5,6],"c":[7]}`},
		{give: `[{"a":3,"b":5},{"a":4,"b":6}]`, when: `(transpose)`, then: `{"a":[3,4],"b":[5,6]}`},
		{give: `{"a":[3,4],"b":[5,6,7,8]}`, when: `(transpose)`, then: `[{"a":3,"b":5},{"a":4,"b":6},{"b":7},{"b":8}]`},
		{give: `{"a":[3,4],"b":[5,6],"c":[7]}`, when: `(transpose)`, then: `[{"a":3,"b":5,"c":7},{"a":4,"b":6}]`},
		{give: `{"a":[3,4],"b":[5,6],"c":7}`, when: `(transpose)`, then: `[{"a":3,"b":5},{"a":4,"b":6}]`},
		{give: `{"a":[3,4],"b":[5,6]}`, when: `(transpose)`, then: `[{"a":3,"b":5},{"a":4,"b":6}]`},
		// (unwind)
		{give: `{"a":3}`, when: `(unwind b)`, then: `[]`},
		{give: `{"a":3,"b":null}`, when: `(unwind b)`, then: `[]`},
		{give: `{"a":3,"b":[]}`, when: `(unwind b)`, then: `[]`},
		{give: `{"a":3,"b":4}`, when: `(unwind b)`, then: `[]`},
		{give: `[{"a":3,"b":[4,5]},{"a":6,"b":[7,8]}]`, when: `(unwind b)`, then: `[{"a":3,"b":4},{"a":3,"b":5},{"a":6,"b":7},{"a":6,"b":8}]`},
		{give: `{"a":3,"b":[4,5]}`, when: `(unwind b -r x)`, then: `[{"a":3,"x":4},{"a":3,"x":5}]`},
		{give: `{"a":3,"b":[4,5]}`, when: `(unwind b)`, then: `[{"a":3,"b":4},{"a":3,"b":5}]`},
		// (expr)
		{give: ``, when: `(expr 1 + 2 ? 3)`, then: `3`},
		{give: ``, when: `(expr -4 * -5 + -(raw 1) - -2 + -3)`, then: `18`},
		{give: `{"a":3}`, when: `(expr 4 * (expr 5 + (get a)))`, then: `32`},
		{give: `{"a":3}`, when: `(expr 4 * 5 + (get a))`, then: `23`},
		{give: ``, when: `(expr 1 + 1 / 2)`, then: `1.5`},
		{give: ``, when: `(expr 1 + 1 % 2)`, then: `2`},
		{give: ``, when: `(expr 1 % 2)`, then: `1`},
		{give: ``, when: `(expr (raw 2) * 3 + 4 * 3 + 4 - 5)`, then: `17`},
		{give: ``, when: `(expr 2 * 3 + 4 * 3 + 4 - 5)`, then: `17`},
		{give: ``, when: `(expr 2 + 3 * 4 + 3 * 4 - 5)`, then: `21`},
		{give: ``, when: `(expr 3 * (expr 4 * 2 - 5 + (expr 2 + 3) * 4 * 2 - 5) + 2)`, then: `116`},
		{give: ``, when: `(expr 3 * 4 * 2 - 5 + (expr 2 + 3) * 4 * 2 - 5 + 2)`, then: `56`},
		{give: ``, when: `(expr 3 * 4 * 2 - 5 + 2 + 3 * 4 * 2 - 5 + 2)`, then: `42`},
		{give: ``, when: `(expr 3 * 4 - 5 + 2)`, then: `9`},
		{give: ``, when: `(expr 3 * 4 * 2 - 5 + 2)`, then: `21`},
		{give: ``, when: `(expr 2 + 3 * 4 - 5 + 2 + 3 * 4 - 5)`, then: `18`},
		{give: ``, when: `(expr 2.1 + 2.5)`, then: `4.6`},
		{give: ``, when: `(expr 2 - 5)`, then: `-3`},
		{give: ``, when: `(expr 2 - -2)`, then: `4`},
		{give: ``, when: `(expr 2 - 2)`, then: `0`},
		{give: ``, when: `(expr 3 * 4)`, then: `12`},
		{give: ``, when: `(expr 3 + 4)`, then: `7`},
		{give: ``, when: `(expr 3)`, then: `3`},
		// (at)
		{give: `[3,4]`, when: `(at (raw 1))`, then: `4`},
		{give: `[3,4]`, when: `(at 0)`, then: `3`},
		// (group)
		{give: `[{"a":3},{"a":3},{"a":4}]`, when: `(group a (val)) (sort)`, then: `{"3":[{"a":3},{"a":3}],"4":[{"a":4}]}`},
		{give: `[3,4,3,4,5]`, when: `(group (val) (key)) (sort)`, then: `{"3":[0,2],"4":[1,3],"5":[4]}`},
		// (unique)
		{give: `[{"a":3},{"a":3},{"a":4}]`, when: `(unique a)`, then: `[3,4]`},
		{give: `[3,4,3,4,5]`, when: `(unique)`, then: `[3,4,5]`},
		// (in)
		{give: `[3,4,5,6,7]`, when: `(collect (not (in [4,6])))`, then: `[3,5,7]`},
		{give: `[3,4,5,6,7]`, when: `(collect (in [4,6]))`, then: `[4,6]`},
		// (slice)
		{give: `3`, when: `(slice 0)`, then: `3`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 2 -1)`, then: `[5,6,7]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice -2)`, then: `[7,8]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 2 4)`, then: `[5,6]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 0 3)`, then: `[3,4,5]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 0 0)`, then: `[3,4,5,6,7,8]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 6)`, then: `[]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 3)`, then: `[6,7,8]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 1)`, then: `[4,5,6,7,8]`},
		{give: `[3,4,5,6,7,8]`, when: `(slice 0)`, then: `[3,4,5,6,7,8]`},
		// (match)
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -kk x)`, then: ``},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -kk a)`, then: `"a"`},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -kk -s ab)`, then: `"aab"`},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -kk -p aa)`, then: `"aab"`},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -kk -r ^aa)`, then: `"aab"`},
		{give: `{"a":3,"b":4}`, when: `(match -k b)`, then: `4`},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -k -r ^aa)`, then: `4`},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -k -s ac)`, then: `5`},
		{give: `{"a":3,"aab":4,"aac":5,"bac":6}`, when: `(match -k -p aa)`, then: `4`},
		{give: `[{"a":3},{"a":4},{"a":5}]`, when: `(get * (match -v (match -k -r ^a) -r [35]))`, then: `[{"a":3},{"a":5}]`},
		{give: `[{"a":3},{"a":4},{"a":5}]`, when: `(get * (match -v a -r [35]))`, then: `[{"a":3},{"a":5}]`},
		{give: `{"a_b":3,"a_c":4,"b_c":5}`, when: `(iterate (match -s _c) (this))`, then: `{"a_c":4,"b_c":5}`},
		{give: `{"a_b":3,"a_c":4,"b_c":5}`, when: `(iterate (match -p a_) (this))`, then: `{"a_b":3,"a_c":4}`},
		{give: `{"a_b":3,"a_c":4,"b_c":5}`, when: `(iterate (match -r ^a_) (this))`, then: `{"a_b":3,"a_c":4}`},
		{give: `{"a_b":3,"a_c":4,"b_c":5}`, when: `(iterate (match a_c) (this))`, then: `{"a_c":4}`},
		// (set)
		{give: `[{"a":[{"a":3},{"a":4}]},{"a":[{"a":5}]}]`, when: `(set * a -r y * a -r x (this))`, then: `[{"y":[{"x":3},{"x":4}]},{"y":[{"x":5}]}]`},
		{give: `[{"a":[{"a":3},{"a":4}]},{"a":[{"a":5}]}]`, when: `(set * a -r y * a (this))`, then: `[{"y":[{"a":3},{"a":4}]},{"y":[{"a":5}]}]`},
		{give: `[{"a":[{"a":3},{"a":4}]},{"a":[{"a":5}]}]`, when: `(set * a * a -r x (this))`, then: `[{"a":[{"x":3},{"x":4}]},{"a":[{"x":5}]}]`},
		{give: `[{"a":3},{"a":3}]`, when: `(set * a -r x (this))`, then: `[{"x":3},{"x":3}]`},
		{give: `{"a":3,"b":4}`, when: `(set c -r x (this))`, then: `{"a":3,"b":4}`},
		{give: `{"a":3,"b":4}`, when: `(set b -r x (this))`, then: `{"a":3,"x":4}`},
		{give: `{"a":3,"b":4}`, when: `(set a -r x (this))`, then: `{"x":3,"b":4}`},
		{give: `3`, when: `(set a 7)`, then: `3`},
		{give: `{"a":{"b":3}}`, when: `(set -i a a 0 a 3)`, then: `{"a":{"b":3,"a":[{"a":3}]}}`},
		{give: `{}`, when: `(set -i 0 a 0 a 3)`, then: `[{"a":[{"a":3}]}]`},
		{give: `{}`, when: `(set -i 0 a 0 0 3)`, then: `[{"a":[[3]]}]`},
		{give: `{}`, when: `(set -i 0 a 0 3)`, then: `[{"a":[3]}]`},
		{give: `{}`, when: `(set -i 0 0 a 3)`, then: `[[{"a":3}]]`},
		{give: `{}`, when: `(set -i 0 0 0 3)`, then: `[[[3]]]`},
		{give: `{}`, when: `(set -i 0 0 3)`, then: `[[3]]`},
		{give: `{}`, when: `(set -i 0 3)`, then: `[3]`},
		{give: `{}`, when: `(set -i 0)`, then: `0`},
		{give: `{}`, when: `(set -i "0" 3)`, then: `{"0":3}`},
		{give: `{}`, when: `(set -i a 0 b 0 3)`, then: `{"a":[{"b":[3]}]}`},
		{give: `{}`, when: `(set -i a 1 b 3)`, then: `{"a":[{"b":3}]}`},
		{give: `{}`, when: `(set -i a 0 b 3)`, then: `{"a":[{"b":3}]}`},
		{give: `{"1":[{"1":[{"1":3},{"1":4}]},{"1":[{"1":5},{"1":6}]}]}`, when: `(set 1 1 1 1 1 7)`, then: `{"1":[{"1":[{"1":3},{"1":4}]},{"1":[{"1":5},{"1":7}]}]}`},
		{give: `{"1":[{"1":[{"1":3},{"1":4}]},{"1":[{"1":5},{"1":6}]}]}`, when: `(set 1 * 1 * 1 7)`, then: `{"1":[{"1":[{"1":7},{"1":7}]},{"1":[{"1":7},{"1":7}]}]}`},
		{give: `{"*":[{"*":[{"*":3},{"*":4}]},{"*":[{"*":5},{"*":6}]}]}`, when: `(set * * * * * 7)`, then: `{"*":[{"*":[{"*":7},{"*":7}]},{"*":[{"*":7},{"*":7}]}]}`},
		{give: `{"a":[{"1":[{"c":3},{"c":4}]},{"1":[{"c":5},{"c":6}]}]}`, when: `(set a * 1 * c 7)`, then: `{"a":[{"1":[{"c":7},{"c":7}]},{"1":[{"c":7},{"c":7}]}]}`},
		{give: `{"a":[{"*":[{"c":3},{"c":4}]},{"*":[{"c":5},{"c":6}]}]}`, when: `(set a * * * c 7)`, then: `{"a":[{"*":[{"c":7},{"c":7}]},{"*":[{"c":7},{"c":7}]}]}`},
		{give: `{"a":[{"b":[{"c":3},{"c":4}]},{"b":[{"c":5},{"c":6}]}]}`, when: `(set a * b * c 7)`, then: `{"a":[{"b":[{"c":7},{"c":7}]},{"b":[{"c":7},{"c":7}]}]}`},
		{give: `[[3],[4,5]]`, when: `(set -i 0 1 x 7)`, then: `[[3,{"x":7}],[4,5]]`},
		{give: `{"a":[{"b":[{"c":3},{"c":4}]},{"b":[{"c":5},{"c":6}]}]}`, when: `(set -i a * b 100 {})`, then: `{"a":[{"b":[{"c":3},{"c":4},{}]},{"b":[{"c":5},{"c":6},{}]}]}`},
		{give: `{"a":[{"b":[{"c":3},{"c":4}]},{"b":[{"c":5},{"c":6}]}]}`, when: `(set -i a * b 7)`, then: `{"a":[{"b":7},{"b":7}]}`},
		{give: `{"a":[{"b":[{"c":3},{"c":4}]},{"b":[{"c":5},{"c":6}]}]}`, when: `(set -i a * b * 7)`, then: `{"a":[{"b":[7,7]},{"b":[7,7]}]}`},
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(set -i a * x c (raw 7))`, then: `{"a":[{"b":3,"x":{"c":7}},{"c":4,"x":{"c":7}}]}`},
		{give: `{"a":{"b":3,"c":4}}`, when: `(set -i a x c (raw 7))`, then: `{"a":{"b":3,"c":4,"x":{"c":7}}}`},
		{give: `[3]`, when: `(set -i 1 (raw 7))`, then: `[3,7]`},
		{give: `{}`, when: `(set -i x 5)`, then: `{"x":5}`},
		{give: `[{"a":[3,4]},{"a":[3,4]},{"b":[3,4]}]`, when: `(set * a 1 (raw 7))`, then: `[{"a":[3,7]},{"a":[3,7]},{"b":[3,4]}]`},
		{give: `{"a":[{"b":3,"a":[{"b":4},{"b":5}]},{"b":6}]}`, when: `(set a 0 a * b (raw 7))`, then: `{"a":[{"b":3,"a":[{"b":7},{"b":7}]},{"b":6}]}`},
		{give: `{"a":[{"b":3,"a":[{"b":4},{"b":5}]},{"b":6}]}`, when: `(set a 0 a 1 b (raw 7))`, then: `{"a":[{"b":3,"a":[{"b":4},{"b":7}]},{"b":6}]}`},
		{give: `{"a":[{"b":3,"a":[{"b":4},{"b":5}]},{"b":6}]}`, when: `(set a 0 a 0 b (raw 7))`, then: `{"a":[{"b":3,"a":[{"b":7},{"b":5}]},{"b":6}]}`},
		{give: `{"e":3,"a":{"b":{"c":[{"d":4},{"d":5}],"f":6},"f":8},"f":9}`, when: `(set a b c * d (raw 7))`, then: `{"e":3,"a":{"b":{"c":[{"d":7},{"d":7}],"f":6},"f":8},"f":9}`},
		{give: `[{"a":3},{"a":4},{"b":5}]`, when: `(set * a (raw 7))`, then: `[{"a":7},{"a":7},{"b":5}]`},
		{give: `{"a":[{"b":3},{"b":4}],"b":5}`, when: `(set a 1 b (raw 7))`, then: `{"a":[{"b":3},{"b":7}],"b":5}`},
		{give: `{"a":[{"b":3},{"b":4}],"b":5}`, when: `(set a 0 b (raw 7))`, then: `{"a":[{"b":7},{"b":4}],"b":5}`},
		{give: `{"a":{"0":[3,4]}}`, when: `(set a 0 1 (raw 7))`, then: `{"a":{"0":[3,7]}}`},
		{give: `{"a":{"b":[3,4]}}`, when: `(set a b 1 (raw 7))`, then: `{"a":{"b":[3,7]}}`},
		{give: `[[3],[4,5],[[8]]]`, when: `(set * 0 (nothing))`, then: `[[],[5],[]]`},
		{give: `[[3],[4,5],[[8]]]`, when: `(set 0 (nothing))`, then: `[[4,5],[[8]]]`},
		{give: `[[3],[4,5]]`, when: `(set 1 (nothing))`, then: `[[3]]`},
		{give: `[[3],[4,5]]`, when: `(set 0 (nothing))`, then: `[[4,5]]`},
		{give: `[[3],[4,5]]`, when: `(set 0 (raw 7))`, then: `[7,[4,5]]`},
		{give: `[[3],[4,5]]`, when: `(set 1 (raw 7))`, then: `[[3],7]`},
		{give: `[[3],[4,5]]`, when: `(set 1 1 (raw 7))`, then: `[[3],[4,7]]`},
		{give: `[[3],[4,5]]`, when: `(set 1 0 (raw 7))`, then: `[[3],[7,5]]`},
		{give: `[[3],[4,5]]`, when: `(set 0 0 (raw 7))`, then: `[[7],[4,5]]`},
		{give: `[3,4,5]`, when: `(set 3 (nothing))`, then: `[3,4,5]`},
		{give: `[3,4,5]`, when: `(set 2 (nothing))`, then: `[3,4]`},
		{give: `[3,4,5]`, when: `(set 1 (nothing))`, then: `[3,5]`},
		{give: `[3,4,5]`, when: `(set 0 (nothing))`, then: `[4,5]`},
		{give: `[3,4,5]`, when: `(set 3 (raw 7))`, then: `[3,4,5]`},
		{give: `[3,4,5]`, when: `(set 2 (raw 7))`, then: `[3,4,7]`},
		{give: `[3,4,5]`, when: `(set 1 (raw 7))`, then: `[3,7,5]`},
		{give: `[3,4,5]`, when: `(set 0 (raw 7))`, then: `[7,4,5]`},
		{give: `{"0":3}`, when: `(set 0 (raw 7))`, then: `{"0":7}`},
		{give: `[[[3]]]`, when: `(set 0 (raw 7))`, then: `[7]`},
		{give: `[[3]]`, when: `(set 0 (raw 7))`, then: `[7]`},
		{give: `[3]`, when: `(set 0 (raw 7))`, then: `[7]`},
		{give: `{"a":{"b":"c"},"d":3}`, when: `(set a (nothing))`, then: `{"d":3}`},
		{give: `{"a":3,"b":{"c":4,"d":5,"e":{}},"f":6}`, when: `(set b e (nothing))`, then: `{"a":3,"b":{"c":4,"d":5},"f":6}`},
		{give: `{"a":{"b":"c"},"d":3}`, when: `(set d (raw 7))`, then: `{"a":{"b":"c"},"d":7}`},
		{give: `{"a":{"b":"bb","c":"cc"}}`, when: `(set a c (raw 7))`, then: `{"a":{"b":"bb","c":7}}`},
		{give: `{"a":3,"b":{"c":4,"d":5},"e":6}`, when: `(set b d (raw 7))`, then: `{"a":3,"b":{"c":4,"d":7},"e":6}`},
		{give: `{"a":3,"b":4}`, when: `(set a (nothing))`, then: `{"b":4}`},
		{give: `{"a":3,"b":4}`, when: `(set b (nothing))`, then: `{"a":3}`},
		{give: `{"a":3,"b":4}`, when: `(set b (raw 7))`, then: `{"a":3,"b":7}`},
		{give: `{"a":3,"b":4}`, when: `(set a (raw 7))`, then: `{"a":7,"b":4}`},
		{give: `{"a":3,"b":4}`, when: `(set a 7)`, then: `{"a":7,"b":4}`},
		// (key) (val)
		{give: `[3,4]`, when: `(collect (arr (key) (val)))`, then: `[[0,3],[1,4]]`},
		{give: `{"a":3,"c":4}`, when: `(iterate (concat (key) (val)) (arr (key) (val)))`, then: `{"a3":["a",3],"c4":["c",4]}`},
		{give: `{"a":3,"c":4}`, when: `(iterate -r (key) (arr (key) (val)))`, then: `[null,{"a":["a",3],"c":["c",4]}]`},
		{give: `{"a":3}`, when: `(iterate -k (key))`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(iterate -v (val))`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(iterate -kv (concat (key) (val)))`, then: `{"aa":"33"}`},
		// (save) (load)
		{give: `{"a":3,"b":4}`, when: `(save -k a b c -v (raw 7)) (obj x (load a) y (load b) z (load c)`, then: `{"x":3,"y":4,"z":7}`},
		{give: `{"a":3}`, when: `(save (get a)) (arr (load))`, then: `[3]`},
		{give: `{"a":3}`, when: `(save a) (arr (load))`, then: `[3]`},
		{give: `{"a":3}`, when: `a (save) (arr (load))`, then: `[3]`},
		// (iterate)
		{give: `3`, when: `(iterate -r (key) (stringify))`, then: `"3"`},
		{give: `{"a":{"b":{"c":3}}}`, when: `(iterate -d 3 (upper) (val))`, then: `{"A":{"B":{"C":3}}}`},
		{give: `{"a":{"b":{"c":3}}}`, when: `(iterate -d 2 (upper) (val))`, then: `{"A":{"B":{"c":3}}}`},
		{give: `{"a":{"b":{"c":3}}}`, when: `(iterate -d 1 (upper) (val))`, then: `{"A":{"b":{"c":3}}}`},
		{give: `{"a":3,"b":{},"c":{"d":4,"e":[]},"f":{"g":[]},"h":[5,[6],{}]}`, when: `(iterate (upper) (not (is-empty)))`, then: `{"A":3,"C":{"D":4},"H":[5,[6]]}`},
		{give: `{"a":{"a":3,"b":4,"c":5},"b":{"a":6,"b":7,"c":8},"c":9}`, when: `(iterate (!= "c") (this))`, then: `{"a":{"a":3,"b":4},"b":{"a":6,"b":7}}`},
		{give: `{"a":{"a":3,"b":4,"c":5},"b":{"a":6,"b":7,"c":8},"c":9}`, when: `(iterate (this) (pluck c))`, then: `{"a":{"a":3,"b":4},"b":{"a":6,"b":7},"c":9}`},
		{give: `{"a":{"a":3,"b":4,"c":{"a":5,"c":6}},"b":{"a":6,"b":7,"c":8},"c":9}`, when: `(iterate -r (this) (pluck c))`, then: `{"a":{"a":3,"b":4},"b":{"a":6,"b":7}}`},
		{give: `{"a":{"a":3,"b":4,"c":5},"b":{"a":6,"b":7,"c":8},"c":9}`, when: `(iterate -r (this) (pluck c))`, then: `{"a":{"a":3,"b":4},"b":{"a":6,"b":7}}`},
		{give: `{"a":3}`, when: `(iterate (raw "a b") (this))`, then: `{"a b":3}`},
		{give: `{"a":3}`, when: `(iterate (if (== "a") (nothing)) (this))`, then: `{}`},
		{give: `[4,5]`, when: `(iterate (this) (if (is-num) (nothing)))`, then: `[]`},
		// Test empty context input for many functions.
		{give: `["3 4","5 6","7 8"]`, when: `(collect (!= "5 6") (replace " " "_"))`, then: `["3_4","7_8"]`},
		{give: `[[3,4],[5,6],[7,8]]`, when: `(collect (!= 6 1) (reverse))`, then: `[[4,3],[8,7]]`},
		{give: `[[3,4],[5,6],[7,8]]`, when: `(collect (!= 6 1) (sort desc))`, then: `[[4,3],[8,7]]`},
		// (def)
		{give: `{"a":3,"b":4}`, when: `(def b (a)) (def a (stringify)) (arr (get a (a)) (get b (b)))`, then: `["3","4"]`},
		{give: `{"a":3,"b":4}`, when: `(def a (stringify)) (def b (a)) (arr (get a (a)) (get b (b)))`, then: `["3","4"]`},
		{give: `{"a":3,"b":4}`, when: `(def a (get a)) (def b (load)) (save (raw 5)) (arr (a) (b))`, then: `[3,5]`},
		{give: `{"a":3,"b":4}`, when: `(def a (get a)) (def b (get b)) (arr (a) (b))`, then: `[3,4]`},
		// (upsert)
		{give: `3`, when: `(upsert b 5)`, then: `3`},
		{give: `{"a":3}`, when: `(upsert (nothing) a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(upsert x (nothing))`, then: `{"a":3}`},
		{give: `{"a":3,"b":4,"c":5}`, when: `(upsert x 5 y (get a) b 7 a (nothing))`, then: `{"x":5,"y":3,"b":7,"c":5}`},
		{give: `{"a":3,"b":4}`, when: `(upsert b 5)`, then: `{"b":5,"a":3}`},
		// (pluck)
		{give: `{"a":3,"b":4,"c":5}`, when: `(pluck b c)`, then: `{"a":3}`},
		// (pick)
		{give: `{"a":3,"b":4,"c":5}`, when: `(pick b -r x -m (concat (raw "b") (this)) c)`, then: `{"x":"b4","c":5}`},
		{give: `{"a":3,"b":4,"c":5}`, when: `(pick b -r x c)`, then: `{"x":4,"c":5}`},
		{give: `{"a":{"b":3,"c":4,"d":{"e":3,"f":4}},"c":5}`, when: `(pick a -m (pick b d -m (pick f)))`, then: `{"a":{"b":3,"d":{"f":4}}}`},
		{give: `{"a":{"b":3,"c":4},"b":{"d":3,"e":4},"c":5}`, when: `(pick a -m (pick c) c -m (nothing))`, then: `{"a":{"c":4}}`},
		{give: `{"a":{"b":3,"c":4},"b":{"d":3,"e":4},"c":5}`, when: `(pick a -m (pick c) b -m (pick d))`, then: `{"a":{"c":4},"b":{"d":3}}`},
		{give: `{"a":3,"b":4,"c":5}`, when: `(pick b c)`, then: `{"b":4,"c":5}`},
		{give: `3`, when: `(pick)`, then: `3`},
		// (this)
		{give: `[3,4]`, when: `(collect (> 3))`, then: `[4]`},
		// (split)
		{give: `{"a":"one,two"}`, when: `(split ";" a)`, then: `["one,two"]`},
		{give: `"one,two,three"`, when: `(split ";")`, then: `["one,two,three"]`},
		{give: `"one,\"two three\""`, when: `(split ",")`, then: `["one","\"two three\""]`},
		{give: `"one,two,three"`, when: `(split ",")`, then: `["one","two","three"]`},
		// (join)
		{give: `{"x":["a","b","c"]}`, when: `(join "_" x)`, then: `"a_b_c"`},
		{give: `["a","b","c"]`, when: `(join "_")`, then: `"a_b_c"`},
		// (reverse)
		{give: `[{"a":3},{"b":4}]`, when: `(reverse)`, then: `[{"b":4},{"a":3}]`},
		{give: `[3,4,5,6,2,7]`, when: `(reverse)`, then: `[7,2,6,5,4,3]`},
		// (sort)
		{give: `[1, 10, 2]`, when: `(sort)`, then: `[1,2,10]`},
		{give: `[1, 10, 2]`, when: `(sort desc)`, then: `[10,2,1]`},
		{give: `{"c":1,"b":10,"a":2}`, when: `(sort)`, then: `{"a":2,"b":10,"c":1}`},
		{give: `{"a":1,"b":10,"c":2}`, when: `(sort desc)`, then: `{"c":2,"b":10,"a":1}`},
		{give: `[{"a":1},{"a":10},{"a":2}]`, when: `(sort desc a)`, then: `[{"a":10},{"a":2},{"a":1}]`},
		{give: `[{"a":1},{"a":10},{"a":2}]`, when: `(sort a)`, then: `[{"a":1},{"a":2},{"a":10}]`},
		// (either) should skip its arguments.
		{give: `{"a":"","b":"B","c":""}`, when: `(get (either a b c) (lower))`, then: `"b"`},
		// (either)
		{give: `{"b":"B"}`, when: `(either a b c)`, then: `"B"`},
		{give: `{"a":"","b":"","c":""}`, when: `(either a b c)`, then: `""`},
		{give: `{"a":"A","b":"","c":"C"}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"A","b":"B","c":""}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"A","b":"","c":""}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"","b":"","c":"C"}`, when: `(either a b c)`, then: `"C"`},
		{give: `{"a":"","b":"B","c":""}`, when: `(either a b c)`, then: `"B"`},
		{give: `{"a":"","b":"B","c":"C"}`, when: `(either a b c)`, then: `"B"`},
		{give: `{"a":"A","b":"B","c":"C"}`, when: `(either a b c)`, then: `"A"`},
		{give: `{"a":"","b":""}`, when: `(either a b)`, then: `""`},
		{give: `{"a":"","b":"B"}`, when: `(either a b)`, then: `"B"`},
		{give: `{"a":"A","b":""}`, when: `(either a b)`, then: `"A"`},
		{give: `{"a":"A","b":"B"}`, when: `(either a b)`, then: `"A"`},
		// (concat)
		{give: `{"a":"x","b":3,"c":true,"d":null,"e":{},"f":[],"g":{"h": 4},"i":false}`, when: `(concat a b c d e f g i)`, then: `"x3truenullfalse"`},
		{give: `{ "a": "hello" }`, when: `(concat a (raw " \"world\""))`, then: `"hello \"world\""`},
		{give: `{ "a": "hello", "b": "world" }`, when: `(concat a (raw " ") b)`, then: `"hello world"`},
		//(replace)
		{give: `"a b"`, when: `(replace " " "_")`, then: `"a_b"`},
		// (upper)
		{give: `"a"`, when: `(upper)`, then: `"A"`},
		// (lower)
		{give: `"A"`, when: `(lower)`, then: `"a"`},
		// (stringify)
		{give: `3`, when: `(stringify)`, then: `"3"`},
		{give: `-3`, when: `(stringify)`, then: `"-3"`},
		{give: `[]`, when: `(stringify)`, then: `"[]"`},
		{give: `{}`, when: `(stringify)`, then: `"{}"`},
		{give: `""`, when: `(stringify)`, then: `"\"\""`},
		{give: `{"a":[{"b":3},4,"5"]}`, when: `(stringify)`, then: `"{\"a\":[{\"b\":3},4,\"5\"]}"`},
		// (jsonify)
		{give: `"{\"a\":[{\"b\":3},4,\"5\"]}"`, when: `(jsonify)`, then: `{"a":[{"b":3},4,"5"]}`},
		{give: `"3"`, when: `(jsonify)`, then: `3`},
		{give: `"{}"`, when: `(jsonify)`, then: `{}`},
		// (bool)
		{give: `{"a":3}`, when: `(bool b)`, then: `false`},
		{give: `{"a":3}`, when: `(bool a)`, then: `true`},
		{give: `3`, when: `(is-num) (bool)`, then: `true`},
		{give: `{}`, when: `(is-num) (bool)`, then: `false`},
		// (or) (and) (not)
		{give: `{"a":3,"b":4}`, when: `(not a)`, then: ``},
		{give: `{"a":3,"b":4}`, when: `(and a b)`, then: `{"a":3,"b":4}`},
		{give: `{"a":3,"b":4}`, when: `(and a x)`, then: ``},
		{give: `{"a":3,"b":4}`, when: `(or a x)`, then: `{"a":3,"b":4}`},
		{give: `[3,"",4,"5"]`, when: `(collect (not (is-str)))`, then: `[3,4]`},
		{give: `[{"a":3},{"a":4},{"a":5},{"a":6}]`, when: `(collect (or (< 4 a) (> 5 a)) a)`, then: `[3,6]`},
		{give: `[{"a":3},{"a":4},{"a":5},{"a":6}]`, when: `(collect (and (>= 4 a) (<= 5 a)) a)`, then: `[4,5]`},
		// (objectify)
		{give: `[[1,3],[2,4]]`, when: `(objectify)`, then: `{"1":3,"2":4}`},
		{give: `[["a",3],["b",4]]`, when: `(objectify)`, then: `{"a":3,"b":4}`},
		// (entries)
		{give: `{"a":3,"b":4}`, when: `(entries) (collect (flatten))`, then: `["a",3,"b",4]`},
		{give: `{"a":3,"b":4}`, when: `(entries)`, then: `[["a",3],["b",4]]`},
		//(values)
		{give: `{"a":3,"b":4}`, when: `(values)`, then: `[3,4]`},
		// (key)
		{give: `{"a":3,"b":4}`, when: `(keys)`, then: `["a","b"]`},
		// (if)
		{give: `[3,4,5,6]`, when: `(iterate (key) (if (== 4) 7 (== 5) 8 (raw 9)))`, then: `[9,7,8,9]`},
		{give: `[3,4,5,6]`, when: `(iterate (key) (if (== 4) 7 (== 5) 8))`, then: `[3,7,8,6]`},
		{give: `[33,44]`, when: `(if -n (is-arr) "array" (raw "skip array"))`, then: `"skip array"`},
		{give: `[33,44]`, when: `(if (not (is-arr)) "array" (raw "skip array"))`, then: `"skip array"`},
		{give: `[33,44]`, when: `(if (is-arr) "array" (raw "not array"))`, then: `"array"`},
		{give: `[33,44]`, when: `(if (is-arr) 4)`, then: `4`},                             // Can omit Else.
		{give: `{"b":3}`, when: `(if (is-num) 4)`, then: `{"b":3}`},                       // Can omit Else.
		{give: `{"a":""}`, when: `(get a (if (is-str) {} (raw 3)))`, then: `{}`},          // Then.
		{give: `{"a":{"b":3}}`, when: `(get a (if (is-str) {} (this)))`, then: `{"b":3}`}, // Else.
		{give: `3`, when: `(if (is-num) (obj b (this)) (raw 3))`, then: `{"b":3}`},        // Then.
		{give: `{"b":3}`, when: `(if (is-num) 3 (this))`, then: `{"b":3}`},                // Else.
		// (is-void)
		{give: `{"a":[]}`, when: `(is-void a)`, then: `{"a":[]}`},
		{give: `{"a":[3]}`, when: `(is-void a)`, then: ``},
		{give: `{}`, when: `(is-void)`, then: `{}`},
		{give: `[]`, when: `(is-void)`, then: `[]`},
		{give: `""`, when: `(is-void)`, then: ``},
		// (is-blank)
		{give: `{"a":null}`, when: `(is-blank a)`, then: `{"a":null}`},
		{give: `{"a":3}`, when: `(is-blank a)`, then: ``},
		{give: `{}`, when: `(is-blank)`, then: `{}`},
		{give: `[]`, when: `(is-blank)`, then: `[]`},
		{give: `null`, when: `(is-blank)`, then: `null`},
		{give: `""`, when: `(is-blank)`, then: ``},
		// (is-nully)
		{give: `{"a":""}`, when: `(is-nully a)`, then: `{"a":""}`},
		{give: `{"a":3}`, when: `(is-nully a)`, then: ``},
		{give: `{}`, when: `(is-nully)`, then: `{}`},
		{give: `[]`, when: `(is-nully)`, then: `[]`},
		{give: `null`, when: `(is-nully)`, then: `null`},
		{give: `""`, when: `(is-nully)`, then: `""`},
		{give: `3`, when: `(is-nully)`, then: ``},
		// (is-some)
		{give: `{"a":0}`, when: `(is-some a)`, then: `{"a":0}`},
		{give: `{"a":null}`, when: `(is-some a)`, then: ``},
		{give: `3`, when: `(is-some)`, then: `3`},
		{give: `""`, when: `(is-some)`, then: `""`},
		{give: `null`, when: `(is-some)`, then: ``},
		// (nothing) - Undocumented.
		{give: `3`, when: `(nothing)`, then: ``},
		// (exists)
		{give: `{"a":3}`, when: `(exists a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(exists b)`, then: ``},
		{give: `3`, when: `(exists)`, then: `3`},
		{give: ``, when: `(exists)`, then: ``},
		// (truthy)
		{give: `{"a":1}`, when: `(truthy a)`, then: `{"a":1}`},
		{give: `{"a":0}`, when: `(truthy a)`, then: ``},
		{give: `{}`, when: `(truthy)`, then: ``},
		{give: `[]`, when: `(truthy)`, then: ``},
		{give: `0`, when: `(truthy)`, then: ``},
		{give: `""`, when: `(truthy)`, then: ``},
		{give: `null`, when: `(truthy)`, then: ``},
		{give: `[0]`, when: `(truthy)`, then: `[0]`},
		{give: `3`, when: `(truthy)`, then: `3`},
		{give: `{"a":3}`, when: `(truthy)`, then: `{"a":3}`},
		{give: `true`, when: `(truthy)`, then: `true`},
		// (falsy)
		{give: `{"a":0}`, when: `(falsy a)`, then: `{"a":0}`},
		{give: `{"a":1}`, when: `(falsy a)`, then: ``},
		{give: `{}`, when: `(falsy)`, then: `{}`},
		{give: `[]`, when: `(falsy)`, then: `[]`},
		{give: `0`, when: `(falsy)`, then: `0`},
		{give: `""`, when: `(falsy)`, then: `""`},
		{give: `false`, when: `(falsy)`, then: `false`},
		{give: `[0]`, when: `(falsy)`, then: ``},
		{give: `3`, when: `(falsy)`, then: ``},
		{give: `null`, when: `(falsy)`, then: `null`},
		// (is-empty-obj)
		{give: `{"a":{}}`, when: `(is-empty-obj a)`, then: `{"a":{}}`},
		{give: `{"a":{"b":3}}`, when: `(is-empty-obj a)`, then: ``},
		{give: `{}`, when: `(is-empty-obj)`, then: `{}`},
		{give: `{"a":3}`, when: `(is-empty-obj)`, then: ``},
		{give: `[]`, when: `(is-empty-obj)`, then: ``},
		// (is-empty-arr)
		{give: `{"a":[]}`, when: `(is-empty-arr a)`, then: `{"a":[]}`},
		{give: `{"a":[3]}`, when: `(is-empty-arr a)`, then: ``},
		{give: `{}`, when: `(is-empty-arr)`, then: ``},
		{give: `[0]`, when: `(is-empty-arr)`, then: ``},
		{give: `[]`, when: `(is-empty-arr)`, then: `[]`},
		// (is-empty-str)
		{give: `{"a":""}`, when: `(is-empty-str a)`, then: `{"a":""}`},
		{give: `{"a":"3"}`, when: `(is-empty-str a)`, then: ``},
		{give: `3`, when: `(is-empty-str)`, then: ``},
		{give: `""`, when: `(is-empty-str)`, then: `""`},
		// (is-empty)
		{give: `{"a":""}`, when: `(is-empty a)`, then: `{"a":""}`},
		{give: `{"a":"3"}`, when: `(is-empty a)`, then: ``},
		{give: `3`, when: `(is-empty)`, then: ``},
		{give: `{}`, when: `(is-empty)`, then: `{}`},
		{give: `[]`, when: `(is-empty)`, then: `[]`},
		{give: `""`, when: `(is-empty)`, then: `""`},
		// (is-null)
		{give: `{"a":null}`, when: `(is-null a)`, then: `{"a":null}`},
		{give: `{"a":3}`, when: `(is-null a)`, then: ``},
		{give: `3`, when: `(is-null)`, then: ``},
		{give: `null`, when: `(is-null)`, then: `null`},
		// (is-bool)
		{give: `{"a":false}`, when: `(is-bool a)`, then: `{"a":false}`},
		{give: `{"a":true}`, when: `(is-bool a)`, then: `{"a":true}`},
		{give: `{"a":3}`, when: `(is-bool a)`, then: ``},
		{give: `3`, when: `(is-bool)`, then: ``},
		{give: `true`, when: `(is-bool)`, then: `true`},
		{give: `false`, when: `(is-bool)`, then: `false`},
		// (is-str)
		{give: `{"a":"3"}`, when: `(is-str a)`, then: `{"a":"3"}`},
		{give: `{"a":3}`, when: `(is-str a)`, then: ``},
		{give: `3`, when: `(is-str)`, then: ``},
		{give: `"3"`, when: `(is-str)`, then: `"3"`},
		// (is-arr)
		{give: `{"a":[3]}`, when: `(is-arr a)`, then: `{"a":[3]}`},
		{give: `{"a":3}`, when: `(is-arr a)`, then: ``},
		{give: `3`, when: `(is-arr)`, then: ``},
		{give: `[]`, when: `(is-arr)`, then: `[]`},
		// (is-obj)
		{give: `{"a":{"b":3}}`, when: `(is-obj a)`, then: `{"a":{"b":3}}`},
		{give: `{"a":3}`, when: `(is-obj a)`, then: ``},
		{give: `3`, when: `(is-obj)`, then: ``},
		{give: `{}`, when: `(is-obj)`, then: `{}`},
		// (is-num)
		{give: `{"a":3}`, when: `(is-num a)`, then: `{"a":3}`},
		{give: `{"a":"3"}`, when: `(is-num a)`, then: ``},
		{give: `"3"`, when: `(is-num)`, then: ``},
		{give: `3`, when: `(is-num)`, then: `3`},
		// (ugly) (pretty)
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5, "d": "e f" } ], "c": [ ], "d": { } } ]`, when: `(pretty)`, then: "[\n    {\n        \"a\": 3,\n        \"b\": [\n            4,\n            {\n                \"c\": 5,\n                \"d\": \"e f\"\n            }\n        ],\n        \"c\": [],\n        \"d\": {}\n    }\n]"},
		{give: `[ { "a" : 3 , "b" : [ 4 , { "c" : 5, "d": "e f" } ], "c": [ ], "d": { } } ]`, when: `(ugly)`, then: `[{"a":3,"b":[4,{"c":5,"d":"e f"}],"c":[],"d":{}}]`},
		// (==)
		{give: `{"a":3}`, when: `(== 2 a)`, then: ``},
		{give: `{"a":3}`, when: `(== 3 a)`, then: `{"a":3}`},
		{give: `3`, when: `(== 2 (this))`, then: ``},
		{give: `3`, when: `(== 3 (this))`, then: `3`},
		{give: `3`, when: `(== 2)`, then: ``},
		{give: `3`, when: `(== 3)`, then: `3`},
		// (!=)
		{give: `{"a":3}`, when: `(!= 2 a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(!= 3 a)`, then: ``},
		{give: `3`, when: `(!= 3 (this))`, then: ``},
		{give: `3`, when: `(!= 2 (this))`, then: `3`},
		{give: `3`, when: `(!= 3)`, then: ``},
		{give: `3`, when: `(!= 2)`, then: `3`},
		// (>=)
		{give: `"bear"`, when: `(>= "ant")`, then: `"bear"`},
		{give: `"ant"`, when: `(>= "bear")`, then: ``},
		{give: `{"a":3}`, when: `(>= 2 a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(>= 3 a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(>= 4 a)`, then: ``},
		{give: `3`, when: `(>= 2 (this))`, then: `3`},
		{give: `3`, when: `(>= 3 (this))`, then: `3`},
		{give: `3`, when: `(>= 4 (this))`, then: ``},
		{give: `3`, when: `(>= 2)`, then: `3`},
		{give: `3`, when: `(>= 3)`, then: `3`},
		{give: `3`, when: `(>= 4)`, then: ``},
		// (<=)
		{give: `"bear"`, when: `(<= "ant")`, then: ``},
		{give: `"ant"`, when: `(<= "bear")`, then: `"ant"`},
		{give: `{"a":3}`, when: `(<= 2 a)`, then: ``},
		{give: `{"a":3}`, when: `(<= 3 a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(<= 4 a)`, then: `{"a":3}`},
		{give: `3`, when: `(<= 2 (this))`, then: ``},
		{give: `3`, when: `(<= 3 (this))`, then: `3`},
		{give: `3`, when: `(<= 4 (this))`, then: `3`},
		{give: `3`, when: `(<= 2)`, then: ``},
		{give: `3`, when: `(<= 3)`, then: `3`},
		{give: `3`, when: `(<= 4)`, then: `3`},
		// (>)
		{give: `{"a":3}`, when: `(> 2 a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(> 3 a)`, then: ``},
		{give: `3`, when: `(> 2 (this))`, then: `3`},
		{give: `3`, when: `(> 3 (this))`, then: ``},
		{give: `3`, when: `(> 2)`, then: `3`},
		{give: `3`, when: `(> 3)`, then: ``},
		// (<)
		{give: `{"a":3}`, when: `(< 4 a)`, then: `{"a":3}`},
		{give: `{"a":3}`, when: `(< 3 a)`, then: ``},
		{give: `3`, when: `(< 4 (this))`, then: `3`},
		{give: `3`, when: `(< 3 (this))`, then: ``},
		{give: `3`, when: `(< 4)`, then: `3`},
		{give: `3`, when: `(< 3)`, then: ``},
		// (comparison)
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (== 7 b) a)`, then: `[4]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (!= 7 b) a)`, then: `[3,5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (>= 7 b) a)`, then: `[4,5]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (<= 7 b) a)`, then: `[3,4]`},
		{give: `[{"a":3,"b":6},{"a":4,"b":7},{"a":5,"b":8}]`, when: `(collect (> 7 b) a)`, then: `[5]`},
		{give: `[{"a":3,"b":11},{"a":4,"b":100},{"a":5,"b":8}]`, when: `(collect (< 77 b) a)`, then: `[3,5]`},
		// (iterate -kv)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate -kv (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3, "b": [ 3 , { "c": "d" } ] }`, when: `(iterate -kv (upper))`, then: `{"A":3,"B":[3,{"C":"D"}]}`},
		{give: `{ "a": "b" }`, when: `(iterate -kv (upper))`, then: `{"A":"B"}`},
		{give: `3`, when: `(iterate -kv (stringify))`, then: `"3"`},
		{give: `3`, when: `(iterate -kv (this))`, then: `3`},
		// (iterate -v)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate -v (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3, "b": [ 3 , { "a" : 3 } ] }`, when: `(iterate -v (if (== 3) (raw 4) (this))))`, then: `{"a":4,"b":[4,{"a":4}]}`},
		{give: `{ "a": 3 }`, when: `(iterate -v (if (== 3) (raw 4) (this)))`, then: `{"a":4}`},
		{give: `3`, when: `(iterate -v (stringify))`, then: `"3"`},
		{give: `3`, when: `(iterate -v (this))`, then: `3`},
		// (iterate -k)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate -k (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3 , "b" : [ 3, { "a": 3 } ] }`, when: `(iterate -k (if (== "a") (raw "x") (this))))`, then: `{"x":3,"b":[3,{"x":3}]}`},
		{give: `{ "a" : 3 }`, when: `(iterate -k (if (== "a") (raw "x") (this)))`, then: `{"x":3}`},
		{give: `3`, when: `(iterate -k (stringify))`, then: `3`},
		{give: `3`, when: `(iterate -k (this))`, then: `3`},
		// (iterate -f)
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate -f (key) (val))`, then: `{"a":3}`},
		{give: "{ \"a\"\t:\t3\t}", when: `(iterate -f (this) (this))`, then: `{"a":3}`},
		{give: `{ "a" : 3 , "b" : [ { "c" : 4 } , { "c" : 5 } ] , "d" : [ 6 , true ] }`, when: `(iterate -f (key) (if (is-num) (stringify) (val)))`, then: `{"a":"3","b":[{"c":"4"},{"c":"5"}],"d":["6",true]}`},
		{give: `{ "a" : 3 , "b" : [ { "c" : 4 } , { "c" : 5 } ] , "d" : [ 6 , true ] }`, when: `(iterate -f (this) (if (is-num) (stringify) (this)))`, then: `{"a":"3","b":[{"c":"4"},{"c":"5"}],"d":["6",true]}`},
		{give: `{ "a" : 3, "b" : 4}`, when: `(iterate -f (this) (if (is-num) (stringify) (this)))`, then: `{"a":"3","b":"4"}`},
		{give: `[3,4]`, when: `(iterate -f (key) (if (is-num) (stringify) (val)))`, then: `["3","4"]`},
		{give: `[3,4]`, when: `(iterate -f (this) (if (is-num) (stringify) (this)))`, then: `["3","4"]`},
		{give: `["30"]`, when: `(iterate -f (this) (jsonify))`, then: `[30]`},
		{give: `3`, when: `(iterate -f (this) (stringify))`, then: `"3"`},
		{give: `3`, when: `(iterate -f (key) (val))`, then: `3`},
		{give: `3`, when: `(iterate -f (this) (this))`, then: `3`},
		// (default)
		{give: `[{"b":3},{"c":4},{"b":5}]`, when: `(collect b (default 0))`, then: `[3,0,5]`},
		// (size)
		{give: `{"a":3,"b":4}`, when: `(size)`, then: `2`},
		{give: `"abc"`, when: `(size)`, then: `3`},
		{give: `[3,4]`, when: `(size)`, then: `2`},
		// (merge)
		{give: `[{"a":3},{"b":4}]`, when: `(merge)`, then: `{"a":3,"b":4}`},
		// (flatten)
		{give: `{"a":{}}`, when: `(flatten -k a)`, then: `{}`},
		{give: `{"a":null}`, when: `(flatten -k a)`, then: `{"a":null}`},
		{give: `[{"a":{"b":3}},{"a":{"b":4}}]`, when: `(flatten -k a)`, then: `[{"b":3},{"b":4}]`},
		{give: `{"a":{"d":3,"e":4},"b":{"f":5},"c":{"g":6}}`, when: `(flatten -k a b)`, then: `{"d":3,"e":4,"f":5,"c":{"g":6}}`},
		{give: `{"a":{},"c":{"d":{}},"f":5,"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten)`, then: `{"f":5,"h":3,"l":7,"j":5}`},
		{give: `{"a":[{"c":{"d":3,"e":{"k":6}},"f":5}],"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten 0)`, then: `{"a":[{"c":{"d":3,"e":{"k":6}},"f":5}],"h":3,"l":7,"j":5}`},
		{give: `{"a":{"c":{"d":3,"e":{"k":6}},"f":5},"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten 3)`, then: `{"d":3,"k":6,"f":5,"h":3,"l":7,"j":5}`},
		{give: `{"a":{"c":{"d":3,"e":{"k":6}},"f":5},"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten 2)`, then: `{"d":3,"e":{"k":6},"f":5,"h":3,"i":{"l":7},"j":5}`},
		{give: `{"a":{"c":{"d":3,"e":{"k":6}},"f":5},"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten 1)`, then: `{"c":{"d":3,"e":{"k":6}},"f":5,"g":{"h":3,"i":{"l":7}},"j":5}`},
		{give: `{"a":{"c":{"d":3,"e":{"k":6}},"f":5},"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten 0)`, then: `{"d":3,"k":6,"f":5,"h":3,"l":7,"j":5}`},
		{give: `{"a":{"c":{"d":3,"e":{"k":6}},"f":5},"b":{"g":{"h":3,"i":{"l":7}},"j":5}}`, when: `(flatten)`, then: `{"d":3,"k":6,"f":5,"h":3,"l":7,"j":5}`},
		{give: `[ 3 , [ 4 , [ [ 5 ] ] ] , [ [ { "a" : [ [ 7 ] ] } ] ] , { "a" : [ [ 8 ] ] } ]`, when: `(flatten 0)`, then: `[3,4,5,{ "a" : [ [ 7 ] ] },{ "a" : [ [ 8 ] ] }]`},
		{give: `[[],3,[4,[5]],[[[]],7],8,[]]`, when: `(flatten 0)`, then: `[3,4,5,7,8]`},
		{give: `[3,[4,[5]],[[[6]],7],8]`, when: `(flatten 3)`, then: `[3,4,5,6,7,8]`},
		{give: `[3,[4,[5]],[[[6]],7],8]`, when: `(flatten 2)`, then: `[3,4,5,[6],7,8]`},
		{give: `[3,[4,[5]],[[[6]],7],8]`, when: `(flatten 1)`, then: `[3,4,[5],[[6]],7,8]`},
		{give: `[3,[4,[5]],[[[6]],7],8]`, when: `(flatten 0)`, then: `[3,4,5,6,7,8]`},
		{give: `[3,4]`, when: `(flatten 0)`, then: `[3,4]`},
		{give: `[[3]]`, when: `(flatten)`, then: `[3]`},
		{give: `[3]`, when: `(flatten)`, then: `3`},
		{give: `3`, when: `(flatten)`, then: `3`},
		// (first) (last)
		{give: `{"a":[]}`, when: `a (last)`, then: ``},
		{give: `{"a":[]}`, when: `(first a)`, then: ``},
		{give: `{"a":[3,4,5]}`, when: `a (last)`, then: `5`},
		{give: `{"a":[3,4,5]}`, when: `a (first)`, then: `3`},
		{give: `[]`, when: `(last)`, then: ``},
		{give: `[]`, when: `(first)`, then: ``},
		{give: `{"a":[{"b":{"c":3}},{"b":{}}]}`, when: `a (last (== 3 (get b c)) b c)`, then: `3`},
		{give: `{"a":[{"b":{"c":3}},{"b":{}}]}`, when: `a (first (== 3 (get b c)) b c)`, then: `3`},
		{give: `{"a":[{"b":3},{"b":4}]}`, when: `a (last) b`, then: `4`},
		{give: `{"a":[{"b":3},{"b":4}]}`, when: `a (first) b`, then: `3`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(last  (== 4 a) b)`, then: `3`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(first (== 4 a) b)`, then: `2`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(last  (== 4 a))`, then: `{"a":4,"b":3}`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(first (== 4 a))`, then: `{"a":4,"b":2}`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(last a)`, then: `5`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(first a)`, then: `3`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(last)`, then: `{"a":5,"b":4}`},
		{give: `[{"a":3,"b":1},{"a":4,"b":2},{"a":4,"b":3},{"a":5,"b":4}]`, when: `(first)`, then: `{"a":3,"b":1}`},
		// (arr)
		{give: `{"a":3,"b":4}`, when: `(arr a b a (raw "hi"))`, then: `[3,4,3,"hi"]`},
		{give: `{"a":3,"b":4}`, when: `(arr a b c)`, then: `[3,4]`},
		{give: `{"a":3,"b":4}`, when: `(arr a b a)`, then: `[3,4,3]`},
		{give: `{"a":3,"b":4}`, when: `(arr (get a) (get b) (get a))`, then: `[3,4,3]`},
		{give: ``, when: `(arr)`, then: `[]`},
		// (obj)
		{give: `{"a":3,"b":4}`, when: `(obj -each (concat (raw "key_") (key)) (this))`, then: `{"key_a":3,"key_b":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj -each (concat (raw "key_") (key)) (val))`, then: `{"key_a":3,"key_b":4}`},
		{give: `[{"a":3,"b":30},{"a":4,"b":40}]`, when: `(obj -each (get a) (val))`, then: `{"3":{"a":3,"b":30},"4":{"a":4,"b":40}}`},
		{give: `{"a":3}`, when: `(obj (raw 4) a)`, then: `{"4":3}`},
		{give: `{"a":3}`, when: `(obj "4" a)`, then: `{"4":3}`},
		{give: `{"a":3}`, when: `(obj 4 a)`, then: `{"4":3}`},
		{give: `{"a":3}`, when: `(obj (nothing) a)`, then: `{}`},
		{give: `{"a":3}`, when: `(obj x (nothing))`, then: `{}`},
		{give: `{"a":"aaa","b":"bbb"}`, when: `(obj (get a) (get b))`, then: `{"aaa":"bbb"}`},
		{give: `{"a":3,"b":4}`, when: `(obj "a b" a y b)`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj "a b" (get a) y (get b))`, then: `{"a b":3,"y":4}`},
		{give: `{"a":3,"b":4}`, when: `(obj x (get a) y (get b))`, then: `{"x":3,"y":4}`},
		{give: `{"a":{"b":{"c":3}}}`, when: `(get a b (obj x c))`, then: `{"x":3}`},
		// (get *)
		{give: `[{"a":[{"a":3},{"a":4}]},{"a":[{"a":5},{"a":6}]}]`, when: `(get * a (get * a) (flatten))`, then: `[3,4,5,6]`},
		{give: `[{"1":[{"1":3},{"1":4}]},{"1":[{"1":5},{"1":6}]}]`, when: `(get 1 1 1 1)`, then: `6`},
		{give: `[{"1":[{"1":3},{"1":4}]},{"1":[{"1":5},{"1":6}]}]`, when: `(get * 1 * 1)`, then: `[[3,4],[5,6]]`},
		{give: `[{"*":[{"*":3},{"*":4}]},{"*":[{"*":5},{"*":6}]}]`, when: `(get * "*" * "*")`, then: `[[3,4],[5,6]]`},
		{give: `[{"a":[{"a":3},{"a":4}]},{"a":[{"a":5},{"a":6}]}]`, when: `(get * a * a)`, then: `[[3,4],[5,6]]`},
		{
			give: `{"a":{"b":{"c":[{"d":"one","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":3}},{"k":{"l":4}}]}}}}]}},{"d":"two","e":{"f":[{"g":{"h":{"i":{"j":[{"k":{"l":5}}]}}}}]}}]}}}`,
			when: `(get a b c * (obj x d e (get e f * g h i j (get * k l) (flatten)))))`,
			then: `[{"x":"one","e":[3,4]},{"x":"two","e":[5]}]`,
		},
		{give: `[3,4,5]`, when: `(get *)`, then: `[3,4,5]`},
		{give: `[{"a":3}]`, when: `(get * a)`, then: `[3]`},
		{give: ``, when: `(get *)`, then: `[]`},
		{give: `3`, when: `(get *)`, then: `[]`},
		{give: `[]`, when: `(get * a)`, then: `[]`},
		{give: `[]`, when: `(get *)`, then: `[]`},
		{give: `[{"b":{"c":3}},{"b":{}}]`, when: `(get * b c)`, then: `[3]`},
		{give: `[{"a":3},{"b":4},{"a":5}]`, when: `(get * a)`, then: `[3,5]`},
		// (get) order should not matter.
		{give: `{"a":3,"b":4]}`, when: `(obj x b y a z a w b)`, then: `{"x":4,"y":3,"z":3,"w":4}`},
		// (get)
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(get a 1 c)`, then: `4`},
		{give: `{"a":[{"b":3},{"c":4}]}`, when: `(get a 0 b)`, then: `3`},
		{give: `[{"a":3},{"a":4}]`, when: `(get 0 a)`, then: `3`},
		{give: `[2,3,4]`, when: `(get 3)`, then: ``},
		{give: `[2,3,4]`, when: `(get 2)`, then: `4`},
		{give: `[2,3,4]`, when: `(get 1)`, then: `3`},
		{give: `[2,3,4]`, when: `(get 0)`, then: `2`},
		{give: `{"aa":{"bb":{"cc":3}}}`, when: `(get aa bb cc)`, then: `3`},
		{give: `{"a b":{"c d":3}}`, when: `(get "a b" "c d")`, then: `3`},
		{give: `{"a b":3}`, when: `(get "a b")`, then: `3`},
		{give: `{"a":{"b":3}}`, when: `(get a b)`, then: `3`},
		{give: `{"a":3}`, when: `(get a)`, then: `3`},
		// (root)
		{give: `3`, when: `(root)`, then: `3`},
		{give: ``, when: `(root)`, then: ``},
		// (raw)
		{give: ``, when: `(raw [ 3, 4 ])`, then: `[ 3, 4 ]`},
		{give: ``, when: `(raw { "a b": "c d" })`, then: `{ "a b": "c d" }`},
		{give: ``, when: `(raw {})`, then: `{}`},
		{give: ``, when: `(raw [])`, then: `[]`},
		{give: ``, when: `(raw null)`, then: `null`},
		{give: ``, when: `(raw true)`, then: `true`},
		{give: ``, when: `(raw false)`, then: `false`},
		{give: ``, when: `(raw 3e2)`, then: `3e2`},
		{give: ``, when: `(raw 3)`, then: `3`},
		{give: ``, when: `(raw "a b")`, then: `"a b"`},
		{give: ``, when: `(raw "a")`, then: `"a"`},
		{give: ``, when: `(raw -3)`, then: `-3`},
		{give: ``, when: `(raw 1.2)`, then: `1.2`},
	}
	for _, tc := range tt {
		r := Get(tc.give, tc.when)
		assertEqual(t, tc.then, r.String(), tc)
	}
}

func Test_Invalid_Query(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `[3,4]`, when: `(collect (this)`, then: `[3,4]`},
	}
	for _, tc := range tt {
		r := Get(tc.give, tc.when)
		assertEqual(t, tc.then, r.String(), tc)
	}
}

func TestGetWith(t *testing.T) {

	tt := []struct {
		give string
		when string
		args []any
		then string
	}{
		{
			give: ``,
			when: `(obj a (arg 0) b (arg 1) c (arg 2) d (arg 3) e (arg 4))`,
			args: []any{"Hello \"World\"", 4, float32(1.2), float64(1.5), []any{3, "4"}},
			then: `{"a":"Hello \"World\"","b":4,"c":1.2,"d":1.5,"e":[3,"4"]}`,
		},
	}
	for _, tc := range tt {
		r := GetWith(tc.give, tc.when, tc.args)
		assertEqual(t, tc.then, r.String(), tc)
	}
}

func BenchmarkGetWith(b *testing.B) {
	a := []any{3}
	for i := 0; i < b.N; i++ {
		_ = GetWith("", "(arg 0)", a)
	}
}

func TestJsonWS(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `  {"a":3,"b":4}`, when: `(pick b)`, then: `{"b":4}`},
		{give: `  {  "a"  :	3,  "b"  :  [  {  "c"  :  4  }  ]  }`, when: `(get b 0 c)`, then: `4`},
		{give: `[3,4 ]`, when: `(get 1)`, then: `4`},
		{give: `[3, 4]`, when: `(get 1)`, then: `4`},
		{give: `[3 ,4]`, when: `(get 1)`, then: `4`},
		{give: `[ 3,4]`, when: `(get 1)`, then: `4`},
		{give: `{"a":3, "b":4}`, when: `(get b)`, then: `4`},
		{give: `{"a":3 ,"b":4}`, when: `(get a)`, then: `3`},
		{give: `{"a": 3}`, when: `(get a)`, then: `3`},
		{give: `{"a" :3}`, when: `(get a)`, then: `3`},
		{give: `{ "a":3}`, when: `(get a)`, then: `3`},
	}
	for _, tc := range tt {
		r := Get(tc.give, tc.when)
		assertEqual(t, tc.then, r.String(), tc)
	}
}

func TestJsonBytes(t *testing.T) {

	tt := []struct {
		give string
		then []byte
	}{
		{give: `[3,4]`, then: []byte{'[', '3', ',', '4', ']'}},
	}
	for _, tc := range tt {
		r := Get(tc.give, `(this`).Bytes()
		assertEqual(t, tc.then, r, tc)
	}
}

func TestJsonGet(t *testing.T) {

	tt := []struct {
		give string
		when string
		then string
	}{
		{give: `{"a":2,"b":3}`, when: `a`, then: `2`},
		{give: `{"a":2,"b":3}`, when: `b`, then: `3`},
		{give: `{"a":2,"b":3}`, when: `c`, then: ``},
		{give: `{"a":[2,3]}`, when: `a`, then: `[2,3]`},
		{give: `[2,3]`, when: `0`, then: `2`},
		{give: `[2,3]`, when: `1`, then: `3`},
		{give: `[2,3]`, when: `2`, then: ``},
		{give: `{"a":{"b":2}}`, when: `a`, then: `{"b":2}`},
	}

	for _, tc := range tt {
		j := JSON(tc.give)
		r := j.Get(tc.when)
		assertEqual(t, tc.then, r.String(), tc)
	}
}

func BenchmarkJson_Get(b *testing.B) {
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Get("age")
	}
}

func TestJsonGet_Order(t *testing.T) {

	j := JSON(`{"a":3,"b":4}`)
	assertEqual(t, "4", j.Get("b").String())
	assertEqual(t, "3", j.Get("a").String())

	j = JSON(`[3,4]`)
	assertEqual(t, "4", j.Get("1").String())
	assertEqual(t, "3", j.Get("0").String())
}

func TestJsonForEachKeyVal(t *testing.T) {
	tt := []struct {
		give string
		then []string
	}{
		{give: `{}`, then: nil},
		{give: `{"a":2}`, then: []string{"a", "2"}},
		{give: `{"a":2,"b":3}`, then: []string{"a", "2", "b", "3"}},
		{give: `{"a":{"b":2}}`, then: []string{"a", `{"b":2}`}},
		{give: `{"a":[2]}`, then: []string{"a", "[2]"}},
	}
	for _, tc := range tt {
		var r []string
		j := JSON(tc.give)
		j.ForEachKeyVal(func(k, v Json) bool {
			r = append(r, k.Str(), v.String())
			return false
		})
		assertEqual(t, tc.then, r, tc.give)
	}
}

func TestJsonForEach(t *testing.T) {
	tt := []struct {
		give string
		then []string
	}{
		{give: "[]", then: nil},
		{give: "[10]", then: []string{"0", "10"}},
		{give: "[10,20]", then: []string{"0", "10", "1", "20"}},
		{give: "[10,20,30]", then: []string{"0", "10", "1", "20", "2", "30"}},
		{give: "[{},{},[10]]", then: []string{"0", "{}", "1", "{}", "2", "[10]"}},
	}
	for _, tc := range tt {
		var r []string
		j := JSON(tc.give)
		j.ForEach(func(i, v Json) bool {
			r = append(r, i.String(), v.String())
			return false
		})
		assertEqual(t, tc.then, r, tc.give)
	}
}

func TestJsonStr(t *testing.T) {
	tt := []struct {
		give string
		then string
	}{
		{give: ``, then: ``},
		{give: `""`, then: ``},
		{give: `"a"`, then: `a`},
		{give: `"\"a\""`, then: `"a"`},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Str(), tc.give)
	}
}

func TestJsonInt(t *testing.T) {
	tt := []struct {
		give string
		then int
	}{
		{give: ``, then: 0},
		{give: `0`, then: 0},
		{give: `1`, then: 1},
		{give: `-2`, then: -2},
		{give: `100`, then: 100},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Int(), tc.give)
	}
}

func TestJsonInt64(t *testing.T) {
	tt := []struct {
		give string
		then int64
	}{
		{give: ``, then: 0},
		{give: `0`, then: 0},
		{give: `1`, then: 1},
		{give: `-2`, then: -2},
		{give: `100`, then: 100},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Int64(), tc.give)
	}
}

func TestJsonUint64(t *testing.T) {
	tt := []struct {
		give string
		then uint64
	}{
		{give: ``, then: 0},
		{give: `0`, then: 0},
		{give: `1`, then: 1},
		{give: `100`, then: 100},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Uint64(), tc.give)
	}
}

func TestJsonFloat(t *testing.T) {
	tt := []struct {
		give string
		then float64
	}{
		{give: ``, then: 0},
		{give: `0`, then: 0},
		{give: `1`, then: 1},
		{give: `-2.2`, then: -2.2},
		{give: `1e2`, then: 100},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Float(), tc.give)
	}
}

func TestJsonBool(t *testing.T) {
	tt := []struct {
		give string
		when bool
	}{
		{give: ``, when: false},
		{give: `false`, when: false},
		{give: `true`, when: true},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.when, j.Bool(), tc.give)
	}
}

func TestJsonIsEmpty(t *testing.T) {
	tt := []struct {
		give string
		when bool
	}{
		{give: ``, when: false},
		{give: `{}`, when: true},
		{give: `[]`, when: true},
		{give: `""`, when: true},
		{give: `null`, when: false},
		{give: `0`, when: false},
		{give: `false`, when: false},
		{give: `true`, when: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.when, j.IsEmpty(), tc.give)
	}
}

func TestJsonIsEmptyString(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `""`, then: true},
		{give: `"a"`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsEmptyString(), tc.give)
	}
}

func TestJsonIsEmptyObject(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `{   }`, then: true},
		{give: `[]`, then: false},
		{give: `""`, then: false},
		{give: `3`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsEmptyObject(), tc.give)
	}
}

func TestJsonIsEmptyArray(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `[]`, then: true},
		{give: `[  ]`, then: true},
		{give: `{}`, then: false},
		{give: `""`, then: false},
		{give: `3`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsEmptyArray(), tc.give)
	}
}

func TestJsonIsFalsy(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `[]`, then: true},
		{give: `{}`, then: true},
		{give: `""`, then: true},
		{give: `false`, then: true},
		{give: `null`, then: true},
		{give: `0`, then: true},
		//
		{give: `3`, then: false},
		{give: `[0]`, then: false},
		{give: `{"a":0}`, then: false},
		{give: `"a"`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsFalsy(), tc)
	}
}

func TestJsonIsTruthy(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `[]`, then: false},
		{give: `{}`, then: false},
		{give: `""`, then: false},
		{give: `false`, then: false},
		{give: `null`, then: false},
		{give: `0`, then: false},
		//
		{give: `3`, then: true},
		{give: `[0]`, then: true},
		{give: `{"a":0}`, then: true},
		{give: `"a"`, then: true},
		{give: `true`, then: true},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsTruthy(), tc)
	}
}

func TestJsonIsVoid(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: false},
		{give: `null`, then: false},
		{give: `0`, then: false},
		{give: `false`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsVoid(), tc.give)
	}
}

func TestJsonIsBlank(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: false},
		{give: `null`, then: true},
		{give: `0`, then: false},
		{give: `false`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsBlank(), tc.give)
	}
}

func TestJsonIsNully(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: true},
		{give: `null`, then: true},
		{give: `0`, then: false},
		{give: `false`, then: false},
		{give: `true`, then: false},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsNully(), tc.give)
	}
}

func TestJsonIsWithExists(t *testing.T) {
	j := JSON("")
	tt := []struct {
		name string
		when func() bool
		then bool
	}{
		{name: "IsArray", when: j.IsArray},
		{name: "IsObject", when: j.IsObject},
		{name: "IsEmptyArray", when: j.IsEmptyArray},
		{name: "IsEmptyObject", when: j.IsEmptyObject},
		{name: "IsEmptyString", when: j.IsEmptyString},
		{name: "IsEmpty", when: j.IsEmpty},
		{name: "IsBlank", when: j.IsBlank},
		{name: "IsNull", when: j.IsNull},
		{name: "IsNully", when: j.IsNully},
		{name: "IsBool", when: j.IsBool},
		{name: "IsTrue", when: j.IsTrue},
		{name: "IsTruthy", when: j.IsTruthy},
		{name: "IsFalse", when: j.IsFalse},
		{name: "IsFalsy", when: j.IsFalsy},
		{name: "IsNumber", when: j.IsNumber},
		{name: "IsString", when: j.IsString},
		{name: "IsSome", when: j.IsSome},
		{name: "IsVoid", when: j.IsVoid},
	}
	for _, tc := range tt {
		assertEqual(t, false, tc.when(), tc.name)
	}
}

func BenchmarkJsonExists(b *testing.B) {
	j := JSON("")
	for i := 0; i < b.N; i++ {
		j.Exists()
	}
}

func TestJsonIsSome(t *testing.T) {
	tt := []struct {
		give string
		then bool
	}{
		{give: ``, then: false},
		{give: `{}`, then: true},
		{give: `[]`, then: true},
		{give: `""`, then: true},
		{give: `null`, then: false},
		{give: `0`, then: true},
		{give: `false`, then: true},
		{give: `true`, then: true},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.IsSome(), tc)
	}
}

func TestJsonStringify(t *testing.T) {
	tt := []struct {
		give string
		then string
	}{
		{give: ``, then: `""`},
		{give: `{}`, then: `"{}"`},
		{give: `{ "hello": "wo\"rld" }`, then: `"{ \"hello\": \"wo\\\"rld\" }"`},
		{give: `[]`, then: `"[]"`},
		{give: `123`, then: `"123"`},
		{give: `null`, then: `"null"`},
		{give: `false`, then: `"false"`},
		{give: `true`, then: `"true"`},
		{give: `""`, then: `"\"\""`},
		{give: `"a b"`, then: `"\"a b\""`},
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Stringify().String(), tc)
	}
}

func BenchmarkJsonToStringify(b *testing.B) {
	j := JSON(`{ "hello": "wo\"rld" }`)
	for i := 0; i < b.N; i++ {
		_ = j.Stringify()
	}
}

func TestJsonJsonify(t *testing.T) {
	tt := []struct {
		give string
		then string
	}{
		{give: ``, then: ``},
		{give: `""`, then: ``},
		{give: `"{}"`, then: `{}`},
		{give: `"{ \"hello\": \"wo\\\"rld\" }"`, then: `{ "hello": "wo\"rld" }`},
		{give: `"[]"`, then: `[]`},
		{give: `"123"`, then: `123`},
		{give: `"null"`, then: `null`},
		{give: `"false"`, then: `false`},
		{give: `"true"`, then: `true`},
		{give: `"a b"`, then: `a b`}, // Invalid JSON.
	}
	for _, tc := range tt {
		j := JSON(tc.give)
		assertEqual(t, tc.then, j.Jsonify().String(), tc)
	}
}

func BenchmarkJsonJsonify(b *testing.B) {
	j := JSON(`"{ \"hello\": \"wo\\\"rld\" }"`)
	for i := 0; i < b.N; i++ {
		_ = j.Jsonify()
	}
}

func Example_funcDebug() {

	j := `[{ "a": { "b": [3] } }, { "a": { "b": [4] } }]`

	v := Get(j, `(collect a (debug) b (debug b_val) (flatten) (debug flatd))`)

	fmt.Println("Result:", v)

	// Output:
	// [debug] { "b": [3] }
	// [b_val] [3]
	// [flatd] 3
	// [debug] { "b": [4] }
	// [b_val] [4]
	// [flatd] 4
	// Result: [3,4]
}

func ExampleJson_ForEachKeyVal() {
	m := func(k, v Json) bool {
		fmt.Println(k, v)
		return false
	}

	j := JSON(TestData1)
	j.ForEachKeyVal(m)

	// Output:
	// "name" "Mary"
	// "last" "Jane"
	// "token" null
	// "settings" {}
	// "posts" []
	// "address" {"city":"Place","country":"USA"}
	// "contacts" [{"name":"Karen"},{"name":"Michelle","last":"Jane"}]
	// "age" 33
	// "random" [3,null,{},[],"",false]
}

func BenchmarkJson_ForEachKeyVal(b *testing.B) {
	m := func(k, v Json) bool {
		return false
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		j.ForEachKeyVal(m)
	}
}

func ExampleJson_ForEach() {
	m := func(k, v Json) bool {
		fmt.Println(k, v)
		return false
	}

	j := JSON(TestData2)
	j.ForEach(m)

	// Output:
	// 0 {"name":"Karen"}
	// 1 {"name":"Michelle","last":"Jane"}
}

func BenchmarkJson_ForEach(b *testing.B) {
	m := func(k, v Json) bool {
		return false
	}
	j := JSON(TestData2)
	for i := 0; i < b.N; i++ {
		j.ForEach(m)
	}
}

func ExampleJson_IterateFast() {

	j := JSON(`{"a":1,"b":2,"c":{"a":3,"b":{"a":4,"b":[{"a":5},{"a":6,"b":7,"c":[8,9,0,{},[]]}]}},"d":1}`)
	v := j.IterateFast(func(k, v Json) (Json, Json) {
		k = JSON(strings.ToUpper(k.String()))
		if v.IsNumber() {
			return k, v.Stringify()
		}
		return k, v
	})

	fmt.Println(v)

	// Output:
	// {"A":"1","B":"2","C":{"A":"3","B":{"A":"4","B":[{"A":"5"},{"A":"6","B":"7","C":["8","9","0",{},[]]}]}},"D":"1"}
}

func BenchmarkJson_Iterate(b *testing.B) {
	m := func(k, v Json) (Json, Json) {
		return k, v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateFast(m)
	}
}

func Benchmark_QueryFunction_Iterate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(iterate (this) (this))`)
	}
}

func ExampleJson_IterateKeysValues() {
	m := func(v Json) Json {
		fmt.Println(v)
		return v
	}

	j := JSON(TestData1)
	r := j.IterateKeysValues(m)

	fmt.Println(r)

	// Output:
	// "name"
	// "Mary"
	// "last"
	// "Jane"
	// "token"
	// null
	// "settings"
	// "posts"
	// "address"
	// "city"
	// "Place"
	// "country"
	// "USA"
	// "contacts"
	// "name"
	// "Karen"
	// "name"
	// "Michelle"
	// "last"
	// "Jane"
	// "age"
	// 33
	// "random"
	// 3
	// null
	// ""
	// false
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_IterateKeysValues(b *testing.B) {
	m := func(v Json) Json {
		return v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateKeysValues(m)
	}
}

func Benchmark_QueryFunction_IterateKV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(iterate -kv (this) (this))`)
	}
}

func ExampleJson_IterateKeys() {
	m := func(v Json) Json {
		fmt.Println(v)
		return v
	}

	j := JSON(TestData1)
	r := j.IterateKeys(m)

	fmt.Println(r)

	// Output:
	// "name"
	// "last"
	// "token"
	// "settings"
	// "posts"
	// "address"
	// "city"
	// "country"
	// "contacts"
	// "name"
	// "name"
	// "last"
	// "age"
	// "random"
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_IterateKeys(b *testing.B) {
	m := func(v Json) Json {
		return v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateKeys(m)
	}
}

func Benchmark_QueryFunction_IterateK(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(iterate -k (this))`)
	}
}

func ExampleJson_IterateValues() {
	m := func(v Json) Json {
		fmt.Println(v)
		return v
	}

	j := JSON(TestData1)
	r := j.IterateValues(m)

	fmt.Println(r)

	// Output:
	// "Mary"
	// "Jane"
	// null
	// "Place"
	// "USA"
	// "Karen"
	// "Michelle"
	// "Jane"
	// 33
	// 3
	// null
	// ""
	// false
	// {"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}
}

func BenchmarkJson_IterateValues(b *testing.B) {
	m := func(v Json) Json {
		return v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.IterateValues(m)
	}
}

func Benchmark_QueryFunction_IterateV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(iterate -v (this))`)
	}
}

func ExampleJson_Iterate() {

	m := func(k, v Json) (Json, Json) {
		fmt.Println(k, v)
		return k, v
	}

	j := JSON(`{ "a": 3, "b": { "c": 4 }, "d": [ 5, { "e": 6 } ] }`)
	v := j.Iterate(0, m)

	fmt.Println("Final:", v)

	// Output:
	// "a" 3
	// "c" 4
	// "b" {"c":4}
	// 0 5
	// "e" 6
	// 1 {"e":6}
	// "d" [5,{"e":6}]
	// null {"a":3,"b":{"c":4},"d":[5,{"e":6}]}
	// Final: {"a":3,"b":{"c":4},"d":[5,{"e":6}]}
}

func BenchmarkJson_IterateAll(b *testing.B) {
	m := func(k, v Json) (Json, Json) {
		return k, v
	}
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		j.Iterate(0, m)
	}
}

func Benchmark_QueryFunction_IterateAll(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(iterate (this) (this))`)
	}
}

func ExampleJson_Flatten() {

	j := JSON(`[3,[4,[[5]]],[[{"a":6},{"a":[7,[8]]}]],{"a":[[9]]}]`)
	j = j.Flatten(0)

	fmt.Println(j)

	// Output:
	// [3,4,5,{"a":6},{"a":[7,[8]]},{"a":[[9]]}]
}

func ExampleJson_Uglify() {

	j := JSON(`{ "a": "b" }`)
	j = j.Uglify()

	fmt.Println(j)

	// Output:
	// {"a":"b"}
}

func BenchmarkJson_Uglify(b *testing.B) {
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Uglify()
	}
}

func ExampleJson_Prettify() {

	j := JSON(`{ "a": "b" }`)
	j = j.Prettify()

	fmt.Println(j)

	// Output:
	// {
	//     "a": "b"
	// }
}

func BenchmarkJson_Prettify(b *testing.B) {
	j := JSON(TestData1)
	for i := 0; i < b.N; i++ {
		_ = j.Prettify()
	}
}

func Benchmark_QueryFunction_Sort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`[{ "a": 5 }, { "a": 4 }, { "a": 3 }]`, `(sort a)`)
	}
}

func Benchmark_QueryFunction_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(set address city "xxx")`)
	}
}

func Benchmark_QueryFunction_Pick(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(pick name age)`)
	}
}

func Benchmark_QueryFunction_Pluck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(pluck address contacts)`)
	}
}

func Benchmark_QueryFunction_Keys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(keys)`)
	}
}

func Benchmark_QueryFunction_Values(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(values)`)
	}
}

func Benchmark_QueryFunction_Entries(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(TestData1, `(entries)`)
	}
}

func Benchmark_QueryFunction_Objectify(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`[["a","Hello"],["b","World"],["count",123]]`, `(objectify)`)
	}
}

func Benchmark_QueryFunction_Merge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`[{"a","Hello"},{"b","World"},{"count",123}]`, `(merge)`)
	}
}

func Benchmark_QueryFunction_Arr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(``, `(arr (raw 3) (raw 4) (raw 5))`)
	}
}

func Benchmark_QueryFunction_Obj(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(``, `(obj a (raw 3) b (raw 4) c (raw 5))`)
	}
}

func Benchmark_QueryFunction_Collect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`[0,1,2,3,4,5,6,7,8,9]`, `(collect (this))`)
	}
}

func Benchmark_QueryFunction_Concat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(``, `(concat (raw "Hello") (raw "World"))`)
	}
}

func Benchmark_QueryFunction_FlattenDeep(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get(`[3,[4],[[5]]]`, `(flatten 0)`)
	}
}

const TestData1 = `{"name":"Mary","last":"Jane","token":null,"settings":{},"posts":[],"address":{"city":"Place","country":"USA"},"contacts":[{"name":"Karen"},{"name":"Michelle","last":"Jane"}],"age":33,"random":[3,null,{},[],"",false]}`
const TestData2 = `[{"name":"Karen"},{"name":"Michelle","last":"Jane"}]`

func assertEqual(t *testing.T, exp, got any, msgs ...any) {
	t.Helper()
	if !reflect.DeepEqual(exp, got) {
		t.Errorf("\nExp:\n%v\nGot:\n%v\nMsg: %v", exp, got, fmt.Sprint(msgs...))
	}
}
