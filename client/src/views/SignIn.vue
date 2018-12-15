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
import AWS from 'aws-sdk';
AWS.config.region = 'ap-northeast-1';

export default {
  components: {
    GSignIn,
  },
  methods: {
    async saveUser (user) {
      const userRef = firestore.collection('users').doc(user.uid);
      const userDoc = await userRef.get();
      const userData = {};

      if (!userDoc.exists) {
        userData.createdAt = firestore.FieldValue.serverTimestamp();
      }
      userRef.set(userData, { merge: true });
    },
    async onSignIn (googleUser) {
      console.log(googleUser.getAuthResponse().id_token);
    },
  },
  async mounted () {
  },
}
</script>

