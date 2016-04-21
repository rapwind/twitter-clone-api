#!/usr/bin/perl

##
## Usage:
##  perl gen-test-data.pl
##  I="mongoimport -h ds023560.mlab.com:23560 -u USERNAME -p PASSWORD -d poppo-mongo-dev --type json"
##  $I -c user --file users.json
##  $I -c tweet --file tweets.json
##  $I -c follow --file follows.json
##  $I -c like --file likes.json

use utf8;
use strict;
use warnings;
use BSON;

##
## Settings
##

my $USERS = 400000; # the number of users
my $MAX_TWEETS = 99; # the max number of tweets per a user
my $MIN_TWEETS = 0; # the min number of tweets per a user
my $MAX_FOLLOWS = 99; # the max number of following users per a user
my $MIN_FOLLOWS = 0; # the min number of following users per a user
my $MAX_LIKES = 99; # the max number of liked tweets per a user
my $MIN_LIKES = 0; # the min number of liked tweets per a user
my $MAX_TIME = time();
my $MIN_TIME = $MAX_TIME - 5 * 365 * 24 * 60 * 60; # 5 years ago

##
## Source data
##

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
my $SITUATION = [qw(竹ぼうきで スカイリムで ポケモンで パソコンで ラケットで フォークで スプーンで バールのようなもので
                    愛車の中で 愛馬に跨って 丸亀製麺で ジェットコースターに乗りながら 人混みの中で
                    製鉄所で 鉱山で 会社の中で 東京駅で 渋谷駅で バス停で スタバで カフェで)];
my $MODIFIER = [qw(勇敢に 元気に 少し 大きく 豪快に 繊細に 愉快に 優雅に エレガントに ゆっくりと 急いで モリモリと さっそうと あえて
                   断腸の思いで まるっと グイッと)];
my $VERB = [qw(遊んだ 立ち上がるぞ 驚きを隠せない 考えまーす 感じた ひらめいた 思いついた 讃岐うどんを食べた 玉子焼きを作ったなう
               腹パンなう ゴミ箱を蹴った 殴られた ジャンプしてる〜)];
my $KAOMOJI = ["(*´ω｀*)", "(´・ω・｀)", "(*´∀｀)", "٩(ˊᗜˋ*)و", "╭( ･ㅂ･)و", "✌(՞ਊ՞✌三✌՞ਊ՞)✌", "(๑◔‿◔๑)", "~(  ~´･_･`)~", "ヽ(･ω･)ﾉ♡", "₍₍ (̨̡ ‾᷄⌂‾᷅)̧̢ ₎₎", "૮(꒦ິཅ꒦ິ)ა"];
my $IMAGES = [
    'https://img.esa.io/uploads/production/attachments/15/2016/04/15/9376/6d8296c6-b702-49f3-b655-777fce676214.png',
    'https://media.githubusercontent.com/media/hidetomo-watanabe/test/master/kannna1_lgtm.jpg',
    'http://livedoor.blogimg.jp/macky8162/imgs/b/9/b9a8ecce.jpg',
    'http://i.gyazo.com/9f18122c37fdfa3af5cc10a45df10404.png',
    'http://optipng.sourceforge.net/pngtech/img/lena.png'];

$| = 1;
my $users = genUsers('users.json');
my $userTweets = genTweets('tweets.json', $users);
genFollows('follows.json', $users);
genLikes('likes.json', $users, $userTweets);

exit;

sub genUsers {
    my $fname = shift;
    my %users;

    open(my $fh, '>:utf8', $fname) or die "Cannot open $fname: $!";

    foreach my $i (0 .. $USERS-1) {
        my $userId = BSON::ObjectId->new;
        my $screenName = mkScreenName($i);
        my $bgUrl = sprintf("http://www.beiz.jp/web/images_P/paint/xpaint_00%03d.jpg.pagespeed.ic.59qDQNa_Nn.webp", $i % 391);
        my $imageUrl = sprintf("http://pokemon.symphonic-net.com/%03d.gif", $i % 650);
        my $createdAt = genTime();

        $users{$userId} = $screenName;

        printf "%s: Generating users ... %.2f %%\n", ptime(), (100.0 * $i / $USERS) if $i % 10000 == 0;

        print {$fh} JSON(
            '_id'                       => ['objectId', $userId],
            'name'                      => ['string', mkUserName()],
            'screenName'                => ['string', $screenName],
            'passwordHash'              => ['string', '7c1e081b170becf92e33fd001769afa73307a5c0889498671ccdf0ab0ff35646'], # password = 1234
            'email'                     => ['string', $screenName . '@example.com'],
            'profileImageUrl'           => ['string', $imageUrl],
            'profileBackgroundImageUrl' => ['string', $bgUrl],
            'url'                       => ['string', 'http://example.com/'],
            'biography'                 => ['string', mkSentence()],
            'createdAt'                 => ['date', $createdAt],
            'updatedAt'                 => ['date', $createdAt]
        ), "\n";
    }

    close($fh);

    return \%users;
}

