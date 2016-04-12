#!/usr/bin/perl

# 使い方
#
# $ cd fixtures/
# $ ls
# load.pl follows.csv tweets.csv users.csv
# $ perl load.pl

use strict;
use warnings;

load('users.csv',
    sub {
      my($uid, $screenName, $name) = @_;
      $uid = sprintf("%024x", $uid);
      my $CreatedAt = "2016-03-20T00:00:00+09:00";
      my $UpdatedAt = "2016-03-20T00:00:00+09:00";
      return
        "db.user.insert({_id: ObjectId(\"$uid\"), " .
        "screenName: \"$screenName\", " .
        "name: \"$name\", " .
        "createdAt: \"$CreatedAt\", " .
        "updatedAt: \"$UpdatedAt\"})";
    });

load('follows.csv',
    sub {
      my($userId, $targetId) = @_;
      return sprintf("db.follow.insert({userId: ObjectId(\"%024x\"), targetId: ObjectId(\"%024x\")})", $userId, $targetId);
    });

load('tweets.csv',
    sub {
      my($n, $uid, $text) = @_;
      my $createdAt = "2016-03-20T00:00:00+09:00";
      return sprintf("db.tweet.insert({userId:ObjectId(\"%024x\"), text:\"%s\", createdAt:\"%s\"})", $uid, $text, $createdAt);
    });

sub load {
    my($infile, $f) = @_;

    $| = 1;
    print "Loading $infile... ";

    open(my $ifh, '<', $infile);
    open(my $ofh, '| mongo 1>/dev/null');

    print {$ofh} "use poppo\n";

    while(defined(my $line = <$ifh>)) {
        chomp $line;
        my $opr = &$f(split(/,/, $line));
        print "$opr\n";
        print {$ofh} "$opr\n";
    }

    close($ifh);
    close($ofh);
}