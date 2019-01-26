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
          <v-btn color="red" dark>Googleでアカウント作成</v-btn>
          <br />
          <v-btn color="light-blue" dark>Twitterでアカウント作成</v-btn>
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
      tab: 0,
      signInError: '',
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
        console.error(err.response.data);
        return;
      }
    },
    async signInWithTwitter () {
      const twitterAuthURL = await axios.post('https://ibsrd4lyxk.execute-api.ap-northeast-1.amazonaws.com/dev/auth/twitter');
      location.href = twitterAuthURL.data;

      // Jump to mounted.twitter-callback
    },
    async signInWithTwitterAfter (token) {
      try {
        const result = (await sdk.signIn({
          twitter: `${token.oauth_token}.${token.oauth_verifier}`,
        })).data;
        await this.signIn(result);
      } catch (err) {
        this.signInError = 'LoginError';
        this.$router.push('/signin');
        return;
      }
    },
    async signIn ({ id_token, user }) {
      localStorage.setItem('id_token', id_token);
      localStorage.setItem('user', user);
      this.$router.push('/dashboard');
    },
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
  },
}
</script>

