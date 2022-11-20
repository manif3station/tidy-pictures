package System;

use strict;
use warnings;

use Capture::Tiny 'capture';

sub Run {
    my (undef, $c) = @_;
    print "> $c\n" if $ENV{VERBOSE};
    capture { system $c };
}

sub RunQW {
    my (undef, @parts) = @_;
    Run undef, join ' ', map { quotemeta } @parts;
}

sub Exec {
    my (undef, $c) = @_;
    print "> $c\n" if $ENV{VERBOSE};
    exec $c;
}

1;
