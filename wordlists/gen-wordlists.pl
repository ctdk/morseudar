#!/usr/bin/perl

# TODO: This could probably stand to be a little go program rather than perl.

print "package wordlists\n";
print "import (\n)\n";
print "var $ARGV[0]Words = []string{";

while (my $l = <STDIN>) {
	chomp $l;
	print qq|"$l",\n|;
}

print "}\n";
