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
import firebase,{ auth } from 'firebase/app';
import 'firebase/auth';
import firebaseui from 'firebaseui';
import 'firebaseui/dist/firebaseui.css';
import firestore from '@/instance/firestore';

export default {
  name: 'signin',
  methods: {
    signInGoogle () {
      const provider = new firebase.auth.GoogleAuthProvider();

      firebase.auth().signInWithPopup(provider).then((result) => {
        const token = result.credential.accessToken;
        const user = result.user;

        console.log(token, user);
      });
    },
    async saveUser (user) {
      const userRef = firestore.collection('users').doc(user.uid);
      const userDoc = await userRef.get();
      const userData = {};

      if (!userDoc.exists) {
        userData.createdAt = firebase.firestore.FieldValue.serverTimestamp();
      }
      userRef.set(userData, { merge: true });
    },
  },
  mounted () {
    let ui = firebaseui.auth.AuthUI.getInstance();
    if (!ui) {
      ui = new firebaseui.auth.AuthUI(firebase.auth());
    }

    ui.start('#firebaseui-auth-container', {
      callbacks: {
        signInSuccessWithAuthResult: (authResult, redirectUrl) => {
          this.saveUser(authResult.user).then(() => {
            this.$router.push('/');
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
}
</script>

