<template>
  <v-container>
    <v-layout
      text-xs-center
      wrap
    >
      <v-flex mb-5 xs12>
        <div id="firebaseui-auth-container" />
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import 'firebaseui/dist/firebaseui.css';
import firebase from 'firebase';
import firebaseui from 'firebaseui';
import firestore from '@/instance/firestore';

export default {
  name: 'signin',
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
    async onMount () {
      let ui = firebaseui.auth.AuthUI.getInstance();
      if (!ui) {
        ui = new firebaseui.auth.AuthUI(firebase.auth());
      }

      ui.start('#firebaseui-auth-container', {
        callbacks: {
          signInSuccessWithAuthResult: (authResult, redirectUrl) => {
            this.saveUser(authResult.user).then(() => {
              this.$router.push(this.$route.query.redirect ? this.$route.query.redirect : '/');
            });
            return false;
          },
        },
        signInFlow: 'redirect',
        signInOptions: [
          firebase.auth.GoogleAuthProvider.PROVIDER_ID,
        ],
      });
    },
  },
  async mounted () {
    await this.onMount();
  },
}
</script>

