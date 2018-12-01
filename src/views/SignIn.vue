<template>
  <v-container>
    <v-layout
      text-xs-center
      wrap
    >
      <v-flex mb-5 xs12>
        <v-btn @click="signInGoogle">
          Sign in with Google
        </v-btn>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
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
  },
}
</script>

