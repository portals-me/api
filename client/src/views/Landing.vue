<template>
  <v-container grid-list-xl>
    <v-layout row wrap text-xs-center>
      <v-flex xs4 offset-xs4>
        <h1>Portals@me</h1>
        <p>クリエイターの発信の場を提供します</p>
      </v-flex>
    </v-layout>

    <v-layout row wrap>
      <v-flex xs3>
        <v-icon x-large color="pink">share</v-icon>
        <h2>作品の共有</h2>
        <span>Portals@meから、あるいは他のサイトから作品を共有します。創作物はイラスト・文章・動画・ゲーム・ソフトウェアなど多岐にわたります。</span>
      </v-flex>

      <v-flex xs3>
        <v-icon x-large color="orange">insert_emoticon</v-icon>
        <h2>リアクション</h2>
        <span>素敵だと思った作品に対してはコメントやリアクションなどを付けてみましょう。作品への感想は新たな創作のきっかけとなります。</span>
      </v-flex>

      <v-flex xs3>
        <v-icon x-large color="teal">collections</v-icon>
        <h2>コレクション</h2>
        <span>クリエイターの作品群はコレクションとして整理され公開されます。クリエイターのマルチな創作を支援し気軽な発信を手助けします。</span>
      </v-flex>

      <v-flex xs3>
        <v-icon x-large color="indigo">public</v-icon>
        <h2>ソーシャル</h2>
        <span>他のクリエイターの活動を覗いてみましょう。新たな作品の発見だけでなく、皆がどんなことに興味を持っているかも分かるかもしれません。</span>
      </v-flex>
    </v-layout>

    <v-layout row wrap text-xs-center>
      <v-container>
        <v-flex xs12>
          <v-btn color="indigo" dark>Sign Up with Google</v-btn>
          <p>あるいは <router-link to="/signin">すでにアカウントを持っています</router-link></p>
        </v-flex>
      </v-container>
    </v-layout>
  </v-container>
</template>

<script>
export default {
}
</script>

