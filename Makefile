ANTLR4=java -jar antlr-4.8-complete.jar

parser: Matcher.g4
	$(ANTLR4) -Dlanguage=Go -visitor -o parser Matcher.g4

clean:
	rm -rf parser
