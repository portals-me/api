<template>
  <v-layout justify-center>
    <v-flex xs4>
      <v-card>
        <v-card-title class="justify-center" primary-title>
          <h2>Sign In</h2>
        </v-card-title>
        <v-card-text>
          <div>
            <p style="color: red">{{ signInError }}</p>
            <v-btn color="red" dark @click="signInWithGoogle">Googleでログイン</v-btn>
            <br />
            <v-btn color="light-blue" dark @click="signInWithTwitter">Twitterでログイン</v-btn>
          </div>
          <br />
          <router-link to="/signup">アカウントを持っていない場合はこちら</router-link>
        </v-card-text>
      </v-card>
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
      localStorage.setItem('user', JSON.stringify(user));
      this.$router.push('/dashboard');
    }
  },
  mounted () {
    // twitter-callback
    // Continue to signInWithTwitter
    if (this.$route.path.startsWith('/signin/twitter-callback')) {
      this.signInWithTwitterAfter(this.$route.query);
    }
  },
}
</script>

