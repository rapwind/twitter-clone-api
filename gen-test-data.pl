#!/usr/bin/perl

##
## Usage:
## $ perl gen-test-data.pl
##

use strict;
use warnings;

my $USERS = 100; # the number of users
my $TWEETS_PER_USER = 100;
my $FOLLOWS_PER_USER = 10;
my $LIKES_PER_USER = 100;

my $SCREENS = [qw(foo baz bar hoge fuga xyz abc)];
my $NAMES = [qw(鈴木 佐藤 田中 中村 山本 伊藤 小林 高橋 吉田 加藤 山田 松本 井上 山口 林 木村 佐々木 橋本 清水 山下 森 石川 阿部 齋藤
                池田 前田 岡田 渡辺 中島 村上 小川 藤田 長谷川 坂本 松田 渡邉 石井 後藤 青木 藤原 岡本 中川 近藤 西村 中野 藤井 中山
                遠藤 原田 太田 小野 上田 竹内 三浦 金子 田村 和田 森田 山崎 石田 福田 原 渡邊 宮本 内田 上野 酒井 安藤 柴田 武田 藤本
                今井 松尾 木下 大野 横山 野村 松井 谷口 工藤 丸山 河野 杉山 佐野 村田 久保 小山 小島 古川 増田 水野 野口 平野 新井
                髙橋 東 西田 山内 市川 櫻井 杉本 宮崎 千葉 松下 矢野 齊藤 北村 本田 大塚 吉村 渡部 安田 吉川 西川 松岡 辻 菊池 中田
                小松 飯田 川口 川上 菅原 平田 大西 樋口 田口 石原 福島 久保田 山中 岩田 土屋 吉岡 五十嵐 大久保 永井 熊谷 高木 森本
                菊地 岩本 堀 関 浅野 島田 永田 野田 服部 斉藤 高田 大石 黒田 馬場 大島 南 三上 西山 斎藤 松浦 岡 小池 松永 高野 篠原
                川村 菅野 大橋 大谷 小田 星野 石橋 秋山 中西 本間 平井 鎌田 内藤 松村 小西 荒木 白石 望月 栗原 澤田 福井 早川 松原
                上原 三宅 宮田 伊東 岩崎 大森 荒井 横田 内山 片山 中嶋 小澤)];
my $NOUN = [@$NAMES, qw(うさぎ カメ 馬)];
my $SITUATION = [qw(竹ぼうきで スカイリムで ポケモンで パソコンで ラケットで バールのようなもので 愛車の中で 愛馬に跨って 丸亀製麺で 製鉄所で 鉱山で 会社の中で 東京駅で 渋谷駅で)];
my $MODIFIER = [qw(勇敢に 元気に 少し 大きく 豪快に 繊細に 愉快に 優雅に エレガントに ゆっくりと 急いで モリモリと さっそうと あえて 断腸の思いで)];
my $VERB = [qw(戦う 遊ぶ 立ち上がる おどろく 考える 熟慮する ひらめく 思いつく 讃岐うどんを食べる 玉子焼きを作る お腹を叩く ゴミ箱を蹴る 殴られる ジャンプする)];
my $KAOMOJI = ["(*´ω｀*)", "(´・ω・｀)", "(*´∀｀)", "٩(ˊᗜˋ*)و", "╭( ･ㅂ･)و", "✌(՞ਊ՞✌三✌՞ਊ՞)✌", "(๑◔‿◔๑)", "~(  ~´･_･`)~", "ヽ(･ω･)ﾉ♡", "₍₍ (̨̡ ‾᷄⌂‾᷅)̧̢ ₎₎", "૮(꒦ິཅ꒦ິ)ა"];
my $IMAGES = [
    'https://img.esa.io/uploads/production/attachments/15/2016/04/15/9376/6d8296c6-b702-49f3-b655-777fce676214.png',
    'https://media.githubusercontent.com/media/hidetomo-watanabe/test/master/kannna1_lgtm.jpg',
    'http://livedoor.blogimg.jp/macky8162/imgs/b/9/b9a8ecce.jpg',
    'http://i.gyazo.com/9f18122c37fdfa3af5cc10a45df10404.png',
    'http://optipng.sourceforge.net/pngtech/img/lena.png'];

open(my $mongo, '| mongo 1>/dev/null');
# open(my $mongo, '>/dev/stdout');

print {$mongo} <<"__HEADER__";
use poppo
db.user.drop()
db.tweet.drop()
db.follow.drop()
db.like.drop()
__HEADER__

##
## Generate users & user tweets
##
my $users = {};
my $tweets = [];
foreach my $i (0 .. $USERS-1) {
    print "Generating user #$i and tweets...\n";

    my($uid, $screenName) = genUser($mongo, $i);
    $$users{$uid} = $screenName;
    genUserTweet($mongo, $uid, $users, $tweets);
}

##
## Generate follows
##
print "Generating follows...\n";
genFollows($mongo, $users);

##
## Generate likes
##
print "Generating likes...\n";
genLikes($mongo, $users, $tweets);

close($mongo);

exit;

