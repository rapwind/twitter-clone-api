#!/usr/bin/perl

# 使い方
#
# $ cd fixtures/
# $ ls
# load.pl follows.csv tweets.csv users.csv
# $ perl load.pl

use strict;
use warnings;

load('users.csv', 'users.tsv',
    'user',
    '_id,screenName,name,profileImageUrl,profileBackgroundImageUrl,biography,locationText,Url,Birthday,createdAt,updatedAt',
    sub {
      my($uid, $screenName, $name) = @_;
      $uid = sprintf("%024x", $uid);
      my $ProfileImageURL = "";
    my $ProfileBackgroundImageURL = "";
      my $Biography = "";
      my $LocationText = "";
      my $URL = "";
      my $Birthday = "";
      my $CreatedAt = "2016-03-20T00:00:00+09:00";
      my $UpdatedAt = "2016-03-20T00:00:00+09:00";
      return "ObjectId($uid)\t$screenName\t$name\t$ProfileImageURL\t$ProfileBackgroundImageURL\t$Biography\t$LocationText\t$URL\t$Birthday\t$CreatedAt\t$UpdatedAt\n";
    });

load('follows.csv', 'follows.tsv',
    'follow',
    'userId,targetId',
    sub {
      my($userId, $targetId) = @_;
      return sprintf("ObjectId(%024x)\tObjectId(%024x)\n", $userId, $targetId);
    });

load('tweets.csv', 'tweets.tsv',
    'tweet',
    'userId,text,createdAt', # inReplyToUserId,inReplyToTweetId,
    sub {
      my($n, $uid, $text) = @_;
      my $createdAt = "2016-03-20T00:00:00+09:00";
      return sprintf("ObjectId(%024x)\t%s\t%s\n", $uid, $text, $createdAt);
    });

sub load {
    my($infile, $outfile, $coll, $fields, $f) = @_;

    $| = 1;
    print "Loading $infile... ";

    open(my $ifh, '<', $infile);
    open(my $ofh, '>', $outfile);

    my $n = 0;
    while(defined(my $line = <$ifh>)) {
        chomp $line;
        print {$ofh} &$f(split(/,/, $line));
        ++$n;
        if ($n % 100000 == 0) { print "[$n]"; }
    }

    print "\n";

    close($ifh);
    close($ofh);

    system("mongoimport --db poppo --collection $coll --type tsv --file $outfile -f $fields");
    unlink($outfile);
}