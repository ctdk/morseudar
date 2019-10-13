#!/usr/bin/perl

print "package wordlists\n";
print "import (\n)\n";
print "var $ARGV[0]Words = []string{";

while (my $l = <STDIN>) {
	chomp $l;
	print qq|"$l",\n|;
}

print "}\n";
