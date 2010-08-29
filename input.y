%package main
%import scanner fmt os strconv

%union {
	fval float
	op func(float)float
}

%token<fval> floating integer
%type<fval> Calc Num Mult Add
%type<op> MultA AddA

%%

Calc : Add				{ fmt.Println("0. Found Calc->Num Mult. Calculated", $$); $$ = $1 }
	 ;
	
Mult : Num MultA		{ fmt.Println("1. Found Mult->'*' Num."); $$ = $2($1) }
	 ;
	
MultA : '*' Mult		{ fmt.Println("2. Found MultA->'*' Mult."); $$ = mult($2) }
	  | '/' Mult		{ fmt.Println("3. Found MultA->'/' Mult."); $$ = div($2) }
	  |					{ fmt.Println("4. Found MultA->{}."); $$ = noop }
	  ;

Add : Mult AddA			{ fmt.Println("5. Found Add->Mult AddA."); $$ = $2($1) }
	;

AddA : '+' Add			{ fmt.Println("6. Found AddA->'+' Add"); $$ = plus($2) }
	 | '-' Add			{ fmt.Println("7. Found AddA->'-' Add"); $$ = minus($2) }
	 |					{ fmt.Println("8. Found AddA->{}"); $$ = noop}
	 ;

Num : floating		{ fmt.Println("9. Found Num->floating. Forwarding value", $1); $$ = $1 }
	| integer		{ fmt.Println("10. Found Num->integer. Forwarding value", $1); $$ = float($1) }
	;

%%

const (
	EOF = USER + 1
)

func mult(m float) func(float)float {
	return func(f float)float {
		return f * m
	}
}

func div(m float) func(float)float {
	return func(f float)float {
		return f / m
	}
}

func plus(m float) func(float)float {
	return func(f float)float {
		return f + m
	}
}

func minus(m float) func(float)float {
	return func(f float)float {
		return f - m
	}
}

func noop(m float) float {
	return m
}

func main() {
	reader := os.Stdin
	for true {
		var s scanner.Scanner
		s.Init(reader)
		nextWord := func(v *yystype)int {
			i := s.Scan()
			switch i {
			case scanner.Float:
				v.fval, _ = strconv.Atof(s.TokenText())
				return floating
			case scanner.Int:
				i, _ := strconv.Atoi(s.TokenText())
				v.fval = float(i)
				return integer
			case scanner.Ident:
				if s.TokenText() == "eof" {return EOF}
				return -1
			case scanner.String:
				return -1
			case scanner.Char:
				return -1
			case scanner.RawString:
				return -1
			case scanner.Comment:
				return -1
			case scanner.EOF:
				return EOF
			default:
				return i
			}
			return -1
		}
		if parse(EOF, nextWord) {
			fmt.Println("Result: ", res[0].(*yystype).fval)
		} else {
			fmt.Println("Can't parse that, dude.")
		}
	}
}
