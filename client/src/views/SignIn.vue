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
import GSignIn from '@/components/GSignIn';
import sdk from '@/app/sdk';

export default {
  components: {
    GSignIn,
  },
  methods: {
    async onSignIn (googleUser) {
      const result = await sdk.signIn(googleUser.getAuthResponse().id_token);
      localStorage.setItem('id_token', result.id_token);
      this.$store.commit('setUser', result.user);
      this.$router.push('/');
    },
  },
  async mounted () {
  },
}
</script>

