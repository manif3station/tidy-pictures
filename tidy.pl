#!/usr/bin/env perl

use strict;
use warnings;
use FindBin '$Bin';
use lib ( '/app/lib', $Bin );

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
    my $dt = eval { DateTime->new(%date) };
    return $dt;
}

sub _mkdir {
    my ($dir) = @_;
    System->RunQW( mkdir => -p => $dir ) if !-d $dir;
    return $dir;
}

sub _move_file {
    my ( $from, $to ) = @_;

    next if $from eq $to;

    if ( -d $to ) {
        my ($filename) = ( $from =~ m {.*/(.*)} );
        $to .= "/$1";
    }

    if ( -f $to ) {
        my $num = 1;
        my $_to;
        do {
            $_to = $to;
            $num++;
            if ( $_to =~ m/(.+)\.(.+)/ ) {
                $_to = "$1 (Copy $num).$2";
            }
            else {
                $_to .= " (Copy $num)";
            }
        } while ( -f $_to );
        $to = $_to;
    }

    System->RunQW( mv => -v => $from, $to );
}

my $_1m = 60;
my $_1h = 60 * $_1m;

sub _timer {
    my ($started) = @_;
    my $now       = time;
    my $elsped    = ( $started - $now ) / $_1h;

    my ( $hours, $minutes_10 ) = split /\n/, $elsped;

    my $h_display = sprintf '%02d', $hours // 0;

    $minutes_10 //= 0;

    my $minutes_60 = 60 * "0.$minutes_10";

    my $m_display = sprintf '%02d', $minutes_60 // 0;

    return "$h_display:$m_display";
}

sub _progress {
    my ( $processed, $total ) = @_;
    return ( ( $processed / $total ) * 100, $processed, $total );
}

sub _filename {
    my ($path)     = @_;
    my ($filename) = ( $path =~ m{.*/(.*)} );
    return $filename;
}

sub main {
    my $from = $ENV{FROM_LOCATION} // '/pictures';
    my $to   = $ENV{TO_LOCATION}   // '/pictures';

    my @files = sort { _filename($a) cmp _filename($b) } split /\n/,
        System->RunQW( find => $from, -type => 'f' );

    my @old_files = sort { _filename($a) cmp _filename($b) } split /\n/,
        System->RunQW( find => $to, -type => 'f' );

    my $started   = time;
    my $total     = scalar @files;
    my $processed = 0;

    my %seen = ();

    $| = 1;

    printf "Indexing ...\n";

    FILE: foreach my $old_file(@old_files) {
        open my $fh, '<', $old_file;
        my $ctx = Digest::MD5->new;
        $ctx->addfile($fh);
        $seen{$ctx->hexdigest}++;
    }

    printf "[Done]\n\nStarted @ %s\n\n", DateTime->now;

    FILE: foreach my $from_file (@files) {
        my $rework_count = sub {
            $total--;
            $processed--;
        };

        $processed++;

        if ($from_file =~ /DS_Store/) {
            $rework_count->();
            next;
        }

        printf "%s In Progress %.01f%% %d of %d\n", _timer($started),
            _progress( $processed, $total );

        if ( !-s $from_file ) {
            _move_file $from_file, _mkdir "$to/Empty-Files/";
            next;
        }

        my $info = ImageInfo($from_file);

        my $mime = $info->{MIMEType} // '';

        if ( !%$info || $mime !~ m/(image|video)/ ) {
            _move_file $from_file, _mkdir "$to/Non-Picture/";
            next;
        }

        my $md5 = do {
            open my $fh, '<', $from_file;
            my $ctx = Digest::MD5->new;
            $ctx->addfile($fh);
            $ctx->hexdigest;
        };

        my $date = $info->{CreateDate} // $info->{FileModifyDate};

        $date = _date($date);

        if ( !$date ) {
            _move_file $from_file, _mkdir "$to/Files-Have-No-Date/";
            next;
        }

        my $dir;

        if ($seen{$md5}++) {
            $dir = _mkdir "$to/Duplicated-Files";
        }
        else {
            $dir = _mkdir sprintf "$to/%04d/%02d", $date->year, $date->month;
        }

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

    map {
        system sprintf "cd %s; find . -name .DS_Store -exec rm {} \\;; find . -empty -type d -delete", quotemeta;
    } ( $from, $to );

    printf "\nDone @ %s\n", DateTime->now;
}

main();
