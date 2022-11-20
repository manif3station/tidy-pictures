#!/usr/bin/env perl

use strict;
use warnings;
use lib '/app/lib';

use DateTime;
use Digest::MD5;
use Image::ExifTool ':Public';
use System;

sub _date {
    my ($date) = @_;
    return if !$date;
    my @date = split /\D/, $date;
    my %date = ();
    @date{qw(year month day hour minute second)} = @date;
    return DateTime->new(%date);
}

sub _mkdir {
    my ($dir) = @_;
    System->RunQW( mkdir => -p => $dir ) if !-d $dir;
    return $dir;
}

sub _move_file {
    my ( $from, $to ) = @_;

    next if $from eq $to;

    System->RunQW( mv => -v => $from, $to );
}

sub main {
    my ( $from, $to ) = @ARGV;

    $from //= '.';
    $to   //= '.';

    die "From where?" if !-d $from;
    die "To where?"   if !-d $to;

    my @files = split /\n/, System->RunQW( find => $from, -type => 'f' );

    my %seen = ();

    foreach my $from_file (@files) {
        next if $from_file =~ /DS_Store/;

        if ( !-s $from_file ) {
            _move_file $from_file, _mkdir "$to/Empty-Files/";
            next;
        }

        my $info = ImageInfo($from_file);

        my $mime = $info->{MIMEType} // '';

        if ( !%$info || $mime !~ m/(image|video)/ ) {
            _move_file $from_file, _mkdir "$to/Non-A-Picture/";
            next;
        }

        my $md5 = do {
            open my $fh, '<', $from_file;
            my $ctx = Digest::MD5->new;
            $ctx->addfile($fh);
            $ctx->hexdigest;
        };

        if ( $seen{$md5}++ ) {
            _move_file $from_file, _mkdir "$to/Duplicated-Files/";
            next;
        }

        my $date = $info->{CreateDate};# // $info->{FileModifyDate};

        $date = _date($date);

        if ( !$date ) {
            _move_file $from_file, _mkdir "$to/Files-Have-No-Date/";
            next;
        }

        my $dir = _mkdir sprintf "$to/%04d/%02d", $date->year, $date->month;

        $date =~ s/\D/-/g;
        $date =~ s/T/-/g;

        my $ext = $info->{FileTypeExtension};

        ($ext) = ( $from_file =~ m/([^\.]+)$/ ) if !$ext;

        my $to_file = sprintf '%s/%s.%s', $dir, $date, $ext;

        my $first_file = $to_file;

        my $next_id = 0;

        while ( -f $to_file ) {
            $to_file = sprintf '%s/%s-%03d.%s', $dir, $date, $next_id++, $ext;
        }

        if ( $next_id == 1 ) {
            my $to_file = sprintf '%s/%s-%03d.%s', $dir, $date, 0, $ext;
            _move_file $first_file, $to_file;
        }

        _move_file $from_file, $to_file;
    }
}

main();
