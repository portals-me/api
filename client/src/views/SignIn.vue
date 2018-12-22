<template>
  <v-container>
    <v-layout
      text-xs-center
      wrap
    >
      <v-flex mb-5 xs12>
        <v-btn @click="signIn">Sign In with Google</v-btn>
      </v-flex>
    </v-layout>
  </v-container>
</template>

<script>
import sdk from '@/app/sdk';

export default {
  methods: {
    async signIn () {
      const user = await this.$gAuth.signIn();
      const result = await sdk.signIn(user.getAuthResponse().id_token);
      localStorage.setItem('id_token', result.id_token);
      this.$store.commit('setUser', result.user);
      this.$router.push('/');
    },
  },
  async mounted () {
  },
}
</script>

