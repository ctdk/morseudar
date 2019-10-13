#!/usr/bin/perl

print "package topwords\n";
print "import (\n)\n";
print "var topWds = []string{";

while (my $l = <STDIN>) {
	chomp $l;
	print qq|"$l",\n|;
}

print "}\n";
