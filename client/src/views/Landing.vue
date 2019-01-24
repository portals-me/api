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
          <v-btn color="indigo" @click="openDialog(0)" dark>アカウントを作成</v-btn>
          <p>あるいは <a @click="openDialog(1)">すでにアカウントを持っていますか？</a></p>
        </v-flex>

        <v-dialog
          v-model="dialog"
          max-width="600px"
        >
          <v-card>
            <v-tabs
              v-model="dialogTab"
            >
              <v-tab
                ripple
              >
                サインアップ
              </v-tab>
              <v-tab
                ripple
              >
                サインイン
              </v-tab>

              <v-tab-item>
                <v-stepper v-model="signupStep">
                  <v-stepper-content step="1">
                    <v-btn color="red" @click="authWithGoogle" dark>Googleでサインアップ</v-btn>
                    <br />
                    <v-btn color="blue" dark>Facebookでサインアップ</v-btn>
                    <br />
                    <v-btn color="grey" dark>GitHubでサインアップ</v-btn>
                  </v-stepper-content>

                  <v-stepper-content step="2">
                    <v-container>
                      <form>
                        <v-flex xs6>
                          <v-text-field
                            v-model="form.name"
                            label="ユーザーID"
                            xs6
                          />
                          <v-text-field
                            v-model="form.display_name"
                            label="表示される名前"
                          />
                        </v-flex>
                        <v-avatar color="orange" size="32px">
                          <v-img :src="form.photo" />
                        </v-avatar>
                        <v-btn depressed>アイコンをアップロード</v-btn>

                        <br />
                        <p>{{ this.errorMessage }}</p>

                        <v-btn color="success" @click="createAccount">送信</v-btn>
                        <v-btn depressed @click="signupStep --;">キャンセル</v-btn>
                      </form>
                    </v-container>
                  </v-stepper-content>
                </v-stepper>
              </v-tab-item>
              <v-tab-item>
                <v-container>
                  <v-btn color="red" @click="signInWithGoogle" dark>Googleでサインイン</v-btn>
                </v-container>
              </v-tab-item>
            </v-tabs>
          </v-card>
        </v-dialog>
      </v-container>
    </v-layout>
  </v-container>
</template>

<script>
import sdk from '@/app/sdk';

export default {
  data () {
    return {
      dialog: false,
      dialogTab: null,
      signupStep: 1,
      form: {
        name: '',
        display_name: '',
        photo: '',
      },
      google_token: '',
      errorMessage: '',
    };
  },
  methods: {
    openDialog (tabIndex) {
      this.dialog = true;
      this.dialogTab = tabIndex;
    },
    async authWithGoogle () {
      const user = await this.$gAuth.signIn();
      const profile = user.getBasicProfile();

      this.google_token = user.getAuthResponse().id_token;
      this.form = {
        name: profile.getId(),
        display_name: profile.getName(),
        photo: profile.getImageUrl(),
      };

      this.signupStep ++;
    },
    async createAccount () {
      try {
        const result = (await sdk.signUp({
          google_token: this.google_token,
          name: this.form.name,
          display_name: this.form.display_name,
          photo: this.form.photo,
        })).data;

        localStorage.setItem('id_token', result.id_token);
        localStorage.setItem('user', JSON.stringify(result.user));
        this.$router.push('/dashboard');
      } catch (err) {
        this.errorMessage = err.response.data;
        return;
      }
    },
    async signInWithGoogle () {
      const user = await this.$gAuth.signIn();

      try {
        const result = (await sdk.signIn(user.getAuthResponse().id_token)).data;

        localStorage.setItem('id_token', result.id_token);
        localStorage.setItem('user', result.user);
        this.$router.push('/dashboard');
      } catch (err) {
        console.error(err.response.data);
        return;
      }
    },
  }
}
</script>
