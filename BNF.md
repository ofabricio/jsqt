# Query BNF

```
fun =
	'(' ')'
	'(' fname farg+ ')'

farg =
	' ' ( fun | str | raw )

fname =
	anything

raw =
	anything
```