sub genTweets {
    my($fname, $users) = @_;
    my @uids = keys %$users;
    my @userTweets;
    #my %tweets;

    $#userTweets = $#uids;

    open(my $fh, '>:utf8', $fname) or die "Cannot open $fname: $!";

    foreach my $i (0 .. $USERS-1) {
        my $userId = $uids[$i];
        my $n = int(rand($MAX_TWEETS - $MIN_TWEETS) + $MIN_TWEETS);

        $userTweets[$i] = $n;

        printf "%s: Generating tweets %.2f %%\n", ptime(), (100.0*$i/$USERS) if $i % 1000 == 0;

        foreach my $j (0 .. $n-1) {
            my $id = getTweetId($i, $j);
            my $text = mkSentence();
            my $createdAt = genTime();

            die "over 140 characters" if length($text) > 140;

            my $contentUrl;
            $contentUrl = choose($IMAGES) if rand() < 0.2; # Tweet+image with prob. 0.2

            my $inReplyToTweetId;
            if (($i != 0 || $j != 0) && rand() < 0.5) { # Reply with prob. 0.5
                my $ui = int(rand($i+1));
                my $ti;
                if ($ui == $i) {
                    $ti = int(rand($j));
                } else {
                    $ti = int(rand($userTweets[$ui]));
                }
                $inReplyToTweetId = getTweetId($ui, $ti);
                my $inReplyToUserId = $uids[$ui];
                my $inReplyToScreenName = $$users{$inReplyToUserId};
                $text = "\@$inReplyToScreenName $text";
            }

            print {$fh} JSON(
                '_id'              => ['objectId', $id],
                'text'             => ['string', $text],
                'contentUrl'       => ['string', $contentUrl],
                'userId'           => ['objectId', $userId],
                'inReplyToTweetId' => ['objectId', $inReplyToTweetId],
                'createdAt'        => ['date', $createdAt]
            ), "\n";
        }
    }

    close($fh);

    return \@userTweets;
}

sub genFollows {
    my($fname, $users) = @_;
    my $uids = [keys %$users];

    open(my $fh, '>:utf8', $fname) or die "Cannot open $fname: $!";

    foreach my $i (0 .. $USERS-1) {
        my $userId = $$uids[$i];
        my $n = int(rand($MAX_FOLLOWS - $MIN_FOLLOWS) + $MIN_FOLLOWS);
        $n = $#{$uids} if $n > $#{$uids};

        printf "%s: Generating folllows %.2f %%\n", ptime(), (100.0*$i/$USERS) if $i % 1000 == 0;

        my(%seen, $count);
        $count = 0;
        while ($count < $n) {
             my $targetId = choose($uids);
             my $key = $userId . $targetId;
             next if $userId eq $targetId || exists $seen{$key};
             $seen{$key} = 1;
             ++$count;

             print {$fh} JSON(
                '_id'       => ['objectId', BSON::ObjectId->new],
                'userId'    => ['objectId', $userId],
                'targetId'  => ['objectId', $targetId],
                'createdAt' => ['date', genTime()]
             ), "\n";
        }
    }

    close($fh);
}

sub genLikes {
    my($fname, $users, $userTweets) = @_;
    my $uids = [keys %$users];

    open(my $fh, '>:utf8', $fname) or die "Cannot open $fname: $!";

    foreach my $i (0 .. $USERS-1) {
        my $userId = $$uids[$i];
        my $n = int(rand($MAX_LIKES - $MIN_LIKES) + $MIN_LIKES);

        printf "%s: Generating likes %.2f %%\n", ptime(), (100.0*$i/$USERS) if $i % 1000 == 0;

        my(%seen, $count);
        $count = 0;
        while ($count < $n) {
            my $ui = int(rand($USERS));
            my $ti = int(rand($$userTweets[$ui]));
            my $tweetId = getTweetId($ui, $ti);

            my $key = $userId . $tweetId;
            next if exists $seen{$key};
            $seen{$key} = 1;
            ++$count;

            print {$fh} JSON(
               '_id'       => ['objectId', BSON::ObjectId->new],
               'userId'    => ['objectId', $userId],
               'tweetId'   => ['objectId', $tweetId],
               'createdAt' => ['date', genTime()]
            ), "\n";
        }
    }

    close($fh);
}

sub genTime {
    my $t = int(rand($MAX_TIME - $MIN_TIME) + $MIN_TIME);
    return time2iso8061($t);
}

sub JSON {
    my %data = @_;
    my @strs;
    foreach my $key (sort keys %data) {
        my($type, $val) = @{$data{$key}};
        if (defined $val) {
            if ($type eq 'objectId') {
                push(@strs, "$key:ObjectId(\"$val\")");
            } elsif ($type eq 'date') {
                push(@strs, "$key:ISODate(\"$val\")");
            } elsif ($type eq 'string') {
                push(@strs, "$key:\"$val\"");
            } else {
                die "Unknown type $type";
            }
        }
    }
    return '{' . join(',', @strs) . '}';
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

sub getTweetId {
    my($ui, $ti) = @_;
    my $id = ($ti + 1) * $USERS + $ui;
    return sprintf('%024x', $id);
}

sub ptime {
    my $t = time() - $^T;
    my $sec = $t % 60;
    my $min = int($t / 60) % 60;
    my $hour = int($t / 3600);
    return sprintf('%d:%02d:%02d', $hour, $min, $sec);
}