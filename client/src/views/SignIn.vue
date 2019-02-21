<template>
  <v-layout justify-center>
    <v-flex xs4>
      <v-tabs
        v-model="tab"
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
          <p style="color: red">{{ signUpError }}</p>
          <v-stepper v-model="signUpStep">
            <v-stepper-content step="1">
              <v-btn color="red" @click="signUpWithGoogle" dark>Googleでアカウント作成</v-btn>
              <br />
              <v-btn color="light-blue" @click="signUpWithTwitter" dark>Twitterでアカウント作成</v-btn>
            </v-stepper-content>

            <v-stepper-content step="2">
              <v-container>
                <form>
                  <v-flex>
                    <v-text-field
                      v-model="form.name"
                      label="ユーザーID"
                      append-outer-icon="check"
                    />
                    <v-text-field
                      v-model="form.display_name"
                      label="表示される名前"
                    />
                  </v-flex>
                  <v-avatar color="orange" size="32px">
                    <v-img :src="form.picture" />
                  </v-avatar>
                  <v-btn depressed>アイコンをアップロード</v-btn>

                  <br />

                  <v-btn color="success" @click="signUp">送信</v-btn>
                  <v-btn depressed @click="signUpStep --;">キャンセル</v-btn>
                </form>
              </v-container>
            </v-stepper-content>
          </v-stepper>
        </v-tab-item>
        <v-tab-item>
          <p style="color: red">{{ signInError }}</p>
          <v-btn color="red" dark @click="signInWithGoogle">Googleでログイン</v-btn>
          <br />
          <v-btn color="light-blue" dark @click="signInWithTwitter">Twitterでログイン</v-btn>
        </v-tab-item>
      </v-tabs>
    </v-flex>
  </v-layout>
</template>

<script>
import axios from 'axios';
import sdk from '@/app/sdk';

export default {
  props: [ 'signin' ],
  data () {
    return {
      form: {
        name: '',
        display_name: '',
        picture: '',
      },
      logins: {},
      tab: 0,
      signUpStep: 1,
      signInError: '',
      signUpError: '',
    };
  },
  methods: {
    async signInWithGoogle () {
      const user = await this.$gAuth.signIn();

      try {
        const result = (await sdk.signIn({
          google: user.getAuthResponse().id_token
        })).data;
        await this.signIn(result);
      } catch (err) {
        this.signInError = err.response.data;
        return;
      }
    },
    async signInWithTwitter () {
      const twitterAuthURL = await axios.post(`${process.env.VUE_APP_API_ENDPOINT}/auth/twitter`);
      location.href = twitterAuthURL.data;

      // Jump to mounted.twitter-callback
    },
    async signInWithTwitterAfter (token) {
      try {
        const credential = (await axios.get(`${process.env.VUE_APP_API_ENDPOINT}/auth/twitter?oauth_token=${token.oauth_token}&oauth_verifier=${token.oauth_verifier}`)).data.credential;
        const result = (await sdk.signIn({
          twitter: credential,
        })).data;
        await this.toDashboard(result);
      } catch (err) {
        this.signInError = 'LoginError';
        this.$router.push('/signin');
        return;
      }
    },
    async toDashboard ({ id_token, user }) {
      localStorage.setItem('id_token', id_token);
      localStorage.setItem('user', user);
      this.$router.push('/dashboard');
    },
    async signUpWithGoogle () {
      const user = await this.$gAuth.signIn();
      const profile = user.getBasicProfile();

      this.form = {
        name: profile.getId(),
        display_name: profile.getName(),
        picture: profile.getImageUrl(),
      };
      this.logins = {
        google: user.getAuthResponse().id_token
      };

      this.signUpStep ++;
    },
    async signUpWithTwitter () {
      const twitterAuthURL = await axios.post(`${process.env.VUE_APP_API_ENDPOINT}/auth/twitter`);
      location.href = twitterAuthURL.data;
    },
    async signUpWithTwitterAfter (token) {
      const result = (await axios.get(`${process.env.VUE_APP_API_ENDPOINT}/auth/twitter?oauth_token=${token.oauth_token}&oauth_verifier=${token.oauth_verifier}`)).data;
      const account = result.account;
      console.log(account)
      console.log(token)

      this.form = {
        name: account.screen_name,
        display_name: account.screen_name,
        picture: account.profile_image_url,
      }
      this.logins = {
        twitter: result.credential,
      };

      this.signUpStep ++;
    },
    async signUp () {
      try {
        const result = (await sdk.signUp({
          form: this.form,
          logins: this.logins,
        })).data;
        await this.toDashboard(result);
      } catch (err) {
        this.signUpError = 'SignUpError';
        this.$router.push('/signup');
        this.signUpStep = 1;
        return;
      }
    }
  },
  mounted () {
    if (this.$route.path.startsWith('/signin')) {
      this.tab = 1;
    }

    // twitter-callback
    // Continue to signInWithTwitter
    if (this.$route.path.startsWith('/signin/twitter-callback')) {
      this.signInWithTwitterAfter(this.$route.query);
    }

    if (this.$route.path.startsWith('/signup/twitter-callback')) {
      this.signUpWithTwitterAfter(this.$route.query);
    }
},
}
</script>

