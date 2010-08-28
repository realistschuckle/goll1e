%package main
%import scanner fmt os

%%

Goal : one						{ fmt.Println("ONE!") }
     | ':' kill					{
									fmt.Println("BANG! kill")
								}
     | '!' Another
	 | ' ' bill
	 | Another AfterAnother
	 | Chair moo
     ;

AfterAnother : StillAnother
			 | moo Something wish
			 ;

Chair : Something Another StillAnother
      ;

Something : kill
		  | bill
		  ;

Another : fighter
		| 
		| myka
		;
		
StillAnother : pillow
			 | fool
			 |
			 ;

%%

const (
	EOF = USER + 1
)

func main() {
	reader := os.Stdin
	for true {
		var s scanner.Scanner
		s.Init(reader)
		nextWord := func()int {
			i := s.Scan()
			switch i {
			case scanner.Ident:
				switch s.TokenText() {
				case "kill":
					return kill
				case "wish":
					return wish
				case "fighter":
					return fighter
				case "moo":
					return moo
				case "one":
					return one
				case "myka":
					return myka
				case "pillow":
					return pillow
				case "fool":
					return fool
				case "bill":
					return bill
				case "eof":
					return EOF
				}
			case scanner.EOF:
				return EOF
			default:
				return i
			}
			return -1
		}
		fmt.Println(parse(EOF, nextWord))
	}

}
