%package main
%import scanner fmt os strconv

%union {
	fval float
	op func(float)float
}

%token<fval> floating integer
%type<fval> Calc Num
%type<op> Mult Add

%%

Calc : Num Mult		{ fmt.Println("0. Found Calc->Num Mult. Calculated", $2($1)); $$ = $2($1) }
	 ;
	
Mult : '*' Calc		{ fmt.Println("1. Found Mult->'*' Calc."); $$ = mult($2) }
	 | '/' Calc		{ fmt.Println("2. Found Mult->'/' Calc."); $$ = div($2) }
	 | Add			{ fmt.Println("3. Found Mult->Add."); $$ = $1 }
	 ;

Add : '+' Calc		{ fmt.Println("4. Found Add->'+' Calc."); $$ = plus($2) }
    | '-' Calc		{ fmt.Println("5. Found Add->'-' Calc."); $$ = minus($2) }
	|				{ fmt.Println("6. Found Add->{}."); $$ = noop }
	;

Num : floating		{ fmt.Println("7. Found Num->floating. Forwarding value", $1); $$ = $1 }
	| integer		{ fmt.Println("8. Found Num->integer. Forwarding value", $1); $$ = float($1) }
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
