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

sub _progress {
    my ( $processed, $total ) = @_;
    return ( ( $processed / $total ) * 100, $processed, $total );
}

sub _filename {
    my ($path)     = @_;
    my ($filename) = ( $path =~ m{.*/(.*)} );
    return $filename;
}

sub _md5 {
    my ($path) = @_;
    open my $fh, '<', $path;
    my $ctx = Digest::MD5->new;
    $ctx->addfile($fh);
    $ctx->hexdigest;
}

my %SEEN;

sub _add_index {
    my ($file, $index_dir) = @_;

    _mkdir $index_dir;

    if ( !-f $file ) {
        die "Unexpected Error: File not found $file\n.";
    }

    my $md5 = $SEEN{$file} //= _md5 $file;

    my $index = "$index_dir/$md5";

    return "seen" if -f $index;

    open my $fh, '>', $index;
    print $fh time;
    close $fh;

    return "new";
}

sub _check_update {
    return if grep {/--skip-check-update/} @ARGV;

    if ( !qx{which curl} ) {
        system qw(apt update);
        system qw(sudo apt install -y curl);
    }

    return if !qx{which curl};

    system
        "curl https://raw.githubusercontent.com/manif3station/tidy-pictures/stable/tidy.pl > $0.new";

    return if !-f "$0.new" || !-s "$0.new" || _md5("$0.new") eq _md5($0);

    system mv => -v => "$0.new", $0;

    exec perl => $0, '--skip-check-update';
}

sub main {
    _check_update;

    printf "\nStarted @ %s\n", DateTime->now;

    my $from = $ENV{FROM_LOCATION} // '/pictures';
    my $to   = $ENV{TO_LOCATION}   // '/pictures';
    my $dup  = _mkdir "$to/Duplicated-Files";

    my $index_dir = "$to/.seen-pictures";

    my $normal_idx = "$index_dir/normal";
    my $dup_idx    = "$index_dir/duplicated";

    my @files = sort { _filename($a) cmp _filename($b) } split /\n/,
        System->RunQW( find => $from, -type => 'f' );

    my @old_files = sort { _filename($a) cmp _filename($b) } split /\n/,
        System->RunQW( find => $to, -path => $dup, '-prune', -type => 'f' );

    my @dup_files = sort { _filename($a) cmp _filename($b) } split /\n/,
        System->RunQW( find => $dup, -type => 'f');

    my $total = scalar @files;

    $| = 1;

    if ( grep { /--reindex/ } @ARGV ) {
        System->RunQW( rm => $index_dir );
    }

    if ( !-d $normal_idx || !-d $dup_idx ) {

        my ($count, $total_files, $printed_init_label);

        $total_files += @old_files if !-d $normal_idx;
        $total_files += @dup_files if !-d $dup_idx;

        foreach my $row ( [ $normal_idx, \@old_files ],
            [ $dup_idx, \@dup_files ] )
        {
            my ( $index, $files ) = @$row;

            next if -d $index;

            foreach my $file (@$files) {
                printf "Indexing ...\n" if $printed_init_label++ == 1;;
                printf "Indexed photos %d of %d\n", ++$count, $total_files;
                _add_index $file, $index;
            }

        }

        printf "\n%s\n", '-=' x 40 if $printed_init_label;
    }

    my $count;

    foreach my $from_file (@files) {
        print "\n";
        printf "Sorting photo %d of %d ", ++$count, $total;

        if ( $from_file =~ /DS_Store/ ) {
            next;
        }

        if ( !-s $from_file ) {
            print "(empty file)";
            _move_file $from_file, _mkdir "$to/Empty-Files/";
            next;
        }

        my $info = ImageInfo $from_file;

        my $mime = $info->{MIMEType} // '';

        if ( !%$info || $mime !~ m/(image|video)/ ) {
            print "(non picture file)";
            _move_file $from_file, _mkdir "$to/Non-Picture/";
            next;
        }

        my $date = $info->{CreateDate} // $info->{FileModifyDate};

        $date = _date($date);

        if ( !$date ) {
            print "(file has no date)";
            _move_file $from_file, _mkdir "$to/Files-Have-No-Date/";
            next;
        }

        print "(Photo date: $date) ";

        my $dir;

        if ( _add_index($from_file, $normal_idx) eq 'seen' ) {
            if ( _add_index($from_file, $dup_idx) eq 'seen' ) {
                print "(seen this file before)";
                System->RunQW( rm => $from_file );
                next;
            }
            else {
                print "(duplicated)";
                $dir = $dup;
            }
        }
        else {
            print "New";
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

        print " (Filename: $date.$ext)";
        _move_file $from_file, $to_file;
    }

    map {
        system sprintf
            "cd %s; find . -name .DS_Store -exec rm {} \\;; find . -empty -type d -delete",
            quotemeta;
    } ( $from, $to );

    printf "\n\nDone @ %s\n", DateTime->now;
}

main
