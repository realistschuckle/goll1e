%package main     # Set the package of the generated file to "main"

%import scanner fmt os strconv  # import (
                                #   scanner
                                #   fmt
                                #   os
                                #   strconv
                                # )

# Replace the default { $$ = $1 } rule code with this custom code.
%defaultcode {
    fmt.Println("Default code. Assigning", $1, " to ", $$, "."); $$ = $1
}

# Define the custom value type for tokens
%union {
    fval float
    ival int
    op func(float)float
}

# Associate the "floating" terminal with the type of fval float
%token<fval> floating

# Associate the "integer" terminal with the type of ival int
%token<ival> integer

# Associate the "Calc", "Num", "Mult", and "Add" nonterminals with the type of fval float 
%type<fval> Calc Num Mult Add

# Associate the "MultA" and "AddA" nonterminals with the type of op func(float)float
%type<op> MultA AddA

%%

Calc : Add        # This will use the code in %defaultcode
     ;

Mult : Num MultA        { fmt.Println("1. Found Mult->'*' Num."); $$ = $2($1) }
     ;
    
MultA : '*' Mult        { fmt.Println("2. Found MultA->'*' Mult."); $$ = mult($2) }
      | '/' Mult        { fmt.Println("3. Found MultA->'/' Mult."); $$ = div($2) }
      |                 { fmt.Println("4. Found MultA->{}."); $$ = noop }
      ;

Add : Mult AddA         { fmt.Println("5. Found Add->Mult AddA."); $$ = $2($1) }
    ;

AddA : '+' Add          { fmt.Println("6. Found AddA->'+' Add"); $$ = plus($2) }
     | '-' Add          { fmt.Println("7. Found AddA->'-' Add"); $$ = minus($2) }
     |                  { fmt.Println("8. Found AddA->{}"); $$ = noop}
     ;

Num : floating          { fmt.Println("9. Found Num->floating. Forwarding value", $1); $$ = float($1) }
    | integer           { fmt.Println("10. Found Num->integer. Forwarding value", $1); $$ = float($1) }
    ;

%%

// Define the end-of-file token used by the scanner
const (
    EOF = yyUSER + 1
)

// Define functions used by the grammar above
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

// Entry point for executable
func main() {
    reader := os.Stdin
    for true {
        var s scanner.Scanner
        s.Init(reader)
        
        // Define a scanner function for the yyparse function
        nextWord := func(v *yystype)int {
            i := s.Scan()
            switch i {
            case scanner.Float:
                // Set the value of the string conversion to the float slot
                v.fval, _ = strconv.Atof(s.TokenText())
                return floating
            case scanner.Int:
                // Set the value of the string conversion to the int slot
                v.ival, _ = strconv.Atoi(s.TokenText())
                return integer
            case scanner.Ident:
                // Return EOF for the token "eof"
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
        
        // Print the result if the parser recognized the input
        // Otherwise, print a colloquial but unhelpful message
        if ok, result := yyparse(EOF, nextWord); ok {
            fmt.Println("Result:", result)
        } else {
            fmt.Println("Can't parse that, dude.")
        }
    }
}
