parser: Matcher.g4
	antlr4 -Dlanguage=Go -visitor -o parser Matcher.g4

clean:
	rm -rf parser
