<template>
  <v-container>
    <v-layout
      text-xs-center
      wrap
    >
      <v-flex mb-5 xs12>
        <g-sign-in @signin="onSignIn" />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import 'firebaseui/dist/firebaseui.css';
import firebase from 'firebase';
import firebaseui from 'firebaseui';
import firestore from '@/instance/firestore';
import GSignIn from '@/components/GSignIn';
import sdk from '@/app/sdk';
import AWS from 'aws-sdk';
AWS.config.region = 'ap-northeast-1';

export default {
  components: {
    GSignIn,
  },
  methods: {
    async onSignIn (googleUser) {
      console.log(googleUser.getAuthResponse().id_token);
      console.log(await sdk.signIn(googleUser.getAuthResponse().id_token));
    },
  },
  async mounted () {
  },
}
</script>