sub genUser {
    my($mongo, $i) = @_;

    my $uid = "uid$i";
    my $name = mkUserName();
    my $screenName = mkScreenName($i);
    my $profileImageUrl = sprintf("http://www.beiz.jp/web/images_P/paint/xpaint_00%03d.jpg.pagespeed.ic.59qDQNa_Nn.webp", $i % 391);
    my $bib = mkSentence();
    my $createdAt = time2iso8061(time() - ($USERS - $i) * 3600);

    # password = 1234 (http://hayam.in/scrypt/)
    print {$mongo} <<"__USER__";
var uid$i = ObjectId()
db.user.insert({
  _id: $uid,
  name:"$name",
  screenName:"$screenName",
  email: "$screenName\@example.com",
  passwordHash:"7c1e081b170becf92e33fd001769afa73307a5c0889498671ccdf0ab0ff35646",
  profileImageUrl:"$profileImageUrl",
  biography:"$bib",
  url:"http://example.com/",
  createdAt:new Date("$createdAt"),
  updatedAt:new Date("$createdAt")
})
__USER__

    return ($uid, $screenName);
}

sub genUserTweet {
    my($mongo, $uid, $users, $tweets) = @_; # Mongo variable name that has a user ID

    foreach my $i (0 .. $TWEETS_PER_USER-1) {
        my $tid = "tid${i}_$uid";
        my $text = mkSentence();
        my $createdAt = time2iso8061(time() - ($TWEETS_PER_USER - $i) * 3600);

        my $contentUrl = "null"; # null in Mongo
        if (rand() < 0.2) { # Tweet with an image with prob. 0.2
            $contentUrl = '"' . choose($IMAGES) . '"';
        }

        if(scalar(@$tweets) > 0 && rand() < 0.5) { # Reply with prob. 0.5
            my $inReplyToTweetId = choose($tweets);
            my $inReplyToUserId = ( split(/_/, $inReplyToTweetId) )[1];
            my $inReplyToScreenName = $$users{$inReplyToUserId};

            print {$mongo} <<"__TWEET__";
var $tid = ObjectId()
db.tweet.insert({
  _id: $tid,
  text: "\@$inReplyToScreenName $text",
  contentUrl: $contentUrl,
  userId: $uid,
  inReplyToTweetId: $inReplyToTweetId,
  createdAt: new Date("$createdAt")
})
__TWEET__
        } else { # No reply
            print {$mongo} <<"__TWEET__";
var $tid = ObjectId()
db.tweet.insert({
  _id: $tid,
  text: "$text",
  contentUrl: $contentUrl,
  userId: $uid,
  createdAt: new Date("$createdAt")
})
__TWEET__
        }

        push(@$tweets, $tid);
    }
}

sub genFollows {
    my($mongo, $users) = @_;

    print {$mongo} <<"__HEADER__";
db.follow.ensureIndex({userId:1, _id:-1}, {unique:false, dropDups:true, background:true, sparse:true})
db.follow.ensureIndex({targetId:1, _id:-1}, {unique:false, dropDups:true, background:true, sparse:true})
db.follow.ensureIndex({userId:1, targetId:1}, {unique:true, dropDups:true, background:true, sparse:true})
__HEADER__

    foreach my $userId (sort keys %$users) {
        print "Generating follows for $userId...\n";

        my %fs = ();
        while (scalar(keys %fs) < $FOLLOWS_PER_USER) {
            my $targetId = choose([keys %$users]);
            my $key = $userId . $targetId;
            next if $userId eq $targetId || exists $fs{$key};

            my $createdAt = time2iso8061(time() - ($FOLLOWS_PER_USER - scalar(keys %fs)) * 3600);

            print {$mongo} <<"__FOLLOW__";
db.follow.insert({ _id: ObjectId(), userId: $userId, targetId: $targetId, createdAt: new Date("$createdAt") });
__FOLLOW__
            $fs{$key} = 1;
        }
    }
}

sub genLikes {
    my($mongo, $users, $tweets) = @_;

    print {$mongo} <<"__HEADER__";
db.like.ensureIndex({userId:1, _id:-1}, {unique:false, dropDups:true, background:true, sparse:true})
db.like.ensureIndex({tweetId:1, _id:-1}, {unique:false, dropDups:true, background:true, sparse:true})
db.like.ensureIndex({userId:1, tweetId:1}, {unique:true, dropDups:true, background:true, sparse:true})
__HEADER__

    foreach my $userId (sort keys %$users) {
        print "Generating likes for $userId...\n";

        my %ls = ();
        while (scalar(keys %ls) < $LIKES_PER_USER) {
            my $tweetId = choose($tweets);
            my $key = $userId . $tweetId;
            next if exists $ls{$key};

            my $createdAt = time2iso8061(time() - ($LIKES_PER_USER - scalar(keys %ls)) * 3600);

            print {$mongo} <<"__LIKE__";
db.like.insert({ _id: ObjectId(), userId: $userId, tweetId: $tweetId, createdAt: new Date("$createdAt") });
__LIKE__
            $ls{$key} = 1;
        }
    }
}

sub mkUserName {
    return choose($NAMES);
}

sub mkScreenName {
    my $i = shift;
    my $s = choose($SCREENS);
    return $s . '-' . $i;
}

sub mkSentence {
    my $s = choose($NOUN);
    my $n = rand(2);
    for (my $i=0; $i < $n; ++$i) {
        $s .= "と" . choose($NOUN);
    }
    $s .= "が" . choose($SITUATION) . choose($MODIFIER) . choose($VERB) . choose($KAOMOJI);
    return $s;
}

# Randomly select an element in a given array ref.
sub choose {
  my $a = shift;
  my $i = rand($#{$a} + 1);
  return $$a[$i];
}

sub time2iso8061 {
    my($sec, $min, $hour, $mday, $mon, $year) = gmtime(shift);
    return sprintf("%04d-%02d-%02dT%02d:%02d:%02d+09:00", $year+1900, $mon+1, $mday, $hour, $min, $sec);
}